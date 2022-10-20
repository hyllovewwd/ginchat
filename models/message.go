package models

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/set"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// Message 消息
type Message struct {
	gorm.Model
	FormId   int64  // 发送者
	TargetId int64  // 接收者
	Type     int    // 发送类型 群聊，私聊，广播
	Media    int    // 消息类型 文字，图片，表情
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// Chat 需要 ；发送者ID，接收者ID，消息类型，发送的内容，发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1.获取参数并校验token等合法性
	query := request.URL.Query() // 获取请求下来的数据
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64) // str转化类型成为int64
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isValida := true // checkToken() 需要去请求数据库进行校验（后期加入）
	conn, err := (&websocket.Upgrader{
		// token校验
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2.获取conn，从Node
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),   // 类似于管道，用户存放数据
		GroupSets: set.New(set.ThreadSafe), // 设置线程安全
	}

	// 3.用户关系
	// 4.userId与node绑定 并加锁
	rwLocker.Lock() // 加锁
	clientMap[userId] = node
	rwLocker.Unlock() // 释放锁

	//5.完成发送逻辑
	go sendProc(node)
	//6.完成接受逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue: // 从通道里面取出数据赋值给data
			fmt.Println("[ws] sendProc >>>>>>> msg:", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data) // 将消息发送出去
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage() // 将来自conn里面的数据接受出来
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		broadMsg(data)                                    // 广播接收到的数据
		fmt.Println("[ws] recvProc <<<<< ", string(data)) // string(data)将接收到的数据进行转换成string类型
	}
}

var udpsendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	fmt.Println("broadMsg data:", string(data))
	udpsendChan <- data
}

func init() {
	//init函数先于main函数自动执行，不能被其他函数调用；
	//init函数没有输入参数、返回值；
	//每个包可以有多个init函数；
	//包的每个源文件也可以有多个init函数，这点比较特殊；
	//同一个包的init执行顺序，golang没有明确定义，编程时要注意程序不要依赖这个执行顺序。
	//不同包的init函数按照包导入的依赖关系决定执行顺序。
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("Init goroutine")
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan: // 从通道里面取出数据赋值给data
			fmt.Println("udpSendProc data:", string(data))
			_, err := con.Write(data) // 将消息发送出去
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

// 完成udp数据接受协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero, // 设置所有端口都可以
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc data:", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg) // json数据转化处理
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: // 私信
		fmt.Println("dispatch data:", string(data))
		sendMsg(msg.TargetId, data)
		//case 2: // 群发
		//	sendGroupMsg()
	}
}
func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg >>> userId:", userId, "msg:", string(msg))
	rwLocker.RLock() // 获取锁
	node, ok := clientMap[userId]
	rwLocker.RUnlock() // 解锁
	if ok {
		node.DataQueue <- msg
	}
}
