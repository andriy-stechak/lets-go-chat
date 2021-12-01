package models

type Token struct {
	Payload string
}

func NewToken(token string) *Token {

	return &Token{
		Payload: token,
	}
}
