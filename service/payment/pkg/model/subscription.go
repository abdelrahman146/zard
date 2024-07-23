package model

import (
	"github.com/abdelrahman146/zard/shared"
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	ID                 string         `json:"id" gorm:"column:id;type:text;primaryKey"`
	WorkspaceID        string         `json:"workspaceId" gorm:"column:workspaceId;index;not null"`
	PlanID             string         `json:"planId" gorm:"column:planId;type:text;index;not null"`
	ServiceName        string         `json:"serviceName" gorm:"column:serviceName;index;type:text;not null"`
	Billing            Billing        `json:"billing,omitempty" gorm:"foreignKey:BillingID;references:ID"`
	BillingID          string         `json:"billingId" gorm:"column:billingId;type:text;index;not null"`
	Status             string         `json:"status" gorm:"column:status;type:text;default:'active'"`
	DaysUntilDue       int            `json:"daysUntilDue" gorm:"column:daysUntilDue;type:integer;default:0"`
	CurrentPeriodEnd   time.Time      `json:"currentPeriodEnd" gorm:"column:currentPeriodEnd;type:timestamp"`
	CurrentPeriodStart time.Time      `json:"currentPeriodStart" gorm:"column:currentPeriodStart;type:timestamp"`
	TrialStartAt       *time.Time     `json:"trialStartAt" gorm:"column:trialStartAt;type:timestamp"`
	TrialEndAt         *time.Time     `json:"trialEndAt" gorm:"column:trialEndAt;type:timestamp"`
	CancelAtPeriodEnd  bool           `json:"cancelAtPeriodEnd" gorm:"column:cancelAtPeriodEnd;type:boolean"`
	CancelAt           *time.Time     `json:"cancelAt" gorm:"column:cancelAt;type:timestamp"`
	CancellationReason *string        `json:"cancellationReason" gorm:"column:cancellationReason;type:text"`
	EndedAt            *time.Time     `json:"endedAt" gorm:"column:endedAt;type:timestamp"`
	CreatedAt          time.Time      `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = "sub_" + shared.Utils.Strings.Cuid()
	return
}
