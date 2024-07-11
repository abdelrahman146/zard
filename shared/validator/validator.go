package validator

type Validator interface {
	ValidateStruct(s interface{}) error
	GetValidationErrors(err error) map[string]string
}
