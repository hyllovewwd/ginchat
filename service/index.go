package service

import (
	"fmt"
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	//ind, err := template.ParseFiles("view/chat/index.html",
	//	"view/chat/head.html",
	//	"view/chat/tabmenu.html",
	//	"view/chat/concat.html",
	//	"view/chat/main.html",
	//	"view/chat/profile.html",
	//	"view/chat/group.html",
	//	"view/chat/foot.html")
	ind, err := template.ParseFiles("view/user/login.html",
		"view/chat/head.html")
	//ind, err := template.ParseFiles("view/chat/main.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		fmt.Println("oooooooooooooooooo", err)
		return
	}
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("view/chat/index.html",
		"view/chat/head.html",
		"view/chat/tabmenu.html",
		"view/chat/concat.html",
		"view/chat/main.html",
		"view/chat/profile.html",
		"view/chat/group.html",
		"view/chat/foot.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	err = ind.Execute(c.Writer, user)
	if err != nil {
		fmt.Println("oooooooooooooooooo", err)
		return
	}
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("view/user/register.html",
		"view/chat/head.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "register")
	if err != nil {
		fmt.Println("oooooooooooooooooo", err)
		return
	}
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
