package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

var (
	directory                 = flag.String("directory", "/tmp/", "Directory to serve files from")
	compressionMethodsAllowed = map[string]bool{
		"gzip": true,
	}
	plainHeaders = map[string]string{"Content-Type": "text/plain"}
)

const (
	compressionHeader = "Accept-Encoding"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	flag.Parse()
	fmt.Println("Logs from your program will appear here!")
	// Uncomment this block to pass the first stage
	//
	fmt.Println("Directory: ", *directory)
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
	case strings.Contains(data[0], "files"):
		HandleFileResponse(data, conn)
	case strings.Contains(data[0], " / "):

		WriteResponse(conn, 200, "OK", plainHeaders, "")
	case strings.Contains(data[0], "user-agent"):
		ReadHeadersAndWriteResponse(data, conn)
	default:
		WriteResponse(conn, 404, "Not Found", plainHeaders, "")
	}
}



func HandleFileResponse(data []string, conn net.Conn) {
	method := strings.Split(data[0], " ")[0]
	switch method {
	case "GET":
		filePath := ExtractFilePath(data)
		log.Printf("File path: %s", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			WriteResponse(conn, 404, "Not Found", plainHeaders, "")
			return
		}
		defer file.Close()
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			WriteResponse(conn, 404, "Not Found", plainHeaders, "")
			return
		}
		log.Printf("File data: %s", string(buf[0:n]))
		headers := makeHeaders(
			"Content-Type", "application/octet-stream",
			"Content-Length", fmt.Sprintf("%d", n))

		WriteResponse(conn, 200, "OK", headers, string(buf[0:n]))
	case "POST":
		log.Printf("POST request: %+#v", data)
		filePath := ExtractFilePath(data)
		requestBody := data[len(data)-1]
		log.Printf("Request body: %s", requestBody)
		file, err := os.Create(filePath)
		if err != nil {
			WriteResponse(conn, 500, "Internal Server Error", plainHeaders, "")
			return
		}
		defer file.Close()
		_, err = file.WriteString(requestBody)
		if err != nil {
			WriteResponse(conn, 500, "Internal Server Error", plainHeaders, "")
			return
		}
		WriteResponse(conn, 201, "Created", plainHeaders, "")
	}
}

func ReadHeadersAndWriteResponse(data []string, conn net.Conn) {
	headerData := strings.Split(data[2], ": ")
	log.Printf("Header data: %+#v", headerData)
	headers := makeHeaders("Content-Type", "text/plain", "Content-Length", fmt.Sprintf("%d", len(headerData[1])))
	WriteResponse(conn, 200, "OK", headers, headerData[1])
}

func ReadEchoPath(data []string, conn net.Conn) {
	headers := ExtractHeaders(data)
	log.Printf("Headers: %+#v", headers)
	routeData := strings.Split(data[0], " ")
	echoPath := strings.Split(routeData[1], "/")
	last := len(echoPath) - 1
	log.Printf("Echo path: %s", echoPath[last])
	var resHeaders map[string]string = makeHeaders(
		"Content-Type", "text/plain",
		"Content-Length", fmt.Sprintf("%d", len(echoPath[last])))
	for _,compressionMethod :=  range strings.Split(headers[compressionHeader]," "){
		if _, ok := compressionMethodsAllowed[compressionMethod]; ok{
			resHeaders["Content-Encoding"] = compressionMethod
			break
		}
	}
	WriteResponse(conn, 200, "OK", resHeaders, echoPath[last])
}