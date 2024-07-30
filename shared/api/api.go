package api

type Struct struct {
	Response Response
}

var Api = Struct{
	Response: Response{},
}
