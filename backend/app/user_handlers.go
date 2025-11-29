package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Role models.Role `json:"role"`
	jwt.RegisteredClaims
}

func (app *App) GetUser(w http.ResponseWriter, r *http.Request) {
	ID, err := app.ExtractPathValue(r, "user_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if userID != ID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
		return
	}
	user, err := app.DB.GetUser(ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting user %d: %w", ID, err))
		return
	}
	cart, err := app.DB.GetCartByUserID(user.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart: %w", err))
		return
	}
	cartItems, err := app.DB.GetCartItems(cart.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart items: %w", err))
		return
	}
	cart.Items = cartItems
	user.Cart = cart
	app.Encode(w, r, user)
}

func (app *App) GetUsers(w http.ResponseWriter, r *http.Request) {
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if role != models.RoleAdmin {
		app.ErrorJSON(w, r, fmt.Errorf("%w: user %d (role %s) does not have required role",
			apperrors.ErrForbidden, userID, role))
		return
	}
	users, err := app.DB.GetUsers()
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting users: %w", err))
		return
	}
	app.Encode(w, r, users)
}

func (app *App) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	email = strings.TrimSpace(email)
	if email == "" {
		app.ErrorJSON(w, r, fmt.Errorf("%w: empty user email", apperrors.ErrBadRequest))
		return
	}
	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting user(%s): %w", email, err))
		return
	}
	app.Encode(w, r, user)
}

func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := app.ExtractPathValue(r, "user_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	updaterID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if role != models.RoleAdmin && userID != updaterID {
		app.ErrorJSON(w, r, fmt.Errorf("%w: user %d (role %s) does not have required role",
			apperrors.ErrForbidden, userID, role))
		return
	}
	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error decoding request body: %w", err))
		return
	}
	if len(updates) == 0 {
		app.ErrorJSON(w, r, fmt.Errorf("%w: no fields to update", apperrors.ErrBadRequest))
		return
	}
	allowedFields := map[string]bool{
		"name":     true,
		"avatar":   true,
		"password": true,
	}
	for field := range updates {
		if !allowedFields[field] {
			app.ErrorJSON(w, r, fmt.Errorf("field %s is not allowed", field))
			return
		}
	}
	pass, ok := updates["password"]
	if ok {
		pass, ok := pass.(string)
		if !ok {
			app.ErrorJSON(w, r, fmt.Errorf("%w: password is not of type string", apperrors.ErrBadRequest))
			return
		}
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			app.ErrorJSON(w, r, fmt.Errorf("error hashing password: %w", err))
			return
		}
		updates["password"] = string(hashedPass)
	}
	updatedUser, err := app.DB.UpdateUser(userID, updates)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error hashing updating user in database: %w", err))
		return
	}
	app.Encode(w, r, updatedUser)
}

func (app *App) Signup(w http.ResponseWriter, r *http.Request) {
	var signup models.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding signup credentials: %w", apperrors.ErrBadRequest, err))
		return
	}
	signup.Email = strings.TrimSpace(signup.Email)
	_, err = mail.ParseAddress(signup.Email)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: incorrect email address: %w", apperrors.ErrBadRequest, err))
		return
	}
	user := models.User{Name: signup.Name, Email: signup.Email}
	newUser, err := app.DB.CreateUser(user)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error creating user: %w", err))
		return
	}
	pass, err := GenerateTmpPassword()
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error generating tmp password: %w", err))
		return
	}
	pass.Email = signup.Email
	_, err = app.DB.DeleteTmpPassword(pass.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		app.ErrorJSON(w, r, fmt.Errorf("error deleting tmp password from database: %w", err))
		return
	}
	pass, err = app.DB.CreateTmpPassword(pass)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error saving tmp password to database: %w", err))
		return
	}
	err = SendTmpPassword(pass)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error sending tmp password to email: %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	app.Encode(w, r, newUser)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var login models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding login credentials: %w", apperrors.ErrBadRequest, err))
		return
	}
	existingUser, err := app.DB.GetUserByEmail(login.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			app.ErrorJSON(w, r, fmt.Errorf("%w: %w", apperrors.ErrInvalidCredentials, err))
			return
		}
		app.ErrorJSON(w, r, fmt.Errorf("error getting user from database: %w", err))
		return
	}
	loggedIn := false
	if existingUser.Password != nil {
		err = bcrypt.CompareHashAndPassword([]byte(*existingUser.Password), []byte(login.Password))
		if err == nil {
			loggedIn = true
		}
	}
	if !loggedIn {
		err = app.LoginWithTmpPassword(login)
		if err != nil {
			app.ErrorJSON(w, r, fmt.Errorf("%w: invalid credentials: %w", apperrors.ErrInvalidCredentials, err))
			return
		}
	}
	claims := Claims{
		Role: existingUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(existingUser.ID),
			Issuer:    "Smartbuy",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	jwt, err := token.SignedString(app.jwtSecret)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error signing token: %w", err))
		return
	}
	cart, err := app.DB.GetCartByUserID(existingUser.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart: %w", err))
		return
	}
	cartItems, err := app.DB.GetCartItems(cart.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart items: %w", err))
		return
	}
	cart.Items = cartItems
	existingUser.Cart = cart
	loginResponse := models.LoginResponse{User: existingUser, Token: jwt}
	app.Encode(w, r, loginResponse)
}

func createContextWithClaims(userID string, role models.Role) context.Context {
	claims := &Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "Smartbuy",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	return context.WithValue(context.Background(), ClaimsKey, claims)
}
