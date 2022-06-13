package entry

/*
#include <stdlib.h>

#include "./exports.h"

// there is probably a cleaner way to set up the C bridge functions...

extern napi_value NapiGoEntryExportBridge0(napi_env env, napi_callback_info info);
*/
import "C"

import (
	"github.com/akshayganeshen/napi-go"
)

var napiGoCallbackTable []napi.Callback

func Export(name string, callback napi.Callback) {
	exportBridgeIndex := len(napiGoCallbackTable)
	napiGoCallbackTable = append(napiGoCallbackTable, callback)

	var exportBridgeFunc C.napi_callback

	// TODO: Only 1 export bridge is supported
	switch exportBridgeIndex {
	default:
		exportBridgeFunc = C.napi_callback(C.NapiGoEntryExportBridge0)
	}

	// TODO: No memory management for name
	C.NapiGoAppendGlobalExport(C.CString(name), exportBridgeFunc)
}

//export NapiGoEntryExportBridge0
func NapiGoEntryExportBridge0(env C.napi_env, info C.napi_callback_info) C.napi_value {
	callback := napiGoCallbackTable[0]

	result := callback(
		napi.Env(env),
		napi.CallbackInfo(info),
	)

	return C.napi_value(result)
}
