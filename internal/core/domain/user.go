package domain

import (
	"fmt"
	"regexp"

	core_error "github.com/Deadcrush-h/ToDo/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUnitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitialzedID,
		UninitialzedVersion,
		fullName,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLength,
			core_error.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_error.ErrInvalidArgument,
			)
		}
	}
	re := regexp.MustCompile(`^\+[0-9]+$`)

	if !re.MatchString(*u.PhoneNumber) {
		return fmt.Errorf(
			"invalid `PhoneNumber` format: %w",
			core_error.ErrInvalidArgument,
		)
	}
	
	return nil
}
