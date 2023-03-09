package cache

import (
	"encoding/gob"
	"github.com/mangenotwork/common/log"
	"testing"
	"time"
)

// go test -test.run Test_Case1 -v
func Test_Case1(t *testing.T) {
	c := NewCache(1*time.Minute, 2*time.Minute)
	err := c.Set("a", 123)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(c.Get("a"))
	log.Info(c.Get("b"))
	c.Delete("a")
	log.Info(c.Get("a"))

	type Case2 struct {
		A string
	}
	caseData2 := make(map[string]Case2)
	caseData2["a"] = Case2{"1"}
	caseData2["b"] = Case2{"2"}
	for k, v := range caseData2 {
		c.Set(k, v)
	}
	log.Info(c.GetAll())
	filePath := "./test1.cache"
	log.Info(c.Save(filePath))
}

// go test -test.run Test_Case2 -v
func Test_Case2(t *testing.T) {
	c := NewCache(1*time.Minute, 2*time.Minute)
	type Case2 struct {
		A string
	}
	gob.Register(Case2{})
	filePath := "./test1.cache"
	err := c.Load(filePath)
	if err != nil {
		log.Error(err)
	}
	log.Info(c.GetAll())
	log.Info(c.Get("a"))
}
