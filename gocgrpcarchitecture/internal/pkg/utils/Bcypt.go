package utils

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const PasswordExpiryDays = 30

func HashPassword(password string) ([]byte, error) {
	const hashCost int = 12
	return bcrypt.GenerateFromPassword([]byte(password), hashCost)
}

func ComparePassword(hashedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}

func HashPasswordIsExpired(lastChangedPassword time.Time) bool {
	expiryDate := lastChangedPassword.AddDate(0, 0, PasswordExpiryDays) // Calculate the expiration date
	return time.Now().After(expiryDate)                                 // Check if the current time is past the expiry date
}
