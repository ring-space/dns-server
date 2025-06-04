package dns

import (
	"log"
	"net"

	"github.com/miekg/dns"
	"github.com/singeol/dns-server/internal/geoip"
	"github.com/singeol/dns-server/internal/metrics"
	"github.com/singeol/dns-server/internal/records"
)

func MakeHandler(cl *records.Client) dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		host, _, _ := net.SplitHostPort(w.RemoteAddr().String())
		clientIP := net.ParseIP(host)

		country := geoip.CountryCode(clientIP)
		metrics.Record(country)

		msg := new(dns.Msg)
		msg.SetReply(r)
		msg.Authoritative = true

		if len(r.Question) == 0 {
			w.WriteMsg(msg)
			return
		}

		q := r.Question[0]
		log.Printf("[QUERY] %s from %s type %s", q.Name, clientIP, dns.TypeToString[q.Qtype])

		switch q.Qtype {
		case dns.TypeA:
			ips := cl.Get(q.Name)
			if len(ips) > 0 {
				sel := geoip.PickNearest(clientIP, ips)
				rr, err := dns.NewRR(q.Name + " A " + sel)
				if err == nil {
					msg.Answer = append(msg.Answer, rr)
					log.Printf("[RESPONSE] %s → %s", q.Name, sel)
				} else {
					log.Printf("[ERROR] при создании RR: %v", err)
				}
			} else {
				msg.Rcode = dns.RcodeNameError
			}
		default:
			msg.Rcode = dns.RcodeNotImplemented
		}

		w.WriteMsg(msg)
	}
}
