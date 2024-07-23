package model

import (
	"github.com/abdelrahman146/zard/shared"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Invoice struct {
	ID                 string                 `json:"id" gorm:"column:id;type:text;primaryKey"`
	Billing            Billing                `json:"organization,omitempty" gorm:"foreignKey:BillingID;references:ID"`
	BillingID          string                 `json:"billingId" gorm:"column:billingId;index;type:text"`
	WorkspaceID        string                 `json:"workspaceId" gorm:"column:workspaceId;index;not null"`
	InvoiceNum         string                 `json:"invoiceNum" gorm:"column:invoiceNum;index;type:text;not null"`
	Name               string                 `json:"name" gorm:"column:name;uniqueIndex;type:text;not null"`
	Email              string                 `json:"email" bson:"column:email;uniqueIndex;type:text;not null"`
	Phone              *string                `json:"phone" bson:"column:phone;type:text"`
	Country            string                 `json:"country" bson:"column:country;type:text"`
	City               string                 `json:"city" bson:"column:city;type:text"`
	Address            string                 `json:"address" bson:"column:address;type:text"`
	AmountDue          float64                `json:"amountDue" gorm:"column:amountDue;type:float"`
	SubTotal           float64                `json:"subTotal" gorm:"column:subTotal;type:float"`
	Tax                float64                `json:"tax" gorm:"column:tax;type:float"`
	Total              float64                `json:"total" gorm:"column:total;type:float"`
	Currency           string                 `json:"currency" gorm:"column:currency;type:text;not null;default:'usd'"`
	Status             string                 `json:"status" gorm:"column:status;type:text;default:'draft'"`
	CollectionMethod   string                 `json:"collectionMethod" gorm:"column:collectionMethod;type:text;default:'charge_automatically'"`
	PaymentDetails     map[string]interface{} `json:"paymentDetails" gorm:"column:paymentDetails;type:jsonb"`
	DueDate            time.Time              `json:"dueDate" gorm:"column:dueDate;type:timestamp;not null"`
	Lines              []InvoiceLine          `json:"lines,omitempty" gorm:"foreignKey:InvoiceID;references:ID"`
	Attempts           int                    `json:"attempts" gorm:"column:attempts;type:integer;default:0"`
	LastAttemptAt      *time.Time             `json:"lastAttemptAt" gorm:"column:lastAttemptAt;type:timestamp"`
	CanceledAt         *time.Time             `json:"canceledAt" gorm:"column:canceledAt;type:timestamp"`
	CancellationReason *string                `json:"cancellationReason" gorm:"column:cancellationReason;type:text"`
	CreatedAt          time.Time              `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt          gorm.DeletedAt         `json:"deletedAt" gorm:"column:deletedAt"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = "inv_" + shared.Utils.Strings.Cuid()
	i.InvoiceNum = "INV-" + strings.ToUpper(shared.Utils.Strings.Cuid())[:6]
	return
}

type InvoiceLine struct {
	ID        string  `json:"id" gorm:"column:id;type:text;primaryKey"`
	Type      string  `json:"type" gorm:"column:type;type:text;not null;default:'invoice_item'"`
	InvoiceID string  `json:"invoiceId" gorm:"column:invoiceId;index;type:text"`
	Invoice   Invoice `json:"invoice,omitempty" gorm:"foreignKey:InvoiceID;references:ID"`
	Name      string  `json:"name" gorm:"column:name;type:text;not null"`
	Quantity  int     `json:"quantity" gorm:"column:quantity;type:integer;not null"`
	UnitPrice float64 `json:"unitPrice" gorm:"column:unitPrice;type:float;not null"`
	Subtotal  float64 `json:"subtotal" gorm:"column:subtotal;type:float;not null"`
	Tax       float64 `json:"tax" gorm:"column:tax;type:float;not null"`
	Total     float64 `json:"total" gorm:"column:total;type:float;not null"`
	Currency  string  `json:"currency" gorm:"column:currency;type:text;not null;default:'usd'"`
}

func (i *InvoiceLine) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = "inl_" + shared.Utils.Strings.Cuid()
	return
}
