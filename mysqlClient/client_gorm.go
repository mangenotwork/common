package mysqlClient

import (
	"fmt"
	"log"
	"time"

	"github.com/mangenotwork/common/conf"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlGorm map[string]*gorm.DB

// InitMysqlGorm 配置文件读取配置
func InitMysqlGorm() {
	MysqlGorm = make(map[string]*gorm.DB, len(conf.Conf.Default.Mysql))
	for _, v := range conf.Conf.Default.Mysql {
		log.Println(v)
		m, err := NewORM(v.DBName, v.User, v.Password, v.Host, v.Port, false)
		if err != nil {
			panic(err)
		}
		MysqlGorm[v.DBName] = m
	}
	log.Print("连接Redis : ", MysqlGorm)
}

// NewORM  连接 orm
func NewORM(database, user, password, host, port string, disablePrepared bool) (*gorm.DB, error) {
	var (
		orm *gorm.DB
		err error
	)
	if database == "" || user == "" || password == "" || host == "" {
		panic("数据库配置信息获取失败")
	}
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database) + "?charset=utf8mb4&parseTime=true&loc=Local"
	if disablePrepared {
		str = str + "&interpolateParams=true"
	}
	for orm, err = gorm.Open("mysql", str); err != nil; {
		log.Println(fmt.Sprintf("[DB]-[%v] 连接异常:%v，正在重试: %v", database, err, str))
		time.Sleep(5 * time.Second)
		orm, err = gorm.Open("mysql", str)
	}
	orm.LogMode(true)
	orm.CommonDB()
	return orm, err
}

// GetGorm 获取gorm对象
func GetGorm(name string) *gorm.DB {
	m, ok := MysqlGorm[name]
	if !ok {
		panic("[DB] 未init")
	}
	return m
}

func SetGorm(c *gorm.DB, name string) {
	MysqlGorm[name] = c
}

func Expr(expression string, args ...interface{}) *gorm.SqlExpr {
	return gorm.Expr(expression, args...)
}
