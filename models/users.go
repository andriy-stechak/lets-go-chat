package models

const NameMinLength = 3

const PasswordMinLength = 6

type User struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
