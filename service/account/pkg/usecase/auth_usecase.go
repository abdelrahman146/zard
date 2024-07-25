package usecase

import (
	"encoding/json"
	"errors"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/cache"
	"github.com/abdelrahman146/zard/shared/config"
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
	AuthenticateToken(token string) (user *UserResponse, err error)
	AuthenticateWorkspaceByApiKey(apiKey string) (id string, err error)
}

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrInvalidApiKey          = errors.New("invalid apikey")
)

type authUseCase struct {
	userRepo    repo.UserRepo
	wrkRepo     repo.WorkspaceRepo
	conf        config.Config
	cacheClient cache.Cache
}

func NewAuthUseCase(userRepo repo.UserRepo, wrkRepo repo.WorkspaceRepo, conf config.Config, cacheClient cache.Cache) AuthUseCase {
	return &authUseCase{
		userRepo:    userRepo,
		wrkRepo:     wrkRepo,
		conf:        conf,
		cacheClient: cacheClient,
	}
}

func (uc *authUseCase) AuthenticateUserByEmailPassword(email, password string) (token string, user *UserResponse, err error) {
	userModel, err := uc.userRepo.GetOneByEmail(email)
	if err != nil {
		return "", nil, ErrInvalidEmailOrPassword
	}
	hashedPassword := shared.Utils.Auth.Encrypt(password, uc.conf.GetString("app.secret"))
	if userModel.Password != nil && *userModel.Password != hashedPassword {
		return "", nil, ErrInvalidEmailOrPassword
	}
	token = shared.Utils.Auth.CreateToken("ztkn", userModel.ID, uc.conf.GetString("app.secret"))
	userJson, _ := json.Marshal(userModel)
	ttl := time.Second * time.Duration(uc.conf.GetInt("app.auth.tokenTTL"))
	if err = uc.cacheClient.Set([]string{"account", "auth", "user", "tokens", token}, userJson, ttl); err != nil {
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

func (uc *authUseCase) AuthenticateToken(token string) (user *UserResponse, err error) {
	userJson, err := uc.cacheClient.Get([]string{"account", "auth", "user", "tokens", token})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err = json.Unmarshal(userJson, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *authUseCase) AuthenticateWorkspaceByApiKey(apiKey string) (id string, err error) {
	secret := uc.conf.GetString("app.secret")
	if ok := shared.Utils.Auth.ValidateToken(apiKey, secret); !ok {
		return "", ErrInvalidApiKey
	}
	if resp, err := uc.cacheClient.Get([]string{"account", "auth", "workspace", "apikeys", apiKey}); err == nil {
		return string(resp), nil
	}
	workspace, err := uc.wrkRepo.GetOneByApiKey(apiKey)
	if err != nil {
		return "", ErrInvalidApiKey
	}
	if err = uc.cacheClient.Set([]string{"account", "auth", "workspace", "apikeys", apiKey}, []byte(workspace.ID), time.Second*time.Duration(uc.conf.GetInt("app.auth.apiKeyTTL"))); err != nil {
		return "", err
	}
	return workspace.ID, nil
}
