package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func generateID(length uint) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Decode JSON request.Body into a struct, return an error if it fails
func decodeJSONBody(writer http.ResponseWriter, request *http.Request, ref interface{}) error {
	err := json.NewDecoder(request.Body).Decode(ref)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	return err
}

func handleDBError(writer http.ResponseWriter, result *gorm.DB, messages ...string) error {
	if result.Error != nil {
		errorMessage := result.Error.Error()
		code := http.StatusInternalServerError
		if len(messages) > 0 {
			errorMessage = messages[0]
			code = http.StatusUnauthorized
		}

		http.Error(writer, errorMessage, code)
	}
	return result.Error
}

func getURLParam(request *http.Request, writer http.ResponseWriter, param string) (string, bool) {
	vars := mux.Vars(request)
	value, ok := vars[param]
	if !ok {
		http.Error(writer, param+" not provided", http.StatusBadRequest)
	}
	return value, ok
}

func initDB() *gorm.DB {
	dsn := "host=db user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}
	db.AutoMigrate(&Chat{}, &Message{})
	return db
}

func handleUsing(f func(*gorm.DB, http.ResponseWriter, *http.Request), db *gorm.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		f(db, writer, request)
	}
}

func prepareChatObject(chat *Chat) {
	if chat.ID == "" {
		chat.ID = generateID(20)
	}
	chat.AdminToken = generateID(64)
	chat.Messages = []Message{}
}
