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

// GetUser gets a single user
// @Summary      Get User Profile
// @Description  Get details of a specific user.
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        user_id path int true "User ID"
// @Success      200  {object}  models.User
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Failure      404  {object}  apperrors.ErrorResponse "Not Found"
// @Router       /users/{user_id} [get]
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

// GetUsers lists all users
// @Summary      Get All Users
// @Description  Admin only. Lists all registered users.
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /users [get]
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

// UpdateUser updates profile
// @Summary      Update User
// @Description  Update name, avatar, or password.
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user_id path int true "User ID"
// @Param        input body models.UpdateRequest true "Update Info"
// @Success      200  {object}  models.User
// @Failure      400  {object}  apperrors.ErrorResponse "Bad Request"
// @Router       /users/{user_id} [patch]
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

// Signup creates a new user
// @Summary      Register User
// @Description  Creates a new user account and sends a tmp password (single use) to the user's email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        name&email body models.SignUpRequest true "Name and Email of the user"
// @Success      201  {object}  models.User
// @Failure      400  {object}  apperrors.ErrorResponse "Bad Request"
// @Router       /signup [post]
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

// Login logs in a user
// @Summary      User Login
// @Description  Authenticates a user and signs a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        email&password body models.LoginRequest true "Login Credentials"
// @Success      200  {object}  models.LoginResponse
// @Failure      400  {object}  apperrors.ErrorResponse "Bad Request"
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Router       /login [post]
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

// DeleteUser deletes a user account
// @Summary      Delete User Account
// @Description  Permanently remove a user account. Users can delete themselves; Admins can delete anyone.
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        user_id path int true "User ID to delete"
// @Success      200  {object}  models.User "Returns the deleted user data"
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden - Not your account"
// @Failure      404  {object}  apperrors.ErrorResponse "User not found"
// @Router       /users/{user_id} [delete]
func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	targetUserID, err := app.ExtractPathValue(r, "user_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	requestorID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if targetUserID != requestorID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, fmt.Errorf("%w: you can only delete your own account", apperrors.ErrForbidden))
		return
	}
	deletedUser, err := app.DB.DeleteUser(targetUserID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error deleting user %d: %w", targetUserID, err))
		return
	}
	app.Encode(w, r, deletedUser)
}

// SetLanguage sets the user's preferred language in a cookie
// @Summary      Set Interface Language
// @Description  Saves the user's preferred language (key="lang") in a cookie for 30 days.
// @Tags         settings
// @Accept       json
// @Produce      json
// @Param        input body models.SetLang true "Language (e.g. ru, en)"
// @Success      201  {object}  models.SetLang "New Language"
// @Failure      400  {object}  apperrors.ErrorResponse "Invalid language code"
// @Router       /language [post]
func (app *App) SetLanguage(w http.ResponseWriter, r *http.Request) {
	var req models.SetLang
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding request: %w", apperrors.ErrBadRequest, err))
		return
	}
	if req.Lang != "ru" && req.Lang != "en" {
		app.ErrorJSON(w, r, fmt.Errorf("%w: unsupported language '%s'", apperrors.ErrBadRequest, req.Lang))
		return
	}
	cookie := &http.Cookie{
		Name:     "lang",
		Value:    req.Lang,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	app.Encode(w, r, models.SetLang{
		Lang: req.Lang,
	})
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
