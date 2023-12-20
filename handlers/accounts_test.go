package handlers

import (
	"strings"
	"testing"
)

func TestHashFunctions(t *testing.T) {
	testPassword := "test123"
	hash, err := hashPassword(testPassword, 12)
	if err != nil {
		t.Fatal("Failed to hash password")
	}

	if !strings.HasPrefix(hash, "$2a") {
		t.Fatal("Token is bad")
	}

	if !checkPassword(testPassword, hash) {
		t.Fatal("Password should have been correct")
	}

	if checkPassword("wrong", hash) {
		t.Fatal("Password should have been incorrect")
	}
}
