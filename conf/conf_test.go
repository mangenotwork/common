package conf

import (
	"github.com/mangenotwork/common/log"
	"testing"
)

// go test -test.run Test_yaml -v
func Test_yaml(t *testing.T) {
	err := NewConf("./app.yaml")
	if err != nil {
		log.Error(err)
	}
	log.Print(Config)
}

// go test -test.run Test_ini -v
func Test_ini(t *testing.T) {
	err := NewConf("./app.ini")
	if err != nil {
		log.Error(err)
	}
	log.Print(Config)
}
