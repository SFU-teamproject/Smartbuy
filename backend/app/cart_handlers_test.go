package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/logger/mocklogger"
	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/sfu-teamproject/smartbuy/backend/storage/mockstorage"
	"github.com/stretchr/testify/mock"
)

func TestGetCart(t *testing.T) {
	ms := new(mockstorage.MockStorage)
	ml := new(mocklogger.MockLogger)
	currTime := time.Now()
	ms.On("GetCartItems", mock.Anything).Return([]models.CartItem{}, nil)
	ms.On("GetCart", 1).Return(models.Cart{1, 1, currTime, currTime, nil}, nil)
	ms.On("GetCart", 2).Return(models.Cart{}, apperrors.ErrNotFound)
	ms.On("GetCart", 3).Return(models.Cart{3, 3, currTime, currTime, nil}, nil)
	ms.On("GetCart", 4).Return(models.Cart{4, 4, currTime, currTime, nil}, nil)
	ml.On("Errorln", mock.Anything, mock.Anything, mock.Anything)

	app := NewApp(ml, nil, ms)
	tests := []struct {
		name   string
		cartID string
		userID string
		role   models.Role
		err    error
	}{
		{"Existing cart, same user", "1", "1", models.RoleUser, nil},
		{"Nonexisting cart", "2", "1", models.RoleUser, apperrors.ErrNotFound},
		{"Existing cart, diff user", "3", "2", models.RoleUser, apperrors.ErrForbidden},
		{"Existing cart, diff user (admin)", "4", "2", models.RoleAdmin, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createContextWithClaims(tt.userID, tt.role)
			r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			r.SetPathValue("cart_id", tt.cartID)
			w := httptest.NewRecorder()
			app.GetCart(w, r)
		})

	}
	ms.AssertExpectations(t)
}
