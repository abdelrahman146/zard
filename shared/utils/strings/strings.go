package strings

import (
	"github.com/lucsky/cuid"
	"strconv"
)

type Struct struct{}

func (Struct) Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < len(runes)/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (Struct) Parse(val string) interface{} {
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(val); err == nil {
		return b
	}
	return val
}

func (Struct) IsEmpty(s string) bool {
	return len(s) == 0
}

func (Struct) Cuid() string {
	return cuid.New()
}
