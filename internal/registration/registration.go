package registration

import (
	"errors"
	"fmt"
	"time"
)

const PasswordMinLen = 5

// RegisterRequest represents fields of registration request.
type RegisterRequest struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

// ConfirmRequest represents fields of confirmation request.
type ConfirmRequest struct {
	Email                string `json:"email"`
	Code                 int    `json:"code"` // confirmation code
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"` // should be the same as password field
}

// ConfirmationData represents information about registration confirmation. It's fetched from database.
type ConfirmationData struct {
	Code      int
	Login     string
	CreatedAt time.Time
}

// EmailConfirmation is used for generating email message.
type EmailConfirmation struct {
	Code int
}

// Validate validates confirmation request.
func (r ConfirmRequest) Validate() error {
	if r.Password == "" || r.PasswordConfirmation == "" {
		return errors.New("password should not be empty")
	}

	if r.Password != r.PasswordConfirmation {
		return errors.New("passwords should match")
	}

	if len(r.Password) < PasswordMinLen {
		return fmt.Errorf("passwords length should be more than %d symbols", PasswordMinLen)
	}

	return nil
}
