package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db := initDB()
	router := mux.NewRouter()

	router.HandleFunc("/ws", handleUsing(handleWebSocketConnection, db))
	router.HandleFunc("/message", handleUsing(handleSendMessage, db)).Methods("POST")
	router.HandleFunc("/chats", handleUsing(handleCreateChat, db)).Methods("POST")
	router.HandleFunc("/chats/{id}/admin", handleUsing(handleGetChatAdmin, db)).Methods("GET")
	router.HandleFunc("/chats/{id}/admin/login", handleUsing(handleChatLogin, db)).Methods("POST")
	router.HandleFunc("/chats/{id}", handleUsing(handleGetChatUser, db)).Methods("GET")
	router.HandleFunc("/message/{id}", handleUsing(handleDeleteMessage, db)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
	})
	corsHandler := c.Handler(router)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
