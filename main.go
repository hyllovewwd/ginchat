package main

import (
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig() // 初始化配置文件
	utils.InitMySql()  // 初始化mysql数据库
	utils.InitRedis()  // 初始化redis数据库
	r := router.Router()
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
