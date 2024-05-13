package models

import (
	"time"

	"github.com/google/uuid"
)

type CatMatch struct {
	ID          uuid.UUID  `db:"id" json:"id" validate:"required,uuid"`
	CatIssuerID uuid.UUID  `db:"cat_issuer_id" json:"userCatId" validate:"required,uuid"`
	CatMatchID  uuid.UUID  `db:"cat_match_id" json:"matchCatId" validate:"required,uuid"`
	Message     string     `db:"message" json:"message" validate:"required,min=5,max=120"`
	Status      string     `db:"status" json:"-"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"-"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`
}

type CatMatchRequest struct {
	CatIssuerID uuid.UUID `db:"cat_issuer_id" json:"userCatId" validate:"required,uuid"`
	CatMatchID  uuid.UUID `db:"cat_match_id" json:"matchCatId" validate:"required,uuid"`
	Message     string    `db:"message" json:"message" validate:"required,min=5,max=120"`
}

type CatMatchUpdateRequest struct {
	ID uuid.UUID `db:"id" json:"matchId" validate:"required,uuid"`
}

type CatMatchDetail struct {
	ID             uuid.UUID  `db:"id" json:"id" validate:"required,uuid"`
	IssuedBy       IssuerUser `db:"issueruser" json:"issuedBy"`
	MatchCatDetail LoveCat    `db:"matchcat" json:"matchCatDetail"`
	UserCatDetail  LoveCat    `db:"issuercat" json:"userCatDetail"`
	Message        string     `db:"message" json:"message"`
	IsApproved     bool       `db:"isapproved" json:"-"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" json:"-"`
	DeletedAt      *time.Time `db:"deleted_at" json:"-"`
}

type IssuerUser struct {
	Name      string    `db:"name" json:"name" validate:"required,min=5,max=50"`
	Email     string    `db:"email" json:"email" validate:"required,email,lte=255"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type LoveCat struct {
	ID          uuid.UUID   `db:"id" json:"id" validate:"required,uuid"`
	Name        string      `db:"name" json:"name" validate:"required,min=1,max=30"`
	Race        string      `db:"race" json:"race" validate:"required"`
	Sex         string      `db:"sex" json:"sex" validate:"required"`
	Description string      `db:"description" json:"description" validate:"required,min=1,max=200"`
	AgeInMonth  int         `db:"ageinmonth" json:"ageInMonth" validate:"required,min=1,max=120082"`
	ImageUrls   []string `db:"imageurls" json:"imageUrls" validate:"required,dive,required"`
	HasMatched  bool        `db:"hasmatched" json:"hasMatched"`
	CreatedAt   time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time  `db:"updated_at" json:"-"`
	DeletedAt   *time.Time  `db:"deleted_at" json:"-"`
}
