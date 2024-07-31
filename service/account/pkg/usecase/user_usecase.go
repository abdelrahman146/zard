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
	UpdateUser(id string, userDto *UpdateUserStruct) (*UserStruct, error)
	UpdateUserEmail(id string, email string) (*UserStruct, error)
	UpdateUserPhone(id string, phone string) (*UserStruct, error)
	UpdateUserPassword(id string, password string) (*UserStruct, error)
	DeactivateUser(id string) (*UserStruct, error)
	ActivateUser(id string) (*UserStruct, error)
	SetUserEmailVerified(id string, verified bool) (*UserStruct, error)
	SetUserPhoneVerified(id string, verified bool) (*UserStruct, error)
	DeleteUser(id string) error
	GetUserByID(id string) (*UserStruct, error)
	GetUserByEmail(email string) (*UserStruct, error)
	GetUserByPhone(phone string) (*UserStruct, error)
	GetAll(page int, limit int) (*shared.List[UserStruct], error)
	GetUsersByOrgID(orgID string, page int, limit int) (*shared.List[UserStruct], error)
	Search(keyword string, page int, limit int) (*shared.List[UserStruct], error)
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

func (uc *userUseCase) ToUserStructList(users []model.User) []UserStruct {
	var userStructList []UserStruct
	for _, user := range users {
		userStructList = append(userStructList, *uc.ToUserStruct(&user))
	}
	return userStructList
}

func (uc *userUseCase) CreateUser(userDto *CreateUserStruct) (*UserStruct, error) {
	if err := uc.toolkit.Validator.ValidateStruct(userDto); err != nil {
		fields := uc.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid user data", fields)
	}
	if userDto.Email != "" {
		if user, _ := uc.userRepo.GetOneByEmail(userDto.Email); user != nil {
			return nil, errs.NewConflictError("email already exists", nil)
		}
	}
	if userDto.Phone != nil || *userDto.Phone != "" {
		if user, _ := uc.userRepo.GetOneByPhone(*userDto.Phone); user != nil {
			return nil, errs.NewConflictError("phone already exists", nil)
		}
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
	if err := uc.userRepo.Create(user); err != nil {
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

func (uc *userUseCase) UpdateUser(id string, userDto *UpdateUserStruct) (*UserStruct, error) {
	if err := uc.toolkit.Validator.ValidateStruct(userDto); err != nil {
		fields := uc.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid user data", fields)
	}
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	if userDto.Name != nil {
		user.Name = *userDto.Name
	}
	if err := uc.userRepo.Save(user); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
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
	user.Email = email
	user.IsEmailVerified = false
	if err := uc.userRepo.Save(user); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) UpdateUserPhone(id string, phone string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	if user.Phone != nil && *user.Phone == phone {
		return uc.ToUserStruct(user), nil
	}
	user.Phone = &phone
	user.IsPhoneVerified = false
	if err := uc.userRepo.Save(user); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) UpdateUserPassword(id string, password string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	if err := uc.userRepo.UpdatePassword(id, password); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) DeactivateUser(id string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	user.Active = false
	if err := uc.userRepo.Save(user); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) ActivateUser(id string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	user.Active = true
	if err := uc.userRepo.Save(user); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) SetUserEmailVerified(id string, verified bool) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	user.IsEmailVerified = verified
	if err := uc.userRepo.Save(user); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) SetUserPhoneVerified(id string, verified bool) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	user.IsPhoneVerified = verified
	if err := uc.userRepo.Save(user); err != nil {
		return nil, errs.NewInternalError("failed to update user", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) DeleteUser(id string) error {
	_, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return errs.NewNotFoundError("User not found", err)
	}
	return uc.userRepo.Delete(id)
}

func (uc *userUseCase) GetUserByID(id string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) GetUserByEmail(email string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByEmail(email)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) GetUserByPhone(phone string) (*UserStruct, error) {
	user, err := uc.userRepo.GetOneByPhone(phone)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	return uc.ToUserStruct(user), nil
}

func (uc *userUseCase) GetAll(page int, limit int) (*shared.List[UserStruct], error) {
	users, total, err := uc.userRepo.GetAll(page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to get users", err)
	}
	userList := uc.ToUserStructList(users)
	return &shared.List[UserStruct]{Items: userList, Total: total, Page: page, Limit: limit}, nil
}

func (uc *userUseCase) GetUsersByOrgID(orgID string, page int, limit int) (*shared.List[UserStruct], error) {
	users, total, err := uc.userRepo.GetAllByOrgID(orgID, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to get users", err)
	}
	userList := uc.ToUserStructList(users)
	return &shared.List[UserStruct]{Items: userList, Total: total, Page: page, Limit: limit}, nil
}

func (uc *userUseCase) Search(keyword string, page int, limit int) (*shared.List[UserStruct], error) {
	users, total, err := uc.userRepo.Search(keyword, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to search users", err)
	}
	userList := uc.ToUserStructList(users)
	return &shared.List[UserStruct]{Items: userList, Total: total, Page: page, Limit: limit}, nil
}
