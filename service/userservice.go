package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0表示成功 ，-1 表示失败
		"message": "获取用户列表成功",
		"data":    data,
	})
}

// FindUserByNameAndPwd
// Summary 所有用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	//name := c.Query("name")
	name := c.Request.FormValue("name")
	//password := c.Query("password")
	password := c.Request.FormValue("password")
	user := models.FindUserByName(name) // 通过Name查找该用户
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0表示成功 ，-1 表示失败
			"message": "该用户不存在",
		})
		return
	}
	flag := utils.VaildPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0表示成功 ，-1 表示失败
			"message": "密码不正确",
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0表示成功 ，-1 表示失败
		"message": "登录成功",
		"data":    data,
	})
}

// CreateUser
// Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31()) // 随机生成一个salt，参与字符串加密
	findData := models.FindUserByName(user.Name)
	fmt.Println(findData)
	if user.Name == "" || password == "" || repassword == "" {

		// 查询到当前用户名不为空，即注册过了
		c.JSON(-1, gin.H{
			"code":    -1, // 0表示成功 ，-1 表示失败
			"message": "用户名为空！",
		})
		return // 直接返回
	}
	if findData.Name != "" {
		// 查询到当前用户名不为空，即注册过了
		c.JSON(-1, gin.H{
			"code":    -1, // 0表示成功 ，-1 表示失败
			"message": "用户名已注册！",
		})
		return // 直接返回
	}
	if repassword != password {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0表示成功 ，-1 表示失败
			"message": "两次输入密码不一致！",
		})
		return
	}
	user.PassWord = utils.MakePassword(password, salt) // 密码加密
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0表示成功 ，-1 表示失败
		"message": "新增用户成功！",
		"data":    user,
	})

}

// DeleteUser
// Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id")) // 转化为整形
	user.ID = uint(id)                   // 转化为uint类型
	models.DeleteUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0表示成功 ，-1 表示失败
		"message": "删除用户成功！",
		"data":    user,
	})
}

// UpdateUser
// Summary 更新用户信息
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param email formData string false "email"
// @param phone formData string false "phone"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id")) // 转化为整形
	user.ID = uint(id)                      // 转化为uint类型
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	_, err := govalidator.ValidateStruct(user) // 校验手机号和邮箱
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0表示成功 ，-1 表示失败
			"message": "更新用户信息失败，参数校验失败！",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0, // 0表示成功 ，-1 表示失败
			"message": "更新用户信息成功！",
			"data":    user,
		})
	}
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息。。", msg)
		tm := time.Now().Format("2006-01-02 15:04:05") // 获取当前的时间
		m := fmt.Sprintf("[ws][%s]:[%s]", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId")) //获取发送随请求发送过来的id
	users := models.SearchFriend(uint(id))               // 查找当前id对应的好友列表
	// 返回当前查询到的信息
	//c.JSON(http.StatusOK, gin.H{
	//	"code":    0, // 0表示成功 ，-1 表示失败
	//	"message": "查找好友信息成功！",
	//	"data":    users,
	//})
	// 将后台的返回进行封装
	utils.RespOkList(c.Writer, users, len(users))
}
