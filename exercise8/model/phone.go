package model

import (
	"fmt"
)

var ErrNoPhoneNumber = fmt.Errorf("no such phone number")

type PhoneRepo interface {
	Setup() error
	ListAll() ([]PhoneNumber, error)
	Find(number string) (*PhoneNumber, error)
	Delete(phone PhoneNumber) error
	Update(phone PhoneNumber) error
}

type PhoneNumber struct {
	ID     int
	Number string
}
