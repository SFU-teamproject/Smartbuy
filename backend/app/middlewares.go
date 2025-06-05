package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func (app *App) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			app.Log.Errorf("error reading request body: %v", err)
		}
		defer r.Body.Close()
		app.Log.Infof("Incoming request:\n%s %s\n%s",
			r.Method, r.RequestURI, string(body))
		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

func (app *App) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.ErrorJSON(w, r, apperrors.ErrInternal)
				app.Log.Errorf("error: %v, stack trace: %s", err, string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *App) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		tokenStrTrim := strings.TrimPrefix(tokenStr, "Bearer ")
		if strings.TrimSpace(tokenStrTrim) == "" {
			app.ErrorJSON(w, r, fmt.Errorf("%w: missing token(%s)", apperrors.ErrUnauthorized, tokenStr))
			return
		}
		token, err := jwt.ParseWithClaims(tokenStrTrim, &Claims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: unexpected signing method: %v",
					apperrors.ErrUnauthorized, t.Header["alg"])
			}
			return app.jwtSecret, nil
		})
		if err != nil {
			app.ErrorJSON(w, r, fmt.Errorf("%w: error parsing token: %w", apperrors.ErrUnauthorized, err))
			return
		}
		claims, ok := token.Claims.(*Claims)
		if !ok {
			app.ErrorJSON(w, r, fmt.Errorf("%w: unsupported claims in token: %v",
				apperrors.ErrUnauthorized, token.Claims))
			return
		}
		if !token.Valid {
			app.ErrorJSON(w, r, fmt.Errorf("%w: invalid token", apperrors.ErrUnauthorized))
			return
		}
		iss, err := claims.GetIssuer()
		if err != nil || iss != "Smartbuy" {
			app.ErrorJSON(w, r, fmt.Errorf("%w: invalid issuer %s", apperrors.ErrUnauthorized, iss))
			return
		}
		exp, err := claims.GetExpirationTime()
		if err != nil {
			app.ErrorJSON(w, r, fmt.Errorf("%w: error getting expiration date: %w", apperrors.ErrUnauthorized, err))
			return
		}
		if time.Now().After(exp.Time) {
			app.ErrorJSON(w, r, fmt.Errorf("%w: token expired at %s", apperrors.ErrUnauthorized, exp.Time.String()))
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
