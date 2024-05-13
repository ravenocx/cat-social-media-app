package models

import (
	"time"

	"github.com/google/uuid"
)

type Cats struct {
	ID          uuid.UUID   `db:"id" json:"id" validate:"required,uuid"`
	UserID      uuid.UUID   `db:"user_id" json:"user_id" validate:"required,uuid"`
	Name        string      `db:"name" json:"name" validate:"required,min=1,max=30"`
	Race        string      `db:"race" json:"race" validate:"required"`
	Sex         string      `db:"sex" json:"sex" validate:"required"`
	AgeInMonth  int         `db:"ageinmonth" json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string      `db:"description" json:"description" validate:"required,min=1,max=200"`
	HasMatched  bool        `db:"hasmatched" json:"hasMatched"`
	ImageUrls   []string `db:"imageurls" json:"imageUrls" validate:"required,dive,required"`
	CreatedAt   time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time  `db:"updated_at" json:"-"`
	DeletedAt   *time.Time  `db:"deleted_at" json:"-"`
}

type CatUpdateRequest struct {
	Name        string     `db:"name" json:"name" validate:"required,min=1,max=30"`
	Race        string     `db:"race" json:"race" validate:"required"`
	Sex         string     `db:"sex" json:"sex" validate:"required"`
	AgeInMonth  int        `db:"ageinmonth" json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string     `db:"description" json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string   `db:"imageurls" json:"imageUrls" validate:"required,min=1,dive,url"`
	UpdatedAt   *time.Time `db:"updated_at" json:"-"`
}

type Cat struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	NewCat
	HasMatched bool      `json:"hasMatched" db:"hasmatched"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"-" db:"updated_at"`
	DeletedAt  time.Time `json:"-" db:"deleted_at"`
}

type NewCat struct {
	Name        string   `json:"name" db:"name" validate:"required,min=1,max=30"`
	Race        string   `json:"race" db:"race" validate:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex         string   `json:"sex" db:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" db:"ageinmonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" db:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" db:"imageurls" validate:"required,min=1,dive,url"`
}

type CatData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Race        string   `json:"race"`
	Sex         string   `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	ImageUrls   []string `json:"imageUrls"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	CreatedAt   string   `json:"createdAt"`
}