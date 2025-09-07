package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"blog/internal/model"
)

// Dbs 变量首字母大写，表示它可以被其他包访问。
var Dbs map[string]*gorm.DB

// InitDB --- 数据库初始化, 在main函数中只执行一次
func InitDB() {

	// 步骤2：初始化全局的数据库连接map
	Dbs = make(map[string]*gorm.DB)

	// 步骤3：配置GORM的日志记录器，以便在控制台看到执行的SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 启用彩色打印
		},
	)

	// 步骤4：从Viper获取 "databases" 配置块
	databases := viper.GetStringMap("settings.databases")
	log.Println("InitDB")
	if len(databases) == 0 {
		log.Fatal("错误: 未在 'setting.yaml' 中找到 'databases' 配置项。")
	}

	log.Println(databases)

	// 步骤5：遍历配置文件中定义的所有数据库，并为每一个创建连接
	for name, value := range databases {
		// 获取当前数据库的详细配置
		log.Println("name:" + name)
		dbConfig, ok := value.(map[string]interface{})

		if !ok {
			fmt.Printf("Error: Value for key '%s' is not a valid map.\n", name)
			continue
		}

		log.Println(dbConfig)

		dsn := dbConfig["dsn"].(string)
		dbType := dbConfig["type"].(string)

		// 校验配置是否存在
		if dsn == "" {
			log.Fatalf("错误: 数据库 '%s' 的 'dsn' 未配置。", name)
		}
		if dbType == "" {
			log.Fatalf("错误: 数据库 '%s' 的 'type' 未配置。", name)
		}

		log.Printf("dsn: %s, dbtype: %s", dsn, dbType)

		if dbType != "mysql" {
			log.Printf("dbtype: %s invalid", dbType)
			continue
		}

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

		if err != nil {
			log.Fatalf(fmt.Sprintf("错误：无法连接到数据库 '%s': %v", name, err))
		}

		log.Printf("数据库 '%s' 连接成功！", name)

		sqlDb, err := db.DB()
		sqlDb.SetMaxOpenConns(1000)
		// SetMaxIdleConns 设置最大的可空闲连接数
		sqlDb.SetMaxIdleConns(300)
		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDb.SetConnMaxLifetime(time.Hour)
		Dbs[name] = db
	}

	primaryDB, ok := Dbs["primary"]
	if !ok {
		log.Fatal("错误: 未找到名为 'primary' 的主数据库配置用于迁移。")
	}

	log.Println("正在 'primary' 数据库上执行自动迁移...")
	err := primaryDB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		log.Fatalf("错误：'primary' 数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移成功！")
}

// CloseDBs 关闭所有已初始化的数据库连接
func CloseDBs() {
	log.Println("正在关闭所有数据库连接...")
	for name, dbInstance := range Dbs {
		sqlDB, err := dbInstance.DB()
		if err != nil {
			log.Printf("错误：无法获取 '%s' 数据库的底层连接: %v", name, err)
			continue
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("错误：关闭 '%s' 数据库连接失败: %v", name, err)
		} else {
			log.Printf("数据库 '%s' 连接已成功关闭。", name)
		}
	}
	log.Println("所有数据库连接均已关闭。")
}
