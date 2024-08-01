package usecase

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
	"strconv"
	"time"
)

type AuthUseCase interface {
	AuthenticateUserByEmailPassword(email, password string) (token string, user *UserStruct, err error)
	CreateUserToken(user *UserStruct) (token string, err error)
	CreateAndSendOTP(target, reason, value string) (maxAge time.Duration, err error)
	VerifyOTP(expectedVal, otp string) (err error)
	AuthenticateToken(token string) (user *UserStruct, err error)
	AuthenticateWorkspaceByApiKey(apiKey string) (id string, err error)
	RevokeToken(token string) (err error)
	RevokeAllUserTokens(userID string) (err error)
}

func NewAuthUseCase(toolkit shared.Toolkit, userRepo repo.UserRepo, wrkRepo repo.WorkspaceRepo) AuthUseCase {
	return &authUseCase{
		toolkit:  toolkit,
		userRepo: userRepo,
		wrkRepo:  wrkRepo,
	}
}

type authUseCase struct {
	toolkit  shared.Toolkit
	userRepo repo.UserRepo
	wrkRepo  repo.WorkspaceRepo
}

func (uc *authUseCase) ToUserStruct(user *model.User) *UserStruct {
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

func (uc *authUseCase) CreateUserToken(user *UserStruct) (token string, err error) {
	userJson, err := json.Marshal(user)
	token = shared.Utils.Auth.CreateToken("ztkn", user.ID, uc.toolkit.Conf.GetString("app.secret"))
	ttl := time.Second * time.Duration(uc.toolkit.Conf.GetInt("app.auth.tokenTTL"))
	if err := uc.toolkit.Cache.Set([]string{"account", "auth", "user", "tokens", token}, userJson, ttl); err != nil {
		return "", errs.NewInternalError("unable to create user session", err)
	}
	return token, nil
}

func (uc *authUseCase) AuthenticateUserByEmailPassword(email, password string) (token string, user *UserStruct, err error) {
	userModel, err := uc.userRepo.GetOneByEmail(email)
	if err != nil {
		return "", nil, errs.NewBadRequestError("invalid email", err)
	}
	if userModel.Password == nil {
		return "", nil, errs.NewBadRequestError("invalid password", nil)
	}
	if userModel.Active == false {
		return "", nil, errs.NewForbiddenError("inactive user", nil)
	}
	if ok := shared.Utils.Auth.Compare(*userModel.Password, password, uc.toolkit.Conf.GetString("app.secret")); !ok {
		return "", nil, errs.NewBadRequestError("invalid password", nil)
	}
	user = uc.ToUserStruct(userModel)
	token, err = uc.CreateUserToken(user)
	return token, user, err
}

func (uc *authUseCase) CreateAndSendOTP(target, reason, value string) (maxAge time.Duration, err error) {
	_, err = uc.toolkit.Cache.Get([]string{"account", "auth", "otp", value})
	if err == nil {
		return 0, errs.NewBadRequestError("otp already exists", nil)
	}
	otpNum, err := shared.Utils.Numbers.GenerateRandomDigits(6)
	if err != nil {
		return 0, errs.NewInternalError("Unable to create otp", err)
	}
	otp := strconv.Itoa(otpNum)
	ttl := time.Duration(uc.toolkit.Conf.GetInt("app.auth.otpTTL"))
	if err = uc.toolkit.Cache.Set([]string{"account", "auth", "otp", value}, []byte(otp), ttl); err != nil {
		return 0, errs.NewInternalError("unable to create otp", err)
	}
	if err := uc.toolkit.PubSub.Publish(&messages.AuthOTPCreated{
		Value:     value,
		Target:    target,
		Reason:    reason,
		Otp:       otp,
		Ttl:       ttl,
		Timestamp: time.Now(),
	}); err != nil {
		return 0, errs.NewInternalError("Unable to publish otp created message", err)
	}
	return ttl, nil
}

func (uc *authUseCase) VerifyOTP(expectedVal, otp string) (err error) {
	res, err := uc.toolkit.Cache.Get([]string{"account", "auth", "otp", expectedVal})
	if err != nil {
		return errs.NewUnauthorizedError("invalid or expired otp", err)
	}
	if string(res) != otp {
		return errs.NewUnauthorizedError("invalid otp", nil)
	}
	if err = uc.toolkit.Cache.Delete([]string{"account", "auth", "otp", expectedVal}); err != nil {
		return errs.NewInternalError("unable to delete otp", err)
	}
	return nil
}

func (uc *authUseCase) AuthenticateToken(token string) (user *UserStruct, err error) {
	userJson, err := uc.toolkit.Cache.Get([]string{"account", "auth", "user", "tokens", token})
	if err != nil {
		return nil, errs.NewUnauthorizedError("invalid or expired token", err)
	}
	if err = json.Unmarshal(userJson, &user); err != nil {
		return nil, errs.NewInternalError("unable to parse user session", err)
	}
	return user, nil
}

func (uc *authUseCase) AuthenticateWorkspaceByApiKey(apiKey string) (id string, err error) {
	secret := uc.toolkit.Conf.GetString("app.secret")
	if ok := shared.Utils.Auth.ValidateToken(apiKey, secret); !ok {
		return "", errs.NewBadRequestError("invalid api key", nil)
	}
	if resp, err := uc.toolkit.Cache.Get([]string{"account", "auth", "workspace", "tokens", apiKey}); err == nil {
		return string(resp), nil
	}
	workspace, err := uc.wrkRepo.GetOneByApiKey(apiKey)
	if err != nil {
		return "", errs.NewUnauthorizedError("invalid api key", err)
	}
	if err = uc.toolkit.Cache.Set([]string{"account", "auth", "workspace", "tokens", apiKey}, []byte(workspace.ID), time.Second*time.Duration(uc.toolkit.Conf.GetInt("app.auth.apiKeyTTL"))); err != nil {
		return "", errs.NewInternalError("unable to create workspace session", err)
	}
	return workspace.ID, nil
}

func (uc *authUseCase) RevokeToken(token string) (err error) {
	if err = uc.toolkit.Cache.Delete([]string{"account", "auth", "user", "tokens", token}); err != nil {
		return errs.NewInternalError("unable to revoke token", err)
	}
	return nil
}

func (uc *authUseCase) RevokeAllUserTokens(userID string) (err error) {
	tokens, err := uc.toolkit.Cache.Keys([]string{"account", "auth", "user", "tokens"})
	if err != nil {
		return errs.NewInternalError("unable to revoke tokens", err)
	}
	for _, token := range tokens {
		userJson, err := uc.toolkit.Cache.Get([]string{"account", "auth", "user", "tokens", token})
		if err != nil {
			continue
		}
		var user UserStruct
		if err = json.Unmarshal(userJson, &user); err != nil {
			continue
		}
		if user.ID == userID {
			if err = uc.toolkit.Cache.Delete([]string{"account", "auth", "user", "tokens", token}); err != nil {
				return errs.NewInternalError("unable to revoke token", err)
			}
		}
	}
	return nil
}
