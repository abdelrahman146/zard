package main

import (
	"fmt"
	"github.com/abdelrahman146/zard/shared/errs"
)

func r() error {
	return errs.NewBadRequestError("Hello Errors", nil)
}

func main() {
	fmt.Println(" Hello From Account Service")
	err := r()
	e := errs.HandleError(err)
	fmt.Println(e.Original)
}
