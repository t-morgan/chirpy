package auth

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "testPassword"
	hashed_password, err := HashPassword(password)
	if err != nil {
		t.Fatalf(`HashPassword returned error: %v`, err)
	}

	err = CheckPasswordHash(password, hashed_password)
	if err != nil {
		t.Fatalf(`CheckPasswordHash returned error: %v`, err)
	}
}
