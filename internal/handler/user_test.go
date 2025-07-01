package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-user-api/internal/model"

	"github.com/stretchr/testify/assert"
)

// ---- 👇 Mock Repository Implementation ----

type mockUserRepo struct{}

func (m *mockUserRepo) Create(_ context.Context, u *model.User) error {
	u.ID = 1 // Simulate DB auto-increment
	return nil
}
func (m *mockUserRepo) Get(_ context.Context, id int) (*model.User, error) { return nil, nil }
func (m *mockUserRepo) Update(_ context.Context, u *model.User) error      { return nil }
func (m *mockUserRepo) Delete(_ context.Context, id int) error             { return nil }

// ---- ✅ Test CreateUser ----

func TestCreateUser(t *testing.T) {
	handler := NewUserHandler(&mockUserRepo{})

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
