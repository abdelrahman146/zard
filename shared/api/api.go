package api

type Struct struct {
	Response Response
	Auth     Auth
}

var Api = Struct{
	Response: Response{},
	Auth:     Auth{},
}
