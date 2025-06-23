package app

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/logger/mocklogger"
	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/sfu-teamproject/smartbuy/backend/storage/mockstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	ms := new(mockstorage.MockStorage)
	ml := new(mocklogger.MockLogger)
	currTime := time.Now().Round(0)
	user2 := models.User{2, "user2", "", models.RoleUser, currTime, models.Cart{}}
	user1 := models.User{1, "user1", "", models.RoleUser, currTime, models.Cart{}}
	user3 := models.User{3, "user3", "", models.RoleUser, currTime, models.Cart{}}
	ms.On("GetUser", 1).Return(user1, nil)
	ms.On("GetUser", 2).Return(user2, nil)
	ms.On("GetUser", 3).Return(models.User{}, apperrors.ErrNotFound)
	ms.On("GetCartByUserID", mock.Anything).Return(models.Cart{}, nil)
	ms.On("GetCartItems", mock.Anything).Return([]models.CartItem(nil), nil)
	ml.On("Errorln", mock.Anything, mock.Anything, mock.Anything)

	app := NewApp(ml, nil, ms)
	tests := []struct {
		name        string
		User        models.User
		ReqUserID   string
		ReqUserRole models.Role
		err         error
	}{
		{"Existing user, same user", user1, "1", models.RoleUser, nil},
		{"Existing user, admin", user2, "1", models.RoleAdmin, nil},
		{"Existing user, diff user", user2, "1", models.RoleUser, apperrors.ErrForbidden},
		{"Non-existing user", user3, "1", models.RoleAdmin, apperrors.ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createContextWithClaims(tt.ReqUserID, tt.ReqUserRole)
			r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			r.SetPathValue("user_id", strconv.Itoa(tt.User.ID))
			w := httptest.NewRecorder()
			app.GetUser(w, r)
			var resp models.User
			err := json.NewDecoder(w.Body).Decode(&resp)
			if tt.err == nil {
				assert.NoError(t, err, "Decoding user failed", r.Body)
				assert.Equal(t, tt.User, resp)
			} else {
				assert.NotEqual(t, 200, w.Result().StatusCode)
			}
		})
	}
	ms.AssertExpectations(t)
}

func TestSignup(t *testing.T) {
	ms := new(mockstorage.MockStorage)
	ml := new(mocklogger.MockLogger)
	currTime := time.Now().Round(0)
	user1 := models.User{1, "user1", "", models.RoleUser, currTime, models.Cart{}}
	ms.On("CreateUser", mock.Anything).Return(user1, nil)
	ml.On("Errorln", mock.Anything, mock.Anything, mock.Anything)
	app := NewApp(ml, nil, ms)
	tests := []struct {
		name string
		user models.User
		err  error
	}{
		{"Signup new user", user1, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.user)
			assert.NoError(t, err, "Marshalling user failed: %v")
			rBody := bytes.NewBuffer(jsonBytes)
			r := httptest.NewRequest(http.MethodPost, "/", rBody)
			w := httptest.NewRecorder()
			app.Signup(w, r)
			var resp models.User
			err = json.NewDecoder(w.Body).Decode(&resp)
			assert.NoError(t, err, "Decoding response user failed: %v")
			if tt.err == nil {
				assert.Equal(t, resp, tt.user)
				assert.Equal(t, w.Result().StatusCode, 201)
			}
		})
	}
	ms.AssertExpectations(t)
}
