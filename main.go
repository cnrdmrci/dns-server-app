package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"time"
)

func resolver(domain string, qtype uint16) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := &dns.Client{Timeout: 5 * time.Second}

	response, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		log.Fatalf("[ERROR] : %v\n", err)
		return nil
	}

	if response == nil {
		log.Fatalf("[ERROR] : no response from server\n")
		return nil
	}

	for _, answer := range response.Answer {
		fmt.Printf("%s\n", answer.String())
	}

	return response.Answer
}

var domainsToAddresses map[string]string = map[string]string{
	"google.com.":     "1.2.3.4",
	"www.google.com.": "1.2.3.4",
}

type dnsHandler struct{}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	switch r.Question[0].Qtype {
	case dns.TypeA:
		domain := msg.Question[0].Name
		address, ok := domainsToAddresses[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
			for _, answer := range msg.Answer {
				fmt.Printf("%s\n", answer.String())
			}
		} else {
			//for _, question := range r.Question {
			//	answers := resolver(question.Name, question.Qtype)
			//	msg.Answer = append(msg.Answer, answers...)
			//}
		}
	}

	w.WriteMsg(msg)
}

func StartDNSServer() {
	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      ":53",
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Println("Starting DNS server on port 53")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}

func main() {
	StartDNSServer()
}
