package userapi

import (
	"github.com/abdelrahman146/zard/service/account/pkg/usecase"
	"github.com/abdelrahman146/zard/shared"
	"github.com/gofiber/fiber/v2"
)

type UserUserApi interface {
	Setup(app *fiber.App)
	Register(ctx *fiber.Ctx) error
	UpdateUser(ctx *fiber.Ctx) error
	UpdateUserEmail(ctx *fiber.Ctx) error
	VerifyUserEmailSendOtp(ctx *fiber.Ctx) error
	VerifyUserEmailSubmitOtp(ctx *fiber.Ctx) error
	UpdateUserPassword(ctx *fiber.Ctx) error
}

func NewUserUserApi(app *fiber.App, toolkit shared.Toolkit, user usecase.UserUseCase) {
	api := &userUserApi{
		toolkit: toolkit,
		user:    user,
	}
	api.Setup(app)
}

type userUserApi struct {
	toolkit shared.Toolkit
	user    usecase.UserUseCase
}

func (api *userUserApi) Setup(app *fiber.App) {
	v1group := app.Group("/v1/user")
	v1group.Post("/", api.Register)
	v1group.Put("/", shared.Api.Auth.AuthorizeUserMiddleware(api.toolkit), api.UpdateUser)
	v1group.Put("/email", shared.Api.Auth.AuthorizeUserMiddleware(api.toolkit), api.UpdateUserEmail)
	v1group.Post("/email/verify/otp", shared.Api.Auth.AuthorizeUserMiddleware(api.toolkit), api.VerifyUserEmail)
	v1group.Post("/email/verify/otp/verify", shared.Api.Auth.AuthorizeUserMiddleware(api.toolkit), api.VerifyUserEmail)
	v1group.Put("/password", shared.Api.Auth.AuthorizeUserMiddleware(api.toolkit), api.UpdateUserPassword)
}
