package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func Log(v ...interface{}) {
	fmt.Println(v...)
	return
}

func Sender(conn *net.TCPConn) {
	defer conn.Close()
	sc := bufio.NewReader(os.Stdin)
	go func() {
		t := time.NewTicker(time.Second) //创建定时器,用来实现定期发送心跳包给服务端
		defer t.Stop()
		for {
			<-t.C
			_, err := conn.Write([]byte("@"))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}()
	name := ""
	fmt.Println("please input username") //用户聊天的昵称
	fmt.Fscan(sc, &name)
	msg := ""
	buffer := make([]byte, 1024)
	//_t := time.NewTimer(time.Second * 5) //创建定时器,每次服务端发送消息就刷新时间
	//defer _t.Stop()
	//
	//go func() {
	//	<-_t.C
	//	fmt.Println("服务器出现故障，断开链接")
	//	return
	//}()
	////_t.Reset(time.Second * 5) //收到消息就刷新_t定时器，如果time.Second*5时间到了，那么就会<-_t.C就不会阻塞，代码会往下走，return结束
	//				//if string(buffer[0:1]) != "1" { //心跳包消息定义为字符串"1",不需要打印出来
	//				//	fmt.Println(string(buffer[0:n]))
	//				//}
	res := false
	for {
		if res{
			break
		}
		go func() {
			for{
				n, err := conn.Read(buffer)
				if err != nil {
				////break Loop
					fmt.Println("Server error")
					res = true
					return
				}
				fmt.Printf("server response: %s\n",string(buffer[0:n]))
				//if q,_ := strconv.Atoi(strings.Trim(string(buffer[0:n])," \r\n\t")); q == 2{
				//	fmt.Printf("%s: Exit.\n",name)
				//	res = true
				//	return
				//}else if q,_ := strconv.Atoi(strings.Trim(string(buffer[0:n])," \r\n\t")); q == 1{
				//	fmt.Printf("%s: Success.\n",name)
				//}else if q,_ := strconv.Atoi(strings.Trim(string(buffer[0:n])," \r\n\t")); q == -1{
				//	fmt.Printf("%s: Fail.\n",name)
				//}else{
				//	fmt.Printf("%s: Unknown.\n",name)
				//}
			}
		}()
		fmt.Fscan(sc, &msg)
		//i := time.Now().Format("2006-01-02 15:04:05")
		//发送消息
		conn.Write([]byte(fmt.Sprintf("%s:%s", name,msg)))
	}
}

func main() {
	server := "127.0.0.1:28086"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		Log(os.Stderr, "Fatal error:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Log("Fatal error:", err.Error())
		os.Exit(1)
	}
	Log(conn.RemoteAddr().String(), "connect success!")
	Sender(conn)
	Log("end")
}
