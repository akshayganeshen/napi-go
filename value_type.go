package napi

/*
#include <node/node_api.h>
*/
import "C"

type ValueType int

const (
	ValueTypeUndefined ValueType = C.napi_undefined
	ValueTypeNull      ValueType = C.napi_null
	ValueTypeBoolean   ValueType = C.napi_boolean
	ValueTypeNumber    ValueType = C.napi_number
	ValueTypeString    ValueType = C.napi_string
	ValueTypeSymbol    ValueType = C.napi_symbol
	ValueTypeObject    ValueType = C.napi_object
	ValueTypeFunction  ValueType = C.napi_function
	ValueTypeExternal  ValueType = C.napi_external
	ValueTypeBigint    ValueType = C.napi_bigint
)
