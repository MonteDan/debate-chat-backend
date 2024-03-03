package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func handleWebSocketConnection(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	adminToken := request.URL.Query().Get("adminToken")
	chatID := request.URL.Query().Get("chatID")

	var chat Chat
	result := db.Where("id = ? AND admin_token = ?", chatID, adminToken).First(&chat)
	if handleDBError(writer, result, "Invalid admin token or chat ID") != nil {
		return
	}

	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	channel := PS.Subscribe(chatID)
	defer PS.Unsubscribe(chatID, channel)

	for item := range channel {
		err := ws.WriteJSON(item.Data)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
}

func handleSendMessage(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	var message Message
	if decodeJSONBody(writer, request, &message) != nil {
		return
	}
	message.ID = generateID(20)

	result := db.Create(&message)
	if handleDBError(writer, result) != nil {
		return
	}

	PS.Publish(message.ChatID, message)

	fmt.Fprintf(writer, "Message sent")
}

func handleCreateChat(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	var chat Chat
	if decodeJSONBody(writer, request, &chat) != nil {
		return
	}

	prepareChatObject(&chat)

	result := db.Create(&chat)
	if handleDBError(writer, result) != nil {
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(chat)
}

func handleGetChatAdmin(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	chatID, ok := getURLParam(request, writer, "id")
	if !ok {
		return
	}

	adminToken, ok := getBearerToken(request, writer)
	if !ok {
		return
	}

	var chat Chat
	result := db.Preload("Messages").Where("id = ? AND admin_token = ?", chatID, adminToken).First(&chat)
	if handleDBError(writer, result, "Invalid admin token or chat ID") != nil {
		return
	}

	json.NewEncoder(writer).Encode(chat)
}

func handleGetChatUser(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	chatID, ok := getURLParam(request, writer, "id")
	if !ok {
		return
	}

	var chat Chat
	result := db.Select("id", "title").Where("id = ?", chatID).First(&chat)
	if handleDBError(writer, result, "Invalid chat ID") != nil {
		return
	}

	json.NewEncoder(writer).Encode(chat)
}

func handleDeleteMessage(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	messageID, ok := getURLParam(request, writer, "id")
	if !ok {
		return
	}

	adminToken, ok := getBearerToken(request, writer)
	if !ok {
		return
	}

	result := db.Where("id = ? AND chat_id IN (SELECT id FROM chats WHERE admin_token = ?)", messageID, adminToken).Delete(&Message{})
	if handleDBError(writer, result) != nil {
		return
	} else if result.RowsAffected == 0 {
		http.Error(writer, "Invalid admin token or message ID", http.StatusNotFound)
		return
	}
}
