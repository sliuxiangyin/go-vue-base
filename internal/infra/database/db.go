package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBType 数据库类型枚举
type DBType string

const (
	SQLite DBType = "sqlite"
	MySQL  DBType = "mysql"
)

// DB 数据库实例
type DB struct {
	*gorm.DB
	Type DBType
}

// getDBType 根据数据库URL判断数据库类型
func getDBType(databaseURL string) DBType {
	if strings.HasPrefix(databaseURL, "file:") {
		return SQLite
	}
	// 检查是否为MySQL连接字符串
	if strings.Contains(databaseURL, "@tcp(") || strings.Contains(databaseURL, "@(") {
		return MySQL
	}
	// 默认使用SQLite
	return SQLite
}

// NewDB 创建新的数据库实例
func NewDB(dsn string) (*DB, error) {
	dbType := getDBType(dsn)
	// 配置 GORM 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)

	var db *gorm.DB
	var err error

	switch dbType {
	case SQLite:

		// 连接SQLite数据库
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	case MySQL:
		// 连接MySQL数据库
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)           // 空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
	log.Printf("Successfully connected to %s database: %s", dbType, dsn)
	return &DB{DB: db, Type: dbType}, nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据库结构
func (db *DB) AutoMigrate(dst ...interface{}) error {
	return db.DB.AutoMigrate(dst...)
}

// GetDB 获取原始的 *gorm.DB 实例
func (db *DB) GetDB() *gorm.DB {
	return db.DB
}
