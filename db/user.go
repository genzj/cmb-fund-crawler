package db

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
)

type userInfo struct {
	Name     string
	Password string
}

type userFunds struct {
	Funds []struct {
		FundOrganization string
		ID               string
		BidRate          string
		Amount           string
		FirstBidDate     time.Time
	}
}

const (
	userInfoKey  = "u:%s:info"
	userFundsKey = "u:%s:fund"
)

func CreateUser(db *badger.DB, username string, password string) error {
	key := fmt.Sprintf(userInfoKey, username)
	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(key)); err == nil {
			return fmt.Errorf("user %s exists", username)
		} else if err != badger.ErrKeyNotFound {
			return err
		}

		if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
			return err
		} else {
			return SaveJSONObject(txn, key, &userInfo{Name: username, Password: string(hashedPassword)})
		}
	})
	return err
}

func GetUser(db *badger.DB, username string) (*userInfo, error) {
	ans := &userInfo{}
	key := fmt.Sprintf(userInfoKey, username)
	err := db.View(func(txn *badger.Txn) error {
		return LoadJSONObject(txn, key, ans)
	})
	return ans, err
}
