package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Pub/Sub "router" instance
var PS = NewPubSub()

// WS Upgrader instance
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// ID generator letters constant
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
