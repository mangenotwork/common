# mysql

有两种，gorm和database/sql

#### 全局 MysqlGorm 客户端
> var MysqlGorm map[string]*gorm.DB

---

#### 配置文件读取配置
> func InitMysqlGorm()

---

#### NewORM  连接 orm
> func NewORM(database, user, password, host, port string, disablePrepared bool) (*gorm.DB, error)

---

#### GetGorm 获取gorm对象
> func GetGorm(name string) *gorm.DB

---

#### SetGorm 设置gorm对象
> func SetGorm(c *gorm.DB, name string)

---

#### MysqlDB 全局对象

---

#### 给mysql对象进行连接
> func NewMysqlDB(host, port, user, password, database string) (err error)

---

#### 创建一个mysql对象
> func NewMysql(host, port, user, password, database string) (*Mysql, error)

---

#### 获取mysql 连接
> func GetMysqlDBConn() (*Mysql, error)

---

#### CloseLog 关闭日志
> func (m *Mysql) CloseLog()

---

#### SetMaxOpenConn 最大连接数
> func (m *Mysql) SetMaxOpenConn(number int) 

---

#### SetMaxIdleConn 最大idle 数
> func (m *Mysql) SetMaxIdleConn(number int)

---

#### Conn 连接mysql
> func (m *Mysql) Conn() (err error)

---

#### IsHaveTable 表是否存在
> func (m *Mysql) IsHaveTable(table string) bool 

---

#### TableInfo 表信息
```shell
type TableInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
}
```

---

#### Describe 获取表结构
> func (m *Mysql) Describe(table string) (*TableDescribe, error)

---

#### Select 查询语句 返回 map
> func (m *Mysql) Select(sql string) ([]map[string]string, error)

---

#### NewTable 创建表
> func (m *Mysql) NewTable(table string, fields map[string]string) error

参数
```shell
fields  字段:类型； name:varchar(10);
```

---

#### NewTableGd 创建新的固定map顺序为字段的表
> func (m *Mysql) NewTableGd(table string, fields *utils.GDMap) error

---

#### Insert 新增数据
> func (m *Mysql) Insert(table string, fieldData map[string]interface{}) error

---

#### InsertAt 新增数据 如果没有表则先创建表
> func (m *Mysql) InsertAt(table string, fieldData map[string]interface{}) error 

---

#### InsertAtGd  固定顺序map写入
> func (m *Mysql) InsertAtGd(table string, fieldData *utils.GDMap) error

---

#### InsertAtJson json字符串存入数据库
> func (m *Mysql) InsertAtJson(table, jsonStr string) error

---

#### Update 更新sql
> func (m *Mysql) Update(sql string) error

---

#### Exec 执行sql
> func (m *Mysql) Exec(sql string) error

---

#### Query 执行selete sql
> func (m *Mysql) Query(sql string) ([]map[string]string, error)

---

#### Delete 执行delete sql
> func (m *Mysql) Delete(sql string) error

---

#### ToVarChar  写入mysql 的字符类型
> func (m *Mysql) ToVarChar(data interface{}) string

---

#### DeleteTable 删除表
> func (m *Mysql) DeleteTable(tableName string) error 

---

#### HasTable 判断表是否存
> func (m *Mysql) HasTable(tableName string) bool

---

#### GetFieldList 获取表字段
> func (m *Mysql) GetFieldList(table string) (fieldList []string)

---

#### ToXls 数据库查询输出到excel
> func (m *Mysql) ToXls(sql, outPath string)

---

#### StringValueMysql 用于mysql字符拼接使用
> func StringValueMysql(i interface{}) string 

---

