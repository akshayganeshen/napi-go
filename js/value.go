package js

import (
	"github.com/akshayganeshen/napi-go"
)

type Value struct {
	Env   Env
	Value napi.Value
}

func (v Value) GetEnv() Env {
	return v.Env
}
