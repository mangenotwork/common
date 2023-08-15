package localdb

import (
	"github.com/mangenotwork/common/log"
	"testing"
)

// go test -test.run Test_Base -v
func Test_Base(t *testing.T) {
	path := "./data.db"
	table := []string{"a", "b"}
	db := NewLocalDB(path, table)
	db.Init()
	_ = db.Set("a", "1", "aa")
	_ = db.Set("a", "2", "bb")
	_ = db.Set("a", "3", "cc")
	var t1 string
	_ = db.Get("a", "1", &t1)
	log.Info(t1)
	keys, _ := db.AllKey("a")
	log.Info("keys = ", keys)
	db.Pg("a")
	log.Info(db.Last("a"))
}
