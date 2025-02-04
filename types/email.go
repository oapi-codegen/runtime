package types

import (
	"encoding/json"
	"errors"
	"net/mail"
)

// ErrValidationEmail is the sentinel error returned when an email fails validation
var ErrValidationEmail = errors.New("email: failed to pass regex validation")

// Email represents an email address.
// It is a string type that must pass regex validation before being marshalled
// to JSON or unmarshalled from JSON.
type Email string

func (e Email) MarshalJSON() ([]byte, error) {
	if _, err := mail.ParseAddress(string(e)); err != nil {
		return nil, ErrValidationEmail
	}

	return json.Marshal(string(e))
}

func (e *Email) UnmarshalJSON(data []byte) error {
	if e == nil {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*e = Email(s)
	if _, err := mail.ParseAddress(s); err != nil {
		return ErrValidationEmail
	}

	return nil
}
