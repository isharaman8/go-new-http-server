package testutils

import (
	"context"
	"go-user-api/internal/model"
)

type MockUserRepo struct{}

func (m *MockUserRepo) Create(_ context.Context, u *model.User) error {
	u.ID = 1 // Simulate DB auto-increment
	return nil
}
func (m *MockUserRepo) Get(_ context.Context, id int) (*model.User, error) { return nil, nil }
func (m *MockUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Update(_ context.Context, u *model.User) error { return nil }
func (m *MockUserRepo) Delete(_ context.Context, id int) error        { return nil }
func (m *MockUserRepo) GetAllUsers(_ context.Context) ([]*model.User, error) {
	return []*model.User{
		{ID: 1, Email: "user1@example.com", Name: "User One"},
		{ID: 2, Email: "user2@example.com", Name: "User Two"},
	}, nil
}
