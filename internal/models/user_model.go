package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Name        string    `db:"name" json:"name" validate:"required,min=5,max=50"`
	Email       string    `db:"email" json:"email" validate:"required,email,lte=255"`
	Password    string    `db:"password" json:"password,omitempty" validate:"required,min=5,max=15"`
	UserStatus  int       `db:"user_status" json:"-" validate:"required,len=1"`
	UserRole    string    `db:"user_role" json:"-" validate:"required"`
	AccessToken string    `db:"-" json:"accessToken"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

type SignUpRequest struct {
	Email    string `db:"email" json:"email" validate:"required,email,lte=255"`
	Name     string `db:"name" json:"name" validate:"required,min=5,max=50"`
	Password string `db:"password" json:"password,omitempty" validate:"required,min=5,max=15"`
}

type AuthResponse struct {
	Email       string `db:"email" json:"email" validate:"required,email,lte=255"`
	Name        string `db:"name" json:"name" validate:"required,min=5,max=50"`
	AccessToken string `db:"token" json:"accessToken"`
}

type SignInRequest struct {
	Email    string `db:"email" json:"email" validate:"required,email,lte=255"`
	Password string `db:"password" json:"password,omitempty" validate:"required,min=5,max=15"`
}
