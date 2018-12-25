package db

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/genzj/cmb-fund-crawler/crawl"
	"github.com/labstack/gommon/log"
)

type FundValue struct {
	PrimaryKey       uint64
	FundOrganization string
	FundID           string
	Detail           crawl.FundDetail
}

const (
	FundValueKeyPrefixFmt      = "f:%s:%s:value:"
	FundValueKeyFmt            = "f:%s:%s:value:%010d"
	FundValueSequenceKey       = "f:%s:%s:key"
	FundValueSequenceBandwidth = 1
)

func SaveFundValueRecord(db *badger.DB, fundOrganization, fundID string, value crawl.FundDetail) error {
	seqKey := fmt.Sprintf(FundValueSequenceKey, fundOrganization, fundID)
	seq, err := db.GetSequence([]byte(seqKey), FundValueSequenceBandwidth)
	if seq != nil {
		defer func() {
			if err := seq.Release(); err != nil {
				log.Errorf("ERROR error at release sequence %s: %s\n", seqKey, err)
			}
		}()
	}
	if err != nil {
		return err
	}
	pk, err := seq.Next()
	if err != nil {
		return err
	}

	key := fmt.Sprintf(FundValueKeyFmt, fundOrganization, fundID, pk)

	v := FundValue{
		PrimaryKey:       pk,
		FundOrganization: fundOrganization,
		FundID:           fundID,
		Detail:           value,
	}

	err = db.Update(func(txn *badger.Txn) error {
		return SaveJSONObject(txn, key, v)
	})
	return err
}

func IterateFundValue(db *badger.DB, fundOrganization string, fundID string, reader func(value FundValue) error) error {
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		prefix := []byte(fmt.Sprintf(FundValueKeyPrefixFmt, fundOrganization, fundID))
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			log.Debugf("reading fund from %s", string(item.Key()))
			v := FundValue{}
			err := LoadJSONFromItem(item, &v)
			if err != nil {
				return err
			}
			if err = reader(v); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
