package geoip

import (
	"log"
	"math"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var db *geoip2.Reader

func Init(path string) {
	var err error
	db, err = geoip2.Open(path)
	if err != nil {
		log.Fatalf("Не удалось открыть GeoIP2 DB: %v", err)
	}
}

func CountryCode(ip net.IP) string {
	rec, err := db.City(ip)
	if err != nil {
		return "unknown"
	}
	if rec.Country.IsoCode != "" {
		return rec.Country.IsoCode
	}
	return "unknown"
}

func PickNearest(clientIP net.IP, candidates []string) string {
	rec, err := db.City(clientIP)
	if err != nil {
		return candidates[0]
	}
	clat, clon := rec.Location.Latitude, rec.Location.Longitude

	nearest := candidates[0]
	minD := math.MaxFloat64

	for _, ipStr := range candidates {
		srvIP := net.ParseIP(ipStr)
		s, err := db.City(srvIP)
		if err != nil {
			continue
		}
		d := haversine(clat, clon, s.Location.Latitude, s.Location.Longitude)
		if d < minD {
			minD = d
			nearest = ipStr
		}
	}
	return nearest
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
