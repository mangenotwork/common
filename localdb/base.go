package localdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

var DBPath = "./data.db"

type LocalDB struct {
	Path   string
	Tables []string
	Conn   *bolt.DB
}

func NewLocalDB(path string, tables []string) *LocalDB {
	return &LocalDB{
		Path:   path,
		Tables: tables,
	}
}

func (ldb *LocalDB) Init() {
	db, err := bolt.Open(ldb.Path, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = db.Close()
	}()
	for _, table := range ldb.Tables {
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(table))
			if b == nil {
				_, err = tx.CreateBucket([]byte(table))
				if err != nil {
					log.Panic(err)
				}
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}
}

func (ldb *LocalDB) open() {
	ldb.Conn, _ = bolt.Open(ldb.Path, 0600, nil)
}

func (ldb *LocalDB) HasTable(table string) bool {
	var rse = false
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	_ = ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b != nil {
			rse = true
		}
		return nil
	})
	return rse
}

func (ldb *LocalDB) Set(table, key string, data interface{}) error {
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	value, err := utils.AnyToJsonB(data)
	if err != nil {
		return err
	}
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b != nil {
			err = b.Put([]byte(key), value)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (ldb *LocalDB) Get(table, key string, data interface{}) error {
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b != nil {
			bt := b.Get([]byte(key))
			err := json.Unmarshal(bt, data)
			return err
		}
		return nil
	})
}

func (ldb *LocalDB) AllKey(table string) ([]string, error) {
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	keys := make([]string, 0)
	err := ldb.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			keys = append(keys, string(k))
		}
		return nil
	})
	return keys, err
}

func (ldb *LocalDB) GetAll(table string, fn func(k, v []byte)) error {
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()

	return ldb.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		forErr := b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			fn(k, v)
			return nil
		})
		return forErr
	})
}

func (ldb *LocalDB) Pg(table string) error {
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table)).Cursor()
		min := []byte("2")
		max := []byte("3")
		for k, v := b.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = b.Next() {
			fmt.Printf("%s: %s\n", k, v)
		}
		return nil
	})
}

func (ldb *LocalDB) Last(table string) error {
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table)).Cursor()
		k, _ := b.Last()
		log.Info("k = ", string(k))
		return nil
	})
}

func (ldb *LocalDB) Delete(table, key string) error {
	ldb.open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return fmt.Errorf("未获取到表")
		}
		if err := b.Delete([]byte(key)); err != nil {
			return err
		}
		return nil
	})
}
