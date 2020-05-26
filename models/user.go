package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User ...
type User struct {
	ID        string `json:"_id",bson:"_id"`
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

// Validate for user struct
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.BirthDate, validation.Required, validation.Date("Monday, January 02, 2006 3:04 PM")),
	)

}
