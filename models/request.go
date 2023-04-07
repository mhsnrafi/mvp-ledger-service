package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"regexp"
)

type AuthRequest struct {
	Email string `json:"email"`
}

func (a AuthRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}

type RefreshRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (a RefreshRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Token,
			validation.Required,
			validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
		),
	)
}

type AddFundsRequest struct {
	Amount float64 `json:"amount"`
}
