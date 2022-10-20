package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//fmt.Println("eteff")
	//fmt.Println(db)
	////// 迁移 schema
	db.AutoMigrate(&models.GroupBasic{})
	// Create
	//user := &models.UserBasic{}
	//user.Name = "hyl"
	//db.Create(user)
	////// Read
	//fmt.Println(db.First(user, 1)) // 按主键值进行查找
	//
	//// Update - 将 product 的 price 更新为 200
	//db.Model(&user).Update("PassWord", "1234")

}
