package model

import (
	"github.com/abdelrahman146/zard/shared"
	"gorm.io/gorm"
	"time"
)

type Billing struct {
	ID               string                 `json:"id" gorm:"column:id;type:text;primaryKey"`
	WorkspaceID      string                 `json:"workspaceId" gorm:"column:workspaceId;index;not null"`
	Name             string                 `json:"name" gorm:"column:name;uniqueIndex;type:text;not null"`
	Email            string                 `json:"email" bson:"column:email;uniqueIndex;type:text;not null"`
	Phone            *string                `json:"phone" bson:"column:phone;type:text"`
	Country          string                 `json:"country" bson:"column:country;type:text"`
	City             string                 `json:"city" bson:"column:city;type:text"`
	Address          string                 `json:"address" bson:"column:address;type:text"`
	Currency         string                 `json:"currency" gorm:"column:currency;type:text;not null;default:'usd'"`
	TaxPercentage    float64                `json:"taxPercentage" gorm:"column:taxPercentage;type:float;default:0.05"`
	BillingMethod    string                 `json:"type" gorm:"column:type;type:text;not null;default:'stripe'"`
	CollectionMethod string                 `json:"collectionMethod" gorm:"column:collectionMethod;type:text;default:'charge_automatically'"`
	PaymentDetails   map[string]interface{} `json:"paymentDetails" gorm:"column:paymentDetails;type:jsonb"`
	Subs             []Subscription         `json:"subs,omitempty" gorm:"foreignKey:BillingID;references:ID"`
	Invoices         []Invoice              `json:"invoices,omitempty" gorm:"foreignKey:BillingID;references:ID"`
	CreatedAt        time.Time              `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt        gorm.DeletedAt         `json:"deletedAt" gorm:"column:deletedAt"`
}

func (b *Billing) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = "bil_" + shared.Utils.Strings.Cuid()
	return
}
