package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `grom:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

// FindUserByNameAndPwd 校验name和密码
func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}                                                       // 声明一个这种类型
	utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user) // 查找当前类型下name为传入值的操作，取第一个

	// token加密
	str := fmt.Sprintf("%d", time.Now().Unix()) // 获取当前时间
	temp := utils.Md5Encode(str)
	utils.DB.Model(&user).Where("id=?", user.ID).Update("identity", temp)

	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}                           // 声明一个这种类型
	utils.DB.Where("name = ?", name).First(&user) // 查找当前类型下name为传入值的操作，取第一个
	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user) // 查询到的值，直接赋值给user
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	// 创建用户返回
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	// 删除用户返回
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	// 更新用户返回
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord,
		Phone: user.Phone, Email: user.Email})
}
