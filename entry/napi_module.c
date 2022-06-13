#include <stdlib.h>

#include "./entry.h"
#include "./exports.h"

napi_value InitializeModule(napi_env env, napi_value exports) {
  EnterCGo();

  if (napi_go_global_exports.len > 0) {
    napi_property_descriptor *props = calloc(
      napi_go_global_exports.len,
      sizeof(napi_property_descriptor)
    );

    for (size_t i = 0; i < napi_go_global_exports.len; i++) {
      props[i] = (napi_property_descriptor){
        napi_go_global_exports.exports[i].name,
        0,
        napi_go_global_exports.exports[i].callback,
        0,
        0,
        0,
        napi_default,
        0,
      };
    }
    napi_status status = napi_define_properties(
      env,
      exports,
      napi_go_global_exports.len,
      props
    );
  }

  return exports;
}

NAPI_MODULE(napiGo, InitializeModule)
