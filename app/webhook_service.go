package app

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/mitoteam/mt-ddns-updater/model"
)

// DDNS logic for webhook

type WebhookService struct {
	webhook *model.DdnsWebhook

	lastMessage string
	lastReply   string
}

func NewWebhookService(webhook *model.DdnsWebhook) *WebhookService {
	return &WebhookService{webhook: webhook}
}

func (ws *WebhookService) GetLastMessage() string {
	return ws.lastMessage
}

func (ws *WebhookService) GetLastReply() string {
	return ws.lastReply
}

func (ws *WebhookService) QuestionIpAddress() (net.IP, error) {
	fullName := ws.webhook.FQDN()

	msg := new(dns.Msg)
	msg.SetQuestion(fullName, dns.TypeA)
	msg.RecursionDesired = true

	// prepare network client
	client := new(dns.Client)

	// save raw-message before sending it
	ws.lastMessage = msg.String()
	ws.lastReply = "" //and clear reply

	// send message to DNS server
	response, _, err := client.Exchange(msg, ws.webhook.DnsServerAddress)
	if err != nil {
		return nil, err
	}

	//save RAW reply
	ws.lastReply = response.String()

	// check server's reply code (Rcode)
	if response.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS Server rejected question. Rcode: %s", dns.RcodeToString[response.Rcode])
	}

	// find and IP in answer
	for _, ans := range response.Answer {
		// Is it an A-record (IPv4) ?
		if aRecord, ok := ans.(*dns.A); ok {
			// aRecord.A has net.IP
			return aRecord.A, nil
		}
	}

	return nil, fmt.Errorf("No IP address found in reply")
}

func (ws *WebhookService) UpdateIpAddress(ip net.IP) error {
	fullName := ws.webhook.FQDN()

	msg := new(dns.Msg)
	msg.SetUpdate(dns.Fqdn(ws.webhook.DnsZoneName))

	// Delete old A-record
	rrRemove := &dns.A{
		Hdr: dns.RR_Header{Name: fullName, Rrtype: dns.TypeA, Class: dns.ClassANY, Ttl: 0},
	}
	msg.RemoveName([]dns.RR{rrRemove})

	// Add new A-Record with new IP and TTL
	rrAdd := &dns.A{
		Hdr: dns.RR_Header{Name: fullName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ws.webhook.Ttl},
		A:   ip,
	}
	msg.Insert([]dns.RR{rrAdd})

	// Add TSIG-key data
	key := ws.webhook.GetKey()
	if key != nil {
		msg.SetTsig(dns.Fqdn(key.Name), dns.Fqdn(key.Algo.String()), 300, time.Now().Unix())
	}

	// prepare network client
	client := new(dns.Client)
	if key != nil {
		client.TsigSecret = map[string]string{dns.Fqdn(key.Name): key.Secret}
	}
	//client.Net = "tcp"

	// save raw-message before sending it
	ws.lastMessage = msg.String()
	ws.lastReply = "" //and clear reply

	// send message to DNS server
	response, _, err := client.Exchange(msg, ws.webhook.DnsServerAddress)
	if err != nil {
		return err
	}

	//save RAW reply
	ws.lastReply = response.String()

	// check server's reply code (Rcode)
	if response.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("DNS Server rejected update. Rcode: %s", dns.RcodeToString[response.Rcode])
	}

	return nil //no error
}
