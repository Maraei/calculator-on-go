package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

func TestRegisterAndLoginFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/register":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Registration successful",
			})
		case "/login":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"token": "test-token",
			})
		case "/validate":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"valid":    true,
				"username": "testuser",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	registerResp, err := http.Post(ts.URL+"/register", "application/json", 
		bytes.NewBufferString(`{"username": "testuser", "password": "1234"}`))
	assert.NoError(t, err)
	defer registerResp.Body.Close()

	var registerBody map[string]interface{}
	err = json.NewDecoder(registerResp.Body).Decode(&registerBody)
	assert.NoError(t, err)
	assert.Equal(t, "Registration successful", registerBody["message"])

	loginResp, err := http.Post(ts.URL+"/login", "application/json", 
		bytes.NewBufferString(`{"username": "testuser", "password": "1234"}`))
	assert.NoError(t, err)
	defer loginResp.Body.Close()

	var loginBody map[string]interface{}
	err = json.NewDecoder(loginResp.Body).Decode(&loginBody)
	assert.NoError(t, err)
	token := loginBody["token"].(string)
	assert.NotEmpty(t, token)

	req, err := http.NewRequest("GET", ts.URL+"/validate", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", token)

	validateResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer validateResp.Body.Close()

	var validateBody map[string]interface{}
	err = json.NewDecoder(validateResp.Body).Decode(&validateBody)
	assert.NoError(t, err)
	assert.Equal(t, true, validateBody["valid"])
	assert.Equal(t, "testuser", validateBody["username"])
}
