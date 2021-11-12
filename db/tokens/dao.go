package tokens

import (
	"github.com/andriystech/lgc/models"
)

func Generate() *models.Token {
	return models.NewToken()
}
