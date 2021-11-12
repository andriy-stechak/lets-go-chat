package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Token struct {
	Url string `json:"url"`
}

const template = "ws://fancy-chat.io/ws&token=%s"

func NewToken() *Token {
	return &Token{
		Url: fmt.Sprintf(template, uuid.New().String()),
	}
}
