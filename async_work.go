package napi

import (
	"unsafe"
)

type AsyncWork struct {
	Handle unsafe.Pointer
	ID     NapiGoAsyncWorkID
}

type AsyncExecuteCallback func(env Env)

type AsyncCompleteCallback func(env Env, status Status)
