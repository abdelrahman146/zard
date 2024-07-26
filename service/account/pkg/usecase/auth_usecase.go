package usecase

import (
	"encoding/json"
	"errors"
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
	"gorm.io/gorm"
	"time"
)

type UserResponse struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	Phone           *string        `json:"phone"`
	IsEmailVerified bool           `json:"isEmailVerified"`
	IsPhoneVerified bool           `json:"isPhoneVerified"`
	Active          bool           `json:"active"`
	OrgID           string         `json:"orgId"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt"`
}

type AuthUseCase interface {
	AuthenticateUserByEmailPassword(email, password string) (token string, user *UserResponse, err error)
	SendMagicLink(email string, withReset bool) error
	AuthenticateWithMagicLink(magiclinkToken string) (token string, user *UserResponse, err error)
	AuthenticateToken(token string) (user *UserResponse, err error)
	AuthenticateWorkspaceByApiKey(apiKey string) (id string, err error)
}

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrSomethingWentWrong     = errors.New("something went wrong")
	ErrUserNotFound           = errors.New("user not found")
	ErrInactiveUser           = errors.New("user is inactive")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrInvalidApiKey          = errors.New("invalid apikey")
)

type authUseCase struct {
	toolkit  shared.Toolkit
	userRepo repo.UserRepo
	wrkRepo  repo.WorkspaceRepo
}

func NewAuthUseCase(toolkit shared.Toolkit, userRepo repo.UserRepo, wrkRepo repo.WorkspaceRepo) AuthUseCase {
	return &authUseCase{
		toolkit:  toolkit,
		userRepo: userRepo,
		wrkRepo:  wrkRepo,
	}
}

func (uc *authUseCase) createUserSession(userModel *model.User) (token string, user *UserResponse, err error) {
	userJson, _ := json.Marshal(userModel)
	token = shared.Utils.Auth.CreateToken("ztkn", userModel.ID, uc.toolkit.Conf.GetString("app.secret"))
	ttl := time.Second * time.Duration(uc.toolkit.Conf.GetInt("app.auth.tokenTTL"))
	if err := uc.toolkit.Cache.Set([]string{"account", "auth", "user", "tokens", token}, userJson, ttl); err != nil {
		return "", nil, err
	}
	user = &UserResponse{
		ID:              userModel.ID,
		Name:            userModel.Name,
		Email:           userModel.Email,
		Phone:           userModel.Phone,
		Active:          userModel.Active,
		IsEmailVerified: userModel.IsEmailVerified,
		IsPhoneVerified: userModel.IsPhoneVerified,
		OrgID:           userModel.OrgID,
		CreatedAt:       userModel.CreatedAt,
		UpdatedAt:       userModel.UpdatedAt,
		DeletedAt:       userModel.DeletedAt,
	}
	return token, user, nil
}

func (uc *authUseCase) AuthenticateUserByEmailPassword(email, password string) (token string, user *UserResponse, err error) {
	userModel, err := uc.userRepo.GetOneByEmail(email)
	if err != nil {
		return "", nil, ErrUserNotFound
	}
	if userModel.Password == nil {
		return "", nil, ErrInvalidEmailOrPassword
	}
	if userModel.Active == false {
		return "", nil, ErrInactiveUser
	}
	if ok := shared.Utils.Auth.Compare(*userModel.Password, password, uc.toolkit.Conf.GetString("app.secret")); !ok {
		return "", nil, ErrInvalidEmailOrPassword
	}
	return uc.createUserSession(userModel)
}

func (uc *authUseCase) AuthenticateToken(token string) (user *UserResponse, err error) {
	userJson, err := uc.toolkit.Cache.Get([]string{"account", "auth", "user", "tokens", token})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err = json.Unmarshal(userJson, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *authUseCase) SendMagicLink(email string, withReset bool) error {
	user, err := uc.userRepo.GetOneByEmail(email)
	if err != nil {
		return err
	}
	if user.Active == false {
		return ErrInactiveUser
	}
	token := shared.Utils.Auth.CreateToken("zml", user.ID, uc.toolkit.Conf.GetString("app.secret"))
	ttl := time.Second * time.Duration(uc.toolkit.Conf.GetInt("app.auth.magicLinkTTL"))
	if err = uc.toolkit.Cache.Set([]string{"account", "auth", "user", "magiclinks", token}, []byte(user.ID), ttl); err != nil {
		return err
	}
	msg := messages.MagicLinkMessage{
		Email:     email,
		Name:      user.Name,
		Token:     token,
		Timestamp: time.Now(),
		WithReset: withReset,
	}

	return uc.toolkit.PubSub.Publish(&msg)
}

func (uc *authUseCase) AuthenticateWithMagicLink(magiclinkToken string) (token string, user *UserResponse, err error) {
	userId, err := uc.toolkit.Cache.Get([]string{"account", "auth", "user", "magiclinks", magiclinkToken})
	if err != nil {
		return "", nil, ErrInvalidToken
	}
	userModel, err := uc.userRepo.GetOneByID(string(userId))
	if err != nil {
		return "", nil, err
	}
	return uc.createUserSession(userModel)
}

func (uc *authUseCase) AuthenticateWorkspaceByApiKey(apiKey string) (id string, err error) {
	secret := uc.toolkit.Conf.GetString("app.secret")
	if ok := shared.Utils.Auth.ValidateToken(apiKey, secret); !ok {
		return "", ErrInvalidApiKey
	}
	if resp, err := uc.toolkit.Cache.Get([]string{"account", "auth", "workspace", "apikeys", apiKey}); err == nil {
		return string(resp), nil
	}
	workspace, err := uc.wrkRepo.GetOneByApiKey(apiKey)
	if err != nil {
		return "", ErrInvalidApiKey
	}
	if err = uc.toolkit.Cache.Set([]string{"account", "auth", "workspace", "apikeys", apiKey}, []byte(workspace.ID), time.Second*time.Duration(uc.toolkit.Conf.GetInt("app.auth.apiKeyTTL"))); err != nil {
		return "", err
	}
	return workspace.ID, nil
}
