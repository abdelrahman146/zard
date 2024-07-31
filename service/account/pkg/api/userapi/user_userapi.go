package userapi

import (
	"github.com/abdelrahman146/zard/service/account/pkg/usecase"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
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

func NewUserUserApi(app *fiber.App, toolkit *shared.Toolkit, usecases *usecase.AccountUseCases) {
	api := &userUserApi{
		toolkit:  toolkit,
		usecases: usecases,
	}
	api.Setup(app)
}

type userUserApi struct {
	toolkit  *shared.Toolkit
	usecases *usecase.AccountUseCases
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

func (api *userUserApi) Register(ctx *fiber.Ctx) error {
	body := &usecase.CreateUserStruct{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided or invalid", err)
	}
	user, err := api.usecases.UserUseCase.CreateUser(body)
	if err != nil {
		return err
	}
	go func() {
		api.usecases.UserUseCase.Em
	}()
	return ctx.Status(fiber.StatusCreated).JSON(shared.Api.Response.NewSuccessResponse(user))
}
