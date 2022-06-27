# napi-go

A Go library for building Node.js Native Addons using Node-API.

## Usage

Use `go get` to install the library:

```sh
go get -u github.com/akshayganeshen/napi-go
```

Then use the library to define handlers:

```go
package handlers

import "github.com/akshayganeshen/napi-go"

func MyHandler(env napi.Env, info napi.CallbackInfo) napi.Value {
  return nil
}
```

Next, create a `main.go` that registers all module exports:

```go
package main

import "github.com/akshayganeshen/napi-go/entry"

func init() {
  entry.Export("myHandler", MyHandler)
}

func main() {}
```

Finally, build the Node.js addon using `go build`:

```sh
go build -buildmode=c-shared -o "example.node" .
```

The output `.node` file can now be imported via `require`:

```js
const example = require("./example.node");

example.myHandler();
```

### JS Helpers

In addition to the Node-API exposed via package `napi`, the `napi-go/js`
package provides functions similar to the `syscall/js` standard library.

```go
package main

import (
  "github.com/akshayganeshen/napi-go/entry"
  "github.com/akshayganeshen/napi-go/js"
)

func init() {
  entry.Export("myCallback", js.AsCallback(MyCallback))
}

func MyCallback(env js.Env, this js.Value, args []js.Value) any {
  return map[string]any{
    "message": "hello world",
    "args":    args,
  }
}

func main() {}
```

## Examples

Check out the example addons in [`docs/examples`](docs/examples).
