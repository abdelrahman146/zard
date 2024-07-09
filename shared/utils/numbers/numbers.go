package numbers

import (
	"errors"
	"math"
	"math/rand"
)

type Struct struct{}

func (Struct) GenerateRandomDigits(maxDigits int) (int, error) {
	if maxDigits <= 0 {
		return 0, errors.New("maxDigits must be greater than 0")
	}
	mn := math.Pow(10, float64(maxDigits-1))
	mx := math.Pow(10, float64(maxDigits)) - 1

	return rand.Intn(int(mx-mn)) + int(mn), nil
}

func (Struct) GenerateRandomInt(min int, max int) (int, error) {
	if min > max {
		return 0, errors.New("min must be less than or equal to max")
	}
	if max-min <= 0 {
		return 0, errors.New("max-min must be greater than 0")
	}
	return rand.Intn(max-min) + min, nil
}

func (Struct) Round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
