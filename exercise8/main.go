package main

import (
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	. "gotraining/exercise8/model"
	"gotraining/exercise8/persist"
	"log"
	"strings"
	"unicode"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("Error during run: ", err)
	}
	log.Println("Numbers updated!")
}

func run() error {
	repo, err := persist.New("pgx", "user=postgres password=mysecretpassword host=localhost port=5432")
	if err != nil {
		return fmt.Errorf("unable to create repo: %w", err)
	}
	if err := repo.Setup(); err != nil {
		return fmt.Errorf("unable to setup db: %w", err)
	}
	numbers, err := repo.ListAll()
	if err != nil {
		return fmt.Errorf("unable to load phone numbers: %w", err)
	}
	return updateTable(repo, normalizePhoneNumbers(numbers))
}

func normalizePhoneNumbers(numbers []PhoneNumber) []PhoneNumber {
	ret := make([]PhoneNumber, 0, len(numbers))
	for _, v := range numbers {
		norm := strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) {
				return r
			}
			return -1
		}, v.Number)
		ret = append(ret, PhoneNumber{ID: v.ID, Number: norm})
	}
	return ret
}

func updateTable(repo PhoneRepo, numbers []PhoneNumber) error {
	for _, v := range numbers {
		if err := processNumber(repo, v); err != nil {
			return fmt.Errorf("updateTable: %w", err)
		}
	}
	return nil
}

func processNumber(repo PhoneRepo, number PhoneNumber) error {
	_, err := repo.Find(number.Number)
	if errors.Is(err, ErrNoPhoneNumber) {
		return repo.Update(number)
	} else {
		return repo.Delete(number)
	}
}
