# `haproxysocket` - a go wrapper around the haproxy unix socket
This covers about 80 of the commandos, and my goal is to support 100%  

## How to use
In your haproxy config under **global** add `stats timeout 2m` and add `stats socket ipv4@127.0.0.1:9999 level admin` or `stats socket /var/run/haproxy.sock mode 666 level admin`  

In your go code:
```go
package main

import (
	"fmt"

	"github.com/mjarkk/haproxysocket"
)

func main() {
	// Create a instace of haproxy
	// Make sure to change /var/sock/haproxy.sock to where your haproxy sock file is
	h := haproxysocket.New("unix", "/var/sock/haproxy.sock")

	// Get the sessions
	sessions, err := h.ShowSess()
	if err != nil {
		panic(err)
	}
	fmt.Println("Sessions:")
	for _, session := range sessions {
		fmt.Println("ID:", session.ID)
	}

	// Set a server to maintenance mode
	s := h.Server("test-backend", "serv1")
	err = s.State("maint")
	if err != nil {
		panic(err)
	}
}
```

## Avaliable functions
Most functions have the same naming sceme as the socket commands, for example `show errors` will become `ShowErrors`   
For docs about the functions see: [mangement.txt > 9.3. Unix Socket commands](http://www.haproxy.org/download/2.0/doc/management.txt)  
- `ShowErrors` 
- `ClearCounters`
- `ShowInfo `
- `ShowStat `
- `ShowSchemaJSON `
- `DisableAgent`
- `DisableHealth`
- `DisableServer`
- `EnableAgent`
- `EnableHealth`
- `EnableServer`
- `SetMaxconnServer`
- `Server`
- `(ServerT)Addr`
- `(ServerT)Agent`
- `(ServerT)AgentAddr`
- `(ServerT)AgentSend`
- `(ServerT)Health`
- `(ServerT)CheckPort`
- `(ServerT)State`
- `(ServerT)Weight`
- `(ServerT)FQDN`
- `GetWeight `
- `SetWeight`
- `ShowSess `
- `ShutdownSession`
- `ShutdownSessionsServer`
- `ClearTable` :x: Not inplemented yet
- `SetTable` :x: Not inplemented yet
- `ShowTable` :x: Not inplemented yet
- `DisableFrontend`
- `EnableFrontend`
- `SetMaxconnFrontend`
- `ShowServersState `
- `ShowBackend `
- `ShutdownFrontend`
- `SetDynamicCookieKeyBackend`
- `DynamicCookieBackend`
- `ShowStatResolvers `
- `SetMaxconnGlobal`
- `SetRateLimit`
- `ShowEnv `
- `ShowCliSockets `
- `AddACL` :x: Not inplemented yet
- `ClearACL` :x: Not inplemented yet
- `DelACL` :x: Not inplemented yet
- `GetACL` :x: Not inplemented yet
- `ShowACL` :x: Not inplemented yet
- `AddMap` :x: Not inplemented yet
- `ClearMap` :x: Not inplemented yet
- `DelMap` :x: Not inplemented yet
- `GetMap` :x: Not inplemented yet
- `SetMap` :x: Not inplemented yet
- `ShowMap`
- `ShowPools `
