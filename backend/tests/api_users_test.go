package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestCart(t *testing.T) {
	base_url := "http://localhost:8080/api/v1"
	user := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{
		Name:     "test_user543",
		Password: "password",
	}
	var token string
	var userID int
	t.Run("signing up", func(t *testing.T) {
		data, err := json.Marshal(user)
		assert.NoError(t, err, "marshalling user failed", err.Error())
		buff := bytes.NewBuffer(data)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/signup", base_url), buff)
		assert.NoError(t, err, "creating request failed", err.Error())
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "sending request failed", err.Error())
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "status code is not 201")
		var newUser models.User
		err = json.NewDecoder(resp.Body).Decode(&newUser)
		assert.NoError(t, err, "decoding user failed", err.Error())
		assert.NotEmpty(t, newUser.ID, "ID is 0")
		userID = newUser.ID
		assert.Equal(t, user.Name, newUser.Name, "user names do not match")
		assert.Equal(t, models.RoleUser, newUser.Role, "user role is not RoleUser")
	})
	t.Run("login", func(t *testing.T) {
		data, err := json.Marshal(user)
		assert.NoError(t, err, "marshalling user failed", err.Error())
		buff := bytes.NewBuffer(data)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login", base_url), buff)
		assert.NoError(t, err, "creating request failed", err.Error())
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "sending request failed", err.Error())
		assert.Equal(t, http.StatusOK, resp.StatusCode, "status code is not 200")
		var response struct {
			User  models.User `json:"user"`
			Token string      `json:"token"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err, "decoding user failed", err.Error())
		assert.NotEmpty(t, response.Token, "token is empty")
		token = response.Token
		assert.Equal(t, userID, response.User.ID, "user's ID is incorrect")
		assert.Equal(t, user.Name, response.User.Name, "user names do not match")
		assert.Equal(t, models.RoleUser, response.User.Role, "user role is not RoleUser")
	})
	t.Run("get user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/%d", base_url, userID), nil)
		req.Header.Add("Authorization", token)
		assert.NoError(t, err, "creating request failed", err.Error())
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "sending request failed", err.Error())
		assert.Equal(t, http.StatusOK, resp.StatusCode, "status code is not 200")
		var newUser models.User
		err = json.NewDecoder(resp.Body).Decode(&newUser)
		assert.NoError(t, err, "decoding user failed", err.Error())
		assert.Equal(t, userID, newUser.ID, "user's ID is incorrect")
		assert.Equal(t, user.Name, newUser.Name, "user names do not match")
		assert.Equal(t, models.RoleUser, newUser.Role, "user role is not RoleUser")
	})
}
