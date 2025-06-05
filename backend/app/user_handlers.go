package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	if userID != ID || role != models.RoleAdmin {
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
	user.Password = ""
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

func (app *App) GetUserByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	name = strings.TrimSpace(name)
	if name == "" {
		app.ErrorJSON(w, r, fmt.Errorf("%w: empty user name", apperrors.ErrBadRequest))
		return
	}
	user, err := app.DB.GetUserByName(name)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting user(%s): %w", name, err))
		return
	}
	user.Password = ""
	app.Encode(w, r, user)
}

func (app *App) Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding user: %w", apperrors.ErrBadRequest, err))
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error hashing password: %w", err))
		return
	}
	user.Password = string(hashedPass)
	newUser, err := app.DB.CreateUser(user)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error creating user: %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	newUser.Password = ""
	app.Encode(w, r, newUser)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding user: %w", apperrors.ErrBadRequest, err))
		return
	}
	existingUser, err := app.DB.GetUserByName(user.Name)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			app.ErrorJSON(w, r, fmt.Errorf("%w: %w", apperrors.ErrInvalidCredentials, err))
			return
		}
		app.ErrorJSON(w, r, fmt.Errorf("error getting user from database: %w", err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: %w", apperrors.ErrInvalidCredentials, err))
		return
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
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	s, err := t.SignedString(app.jwtSecret)
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
	existingUser.Password = ""
	app.Encode(w, r, struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}{Token: s, User: existingUser})
}

func createContextWithClaims(userID string, role models.Role) context.Context {
	claims := Claims{
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
