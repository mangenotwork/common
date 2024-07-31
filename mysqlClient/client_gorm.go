package mysqlClient

import (
	"fmt"
	"log"

	"github.com/mangenotwork/common/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
	orm, err = gorm.Open(mysql.Open(str), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
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

type Page struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func NewPage(page, limit int) *Page {
	return &Page{
		Page:  page,
		Limit: limit,
	}
}

// Pagination 分页
func (page *Page) Pagination() func(db *gorm.DB) *gorm.DB {
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Limit < 1 {
		page.Limit = 10
	}
	offset := (page.Page - 1) * page.Limit
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(page.Limit)
	}
}

func (page *Page) String() string {
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Limit < 1 {
		page.Limit = 10
	}
	offset := (page.Page - 1) * page.Limit
	return fmt.Sprintf("limit %d offset %d ", page.Limit, offset)
}
