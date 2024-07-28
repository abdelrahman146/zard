package usecase

import (
	"gorm.io/gorm"
	"time"
)

type CreateOrgStruct struct {
	Name    string  `json:"name,omitempty" validate:"required,omitempty"`
	Website *string `json:"website,omitempty"`
	Email   string  `json:"email,omitempty" validate:"required,email"`
	Phone   *string `json:"phone,omitempty" validate:"phone"`
	Country string  `json:"country,omitempty" validate:"required,omitempty,iso3166_1_alpha2"`
	City    string  `json:"city,omitempty" validate:"required,omitempty"`
	Address string  `json:"address,omitempty" validate:"required,omitempty"`
}

type UpdateOrgStruct struct {
	Name    *string `json:"name,omitempty"`
	Website *string `json:"website,omitempty"`
	Email   *string `json:"email,omitempty" validate:"email"`
	Phone   *string `json:"phone,omitempty" validate:"phone"`
	Country *string `json:"country,omitempty" validate:"iso3166_1_alpha2"`
	City    *string `json:"city,omitempty"`
	Address *string `json:"address,omitempty"`
}

type CreateWorkspaceStruct struct {
	Name    string  `json:"name,omitempty" validate:"required,omitempty"`
	Website *string `json:"website,omitempty"`
	OrgID   string  `json:"orgId,omitempty" validate:"required,omitempty"`
}

type UpdateWorkspaceStruct struct {
	Name    *string `json:"name,omitempty"`
	Website *string `json:"website,omitempty"`
}

type UserStruct struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	Phone           *string        `json:"phone"`
	IsEmailVerified bool           `json:"isEmailVerified"`
	IsPhoneVerified bool           `json:"isPhoneVerified"`
	Active          bool           `json:"active"`
	OrgID           string         `json:"orgId"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt"`
}

type CreateUserStruct struct {
	Name     string  `json:"name,omitempty" validate:"required,omitempty"`
	Email    string  `json:"email,omitempty" validate:"required,email"`
	Phone    *string `json:"phone,omitempty" validate:"phone"`
	Password *string `json:"password,omitempty" validate:"required,omitempty"`
	OrgID    string  `json:"orgId,omitempty" validate:"required,omitempty"`
}

type UpdateUserStruct struct {
	Name *string `json:"name,omitempty"`
}
