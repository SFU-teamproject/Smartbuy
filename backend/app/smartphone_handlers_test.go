package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/logger/mocklogger"
	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/sfu-teamproject/smartbuy/backend/storage/mockstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSmartphones(t *testing.T) {
	ms := new(mockstorage.MockStorage)
	ml := new(mocklogger.MockLogger)
	sms := []models.Smartphone{
		{ID: 1, Model: "model", Producer: "producer", Memory: 1, Ram: 1,
			DisplaySize: 1, Price: 1, RatingsSum: 1, RatingsCount: 1, ImagePath: ""}}
	ms.On("GetSmartphones").Return(sms, nil)
	ml.On("Errorln", mock.Anything, mock.Anything, mock.Anything)

	app := NewApp(ml, nil, ms)
	tests := []struct {
		name string
		err  error
	}{
		{"Get smartphones", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			app.GetSmartphones(w, r)
			var resp []models.Smartphone
			err := json.NewDecoder(w.Body).Decode(&resp)
			assert.NoError(t, err, "Decoding smartphones failed")
			assert.Equal(t, resp, sms)
		})
	}
	ms.AssertExpectations(t)
}

func TestGetSmartphone(t *testing.T) {
	ms := new(mockstorage.MockStorage)
	ml := new(mocklogger.MockLogger)
	sm1 := models.Smartphone{
		ID: 1, Model: "model1", Producer: "producer1", Memory: 1, Ram: 1,
		DisplaySize: 1, Price: 1, RatingsSum: 1, RatingsCount: 1, ImagePath: ""}
	sm2 := models.Smartphone{ID: 2}
	ms.On("GetSmartphone", 1).Return(sm1, nil)
	ms.On("GetSmartphone", 2).Return(sm2, apperrors.ErrNotFound)
	ms.On("GetReviews", mock.Anything).Return([]models.Review{}, nil)
	ml.On("Errorln", mock.Anything, mock.Anything, mock.Anything)
	app := NewApp(ml, nil, ms)
	tests := []struct {
		name       string
		smartphone models.Smartphone
		err        error
	}{
		{"Get existing smartphone", sm1, nil},
		{"Get non-existing smartphone", sm2, apperrors.ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.SetPathValue("smartphone_id", strconv.Itoa(tt.smartphone.ID))
			w := httptest.NewRecorder()
			app.GetSmartphone(w, r)
			var resp models.Smartphone
			err := json.NewDecoder(w.Body).Decode(&resp)
			if tt.err == nil {
				assert.NoError(t, err, "Decoding smartphone failed", tt.err, err)
				assert.Equal(t, resp, tt.smartphone)
			}
		})
	}
	ms.AssertExpectations(t)
}
