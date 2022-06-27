package js

import (
	"fmt"
	"unsafe"

	"github.com/akshayganeshen/napi-go"
)

type Env struct {
	Env napi.Env
}

type InvalidValueTypeError struct {
	Value any
}

var _ error = InvalidValueTypeError{}

func AsEnv(env napi.Env) Env {
	return Env{
		Env: env,
	}
}

func (e Env) Global() Value {
	v, st := napi.GetGlobal(e.Env)
	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}
	return Value{
		Env:   e,
		Value: v,
	}
}

func (e Env) Null() Value {
	v, st := napi.GetNull(e.Env)
	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}
	return Value{
		Env:   e,
		Value: v,
	}
}

func (e Env) Undefined() Value {
	v, st := napi.GetUndefined(e.Env)
	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}
	return Value{
		Env:   e,
		Value: v,
	}
}

func (e Env) ValueOf(x any) Value {
	var (
		v  napi.Value
		st napi.Status
	)

	switch xt := x.(type) {
	case Value:
		return xt
	case []Value:
		l := len(xt)
		v, st = napi.CreateArrayWithLength(e.Env, l)
		if st != napi.StatusOK {
			break
		}

		for i, xti := range xt {
			// TODO: Use Value.SetIndex helper
			st = napi.SetElement(e.Env, v, i, xti.Value)
			if st != napi.StatusOK {
				break
			}
		}
	case Func:
		return xt.Value
	case Callback:
		return e.FuncOf(xt).Value
	case *Promise:
		v, st = xt.Promise.Value, napi.StatusOK
	case napi.Value:
		v, st = xt, napi.StatusOK

	case nil:
		v, st = napi.GetNull(e.Env)
	case bool:
		v, st = napi.GetBoolean(e.Env, xt)
	case int:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case int8:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case int16:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case int64:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case uint:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case uint8:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case uint16:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case uint64:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case uintptr:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case unsafe.Pointer:
		v, st = napi.CreateDouble(e.Env, float64(uintptr(xt)))
	case float32:
		v, st = napi.CreateDouble(e.Env, float64(xt))
	case float64:
		v, st = napi.CreateDouble(e.Env, xt)
	case string:
		v, st = napi.CreateStringUtf8(e.Env, xt)
	case error:
		msg := e.ValueOf(xt.Error())
		v, st = napi.CreateError(e.Env, nil, msg.Value)
	case []any:
		l := len(xt)
		v, st = napi.CreateArrayWithLength(e.Env, l)
		if st != napi.StatusOK {
			break
		}

		for i, xti := range xt {
			// TODO: Use Value.SetIndex helper
			vti := e.ValueOf(xti)
			st = napi.SetElement(e.Env, v, i, vti.Value)
			if st != napi.StatusOK {
				break
			}
		}
	case map[string]any:
		v, st = napi.CreateObject(e.Env)
		if st != napi.StatusOK {
			break
		}

		for xtk, xtv := range xt {
			// TODO: Use Value.Set helper
			vtk, vtv := e.ValueOf(xtk), e.ValueOf(xtv)
			st = napi.SetProperty(e.Env, v, vtk.Value, vtv.Value)
			if st != napi.StatusOK {
				break
			}
		}

	default:
		panic(InvalidValueTypeError{x})
	}

	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}

	return Value{
		Env:   e,
		Value: v,
	}
}

func (e Env) FuncOf(fn Callback) Func {
	// TODO: Add CreateReference to FuncOf to keep value alive
	v, st := napi.CreateFunction(
		e.Env,
		"",
		AsCallback(fn),
	)

	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}

	return Func{
		Value: Value{
			Env:   e,
			Value: v,
		},
	}
}

func (e Env) NewPromise() *Promise {
	var result Promise
	result.reset(e)
	return &result
}

func (err InvalidValueTypeError) Error() string {
	return fmt.Sprintf("Value cannot be represented in JS: %T", err.Value)
}
