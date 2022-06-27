package js

import (
	"github.com/akshayganeshen/napi-go"
)

type Callback = func(env Env, this Value, args []Value) any

func AsCallback(fn Callback) napi.Callback {
	return func(env napi.Env, info napi.CallbackInfo) napi.Value {
		cbInfo, st := napi.GetCbInfo(env, info)
		if st != napi.StatusOK {
			panic(napi.StatusError(st))
		}

		jsEnv := AsEnv(env)
		this := Value{
			Env:   jsEnv,
			Value: cbInfo.This,
		}
		args := make([]Value, len(cbInfo.Args))
		for i, cbArg := range cbInfo.Args {
			args[i] = Value{
				Env:   jsEnv,
				Value: cbArg,
			}
		}

		result := fn(jsEnv, this, args)
		return jsEnv.ValueOf(result).Value
	}
}
