package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)


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


func makeHeaders(keyVals ...string) map[string]string {
	headers := make(map[string]string)
	for i := 0; i < len(keyVals); i += 2 {
		headers[keyVals[i]] = keyVals[i+1]
	}

	return headers
}


func ExtractFilePath(data []string) string {
	fileName := strings.Split(strings.Split(data[0], " ")[1], "/")[2]
	filePath := *directory + fileName
	return filePath
}


func WriteResponse(conn net.Conn, statusCode int, statusString string, headers map[string]string, body string) {
	log.Printf("Writing response: %d %s %+#v %s", statusCode, statusString, headers, body)
	sb := strings.Builder{}
	sb.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusString)))
	for key, value := range headers {
		sb.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
	}
	sb.Write([]byte("\r\n"))
	sb.Write([]byte(body))
	log.Printf("Response: %s", sb.String())
	_, err := conn.Write([]byte(sb.String()))
	if err != nil {
		log.Printf("Error writing response: %s", err.Error())
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}
}

func ExtractHeaders(data []string) map[string]string {
	headers := make(map[string]string)
	for i := 1; i < len(data)-2; i++ {
		headerData := strings.Split(data[i], ": ")
		if len(headerData) == 2 {
			headers[headerData[0]] = headerData[1]
		}
	}
	return headers
}