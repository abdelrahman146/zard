package model

import (
	"github.com/abdelrahman146/zard/shared"
	"gorm.io/gorm"
	"time"
)

type Organization struct {
	ID         string         `json:"id" gorm:"column:id;type:text;primaryKey"`
	Name       string         `json:"name" gorm:"column:name;uniqueIndex;type:text;not null"`
	Website    *string        `json:"website" gorm:"column:website;type:text"`
	Email      string         `json:"email" bson:"column:email;uniqueIndex;type:text;not null"`
	Phone      *string        `json:"phone" bson:"column:phone;type:text"`
	Country    string         `json:"country" bson:"column:country;type:text"`
	City       string         `json:"city" bson:"column:city;type:text"`
	Address    string         `json:"address" bson:"column:address;type:text"`
	Users      []User         `json:"users,omitempty" gorm:"foreignKey:OrgID;references:ID"`
	Workspaces []Workspace    `json:"workspaces,omitempty" gorm:"foreignKey:OrgID;references:ID"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = "org_" + shared.Utils.Strings.Cuid()
	return
}
