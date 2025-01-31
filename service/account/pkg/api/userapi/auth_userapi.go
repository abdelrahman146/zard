package userapi

import (
	"github.com/abdelrahman146/zard/service/account/pkg/usecase"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
	"github.com/gofiber/fiber/v2"
)

type AuthUserApi interface {
	Setup(app *fiber.App)
	LoginWithEmailAndPassword(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	GenerateOTPForEmail(ctx *fiber.Ctx) error
	VerifyOTP(ctx *fiber.Ctx) error
	LogoutFromAllUserSessions(ctx *fiber.Ctx) error
	GetUserSession(ctx *fiber.Ctx) error
}

func NewAuthUserApi(app *fiber.App, toolkit *shared.Toolkit, auth usecase.AuthUseCase) {
	api := &authUserApi{
		toolkit: toolkit,
		auth:    auth,
	}
	api.SetupV1(app)
}

type authUserApi struct {
	toolkit *shared.Toolkit
	auth    usecase.AuthUseCase
}

func (api *authUserApi) SetupV1(app *fiber.App) {
	v1group := app.Group("/v1/auth")
	cache := api.toolkit.Cache
	v1group.Get("/user", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.GetUserSession)
	v1group.Post("/login", api.LoginWithEmailAndPassword)
	v1group.Post("/otp/email", api.GenerateOTPForEmail)
	v1group.Post("/otp/verify", api.VerifyOTP)
	v1group.Post("/logout", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.Logout)
	v1group.Post("/logout/all", shared.Api.Auth.AuthorizeUserMiddleware(cache), api.LogoutFromAllUserSessions)
}

func (api *authUserApi) LoginWithEmailAndPassword(ctx *fiber.Ctx) error {
	body := &LoginWithEmailAndPasswordRequest{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided or invalid", err)
	}
	if err := api.toolkit.Validator.ValidateStruct(body); err != nil {
		fields := api.toolkit.Validator.GetValidationErrors(err)
		return errs.NewValidationError("Invalid request body", fields)
	}
	token, user, err := api.auth.AuthenticateUserByEmailPassword(body.Email, body.Password)
	if err != nil {
		return err
	}
	maxAge := api.toolkit.Conf.GetInt("app.auth.tokenTTL")
	shared.Api.Auth.InitSession(ctx, token, maxAge)
	resp := shared.Api.Response.NewSuccessResponse(user)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (api *authUserApi) Logout(ctx *fiber.Ctx) error {
	token := ctx.Cookies("token")
	ctx.ClearCookie("token")
	if err := api.auth.RevokeToken(token); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(shared.Api.Response.NewSuccessResponse(nil))
}

func (api *authUserApi) GenerateOTPForEmail(ctx *fiber.Ctx) error {
	body := &GenerateOTPForEmailRequest{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided or invalid", err)
	}
	if err := api.toolkit.Validator.ValidateStruct(body); err != nil {
		fields := api.toolkit.Validator.GetValidationErrors(err)
		return errs.NewValidationError("Invalid request body", fields)
	}
	maxAge, err := api.auth.CreateAndSendOTP("email", "confirm", body.Email)
	if err != nil {
		return err
	}
	resp := shared.Api.Response.NewSuccessResponse(fiber.Map{"maxAge": maxAge})
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (api *authUserApi) VerifyOTP(ctx *fiber.Ctx) error {
	body := &VerifyOTPForEmailRequest{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided or invalid", err)
	}
	if err := api.toolkit.Validator.ValidateStruct(body); err != nil {
		fields := api.toolkit.Validator.GetValidationErrors(err)
		return errs.NewValidationError("Invalid request body", fields)
	}
	if err := api.auth.VerifyOTP(body.Value, body.Otp); err != nil {
		return err
	}
	resp := shared.Api.Response.NewSuccessResponse(nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (api *authUserApi) LogoutFromAllUserSessions(ctx *fiber.Ctx) error {
	user, err := shared.Api.Auth.GetUserFromContext(ctx.UserContext())
	if err != nil {
		return err
	}
	if err := api.auth.RevokeAllUserTokens(user.ID); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(shared.Api.Response.NewSuccessResponse(nil))
}

func (api *authUserApi) GetUserSession(ctx *fiber.Ctx) error {
	user, err := shared.Api.Auth.GetUserFromContext(ctx.UserContext())
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(shared.Api.Response.NewSuccessResponse(user))
}
