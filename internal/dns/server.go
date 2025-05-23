package dns

import (
	"log"
	"time"

	"github.com/miekg/dns"
	"github.com/singeol/dns-server/internal/geoip"
	"github.com/singeol/dns-server/internal/records"
)

// Run инициализирует пакеты и запускает DNS-сервер
func Run(addr, recordsURL, geoDBPath string, refreshIntervalDuration time.Duration, apiKey string) {
	// Инициализация GeoIP и records-клиента
	recordClient := records.NewClient(recordsURL, refreshIntervalDuration, apiKey)
	geoip.Init(geoDBPath)

	// Автообновление записей
	go recordClient.StartAutoRefresh()

	// Регистрируем обработчик
	dns.HandleFunc(".", MakeHandler(recordClient))

	// Запуск
	server := &dns.Server{Addr: addr, Net: "udp"}
	log.Printf("Starting DNS server on %s/udp…", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
