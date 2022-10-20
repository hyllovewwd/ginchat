package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H struct {
	Code  int
	Msg   string
	Data  interface{} // 接收任意类型的消息
	Rows  interface{}
	Total interface{} // 一共多少行
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json") // 设置进行连接的数据类型
	w.WriteHeader(http.StatusOK)                       // 成功的状态码
	// 输出定义好的结构体
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	// 转化成json格式
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	// 将转换得到的json数据“写出去”
	w.Write(ret)
}
func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json") // 设置进行连接的数据类型
	w.WriteHeader(http.StatusOK)                       // 成功的状态码
	// 输出定义好的结构体
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	// 转化成json格式
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	// 将转换得到的json数据“写出去”
	w.Write(ret)

}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}
func RespOk(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}
func RespOkList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
