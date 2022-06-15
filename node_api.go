package napi

/*
#include <node/node_api.h>
*/
import "C"

func GetNodeVersion(env Env) (NodeVersion, Status) {
	var cresult *C.napi_node_version
	status := Status(C.napi_get_node_version(
		C.napi_env(env),
		(**C.napi_node_version)(&cresult),
	))

	if status != StatusOK {
		return NodeVersion{}, status
	}

	return NodeVersion{
		Major:   uint(cresult.major),
		Minor:   uint(cresult.minor),
		Patch:   uint(cresult.patch),
		Release: C.GoString(cresult.release),
	}, status
}

func GetModuleFileName(env Env) (string, Status) {
	var cresult *C.char
	status := Status(C.node_api_get_module_file_name(
		C.napi_env(env),
		(**C.char)(&cresult),
	))

	if status != StatusOK {
		return "", status
	}

	return C.GoString(cresult), status
}
