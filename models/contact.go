package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

// Contact 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint // 谁的关系信息
	TargetID uint // 对应的谁
	Type     int  // 对用的类型 1,好友 2,群组 3，
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

// SearchFriend 查询好友，返回好友列表
func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)                                  // 声明一个长度为零的Contact数组
	objIDs := make([]uint64, 0)                                     // 声明一个unit类型的数组，用于存放查询到的id
	utils.DB.Where("owner_id=? and type=1", userId).Find(&contacts) // 查询数据并且保存在contacts里面
	for _, v := range contacts {
		fmt.Println(v)
		objIDs = append(objIDs, uint64(v.TargetID)) // 好友的id
	}
	fmt.Println(objIDs)
	users := make([]UserBasic, 0) // 声明一个数组存放查询到的好友信息
	utils.DB.Where("id in ?", objIDs).Find(&users)
	return users // 返回查询到的信息

}
