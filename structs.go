package main

type Chat struct {
	ID         string    `gorm:"type:varchar(20);primary_key;" json:"id"`
	AdminToken string    `json:"adminToken"`
	Password   string    `gorm:"type:varchar(64)" json:"password"`
	Title      string    `json:"title"`
	Messages   []Message `gorm:"foreignKey:ChatID" json:"messages"`
}

type Message struct {
	ID      string `gorm:"type:varchar(20);primary_key;" json:"id"`
	Content string `json:"content"`
	ChatID  string `gorm:"type:varchar(20);" json:"chatId"`
}
