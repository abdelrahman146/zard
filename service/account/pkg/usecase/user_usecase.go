package usecase

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
	"github.com/abdelrahman146/zard/shared/logger"
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
	"time"
)

type UserUseCase interface {
	CreateUser(userDto *CreateUserStruct) (*UserStruct, error)
	UpdateUserName(id string, name string) (*UserStruct, error)
	UpdateUserEmail(id string, email string) (*UserStruct, error)
	UpdateUserPhone(id string, phone string) (*UserStruct, error)
	UpdateUserPassword(id string, password string) (*UserStruct, error)
	SetUserActivated(id string, active bool) (*UserStruct, error)
	SetUserEmailVerified(id string, verified bool) (*UserStruct, error)
	SetUserPhoneVerified(id string, verified bool) (*UserStruct, error)
	DeleteUser(id string) error
	GetUserByID(id string) (*UserStruct, error)
	GetUserByEmail(email string) (*UserStruct, error)
	GetUserByPhone(phone string) (*UserStruct, error)
	GetAll(page int, limit int) (*shared.List[model.User], error)
	GetUsersByOrgID(orgID string, page int, limit int) (*shared.List[model.User], error)
}

func NewUserUseCase(toolkit shared.Toolkit, userRepo repo.UserRepo) UserUseCase {
	return &userUseCase{
		toolkit:  toolkit,
		userRepo: userRepo,
	}
}

type userUseCase struct {
	toolkit  shared.Toolkit
	userRepo repo.UserRepo
}

func (uc *userUseCase) ToUserStruct(user *model.User) *UserStruct {
	return &UserStruct{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Phone:           user.Phone,
		Active:          user.Active,
		IsEmailVerified: user.IsEmailVerified,
		IsPhoneVerified: user.IsPhoneVerified,
		OrgID:           user.OrgID,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		DeletedAt:       user.DeletedAt,
	}
}

func (uc *userUseCase) CreateUser(userDto *CreateUserStruct) (*UserStruct, error) {
	if err := uc.toolkit.Validator.ValidateStruct(userDto); err != nil {
		fields := uc.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid user data", fields)
	}
	user := &model.User{
		Name:            userDto.Name,
		Email:           userDto.Email,
		Phone:           userDto.Phone,
		Password:        userDto.Password,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		Active:          true,
		OrgID:           userDto.OrgID,
	}
	if err := uc.userRepo.Insert(user); err != nil {
		return nil, errs.NewInternalError("failed to create user", err)
	}
	userCreatedMessage := &messages.UserCreatedMessage{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Timestamp: time.Now(),
	}
	if err := uc.toolkit.PubSub.Publish(userCreatedMessage); err != nil {
		logger.GetLogger().Error("failed to publish user created message", logger.Field("error", err))
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) UpdateUserName(id string, name string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	if err := uc.userRepo.UpdateName(id, name); err != nil {
		return nil, errs.NewInternalError("failed to update user name", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) UpdateUserEmail(id string, email string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	if user.Email == email {
		return uc.ToUserStruct(user), nil
	}
	if err := uc.userRepo.UpdateEmail(id, email); err != nil {
		return nil, errs.NewInternalError("failed to update user email", err)
	}
	return uc.ToUserStruct(user), nil
}
