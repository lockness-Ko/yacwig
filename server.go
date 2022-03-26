package main

import (
	"fmt"
	"strings"
)

type server struct {
	clients map[string]client
	port    int
}

func bytesToIP(b []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func pad(tobepadded string) string {
	// pad to 4 bytes
	for len(tobepadded)%4 != 0 {
		tobepadded += " "
	}
	return tobepadded
}

func ipencode(str string) []string {
	str = pad(str)
	// split str into chunks of 4 bytes
	chunks := make([][]byte, len(str)/4)
	for i := 0; i < len(str)/4; i++ {
		chunks[i] = []byte(str[i*4 : i*4+4])
	}

	// convert each chunk to an IP
	ips := make([]string, len(chunks))
	for i, chunk := range chunks {
		ips[i] = bytesToIP(chunk)
	}

	return ips
}

func incoming(query string) []string {
	return ipencode(fmt.Sprintf("Hello %s, How are you?", strings.Split(query, ".")[0]))
}

func Serve(port int) server {
	// Start the DNS server
	go DNS(port)

	return server{}
}
