package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db := initDB()
	router := mux.NewRouter()

	router.HandleFunc("/ws", handleUsing(handleWebSocketConnection, db))
	router.HandleFunc("/message", handleUsing(handleSendMessage, db)).Methods("POST")
	router.HandleFunc("/chat", handleUsing(handleCreateChat, db)).Methods("POST")
	router.HandleFunc("/chat/{id}/admin", handleUsing(handleGetChatAdmin, db)).Methods("GET")
	router.HandleFunc("/chat/{id}/user", handleUsing(handleGetChatUser, db)).Methods("GET")
	router.HandleFunc("/message/{id}", handleUsing(handleDeleteMessage, db)).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
