package napi

/*
#include <stdlib.h>
#include <node/node_api.h>

extern void DeleteInstanceData(
	napi_env env,
	void *finalize_data,
	void *finalize_hint
);

extern void DeleteCallbackData(
	napi_env env,
	void *finalize_data,
	void *finalize_hint
);

extern napi_value ExecuteCallback(
	napi_env env,
	napi_callback_info info
);

extern void ExecuteAsyncExecuteCallback(
	napi_env env,
	void *data
);

extern void ExecuteAsyncCompleteCallback(
	napi_env env,
	napi_status status,
	void *data
);
*/
import "C"

import (
	"fmt"
	"runtime"
	"runtime/cgo"
	"sync"
	"unsafe"
)

type NapiGoInstanceData struct {
	UserData      any
	CallbackData  NapiGoInstanceCallbackData
	AsyncWorkData NapiGoInstanceAsyncWorkData
}

type NapiGoInstanceCallbackData struct {
	CallbackMap NapiGoInstanceCallbackMap
	NextID      NapiGoCallbackID
	Lock        sync.RWMutex
}

type NapiGoCallbackID int

type NapiGoInstanceCallbackMap map[NapiGoCallbackID]*NapiGoCallbackMapEntry

type NapiGoCallbackMapEntry struct {
	Callback Callback
	ID       NapiGoCallbackID
}

type NapiGoAsyncWorkID int

type NapiGoInstanceAsyncWorkData struct {
	AsyncWorkMap NapiGoInstanceAsyncWorkMap
	NextID       NapiGoAsyncWorkID
	Lock         sync.RWMutex
}

type NapiGoInstanceAsyncWorkMap map[NapiGoAsyncWorkID]*NapiGoAsyncWorkMapEntry

type NapiGoAsyncWorkMapEntry struct {
	Execute  AsyncExecuteCallback
	Complete AsyncCompleteCallback
	ID       NapiGoAsyncWorkID
}

type InstanceDataProvider interface {
	GetUserData() any
	SetUserData(userData any)

	GetCallbackData() CallbackDataProvider
	GetAsyncWorkData() AsyncWorkDataProvider
}

type CallbackDataProvider interface {
	CreateCallback(env Env, name string, cb Callback) (Value, Status)
	GetCallback(id NapiGoCallbackID) *NapiGoCallbackMapEntry
	DeleteCallback(id NapiGoCallbackID)
}

type AsyncWorkDataProvider interface {
	CreateAsyncWork(
		env Env,
		asyncResource, asyncResourceName Value,
		execute AsyncExecuteCallback,
		complete AsyncCompleteCallback,
	) (AsyncWork, Status)
	GetAsyncWork(id NapiGoAsyncWorkID) *NapiGoAsyncWorkMapEntry
	DeleteAsyncWork(id NapiGoAsyncWorkID)
}

var _ InstanceDataProvider = &NapiGoInstanceData{}
var _ CallbackDataProvider = &NapiGoInstanceCallbackData{}
var _ AsyncWorkDataProvider = &NapiGoInstanceAsyncWorkData{}

const (
	maxStackTraceSize = 8192
)

func InitializeInstanceData(env Env) Status {
	return setInstanceData(env, &NapiGoInstanceData{})
}

//export DeleteInstanceData
func DeleteInstanceData(
	env C.napi_env,
	finalizeData, finalizeHint unsafe.Pointer,
) {
	instanceDataHandle := cgo.Handle(finalizeData)
	instanceDataHandle.Delete()
}

//export DeleteCallbackData
func DeleteCallbackData(
	cEnv C.napi_env,
	finalizeData, finalizeHint unsafe.Pointer,
) {
	env := Env(cEnv)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("napi.DeleteCallbackData: Recovered from panic: %s\n", err)
			reportStackTrace()

			msg := "unknown error"
			if err, ok := err.(error); ok {
				msg = err.Error()
			}
			ThrowError(env, "", msg)
		}
	}()

	instanceData, status := getInstanceData(env)
	if status != StatusOK {
		panic(StatusError(status))
	}

	id := *(*NapiGoCallbackID)(finalizeData)
	instanceData.GetCallbackData().DeleteCallback(id)
}

