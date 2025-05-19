package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/singeol/dns-server/internal/dns"
	"github.com/singeol/dns-server/internal/metrics"
)

// getEnv возвращает значение из ENV или возвращает fallback
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvDuration парсит duration из ENV или возвращает fallback
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
	// Defaults from ENV or hardcoded
	defaultPort := getEnv("DNS_PORT", ":53")
	defaultRecordsURL := getEnv("RECORDS_SERVICE_URL", "")
	defaultGeoDB := getEnv("GEOIP_DB_PATH", "")
	defaultRefresh := getEnvDuration("REFRESH_INTERVAL", 30*time.Second)
	defaultMetricsPort := getEnv("METRICS_PORT", ":9090")
	defaultMetricsFlush := getEnvDuration("METRICS_FLUSH", 5*time.Minute)

	// Flags (override ENV if provided)
	dnsPort := flag.String("port", defaultPort, "адрес и порт DNS-сервера, например :53")
	recordsURL := flag.String("records-url", defaultRecordsURL, "URL сервиса записей, пример http://localhost:8080/records")
	geoDBPath := flag.String("geoip-db", defaultGeoDB, "путь к GeoIP2 DB, например ./GeoLite2-City.mmdb")
	refresh := flag.Duration("refresh", defaultRefresh, "интервал обновления записей")
	metricsPort := flag.String("metrics-port", defaultMetricsPort, "порт для Prometheus /metrics")
	metricsFlush := flag.Duration("metrics-flush", defaultMetricsFlush, "интервал сброса и публикации метрик")
	flag.Parse()

	// Проверяем обязательные параметры
	if *recordsURL == "" {
		log.Fatal("не задан records-url: установите флаг -records-url или ENV RECORDS_SERVICE_URL")
	}
	if *geoDBPath == "" {
		log.Fatal("не задан путь к GeoIP: установите флаг -geoip-db или ENV GEOIP_DB_PATH")
	}

	// Стартуем HTTP-сервер метрик и флаш
	metrics.Init(*metricsPort, *metricsFlush)

	// Запускаем DNS
	dns.Run(*dnsPort, *recordsURL, *geoDBPath, *refresh)
}
