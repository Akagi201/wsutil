# wsutil

WebSocket utils to help debug WebSocket applications

## Tools
* [client](/client) WebSocket client.
* [dump](/dump) WebSocket client with echo test.
* [proxy](/proxy) simple single WebSocket proxy.

## Run
* `./client --ws=ws://echo.websocket.org/`
* `./client --ws=ws://localhost:8327`
* `./proxy --upstream=ws://localhost:8328`
* `./dump --listen=0.0.0.0:8328 --echo`
