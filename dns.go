package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

var serv server

func parseQuery(m *dns.Msg, ip string) {
	// parse query
	for _, q := range m.Question {
		// check if query is in records
		switch q.Qtype {
		// A record
		case dns.TypeA:
			// log.Printf("Query for %s\n", q.Name)
			// IP from records
			ip := incoming(q.Name, ip)
			for _, i := range ip {
				rr, _ := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, i))
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	// handle query
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = true

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, w.RemoteAddr().String())
	}

	w.WriteMsg(m)
}

func DNS(port int, serv server) {

	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
