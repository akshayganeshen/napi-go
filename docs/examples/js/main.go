package main

import (
	"fmt"
	"time"

	"github.com/akshayganeshen/napi-go"
	"github.com/akshayganeshen/napi-go/entry"
	"github.com/akshayganeshen/napi-go/js"
)

func init() {
	entry.Export("getMap", GetMapHandler)
	entry.Export("getCallback", js.AsCallback(GetCallback))
	entry.Export("getArray", js.AsCallback(GetArray))
	entry.Export("getPromiseResolve", js.AsCallback(GetPromiseResolve))
	entry.Export("getPromiseReject", js.AsCallback(GetPromiseReject))
}

func GetMapHandler(env napi.Env, info napi.CallbackInfo) napi.Value {
	jsEnv := js.AsEnv(env)

	return jsEnv.ValueOf(
		map[string]any{
			"string":    "hello world",
			"number":    123,
			"bool":      false,
			"undefined": jsEnv.Undefined(),
			"null":      nil,
			"function": jsEnv.FuncOf(
				func(env js.Env, this js.Value, args []js.Value) any {
					return "hello world"
				},
			),
		},
	).Value
}

func GetCallback(env js.Env, this js.Value, args []js.Value) any {
	return func(env js.Env, this js.Value, args []js.Value) any {
		return map[string]any{
			"this": this,
			"args": args,
		}
	}
}

func GetArray(env js.Env, this js.Value, args []js.Value) any {
	return []any{
		"hello world",
		123,
		true,
		map[string]any{
			"key": "value",
		},
	}
}

func GetPromiseResolve(env js.Env, this js.Value, args []js.Value) any {
	promise := env.NewPromise()

	go func() {
		time.Sleep(time.Second)
		promise.Resolve("resolved")
	}()

	return promise
}

func GetPromiseReject(env js.Env, this js.Value, args []js.Value) any {
	promise := env.NewPromise()

	go func() {
		time.Sleep(time.Second)
		promise.Reject(fmt.Errorf("rejected"))
	}()

	return promise
}

func main() {}
