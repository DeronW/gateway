package teleport

import (
	"fmt"
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

}

func ping() {
	for i := range connections {
		fmt.Printf("%s", i)
		//if i.close_time < "now" { }
	}
	time.AfterFunc(30*time.Second, ping)
}

func Run() {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.PORT))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		err = conn.SetReadDeadline(time.Time{})
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		err = conn.SetWriteDeadline(time.Time{})
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		go handleConn(conn)
	}

}
