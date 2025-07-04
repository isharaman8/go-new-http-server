package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
}
