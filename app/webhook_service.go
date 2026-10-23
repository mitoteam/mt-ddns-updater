package app

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-ddns-updater/model"
)

// DDNS logic for webhook

type WebhookService struct {
	webhook *model.DdnsWebhook
}

func NewWebhookService(webhook *model.DdnsWebhook) *WebhookService {
	return &WebhookService{webhook: webhook}
}

func (ws *WebhookService) QuestionIpAddress() (log *model.DdnsWebhookLog) {
	fullName := ws.webhook.FQDN()

	msg := new(dns.Msg)
	msg.SetQuestion(fullName, dns.TypeA)
	msg.RecursionDesired = true

	// prepare network client and new log entry
	client := new(dns.Client)
	log = model.NewWebhookLog(ws.webhook, model.LogKindQuestionGui)

	// save raw-message before sending it
	log.Message = msg.String()

	// send message to DNS server
	response, _, err := client.Exchange(msg, ws.webhook.DnsServerAddress)
	if err != nil {
		log.Error = err.Error()
		goapp.SaveObject(log)
		return log
	}

	//save RAW reply
	log.Reply = response.String()

	// check server's reply code (Rcode)
	if response.Rcode != dns.RcodeSuccess {
		log.Error = fmt.Sprintf("DNS Server rejected question. Rcode: %s", dns.RcodeToString[response.Rcode])
		goapp.SaveObject(log)
		return log
	}

	// find an IP in answer
	for _, ans := range response.Answer {
		// Is it an A-record (IPv4) ?
		if aRecord, ok := ans.(*dns.A); ok {
			// aRecord.A is net.IP
			log.Ip = aRecord.A.String()
		}
	}

	if log.Ip == "" {
		log.Error = "No IP address found in reply"
	}
	goapp.SaveObject(log)
	return log
}

func (ws *WebhookService) UpdateIpAddress(ip net.IP) (log *model.DdnsWebhookLog) {
	fullName := ws.webhook.FQDN()

	// prepare network client and new log entry
	client := new(dns.Client)
	log = model.NewWebhookLog(ws.webhook, model.LogKindUpdateGui)

	//prepare message
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
		msg.SetTsig(dns.Fqdn(key.Name), dns.Fqdn(key.Algo.Name()), 300, time.Now().Unix())
	}

	if key != nil {
		client.TsigSecret = map[string]string{dns.Fqdn(key.Name): key.Secret}
	}
	//client.Net = "tcp"

	// save raw-message before sending it
	log.Message = msg.String()

	// send message to DNS server
	response, _, err := client.Exchange(msg, ws.webhook.DnsServerAddress)
	if err != nil {
		log.Error = err.Error()
		goapp.SaveObject(log)
		return log
	}

	//save RAW reply
	log.Reply = response.String()

	// check server's reply code (Rcode)
	if response.Rcode == dns.RcodeSuccess {
		log.Ip = ip.String()
	} else {
		log.Error = fmt.Sprintf("DNS Server rejected update. Rcode: %s", dns.RcodeToString[response.Rcode])
	}

	goapp.SaveObject(log)
	return log
}