//export ExecuteCallback
func ExecuteCallback(
	cEnv C.napi_env,
	cInfo C.napi_callback_info,
) C.napi_value {
	env := Env(cEnv)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("napi.ExecuteCallback: Recovered from panic: %s\n", err)
			reportStackTrace()

			msg := "unknown error"
			if err, ok := err.(error); ok {
				msg = err.Error()
			}
			ThrowError(env, "", msg)
		}
	}()

	instanceData, status := getInstanceData(env)
	if status != StatusOK {
		panic(StatusError(status))
	}

	argc := C.size_t(0)
	var cData unsafe.Pointer
	status = Status(C.napi_get_cb_info(
		cEnv,
		cInfo,
		&argc,
		nil,
		nil,
		&cData,
	))

	if status != StatusOK {
		panic(StatusError(status))
	}

	id := *(*NapiGoCallbackID)(cData)
	callbackData := instanceData.GetCallbackData().GetCallback(id)

	info := CallbackInfo(cInfo)
	result := callbackData.Callback(env, info)
	return C.napi_value(result)
}

//export ExecuteAsyncExecuteCallback
func ExecuteAsyncExecuteCallback(cEnv C.napi_env, cData unsafe.Pointer) {
	env := Env(cEnv)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf(
				"napi.ExecuteAsyncExecuteCallback: Recovered from panic: %s\n",
				err,
			)
			reportStackTrace()

			msg := "unknown error"
			if err, ok := err.(error); ok {
				msg = err.Error()
			}
			ThrowError(env, "", msg)
		}
	}()

	instanceData, status := getInstanceData(env)
	if status != StatusOK {
		panic(StatusError(status))
	}

	id := *(*NapiGoAsyncWorkID)(cData)
	asyncWorkData := instanceData.GetAsyncWorkData().GetAsyncWork(id)
	asyncWorkData.Execute(env)
}

//export ExecuteAsyncCompleteCallback
func ExecuteAsyncCompleteCallback(
	cEnv C.napi_env,
	cStatus C.napi_status,
	cData unsafe.Pointer,
) {
	env := Env(cEnv)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf(
				"napi.ExecuteAsyncExecuteCallback: Recovered from panic: %s\n",
				err,
			)
			reportStackTrace()

			msg := "unknown error"
			if err, ok := err.(error); ok {
				msg = err.Error()
			}
			ThrowError(env, "", msg)
		}
	}()

	instanceData, status := getInstanceData(env)
	if status != StatusOK {
		panic(StatusError(status))
	}

	id := *(*NapiGoAsyncWorkID)(cData)
	asyncWorkData := instanceData.GetAsyncWorkData().GetAsyncWork(id)
	asyncWorkData.Complete(env, Status(cStatus))
}

func getInstanceDataHandle(env Env) (cgo.Handle, Status) {
	var result unsafe.Pointer
	status := Status(C.napi_get_instance_data(
		C.napi_env(env),
		&result,
	))

	if status != StatusOK || result == nil {
		return cgo.Handle(0), status
	}

	return cgo.Handle(result), status
}

func getInstanceData(env Env) (InstanceDataProvider, Status) {
	handle, status := getInstanceDataHandle(env)
	if status != StatusOK || handle == 0 {
		return nil, status
	}

	return handle.Value().(InstanceDataProvider), status
}

func setInstanceData(env Env, data *NapiGoInstanceData) Status {
	// check if an existing handle is already set, and clean it up if so
	// (napi won't invoke the finalizer if overwriting instance data)
	handle, status := getInstanceDataHandle(env)
	if status != StatusOK {
		return status
	}

	if handle != 0 {
		handle.Delete()
	}

	dataHandle := cgo.NewHandle(data)
	return Status(C.napi_set_instance_data(
		C.napi_env(env),
		unsafe.Pointer(dataHandle),
		C.napi_finalize(C.DeleteInstanceData),
		nil,
	))
}

func reportStackTrace() {
	stackTraceBuf := make([]byte, maxStackTraceSize)
	stackTraceSz := runtime.Stack(stackTraceBuf, false)
	fmt.Printf("%s\n", string(stackTraceBuf[:stackTraceSz]))
}

