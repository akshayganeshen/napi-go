package js

import (
	"errors"

	"github.com/akshayganeshen/napi-go"
)

type Promise struct {
	Promise            napi.Promise
	ThreadsafeFunction napi.ThreadsafeFunction
	Result             any
	ResultType         PromiseResultType
}

type PromiseResultType string

type PromiseProvider interface {
	Resolve(resolution any)
	Reject(rejection any)
}

var _ PromiseProvider = &Promise{}

const (
	PromiseResultTypeResolved PromiseResultType = "resolved"
	PromiseResultTypeRejected PromiseResultType = "rejected"
)

var ErrPromiseSettled = errors.New(
	"Promise: Cannot resolve/reject a settled promise",
)

func (p *Promise) Resolve(resolution any) {
	p.ensurePending()

	p.Result = resolution
	p.ResultType = PromiseResultTypeResolved

	// function has already been acquired during reset
	defer p.release()
	p.settle()
}

func (p *Promise) Reject(rejection any) {
	p.ensurePending()

	p.Result = rejection
	p.ResultType = PromiseResultTypeRejected

	// function has already been acquired during reset
	defer p.release()
	p.settle()
}

func (p *Promise) reset(e Env) {
	np, st := napi.CreatePromise(e.Env)
	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}

	asyncResourceName := e.ValueOf("napi-go/js-promise")
	fn := e.FuncOf(func(env Env, this Value, args []Value) any {
		value := env.ValueOf(p.Result)

		st := napi.StatusOK
		switch p.ResultType {
		case PromiseResultTypeResolved:
			st = napi.ResolveDeferred(env.Env, p.Promise.Deferred, value.Value)
		case PromiseResultTypeRejected:
			st = napi.RejectDeferred(env.Env, p.Promise.Deferred, value.Value)
		}

		if st != napi.StatusOK {
			panic(napi.StatusError(st))
		}

		return nil
	})

	tsFn, st := napi.CreateThreadsafeFunction(
		e.Env,
		fn.Value.Value,
		nil, asyncResourceName.Value,
		0,
		1, // initialize with 1 acquisition
	)
	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}

	*p = Promise{
		Promise:            np,
		ThreadsafeFunction: tsFn,
	}
}

func (p *Promise) ensurePending() {
	if p.ResultType != "" {
		panic(ErrPromiseSettled)
	}
}

func (p *Promise) settle() {
	st := napi.CallThreadsafeFunction(p.ThreadsafeFunction)
	if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}
}

func (p *Promise) release() {
	st := napi.ReleaseThreadsafeFunction(p.ThreadsafeFunction)
	if st == napi.StatusClosing {
		p.ThreadsafeFunction = nil
	} else if st != napi.StatusOK {
		panic(napi.StatusError(st))
	}
}
