GLSP
====

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/bennypowers/glsp.svg)](https://pkg.go.dev/github.com/bennypowers/glsp)
[![Go Report Card](https://goreportcard.com/badge/github.com/bennypowers/glsp)](https://goreportcard.com/report/github.com/bennypowers/glsp)

[Language Server Protocol](https://microsoft.github.io/language-server-protocol/) SDK for Go.

Fork of [tliron/glsp](https://github.com/tliron/glsp) -- thanks to Tal Liron for the original work. Adds fixes and LSP 3.17 support.

It enables you to more easily implement language servers by writing them in Go. GLSP contains:

1) all the message structures for easy serialization,
2) a handler for all client methods, and
3) a ready-to-run JSON-RPC 2.0 server supporting stdio, TCP, WebSockets, and Node.js IPC.

All you need to do, then, is provide the features for the language you want to support.

Protocol Versions
-----------------

`protocol_3_17` extends `protocol_3_16` via type aliases for unchanged types and native
definitions for 3.17 additions. Import only the version you need:

```go
import protocol "github.com/bennypowers/glsp/protocol_3_17"
```

Types shared between versions are assignment-compatible, so handler functions inherited
from 3.16 accept 3.17 values without casting.

Projects using GLSP:

* [asimonim](https://bennypowers.dev/asimonim) - Design token language server
* [cem](https://bennypowers.dev/cem) - Custom elements manifest language server


Minimal Example
---------------

```go
package main

import (
	"github.com/bennypowers/glsp"
	protocol "github.com/bennypowers/glsp/protocol_3_17"
	"github.com/bennypowers/glsp/server"
	"github.com/tliron/commonlog"

	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "my language"

var (
	version string = "0.0.1"
	handler protocol.Handler
)

func main() {
	// This increases logging verbosity (optional)
	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
```
