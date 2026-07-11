package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := []string{"password123", "admin123"}
	for _, p := range passwords {
		hash, _ := bcrypt.GenerateFromPassword([]byte(p), 10)
		fmt.Printf("Password: %s\nHash:     %s\n\n", p, string(hash))
	}
}