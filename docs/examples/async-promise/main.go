package main

import (
	"fmt"
	"time"

	"github.com/akshayganeshen/napi-go"
	"github.com/akshayganeshen/napi-go/entry"
)

func init() {
	entry.Export("getPromise", GetPromiseHandler)
}

func GetPromiseHandler(env napi.Env, info napi.CallbackInfo) napi.Value {
	result, _ := napi.CreatePromise(env)
	asyncResourceName, _ := napi.CreateStringUtf8(
		env,
		"napi-go/async-promise-example",
	)

	var asyncWork napi.AsyncWork
	asyncWork, _ = napi.CreateAsyncWork(
		env,
		nil, asyncResourceName,
		func(env napi.Env) {
			fmt.Printf("AsyncExecuteCallback(start)\n")
			defer fmt.Printf("AsyncExecuteCallback(stop)\n")
			time.Sleep(time.Second)
		},
		func(env napi.Env, status napi.Status) {
			defer napi.DeleteAsyncWork(env, asyncWork)

			if status == napi.StatusCancelled {
				fmt.Printf("AsyncCompleteCallback(cancelled)\n")
				return
			}

			fmt.Printf("AsyncCompleteCallback\n")
			resolution, _ := napi.CreateStringUtf8(env, "resolved")
			napi.ResolveDeferred(env, result.Deferred, resolution)
		},
	)
	napi.QueueAsyncWork(env, asyncWork)

	return result.Value
}

func main() {}
