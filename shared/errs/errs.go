package errs

type CustomError interface {
	error
	Code() string
	HttpCode() int
	Original() error
}

type ValidationError interface {
	CustomError
	Fields() map[string]string
}
