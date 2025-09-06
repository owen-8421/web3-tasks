package db

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
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
func InitDB() error {
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
	databases := viper.GetStringMap("databases")
	if len(databases) == 0 {
		log.Fatal("错误: 未在 'setting.yaml' 中找到 'databases' 配置项。")
		return errors.New("错误: 未在 'setting.yaml' 中找到 'databases' 配置项。")
	}

	// 步骤5：遍历配置文件中定义的所有数据库，并为每一个创建连接
	for name := range databases {
		// 获取当前数据库的详细配置
		dbConfig := viper.Sub("databases." + name)
		dsn := dbConfig.GetString("dsn")
		dbType := dbConfig.GetString("type")

		// 校验配置是否存在
		if dsn == "" {
			log.Fatalf("错误: 数据库 '%s' 的 'dsn' 未配置。", name)
		}
		if dbType == "" {
			log.Fatalf("错误: 数据库 '%s' 的 'type' 未配置。", name)
		}

		// 根据配置的类型选择合适的GORM驱动
		var dialector gorm.Dialector
		switch dbType {
		case "sqlite":
			dialector = sqlite.Open(dsn)
		case "mysql":
			// 如果要用MySQL，取消这里的注释并导入MySQL驱动
			dialector = mysql.Open(dsn)
		default:
			return errors.New(fmt.Sprintf("错误: 不支持的数据库类型 '%s'", dbType))
		}

		// 使用选定的驱动和配置打开数据库连接
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("错误：无法连接到数据库 '%s': %v", name, err))
		}

		log.Printf("数据库 '%s' 连接成功！", name)
		// 将创建好的连接存入全局map中
		Dbs[name] = db
		primaryDB, ok := Dbs["primary"]
		if !ok {
			log.Fatal("错误: 未找到名为 'primary' 的主数据库配置用于迁移。")
		}

		log.Println("正在 'primary' 数据库上执行自动迁移...")
		err = primaryDB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
		if err != nil {
			log.Fatalf("错误：'primary' 数据库迁移失败: %v", err)
		}
		log.Println("数据库迁移成功！")
	}
	return nil
}
