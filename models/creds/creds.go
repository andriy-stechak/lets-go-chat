package creds

import "github.com/andriystech/lgc/api/errors"

type User struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func (creds *User) Validate() *errors.AppError {
	if len(creds.UserName) == 0 {
		return errors.BadRequest("Field 'userName' was not provided inside body")
	}
	if len(creds.Password) == 0 {
		return errors.BadRequest("Field 'password' was not provided inside body")
	}
	return nil
}
