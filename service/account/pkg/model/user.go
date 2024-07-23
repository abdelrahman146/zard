package model

import (
	"github.com/abdelrahman146/zard/shared"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              string         `json:"id" gorm:"column:id;type:text;primaryKey"`
	Name            string         `json:"name" gorm:"column:name;uniqueIndex;type:text;not null"`
	Email           string         `json:"email" bson:"column:email;uniqueIndex;type:text;not null"`
	Phone           *string        `json:"phone" bson:"column:phone;type:text"`
	Password        *string        `json:"password" gorm:"column:website;type:text"`
	IsEmailVerified bool           `json:"isEmailVerified" gorm:"column:isEmailVerified;type:boolean"`
	IsPhoneVerified bool           `json:"isPhoneVerified" gorm:"column:isPhoneVerified;type:boolean"`
	Active          bool           `json:"active" gorm:"column:active;type:boolean"`
	Organization    Organization   `json:"organization,omitempty" gorm:"foreignKey:OrgID;references:ID"`
	OrgID           string         `json:"orgId" gorm:"column:orgId;type:text"`
	CreatedAt       time.Time      `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = "usr_" + shared.Utils.Strings.Cuid()
	return
}
