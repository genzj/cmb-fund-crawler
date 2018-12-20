package db

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/dgraph-io/badger"
)

func openDatabase() *badger.DB {
	opts := badger.DefaultOptions
	opts.Dir = "./db-data"
	opts.ValueDir = "./db-data"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

var once sync.Once
var instance *badger.DB

func GetDatabaseInstance() *badger.DB {
	once.Do(func() {
		instance = openDatabase()
	})
	return instance
}

func SaveJSONObject(txn *badger.Txn, key string, obj interface{}) error {
	value, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return txn.Set(
		[]byte(key),
		value,
	)
}

func LoadJSONObject(txn *badger.Txn, key string, obj interface{}) error {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return err
	}

	bs, err := item.Value()
	if err != nil {
		return err
	}

	if json.Unmarshal(bs, obj) != nil {
		return err
	}
	return nil
}
