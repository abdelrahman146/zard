package model

type BackofficeUser struct {
	ID       uint   `json:"id" gorm:"column:id;type:autoIncrement;primaryKey"`
	Name     string `json:"name" gorm:"column:name;type:text;not null"`
	Role     string `json:"role" gorm:"column:role;type:text;not null"`
	Email    string `json:"email" gorm:"column:email;type:text;not null"`
	Password string `json:"password" gorm:"column:password;type:text;not null"`
	Active   bool   `json:"active" gorm:"column:active;type:boolean;not null"`
}
