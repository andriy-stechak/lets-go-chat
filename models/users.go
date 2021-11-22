package models

const NameMinLength = 3

const PasswordMinLength = 6

type User struct {
	Id       string `bson:"_id"`
	UserName string `bson:"userName"`
	Password string `bson:"password"`
}

func NewUser(id, name, password string) *User {
	return &User{
		Id:       id,
		UserName: name,
		Password: password,
	}
}
