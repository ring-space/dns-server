package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/singeol/dns-server/internal/dns"
	"github.com/singeol/dns-server/internal/metrics"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Printf("invalid duration for %s: %v; using %s", key, err, fallback)
			return fallback
		}
		return d
	}
	return fallback
}

func main() {
	defaultAPIKey := getEnv("RECORDS_API_KEY", "")
	defaultPort := getEnv("DNS_PORT", ":53")
	defaultRecordsURL := getEnv("RECORDS_SERVICE_URL", "")
	defaultGeoDB := getEnv("GEOIP_DB_PATH", "")
	defaultRefresh := getEnvDuration("REFRESH_INTERVAL", 30*time.Second)
	defaultMetricsPort := getEnv("METRICS_PORT", ":9090")
	defaultMetricsFlush := getEnvDuration("METRICS_FLUSH", 5*time.Minute)

	apiKey := flag.String("api-key", defaultAPIKey, "API-ключ для запроса к сервису записей (X-Api-Key)")
	dnsPort := flag.String("port", defaultPort, "адрес и порт DNS-сервера, например :53")
	recordsURL := flag.String("records-url", defaultRecordsURL, "URL сервиса записей, пример http://localhost:8080/records")
	geoDBPath := flag.String("geoip-db", defaultGeoDB, "путь к GeoIP2 DB, например ./GeoLite2-City.mmdb")
	refresh := flag.Duration("refresh", defaultRefresh, "интервал обновления записей")
	metricsPort := flag.String("metrics-port", defaultMetricsPort, "порт для Prometheus /metrics")
	metricsFlush := flag.Duration("metrics-flush", defaultMetricsFlush, "интервал сброса и публикации метрик")
	flag.Parse()

	if *recordsURL == "" {
		log.Fatal("не задан records-url: установите флаг -records-url или ENV RECORDS_SERVICE_URL")
	}
	if *geoDBPath == "" {
		log.Fatal("не задан путь к GeoIP: установите флаг -geoip-db или ENV GEOIP_DB_PATH")
	}

	metrics.Init(*metricsPort, *metricsFlush)

	dns.Run(*dnsPort, *recordsURL, *geoDBPath, *refresh, *apiKey)
}
