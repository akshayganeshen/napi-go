package napi

import (
	"unsafe"
)

type Deferred unsafe.Pointer

type Promise struct {
	Deferred Deferred
	Value    Value
}
