package rpc

import (
	"github.com/abdelrahman146/zard/shared/rpc/requests"
)

type RPC interface {
	Request(req requests.Request) (resp []byte, err error)
	Handle(req requests.Request, handler func(req []byte) (resp []byte)) error
}
