package main

import (
	"github.com/akshayganeshen/napi-go"
	"github.com/akshayganeshen/napi-go/entry"
)

func init() {
	entry.Export("describeArgs", DescribeArgsHandler)
}

func DescribeArgsHandler(env napi.Env, info napi.CallbackInfo) napi.Value {
	extractedInfo, _ := napi.GetCbInfo(env, info)
	result, _ := napi.CreateArrayWithLength(env, len(extractedInfo.Args))
	for i, arg := range extractedInfo.Args {
		vt, _ := napi.Typeof(env, arg)
		dv, _ := napi.CreateStringUtf8(env, DescribeValueType(vt))
		napi.SetElement(env, result, i, dv)
	}

	return result
}

func DescribeValueType(vt napi.ValueType) string {
	switch vt {
	case napi.ValueTypeUndefined:
		return "undefined"
	case napi.ValueTypeNull:
		return "null"
	case napi.ValueTypeBoolean:
		return "boolean"
	case napi.ValueTypeNumber:
		return "number"
	case napi.ValueTypeString:
		return "string"
	case napi.ValueTypeSymbol:
		return "symbol"
	case napi.ValueTypeObject:
		return "object"
	case napi.ValueTypeFunction:
		return "function"
	case napi.ValueTypeExternal:
		return "external"
	case napi.ValueTypeBigint:
		return "bigint"

	default:
		return "other"
	}
}

func main() {}
