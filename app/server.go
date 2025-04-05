package main

import (
	"fmt"
	"log"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		log.Println("Accepted connection from: ", conn.RemoteAddr())
		buf,err := ReadFromRequest(conn)
		if err!=nil{
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}
		data := strings.Split(string(buf), "\r\n")
		log.Printf("Read data: %+#v", data)
		if strings.Contains(data[0],"echo"){
			ReadEchoPath(data,conn)
		} else if strings.Contains(data[0]," / "){
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		conn.Close()
	}
}

func ReadEchoPath(data []string,conn net.Conn){
	routeData := strings.Split(data[0], " ")
	echoPath := strings.Split(routeData[1],"/")
	last := len(echoPath)-1
	log.Printf("Echo path: %s", echoPath[last])
	conn.Write(fmt.Appendf(nil, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",len(echoPath[last]),echoPath[last]))
}

func ReadFromRequest(conn net.Conn) ([]byte,error){
	var buf []byte = make([]byte,1024)
	n,err := conn.Read(buf[0:])
	if err!=nil{
		return nil,err
	}
	log.Printf("Length: %d\n", n)
	log.Println("Data: ", string(buf[0:n]))
	return buf[0:n],nil
}
