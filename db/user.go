package db

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

type userInfo struct {
	Name string
}

type userFunds struct {
	Funds []struct {
		ID      string
		BidRate string
		Amount  string
	}
}

const (
	userInfoKey  = "u:%s:info"
	userFundsKey = "u:%s:fund"
)

func CreateUser(db *badger.DB, username string) error {
	key := fmt.Sprintf(userInfoKey, username)
	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(key)); err == nil {
			return fmt.Errorf("user %s exists", username)
		} else if err != badger.ErrKeyNotFound {
			return err
		}

		return SaveJSONObject(txn, key, &userInfo{Name: username})
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
