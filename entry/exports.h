#ifndef __ENTRY_EXPORTS_H__
#define __ENTRY_EXPORTS_H__

#include <node/node_api.h>

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct {
  const char     *name;
  napi_callback   callback;
} napi_go_export;

typedef struct {
  napi_go_export *exports;
  size_t          len;
} napi_go_exports;

extern napi_go_exports napi_go_global_exports;

// NapiGoAppendGlobalExport appends an export callback to the global exports.
extern void NapiGoAppendGlobalExport(
  const char     *name,
  napi_callback   callback
);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __ENTRY_EXPORTS_H__ */
