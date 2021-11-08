package token

import (
	"fmt"

	"github.com/google/uuid"
)

const template = "ws://fancy-chat.io/ws&token=%s"

type Token struct {
	Url string `json:"url"`
}

func (t *Token) Generate() {
	t.Url = fmt.Sprintf(template, uuid.New().String())
}
