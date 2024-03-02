package main

import (
	"github.com/gorilla/websocket"
)

// Pub/Sub constant
var PS = NewPubSub()

// WS Upgrader constant
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ID generator letters constant
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
