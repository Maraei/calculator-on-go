package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLoginFlow(t *testing.T) {
	registerResp, err := http.Post("http://localhost:50051/register", "application/json", bytes.NewBufferString(`{"username": "testuser", "password": "1234"}`))
	assert.NoError(t, err)

	var registerBody map[string]interface{}
	json.NewDecoder(registerResp.Body).Decode(&registerBody)
	assert.Equal(t, "Registration successful", registerBody["message"])

	loginResp, err := http.Post("http://localhost:50051/login", "application/json", bytes.NewBufferString(`{"username": "testuser", "password": "1234"}`))
	assert.NoError(t, err)

	var loginBody map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginBody)
	token := loginBody["token"].(string)
	assert.NotEmpty(t, token)

	req, err := http.NewRequest("GET", "http://localhost:50051/validate", nil)
	req.Header.Add("Authorization", token)
	validateResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	var validateBody map[string]interface{}
	json.NewDecoder(validateResp.Body).Decode(&validateBody)
	assert.Equal(t, true, validateBody["valid"])
	assert.Equal(t, "testuser", validateBody["username"])
}
