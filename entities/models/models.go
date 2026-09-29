package models

import (
	"encoding/json"
	"time"
)

type Onboarding struct {
	ID             uint64          `db:"id" json:"id"`
	Name           string          `db:"name" json:"name"`
	Email          string          `db:"email" json:"email"`
	Password       string          `db:"password" json:"-"`
	PhoneNumber    string          `db:"phone_number" json:"phone_number"`
	Address        string          `db:"address" json:"address"`
	RegisterAs     string          `db:"register_as" json:"register_as"`
	Status         string          `db:"status" json:"status"`
	AdditionalInfo json.RawMessage `db:"additional_info" json:"additional_info"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
}

type User struct {
	ID             uint64          `db:"id" json:"id"`
	OnboardingID   uint64          `db:"onboarding_id" json:"onboarding_id"`
	UserName       string          `db:"user_name" json:"user_name"`
	Status         string          `db:"status" json:"status"`
	AdditionalInfo json.RawMessage `db:"additional_info" json:"additional_info"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
}
