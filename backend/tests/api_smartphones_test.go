package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestSmartphones(t *testing.T) {
	base_url := "http://localhost:8080/api/v1"
	t.Run("get smartphones", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/smartphones", base_url), nil)
		assert.NoError(t, err, "creating request failed", err.Error())
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "sending request failed", err.Error())
		assert.Equal(t, http.StatusOK, resp.StatusCode, "status code is not 200")
		var smartphones []models.Smartphone
		err = json.NewDecoder(resp.Body).Decode(&smartphones)
		assert.NoError(t, err, "decoding smartphones failed", err.Error())
		assert.NotEmpty(t, smartphones, "smartphones slice is empty")
	})
	t.Run("get smartphone", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/smartphones/1", base_url), nil)
		assert.NoError(t, err, "creating request failed", err.Error())
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "sending request failed", err.Error())
		assert.Equal(t, http.StatusOK, resp.StatusCode, "status code is not 200")
		var smartphone models.Smartphone
		err = json.NewDecoder(resp.Body).Decode(&smartphone)
		assert.NoError(t, err, "decoding smartphone failed", err.Error())
		assert.NotEqual(t, smartphone, models.Smartphone{}, "smartphone struct is empty")
	})
}
