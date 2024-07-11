package validator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := NewMockValidator(ctrl)
	type User struct {
		Email string `validate:"required,email"`
	}
	v.EXPECT().ValidateStruct(User{Email: "test@example.com"}).Return(nil)
	v.EXPECT().ValidateStruct(User{Email: "invalid-email"}).Return(errors.New("error"))
	v.EXPECT().GetValidationErrors(gomock.Any()).Return(map[string]string{"Email": "Email is Required"})
	tests := []struct {
		name   string
		user   User
		hasErr bool
	}{
		{"Valid email", User{Email: "test@example.com"}, false},
		{"Invalid email", User{Email: "invalid-email"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateStruct(tt.user)
			if tt.hasErr {
				assert.NotNil(t, err)
				validations := v.GetValidationErrors(err)
				assert.Equal(t, "Email is Required", validations["Email"])
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
