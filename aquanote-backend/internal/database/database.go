package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitDB 初始化 GORM 連線池
func InitDB() {
	once.Do(func() {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("取得執行檔路徑失敗: %v", err)
		}
		envPath := filepath.Join(filepath.Dir(exePath), ".env")
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("無法載入 .env 檔: %v", err)
		}

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_DATABASE")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			user, password, host, port, database,
		)

		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err != nil {
			log.Fatalf("connect db failed: %v", err)
		}

		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatalf("get sql.DB failed: %v", err)
		}

		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(3 * time.Minute)

		db = conn
		log.Println("database connected successfully")
	})
}

// GetDB 提供其他套件（repository）取用連線池
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("database not initialized, call InitDB first")
	}
	return db
}

// Close 釋放連線池
func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
