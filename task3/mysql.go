package task3

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // 加载mysql
	"github.com/jinzhu/gorm"
)

// Eloquent TODO
var Eloquent *gorm.DB

var (
	// DbType TODO
	DbType string
	// Host TODO
	Host string
	// Port TODO
	Port int
	// Name TODO
	Name string
	// Username TODO
	Username string
	// Password TODO
	Password string
)

// MysqlDB TODO
type MysqlDB struct {
	DB     *gorm.DB  `json:"-"`
	User   string    `json:"user"`
	Pwd    string    `json:"password"`
	Host   string    `json:"host"`
	Port   string    `json:"port"`
	DbName string    `json:"dbName"`
	Once   sync.Once `json:"-"`
}

// Databases TODO
var Databases = map[string]*MysqlDB{}

// Setup TODO
func Setup() {
	Databases["local"] = &MysqlDB{
		User:   "root",
		Pwd:    "xxx",
		Host:   "localhost",
		Port:   "3306",
		DbName: "web3",
	}

	conn := GetMysqlConnect()
	var err error
	var db Database
	if DbType != "mysql" {
		panic("db type unknow")
	}
	db = new(Mysql)
	Eloquent, err = db.Open(DbType, conn)
	if Eloquent != nil {
		Eloquent.DB().SetMaxOpenConns(100)
		Eloquent.DB().SetMaxIdleConns(15)
	}
	if err != nil {
		log.Fatalf("%s connect error %v", DbType, err)
	} else {
		log.Printf("%s connect success!", DbType)
	}
	if Eloquent.Error != nil {
		log.Fatalf("database error %v", Eloquent.Error)
	}
	Eloquent.LogMode(true)
}

// GetMysqlConnect TODO
func GetMysqlConnect() string {
	var conn bytes.Buffer
	conn.WriteString(Username)
	conn.WriteString(":")
	conn.WriteString(Password)
	conn.WriteString("@tcp(")
	conn.WriteString(Host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(Port))
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(Name)
	conn.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=1000ms")
	return conn.String()
}

// Database TODO
type Database interface {
	Open(dbType string, conn string) (db *gorm.DB, err error)
}

// Mysql TODO
type Mysql struct {
}

// Open TODO
func (*Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

// SqlLite TODO
type SqlLite struct {
}

// Open TODO
func (*SqlLite) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

// GetDataBase TODO
func GetDataBase(key string) (*MysqlDB, error) {
	tmp, ok := Databases[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No Key %s DataBase Exist.", key))
	}
	return tmp, nil
}

// GetDB TODO
func GetDB(connKey string) (*gorm.DB, error) {
	database, err := GetDataBase(connKey)
	if database.DB == nil {
		database.Once.Do(func() {
			err = database.initDB()
		})
	}
	return database.DB, err
}

// SourceName TODO
func (p *MysqlDB) SourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		p.User,
		p.Pwd,
		p.Host,
		p.Port,
		p.DbName,
	)
}

func (p *MysqlDB) initDB() error {
	var err error
	p.DB, err = gorm.Open("mysql", p.SourceName())
	if err != nil {
		log.Printf("Open db failed: %v", err.Error())
		return err
	}

	sqlDb := p.DB.DB()
	sqlDb.SetMaxOpenConns(1000)
	// SetMaxIdleConns 设置最大的可空闲连接数
	sqlDb.SetMaxIdleConns(300)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDb.SetConnMaxLifetime(time.Hour)
	p.DB.LogMode(true)
	return nil
}
