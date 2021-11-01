package main

import (
	"fmt"

	"github.com/andriystech/lgc/pkg/hasher"
)

func main() {
	hash := hasher.HashPassword("hello")
	fmt.Println(hash)
}
