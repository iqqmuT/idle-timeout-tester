# Idle Timeout Tester

Online tester for idle TCP connection timeouts. Made for detecting if a
firewall or something else drops idle TCP connections.

Based heavily on [Gorilla WebSocket echo example](https://github.com/gorilla/websocket/tree/master/examples/echo).

## Test Procedure

1. Create TCP connection using WebSocket
1. Idle
1. Send a message to server and see what happens

## Setup

    go get github.com/gorilla/websocket

## Run

    go run server.go

## Build

    go build server.go

You can change default port: `./server.go -addr 0.0.0.0:9999`
