package main

type Chat struct {
	ID         string    `gorm:"type:varchar(20);primary_key;" json:"id"`
	AdminToken string    `json:"admin_token"`
	Title      string    `json:"title"`
	Messages   []Message `gorm:"foreignKey:ChatID" json:"messages"`
}

type Message struct {
	ID      string `gorm:"type:varchar(20);primary_key;" json:"id"`
	Content string `json:"content"`
	ChatID  string `gorm:"type:varchar(20);" json:"chat_id"`
}