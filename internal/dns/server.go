package dns

import (
	"log"
	"time"

	"github.com/miekg/dns"
	"github.com/singeol/dns-server/internal/geoip"
	"github.com/singeol/dns-server/internal/records"
)

func Run(addr, recordsURL, geoDBPath string, refreshIntervalDuration time.Duration, apiKey string) {
	recordClient := records.NewClient(recordsURL, refreshIntervalDuration, apiKey)
	geoip.Init(geoDBPath)

	go recordClient.StartAutoRefresh()

	dns.HandleFunc(".", MakeHandler(recordClient))

	server := &dns.Server{Addr: addr, Net: "udp"}
	log.Printf("Starting DNS server on %s/udp…", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
