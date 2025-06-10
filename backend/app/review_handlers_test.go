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
	"github.com/stretchr/testify/mock"
)

func TestGetReviews(t *testing.T) {
	ms := new(mockstorage.MockStorage)
	ml := new(mocklogger.MockLogger)
	currTime := time.Now()
	ms.On("GetReviews", 1).Return([]models.Review{{1, 1, 1, "user1", 5, nil, currTime, currTime}}, nil)
	ms.On("GetReviews", 2).Return([]models.Review{}, nil)
	ms.On("GetReviews", 3).Return([]models.Review{}, apperrors.ErrNotFound)
	ml.On("Errorln", mock.Anything, mock.Anything, mock.Anything)

	app := NewApp(ml, nil, ms)
	tests := []struct {
		name         string
		smartphoneID string
		err          error
	}{
		{"Existing smartphone with reviews", "1", nil},
		{"Existing smartphone without reviews", "2", nil},
		{"Nonexisting smartphone", "3", apperrors.ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.SetPathValue("smartphone_id", tt.smartphoneID)
			w := httptest.NewRecorder()
			app.GetReviews(w, r)
		})
	}
	ms.AssertExpectations(t)
}

func TestCreateReview(t *testing.T) {
	ms := new(mockstorage.MockStorage)
	ml := new(mocklogger.MockLogger)
	currTime := time.Now().Round(0)
	r := models.Review{1, 1, 1, "user", 5, nil, currTime, currTime}
	rErr := models.Review{1, 2, 1, "user", 5, nil, currTime, currTime}
	ms.On("CreateReview", r).Return(r, nil)
	ms.On("GetSmartphone", 1).Return(models.Smartphone{}, nil)
	ms.On("GetSmartphone", 2).Return(models.Smartphone{}, apperrors.ErrNotFound)
	ml.On("Errorln", mock.Anything, mock.Anything, mock.Anything)

	app := NewApp(ml, nil, ms)
	tests := []struct {
		name   string
		review models.Review
		err    error
	}{
		{"Create review for existing smartphone", r, nil},
		{"Create review for non-existing smartphone", rErr, apperrors.ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.review)
			if err != nil {
				t.Fatalf("Marshalling review: %v", err)
			}
			rBody := bytes.NewBuffer(jsonBytes)
			ctx := createContextWithClaims(strconv.Itoa(tt.review.UserID), models.RoleUser)
			r := httptest.NewRequestWithContext(ctx, http.MethodPost, "/", rBody)
			r.SetPathValue("smartphone_id", strconv.Itoa(tt.review.SmartphoneID))
			w := httptest.NewRecorder()
			app.CreateReview(w, r)
		})
	}
	ms.AssertExpectations(t)
}
