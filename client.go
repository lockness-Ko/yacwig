package main

import "fmt"

type client struct {
	name  string
	queue []string
	ip    string
}

func (_this client) Send(message string) {
	(&_this).queue = append(_this.queue, message)
	fmt.Println(_this.queue)
}

func (_this client) Info() {
	fmt.Printf("%s: %s\n", _this.name, _this.ip)
}

func (_this client) Kill() {
	_this.Send("kill")
	(&_this).queue = nil
	fmt.Printf("%s: killed\n", _this.name)
}
