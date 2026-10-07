package app

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

type DdnsHelper struct {
	dnsServerAddress string // server address or IP with port, ex. "ns.example.com:53" or "8.8.4.4:53"
	dnsZone          string // full zone, ex. ".example.com."

	TsigKeyName string // key name
	TsigSecret  string // key itself in base64
	TsigAlgo    string // kay algorithm

	lastMessage string
	lastReply   string
}

func NewDdnsHelper(serverAddress, zone, keyName, secret string) DdnsHelper {
	h := DdnsHelper{
		dnsServerAddress: serverAddress,
		dnsZone:          dns.Fqdn(zone),

		TsigKeyName: dns.Fqdn(keyName),
		TsigSecret:  secret,
		TsigAlgo:    dns.HmacSHA256,
	}

	return h
}

func (h *DdnsHelper) GetLastMessage() string {
	return h.lastMessage
}

func (h *DdnsHelper) GetLastReply() string {
	return h.lastReply
}

func (h *DdnsHelper) UpdateARecordIpAddress(recordName string, newIP net.IP, newTtl uint32) error {
	fullName := dns.Fqdn(recordName) + h.dnsZone

	msg := new(dns.Msg)
	msg.SetUpdate(h.dnsZone)

	// Delete old A-record
	rrRemove := &dns.A{
		Hdr: dns.RR_Header{Name: fullName, Rrtype: dns.TypeA, Class: dns.ClassANY, Ttl: 0},
	}
	msg.RemoveName([]dns.RR{rrRemove})

	// Add new A-Record with new IP and TTL
	rrAdd := &dns.A{
		Hdr: dns.RR_Header{Name: fullName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: newTtl},
		A:   newIP,
	}
	msg.Insert([]dns.RR{rrAdd})

	// Add TSIG-key data
	if h.TsigKeyName != "" {
		msg.SetTsig(h.TsigKeyName, h.TsigAlgo, 300, time.Now().Unix())
	}

	// prepare network client
	client := new(dns.Client)
	if h.TsigKeyName != "" {
		client.TsigSecret = map[string]string{h.TsigKeyName: h.TsigSecret}
	}
	//client.Net = "tcp"

	// save raw-message before sending it
	h.lastMessage = msg.String()
	h.lastReply = "" //and clear reply

	// send message to DNS server
	response, _, err := client.Exchange(msg, h.dnsServerAddress)
	if err != nil {
		return err
	}

	//save RAW reply
	h.lastReply = response.String()

	// check server's reply code (Rcode)
	if response.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("DNS Server rejected update. Rcode: %s", dns.RcodeToString[response.Rcode])
	}

	return nil
}

func (h *DdnsHelper) GetARecordIpAddress(recordName string) (ip net.IP, err error) {
	fullName := dns.Fqdn(recordName) + h.dnsZone

	msg := new(dns.Msg)
	msg.SetQuestion(fullName, dns.TypeA)
	msg.RecursionDesired = true

	// prepare network client
	client := new(dns.Client)

	// save raw-message before sending it
	h.lastMessage = msg.String()
	h.lastReply = "" //and clear reply

	// send message to DNS server
	response, _, err := client.Exchange(msg, h.dnsServerAddress)
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
