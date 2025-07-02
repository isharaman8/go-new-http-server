package handler_test

import (
	"bytes"
	"encoding/json"
	"go-user-api/internal/handler"
	"go-user-api/internal/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
	mock := &testutils.MockUserRepo{}
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
