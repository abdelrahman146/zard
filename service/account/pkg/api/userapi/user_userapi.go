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
	cache := api.toolkit.Cache
	v1group.Post("/register", api.RegisterWithEmailAndPassword)
	v1group.Post("/email/verify/otp", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.VerifyUserEmailSendOtp)
	v1group.Post("/email/verify/otp/verify", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.VerifyUserEmailValidateOtp)
	v1group.Get("/email/verify/otp/verify/:hash", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.VerifyUserEmailValidateHash)
	v1group.Put("/", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.UpdateUser)
	v1group.Put("/email", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.UpdateUserEmail)
	v1group.Put("/password", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.UpdateUserPassword)
}

func (api *userUserApi) RegisterWithEmailAndPassword(ctx *fiber.Ctx) error {
	body := &usecase.CreateUserStruct{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided or invalid", err)
	}
	user, err := api.usecases.UserUseCase.CreateUser(body)
	if err != nil {
		return err
	}
	_, _ = api.usecases.AuthUseCase.CreateAndSendOTP("email", "verify", user.Email)
	token, err := api.usecases.AuthUseCase.CreateUserToken(user)
	shared.Api.Auth.InitSession(ctx, token, api.toolkit.Conf.GetInt("app.auth.tokenTTL"))
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(shared.Api.Response.NewSuccessResponse(user))
}

func (api *userUserApi) VerifyUserEmailSendOtp(ctx *fiber.Ctx) error {
	user, err := shared.Api.Auth.GetUserFromContext(ctx.UserContext())
	if err != nil {
		return err
	}
	_, err := api.usecases.AuthUseCase.CreateAndSendOTP("email", "verify", user.Email)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(shared.Api.Response.NewSuccessResponse(nil))
}
