package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-user-api/internal/handler"
	"go-user-api/internal/model"
	"go-user-api/internal/testutils"

	"github.com/stretchr/testify/assert"
)

// ---- ✅ Test CreateUser ----

func TestCreateUser(t *testing.T) {
	handler := handler.NewUserHandler(&testutils.MockUserRepo{})

	// Prepare input user JSON
	inputUser := model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	body, _ := json.Marshal(inputUser)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Call handler
	handler.CreateUser(recorder, req)

	// Assert response
	assert.Equal(t, http.StatusOK, recorder.Code)

	var resp model.User
	err := json.NewDecoder(recorder.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, inputUser.Name, resp.Name)
	assert.Equal(t, inputUser.Email, resp.Email)
}
