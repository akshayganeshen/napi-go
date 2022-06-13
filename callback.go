package napi

type Callback func(env Env, info CallbackInfo) Value
