package mysqlClient

import (
	"github.com/mangenotwork/common/log"
	"testing"
)

// go test -test.run Test_client_gorm -v
func Test_client_gorm(t *testing.T) {
	host := "127.0.0.1"
	port := "3306"
	database := "niu_pp"
	user := "root"
	password := "123456"
	c, err := NewORM(database, user, password, host, port, true)
	log.Print(c, err)

	type AD struct {
		Id        int64  `gorm:"primary_key;column:id" json:"id"`
		ADContent string `gorm:"column:ad_content" json:"ad_content"`
		ADUrl     string `gorm:"column:ad_url" json:"ad_url"`
	}

	a := make([]*AD, 0)
	err = c.Table("tbl_ad").Find(&a).Error
	log.Print(err)
	for k, v := range a {
		log.Print(k, v)
	}
}

// go test -test.run Test_client -v
func Test_client(t *testing.T) {
	host := "127.0.0.1"
	port := "3306"
	database := "niu_pp"
	user := "root"
	password := "123456"
	c, err := NewMysql(host, port, user, password, database)
	log.Print(c, err)

	data, err := c.Select("select * from tbl_ad;")
	log.Print(err)
	for k, v := range data {
		log.Print(k, v)
	}
}
