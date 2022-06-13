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
