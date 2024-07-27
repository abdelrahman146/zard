package shared

import (
	"github.com/abdelrahman146/zard/shared/cache"
	"github.com/abdelrahman146/zard/shared/config"
	"github.com/abdelrahman146/zard/shared/pubsub"
	"github.com/abdelrahman146/zard/shared/rpc"
	"github.com/abdelrahman146/zard/shared/utils"
	"github.com/abdelrahman146/zard/shared/validator"
)

var Utils = utils.Utils

type Toolkit struct {
	Rpc       rpc.RPC
	PubSub    pubsub.PubSub
	Cache     cache.Cache
	Conf      config.Config
	Validator validator.Validator
}

type List[T any] struct {
	Items []T   `json:"items"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
