package main

import (
	"fmt"
	"strconv"
	"net/http"
	"os"
)

func main(){
	var start, end  int
	fmt.Printf("请输入起始页( >= 1) :")
	fmt.Scan(&start)
	fmt.Printf("请输入终止页( >= 起始页) :")
	fmt.Scan(&end)
	DoWork(start,end)
}


func DoWork(start,end int) {
	fmt.Printf("正在爬取 %d 到 %d 的页面\n", start, end)
	page :=make(chan int)
	for i :=start;i<=end ;i++  {
		go SpiderPage(i,page)
	}
	for i :=start;i<=end ;i++  {
		fmt.Printf("第%d个页面爬取完成\n", <-page)
	}
}

func SpiderPage(i int, page chan<- int) {
	url :="http://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn="+strconv.Itoa((i-1)*50)
	fmt.Printf("正在爬第%d页网页: %s\n", i, url)
	
	error,result := HttpGet(url)
	if error != nil {
		fmt.Println("HttpGet err = ", error)
		return
	}
	filename :=strconv.Itoa(i)+".html"
	file , error :=os.Create(filename)
	if error!=nil {
		fmt.Println("os.Create err1 = ", error)
		return
	}
	file.WriteString(result)
	file.Close()
	page<-i

}
func HttpGet(url string) (error error,result string){
	response,error1 := http.Get(url)
	if  error1!= nil{
		error=error1
		return
	}
	defer response.Body.Close()

	buf :=make([]byte,1024*4)
	for   {
		n,_ := response.Body.Read(buf)
		if n==0 {
			break
		}
		result += string(buf[:n])
	}
	return
}