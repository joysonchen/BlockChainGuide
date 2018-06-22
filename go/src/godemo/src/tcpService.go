package main

import "fmt"
import "net"

func main(){
	 listener,err := net.Listen("tcp","127.0.0.1:8000")
	if  err!=nil {
		fmt.Print("err=",err)
		return
	}

	defer listener.Close()

	conn,err := listener.Accept()
	if err!=nil {
		fmt.Print("err=",err)
		return
	}
	buf := make([]byte,1024)
	n,err := conn.Read(buf)
	if err!=nil {
		fmt.Print("err1=",err)
		return
	}
	fmt.Print("buf=",string(buf[0:n]))
	defer conn.Close()

}