package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/logger"
	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/sfu-teamproject/smartbuy/backend/storage"
)

var (
	errBadRequest         = errors.New("bad request")
	errInternal           = errors.New("internal server error")
	errInvalidCredentials = errors.New("invalid credentials")
	errUnauthorized       = errors.New("unauthorized")
	errForbidden          = errors.New("forbidden")
)

type App struct {
	Log       logger.Logger
	Server    *http.Server
	DB        storage.Storage
	jwtSecret []byte
}

func NewApp(logger logger.Logger, server *http.Server, DB storage.Storage) *App {
	jwt := os.Getenv("JWT_SECRET")
	return &App{Log: logger, Server: server, DB: DB, jwtSecret: []byte(jwt)}
}

func (app *App) ErrorJSON(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	app.Log.Errorln(r.Method, r.URL, err.Error())
	var code int
	if errors.Is(err, apperrors.ErrNotFound) {
		err = apperrors.ErrNotFound
		code = http.StatusNotFound
	} else if errors.Is(err, errBadRequest) {
		err = errBadRequest
		code = http.StatusBadRequest
	} else if errors.Is(err, errInvalidCredentials) || errors.Is(err, errUnauthorized) {
		err = errUnauthorized
		code = http.StatusUnauthorized
	} else if errors.Is(err, errForbidden) {
		err = errForbidden
		code = http.StatusForbidden
	} else if errors.Is(err, apperrors.ErrAlreadyExists) {
		err = apperrors.ErrAlreadyExists
		code = http.StatusConflict
	} else {
		err = errInternal
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func (app *App) Encode(w http.ResponseWriter, r *http.Request, obj any) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(obj)
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
}

func (app *App) ExtractPathValue(r *http.Request, pathValue string) (int, error) {
	stringID := r.PathValue(pathValue)
	id, err := strconv.Atoi(stringID)
	if stringID == "" || err != nil {
		return 0, fmt.Errorf("%w: error extracting path value(%s, %s): %w",
			errBadRequest, pathValue, stringID, err)
	}
	return id, nil
}

func (app *App) GetClaims(r *http.Request) (userID int, role models.Role, err error) {
	claims, ok := r.Context().Value(ClaimsKey).(*Claims)
	if !ok || claims == nil {
		err = fmt.Errorf("error getting claims")
		return
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return
	}
	userID, err = strconv.Atoi(sub)
	if err != nil {
		return
	}
	role = claims.Role
	return
}
