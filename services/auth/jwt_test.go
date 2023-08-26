package auth

import (
	"log"
	"testing"
)

func TestTokenGeneration(t *testing.T) {
	token, err := GenerateJWT(1, "username")
	if err != nil {
		t.Fatalf(`Failed to generate the token %v`, err)
	}
	log.Println("[TestTokenGeneration]", token)
}