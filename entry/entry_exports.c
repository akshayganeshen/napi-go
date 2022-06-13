#include <stdlib.h>

#include "./exports.h"

napi_go_exports napi_go_global_exports = {
  .exports  = NULL,
  .len      = 0,
};

void NapiGoAppendGlobalExport(const char *name, napi_callback callback) {
  const int new_len = napi_go_global_exports.len + 1;
  napi_go_export *new_exports = realloc(
    napi_go_global_exports.exports,
    new_len * sizeof(napi_callback)
  );

  new_exports[new_len-1] = (napi_go_export){
    .name     = name,
    .callback = callback,
  };

  napi_go_global_exports = (napi_go_exports){
    .exports  = new_exports,
    .len      = new_len,
  };
}
