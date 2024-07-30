package userapi

import (
	"github.com/abdelrahman146/zard/service/account/pkg/usecase"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
	"github.com/abdelrahman146/zard/shared/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthUserApi interface {
	Setup(app *fiber.App) error
	LoginWithEmailAndPassword(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	GenerateOTPForEmail(ctx *fiber.Ctx) error
	GenerateOTPForPhone(ctx *fiber.Ctx) error
	VerifyOTP(ctx *fiber.Ctx) error
	LogoutFromAllUserSessions(ctx *fiber.Ctx) error
}

func NewAuthUserApi(app *fiber.App, toolkit shared.Toolkit, auth usecase.AuthUseCase, user usecase.UserUseCase) {
	api := &authUserApi{
		toolkit: toolkit,
		auth:    auth,
		user:    user,
	}
	if err := api.Setup(app); err != nil {
		logger.GetLogger().Panic("failed to setup auth user api", logger.Field("error", err))
	}
}

type authUserApi struct {
	toolkit shared.Toolkit
	auth    usecase.AuthUseCase
	user    usecase.UserUseCase
}

func (api *authUserApi) Setup(app *fiber.App) error {
	group := app.Group("/auth")
	group.Post("/login", api.LoginWithEmailAndPassword)
	return nil
}

func (api *authUserApi) LoginWithEmailAndPassword(ctx *fiber.Ctx) error {
	body := &LoginWithEmailAndPasswordRequest{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided", err)
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
	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   maxAge,
		Secure:   true,
		HTTPOnly: true,
	})
	resp := shared.Api.Response.NewSuccessResponse(user)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (api *authUserApi) Logout(ctx *fiber.Ctx) error {
	ctx.ClearCookie("token")
	token := ctx.Cookies("token")
	if err := api.auth.RevokeToken(token); err != nil {
		return err
	}
	resp := shared.Api.Response.NewSuccessResponse(nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (api *authUserApi) GenerateOTPForEmail(ctx *fiber.Ctx) error {
	body := &GenerateOTPForEmailRequest{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided", err)
	}
	if err := api.toolkit.Validator.ValidateStruct(body); err != nil {
		fields := api.toolkit.Validator.GetValidationErrors(err)
		return errs.NewValidationError("Invalid request body", fields)
	}
	if err := api.auth.CreateAndSendOTP("email", body.Email); err != nil {
		return err
	}
	resp := shared.Api.Response.NewSuccessResponse(nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (api *authUserApi) VerifyOTP(ctx *fiber.Ctx) error {
	body := &VerifyOTPForEmailRequest{}
	if err := ctx.BodyParser(body); err != nil {
		return errs.NewBadRequestError("Request body is not provided", err)
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
