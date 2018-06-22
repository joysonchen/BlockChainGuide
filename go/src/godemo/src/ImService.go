package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

var onlineMap map[string]Client
var message = make(chan string)
var isQuit = make(chan bool)
var isHasData = make(chan bool)

func main() {
	listener, error := net.Listen("tcp", ":8000")
	if error != nil {
		fmt.Println("net.Listen,error=", error)
		return
	}
	defer listener.Close()

	go ManagerMessage()
	for {
		conn, error := listener.Accept()
		if error != nil {
			fmt.Println("listener.Accept,error=", error)
			continue
		}
		go handleConn(conn)

	}
}

//协程3：对上线所有用户进行消息通知
func ManagerMessage() {
	onlineMap = make(map[string]Client)
	for {
		msg := <-message
		for _, cli := range onlineMap {
			cli.C <- msg
		}
	}
}

type Client struct {
	C    chan string
	name string
	addr string
}

//处理连接 1.获取addr，且给client发送消息
func handleConn(conn net.Conn) {
	defer conn.Close()
	cliAddr := conn.RemoteAddr().String()
	client := Client{make(chan string), cliAddr, cliAddr}
	onlineMap[cliAddr] = client
	//新开协程1，专门给client发送任务
	go WriteMsgToClient(client, conn)
	//发送给message ，是遍历通知所有client，发送具体client，只有当前client会收到消息
	message <- makeMsg(client, "login")
	client.C <- makeMsg(client,"i am here")


	//新建协程2，用来处理用户发过来的数据
	go handleMsgFromClient(client,conn)
	handleLoginOrTimeout(client)

}
//用于检测用户是否退出,或者是否超时
func handleLoginOrTimeout(client Client)  {
	for {
		select {
		case <-isQuit:
			delete(onlineMap,client.addr)
			message<-makeMsg(client,"login out")
			return
		case <-isHasData:
		case <-time.After(60*time.Second):
			delete(onlineMap,client.addr)
			message<-makeMsg(client,"time out leave out")
			return
		}
	}
}
func handleMsgFromClient(client Client,conn net.Conn) {
	buf := make([]byte,2048)
	for  {
		n,error := conn.Read(buf)
		if error!=nil || n==0 {
			isQuit<-true
			fmt.Println("conn.Read error,",error)
			return
		}
		msg := string(buf[0:n-1])//在win环境下，nc结尾会多一个换行

		printUserList(client,msg,conn)
		isHasData<-true
	}
}

//打印用户列表
func printUserList(client Client,msg string,conn net.Conn) {
	if len(msg) == 3 && msg == "who" {
		conn.Write([]byte("user list:\n"))
		for _,tmp := range onlineMap  {
			msg="userName: " +tmp.name +",Address: "+tmp.addr +"\n"
			conn.Write([]byte(msg))
		}
	}else if len(msg)>8 && msg[0:6]=="rename" {
		name :=strings.Split(msg,"|")[1]
		client.name=name
		onlineMap[client.addr]=client
		conn.Write([]byte("rename is :"+name+"\n"))
	} else {
		message<- makeMsg(client,msg)
	}
}


func WriteMsgToClient(client Client, conn net.Conn) {
	for msg := range client.C {
		conn.Write([]byte(msg + "\n"))
	}
}

func makeMsg(client Client, msg string) (buf string) {
	buf = "[" + client.addr + "]" + client.name + ": " + msg
	return
}
