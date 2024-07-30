package api

import "github.com/abdelrahman146/zard/shared/errs"

type Response struct{}

type ErrorResponseBody struct {
	Message  string `json:"message"`
	HttpCode int    `json:"httpCode"`
	Code     string `json:"code"`
	Reason   string `json:"reason"`
}

type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   ErrorResponseBody `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
}

func (Response) NewErrorResponse(err error) (httpCode int, resp ErrorResponse) {
	customErr := errs.HandleError(err)
	httpCode = customErr.HttpCode
	resp = ErrorResponse{
		Success: false,
		Error: ErrorResponseBody{
			Message:  customErr.Desc,
			HttpCode: customErr.HttpCode,
			Code:     customErr.Code,
			Reason:   customErr.Original.Error(),
		},
	}
	return httpCode, resp
}

func (Response) NewSuccessResponse(result interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Result:  result,
	}
}
