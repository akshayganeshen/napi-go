package napi

/*
#include <stdlib.h>
#include <node/node_api.h>
*/
import "C"

import (
	"unsafe"
)

func GetUndefined(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_get_undefined(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetNull(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_get_null(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateObject(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_create_object(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateArray(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_create_array(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateArrayWithLength(env Env, length int) (Value, Status) {
	var result Value
	status := Status(C.napi_create_array_with_length(
		C.napi_env(env),
		C.size_t(length),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateStringUtf8(env Env, str string) (Value, Status) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	var result Value
	status := Status(C.napi_create_string_utf8(
		C.napi_env(env),
		cstr,
		C.size_t(len([]byte(str))), // must pass number of bytes
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateFunction(env Env, name string, cb Callback) (Value, Status) {
	provider, status := getInstanceData(env)
	if status != StatusOK || provider == nil {
		return nil, status
	}

	return provider.GetCallbackData().CreateCallback(env, name, cb)
}

func CreateError(env Env, code, msg Value) (Value, Status) {
	var result Value
	status := Status(C.napi_create_error(
		C.napi_env(env),
		C.napi_value(code),
		C.napi_value(msg),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func Typeof(env Env, value Value) (ValueType, Status) {
	var result ValueType
	status := Status(C.napi_typeof(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_valuetype)(unsafe.Pointer(&result)),
	))
	return result, status
}

func SetProperty(env Env, object, key, value Value) Status {
	return Status(C.napi_set_property(
		C.napi_env(env),
		C.napi_value(object),
		C.napi_value(key),
		C.napi_value(value),
	))
}

func SetElement(env Env, object Value, index int, value Value) Status {
	return Status(C.napi_set_element(
		C.napi_env(env),
		C.napi_value(object),
		C.uint32_t(index),
		C.napi_value(value),
	))
}

func StrictEquals(env Env, lhs, rhs Value) (bool, Status) {
	var result bool
	status := Status(C.napi_strict_equals(
		C.napi_env(env),
		C.napi_value(lhs),
		C.napi_value(rhs),
		(*C.bool)(&result),
	))
	return result, status
}

type GetCbInfoResult struct {
	Args []Value
	This Value
}

func GetCbInfo(env Env, info CallbackInfo) (GetCbInfoResult, Status) {
	// call napi_get_cb_info twice
	// first is to get total number of arguments
	// second is to populate the actual arguments
	argc := C.size_t(0)
	status := Status(C.napi_get_cb_info(
		C.napi_env(env),
		C.napi_callback_info(info),
		&argc,
		nil,
		nil,
		nil,
	))

	if status != StatusOK {
		return GetCbInfoResult{}, status
	}

	argv := make([]Value, int(argc))
	var thisArg Value

	status = Status(C.napi_get_cb_info(
		C.napi_env(env),
		C.napi_callback_info(info),
		&argc,
		(*C.napi_value)(unsafe.Pointer(&argv[0])), // must pass element pointer
		(*C.napi_value)(unsafe.Pointer(&thisArg)),
		nil,
	))

	return GetCbInfoResult{
		Args: argv,
		This: thisArg,
	}, status
}

func Throw(env Env, err Value) Status {
	return Status(C.napi_throw(
		C.napi_env(env),
		C.napi_value(err),
	))
}

func ThrowError(env Env, code, msg string) Status {
	codeCStr, msgCCstr := C.CString(code), C.CString(msg)
	defer C.free(unsafe.Pointer(codeCStr))
	defer C.free(unsafe.Pointer(msgCCstr))

	return Status(C.napi_throw_error(
		C.napi_env(env),
		codeCStr,
		msgCCstr,
	))
}
