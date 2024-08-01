package api

import (
	"context"
	"encoding/json"
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/usecase"
	"github.com/abdelrahman146/zard/shared/cache"
	"github.com/abdelrahman146/zard/shared/errs"
	"github.com/gofiber/fiber/v2"
)

type Auth struct{}

func Authorize(ctx context.Context, tokenOwner string, token string, cache cache.Cache) (context.Context, error) {
	if token == "" {
		return nil, errs.NewUnauthorizedError("token is not provided", nil)
	}
	resp, err := cache.Get([]string{"account", "auth", tokenOwner, "tokens", token})
	if err != nil {
		return nil, errs.NewUnauthorizedError("invalid or expired token", err)
	}
	userContext := context.WithValue(ctx, tokenOwner, resp)
	return userContext, nil
}

func (Auth) InitSession(ctx *fiber.Ctx, token string, maxAge int) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   maxAge,
		Secure:   true,
		HTTPOnly: true,
	})
}

func (Auth) AuthorizeUserMiddleware(cache cache.Cache) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Cookies("token")
		userContext, err := Authorize(ctx.UserContext(), "user", token, cache)
		if err != nil {
			return err
		}
		ctx.SetUserContext(userContext)
		return ctx.Next()
	}
}

func (Auth) AuthorizeWorkspaceMiddleware(cache cache.Cache) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Cookies("token")

		wsContext, err := Authorize(ctx.UserContext(), "workspace", token, cache)
		if err != nil {
			return err
		}
		ctx.SetUserContext(wsContext)
		return ctx.Next()
	}
}

func (Auth) AuthorizeBackofficeMiddleware(cache cache.Cache) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Cookies("token")
		boContext, err := Authorize(ctx.UserContext(), "backoffice", token, cache)
		if err != nil {
			return err
		}
		ctx.SetUserContext(boContext)
		return ctx.Next()
	}
}

func (Auth) GetUserFromContext(ctx context.Context) (user *usecase.UserStruct, err error) {
	bytes, ok := ctx.Value("user").([]byte)
	if !ok {
		return nil, errs.NewUnauthorizedError("UnAuthorized", nil)
	}
	user = &usecase.UserStruct{}
	if err = json.Unmarshal(bytes, user); err != nil {
		return nil, errs.NewInternalError("unable to parse user session", err)
	}
	return user, nil
}

func (Auth) GetWorkspaceFromContext(ctx context.Context) (workspace *model.Workspace, err error) {
	bytes, ok := ctx.Value("workspace").([]byte)
	if !ok {
		return nil, errs.NewInternalError("unable to parse workspace session", nil)
	}
	workspace = &model.Workspace{}
	if err = json.Unmarshal(bytes, workspace); err != nil {
		return nil, errs.NewInternalError("unable to parse workspace session", err)
	}
	return workspace, nil
}

//func GetBackofficeFromContext(ctx context.Context) (backoffice *usecase.BackofficeStruct, err error) {
//	bytes, ok := ctx.Value("backoffice").([]byte)
//	if !ok {
//		return nil, errs.NewInternalError("unable to parse backoffice session", nil)
//	}
//	backoffice = &usecase.BackofficeStruct{}
//	if err = json.Unmarshal(bytes, backoffice); err != nil {
//		return nil, errs.NewInternalError("unable to parse backoffice session", err)
//	}
//	return backoffice, nil
//}
