package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func InitDB() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("无法获取用户主目录: " + err.Error())
	}

	defaultDBFolder := filepath.Join(homeDir, ".todo-list-mcp")
	defaultDBFile := "todos.sqlite"
	defaultDBPath := filepath.Join(defaultDBFolder, defaultDBFile)

	DB, err = gorm.Open(sqlite.Open(defaultDBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	//// 自动建表
	//err = DB.AutoMigrate(&models.Todo{})
	//if err != nil t
	//	log.Fatalf("自动建表失败: %v", err)
	//}
}
