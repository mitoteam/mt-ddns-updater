package app

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	"github.com/mitoteam/mt-ddns-updater/model"
)

// DDNS logic for webhook

type WebhookHelper struct {
	webhook *model.DdnsWebhook

	lastMessage string
	lastReply   string
}

func NewWebhookHelper(webhook *model.DdnsWebhook) *WebhookHelper {
	return &WebhookHelper{webhook: webhook}
}

func (h *WebhookHelper) GetLastMessage() string {
	return h.lastMessage
}

func (h *WebhookHelper) GetLastReply() string {
	return h.lastReply
}

func (h *WebhookHelper) QuestionIpAddress() (ip net.IP, err error) {
	fullName := h.webhook.FQDN()

	msg := new(dns.Msg)
	msg.SetQuestion(fullName, dns.TypeA)
	msg.RecursionDesired = true

	// prepare network client
	client := new(dns.Client)

	// save raw-message before sending it
	h.lastMessage = msg.String()
	h.lastReply = "" //and clear reply

	// send message to DNS server
	response, _, err := client.Exchange(msg, h.webhook.DnsServerAddress)
	if err != nil {
		return nil, err
	}

	//save RAW reply
	h.lastReply = response.String()

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