func (d *NapiGoInstanceData) GetUserData() any {
	return d.UserData
}

func (d *NapiGoInstanceData) SetUserData(userData any) {
	d.UserData = userData
}

func (d *NapiGoInstanceData) GetCallbackData() CallbackDataProvider {
	return &d.CallbackData
}

func (d *NapiGoInstanceData) GetAsyncWorkData() AsyncWorkDataProvider {
	return &d.AsyncWorkData
}

func (d *NapiGoInstanceCallbackData) CreateCallback(
	env Env,
	name string,
	cb Callback,
) (Value, Status) {
	d.Lock.Lock()
	defer d.Lock.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	callbackState := d.insert(cb)

	var result Value
	status := Status(C.napi_create_function(
		C.napi_env(env),
		cname,
		C.size_t(len([]byte(name))),
		C.napi_callback(C.ExecuteCallback),
		unsafe.Pointer(&callbackState.ID),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))

	if status == StatusOK {
		status = Status(C.napi_add_finalizer(
			C.napi_env(env),
			C.napi_value(result),
			unsafe.Pointer(&callbackState.ID),
			C.napi_finalize(C.DeleteCallbackData),
			nil,
			nil,
		))
	}

	return result, status
}

func (d *NapiGoInstanceCallbackData) GetCallback(
	id NapiGoCallbackID,
) *NapiGoCallbackMapEntry {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	return d.CallbackMap[id]
}

func (d *NapiGoInstanceCallbackData) DeleteCallback(id NapiGoCallbackID) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	delete(d.CallbackMap, id)
}

func (d *NapiGoInstanceCallbackData) insert(
	cb Callback,
) *NapiGoCallbackMapEntry {
	// callers are expected to lock

	if d.CallbackMap == nil {
		d.CallbackMap = NapiGoInstanceCallbackMap{}
	}

	for {
		id := d.NextID
		d.NextID++

		if d.CallbackMap[id] == nil {
			result := &NapiGoCallbackMapEntry{
				Callback: cb,
				ID:       id,
			}
			d.CallbackMap[id] = result
			return result
		}
	}
}

func (d *NapiGoInstanceAsyncWorkData) CreateAsyncWork(
	env Env,
	asyncResource, asyncResourceName Value,
	execute AsyncExecuteCallback,
	complete AsyncCompleteCallback,
) (AsyncWork, Status) {
	d.Lock.Lock()
	defer d.Lock.Unlock()

	asyncWorkState := d.insert(execute, complete)

	result := AsyncWork{
		ID: asyncWorkState.ID,
	}
	status := Status(C.napi_create_async_work(
		C.napi_env(env),
		C.napi_value(asyncResource),
		C.napi_value(asyncResourceName),
		C.napi_async_execute_callback(C.ExecuteAsyncExecuteCallback),
		C.napi_async_complete_callback(C.ExecuteAsyncCompleteCallback),
		unsafe.Pointer(&asyncWorkState.ID),
		(*C.napi_async_work)(unsafe.Pointer(&result.Handle)),
	))

	return result, status
}

func (d *NapiGoInstanceAsyncWorkData) GetAsyncWork(
	id NapiGoAsyncWorkID,
) *NapiGoAsyncWorkMapEntry {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	return d.AsyncWorkMap[id]
}

func (d *NapiGoInstanceAsyncWorkData) DeleteAsyncWork(id NapiGoAsyncWorkID) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	delete(d.AsyncWorkMap, id)
}

func (d *NapiGoInstanceAsyncWorkData) insert(
	execute AsyncExecuteCallback,
	complete AsyncCompleteCallback,
) *NapiGoAsyncWorkMapEntry {
	// callers are expected to lock

	if d.AsyncWorkMap == nil {
		d.AsyncWorkMap = NapiGoInstanceAsyncWorkMap{}
	}

	for {
		id := d.NextID
		d.NextID++

		if d.AsyncWorkMap[id] == nil {
			result := &NapiGoAsyncWorkMapEntry{
				Execute:  execute,
				Complete: complete,
				ID:       id,
			}
			d.AsyncWorkMap[id] = result
			return result
		}
	}
}
