package model

import (
	"github.com/abdelrahman146/zard/shared"
	"gorm.io/gorm"
	"time"
)

type Workspace struct {
	ID        string         `json:"id" gorm:"column:id;type:text;primaryKey"`
	Name      string         `json:"name" gorm:"column:name;uniqueIndex;type:text;not null"`
	Website   *string        `json:"website" gorm:"column:website;type:text"`
	ApiKey    string         `json:"apiKey" gorm:"column:apiKey;type:text;not null"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
}

func (w *Workspace) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = "wrk_" + shared.Utils.Strings.Cuid()
	apikey := shared.Utils.Auth.CreateToken("zky", w.ID, "zard") // TODO get from config
	w.ApiKey = shared.Utils.Auth.Encrypt(apikey, "zard")         // TODO get from config
	return
}
