package utils

import (
	"github.com/abdelrahman146/zard/shared/utils/auth"
	"github.com/abdelrahman146/zard/shared/utils/numbers"
	"github.com/abdelrahman146/zard/shared/utils/strings"
)

type Struct struct {
	Numbers numbers.Struct
	Strings strings.Struct
	Auth    auth.Struct
}

var Utils = Struct{
	Numbers: numbers.Struct{},
	Strings: strings.Struct{},
	Auth:    auth.Struct{},
}
