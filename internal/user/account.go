package user

import (
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID           uuid.UUID
	Name         string
	HashPassword string
	Quota        int64
}

func CreateAccount(username, password string) *Account {
	acc := Account{
		// ID:   [16]byte{},
		Name: username,
	}
	acc.createHash(password)
	acc.MustGenerateID()
	return &acc
}

func (acc *Account) createHash(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	acc.HashPassword = string(hash)
	return nil
}

func (acc *Account) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(acc.HashPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (acc *Account) MustGenerateID() {
	acc.ID = uuid.New()
}




