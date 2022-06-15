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
