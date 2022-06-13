#ifndef __ENTRY_ENTRY_H__
#define __ENTRY_ENTRY_H__

#include <node/node_api.h>

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

// InitializeModule is a N-API module initialization function.
// InitializeModule is suitable for use as a napi_addon_register_func.
extern napi_value InitializeModule(
  napi_env    env,
  napi_value  exports
);

// EnterCGo is a no-op function implemented in Go to be called from C.
// EnterCGo ensures all Go module init functions are executed.
extern void EnterCGo();

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __ENTRY_ENTRY_H__ */
