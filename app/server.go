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
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	log.Println("Accepted connection from: ", conn.RemoteAddr())
	defer conn.Close()
	buf, err := ReadFromRequest(conn)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	data := strings.Split(string(buf), "\r\n")
	log.Printf("Read data: %+#v", data)
	PathMapper(data, conn)
}

func PathMapper(data []string, conn net.Conn) {
	switch {
	case strings.Contains(data[0], "echo"):
		ReadEchoPath(data, conn)
	case strings.Contains(data[0], " / "):
		WriteResponse(conn, 200, "OK", "text/plain", -1, "")
	case strings.Contains(data[0], "user-agent"):
		ReadHeadersAndWriteResponse(data, conn)
	default:
		WriteResponse(conn, 404, "Not Found", "text/plain", -1, "")
	}
}

func WriteResponse(conn net.Conn, statusCode int, statusString string, contentType string, contentLength int, body string) {
	log.Printf("Writing response: %d %s %s %d %s", statusCode, statusString, contentType, contentLength, body)
	if contentLength == -1 || contentType == "" {
		conn.Write(fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n\r\n", statusCode, statusString))
		return
	}
	conn.Write(fmt.Appendf(nil, "HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", statusCode, statusString, contentType, contentLength, body))

}

func ReadHeadersAndWriteResponse(data []string, conn net.Conn) {
	headerData := strings.Split(data[2], ": ")
	log.Printf("Header data: %+#v", headerData)
	WriteResponse(conn, 200, "OK", "text/plain", len(headerData[1]), headerData[1])
}

func ReadEchoPath(data []string, conn net.Conn) {
	routeData := strings.Split(data[0], " ")
	echoPath := strings.Split(routeData[1], "/")
	last := len(echoPath) - 1
	log.Printf("Echo path: %s", echoPath[last])
	WriteResponse(conn, 200, "OK", "text/plain", len(echoPath[last]), echoPath[last])
}

func ReadFromRequest(conn net.Conn) ([]byte, error) {
	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf[0:])
	if err != nil {
		return nil, err
	}
	log.Printf("Length: %d\n", n)
	log.Println("Data: ", string(buf[0:n]))
	return buf[0:n], nil
}
