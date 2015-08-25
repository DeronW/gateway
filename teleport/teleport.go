package teleport

import (
	"fmt"
	"gateway/common"
	"gateway/configs"
	"net"
	"time"
)

type connection struct {
}

var connections []connection

func handleConn(conn net.Conn) {
	//connections = append(connections, conn)
	buf := make([]byte, 1024)
	length, _ := conn.Read(buf)

	if length > 0 {
		buf[length] = 0
	}

	ctx := string(buf[:length])

	fmt.Printf("%s", ctx)
}

func ping() {
	for i := range connections {
		fmt.Printf("%s", i)
		//if i.close_time < "now" { }
	}
	time.AfterFunc(30*time.Second, ping)
}

func Run() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", configs.HOST, configs.PORT))
	common.FatalError(err)

	lt, err := net.ListenTCP("tcp", addr)
	common.FatalError(err)

	for {
		conn, err := lt.Accept()
		common.CheckError(err)
		//err = conn.SetKeepAlive(true)
		//checkError(err)
		//err = conn.SetReadDeadline(time.Time{})
		//common.CheckError(err)
		//err = conn.SetWriteDeadline(time.Time{})
		//common.CheckError(err)
		go handleConn(conn)
	}

}
