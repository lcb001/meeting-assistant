package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("/Users/leyan/.todo-list-mcp/todos.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	//// 自动建表
	//err = DB.AutoMigrate(&models.Todo{})
	//if err != nil t
	//	log.Fatalf("自动建表失败: %v", err)
	//}
}
