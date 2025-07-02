package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go-user-api/internal/handler"
	"go-user-api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepo struct{}

func (m *mockRepo) Create(_ context.Context, u *model.User) error {
	u.ID = 1 // pretend DB inserted user with ID 1
	return nil
}
func (m *mockRepo) Get(_ context.Context, id int) (*model.User, error) {
	return &model.User{ID: id, Email: "test@example.com", Name: "Test User"}, nil
}
func (m *mockRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	return &model.User{ID: 1, Email: email, Password: "$2a$10$dummyhash"}, nil // bcrypt dummy
}
func (m *mockRepo) Update(_ context.Context, u *model.User) error {
	return nil
}

func (m *mockRepo) Delete(_ context.Context, id int) error {
	return nil // or return an error for testing failure cases
}

func TestSignup(t *testing.T) {
	mock := &mockRepo{}
	authHandler := handler.NewAuthRouteHandler(mock)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "test123",
		"name":     "Test User",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	authHandler.Signup(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 ok, got %d", rr.Code)
	}
}
