package offilne_geocoding_db

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoLite2Geocoder struct {
	db *geoip2.Reader
}

// NewGeoLite2Geocoder открывает базу GeoLite2 по указанному пути.
func NewGeoLite2Geocoder(dbPath string) (*GeoLite2Geocoder, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening GeoLite2 database at %s: %v", dbPath, err)
	}
	return &GeoLite2Geocoder{db: db}, nil
}

// LocalGeocode возвращает координаты (lat, lon) для данного IP.
// Если база не содержит записи или возникает ошибка, возвращается ошибка.
func (g *GeoLite2Geocoder) LocalGeocode(ip string) (float64, float64, error) {
	if ip == "" {
		return 0, 0, fmt.Errorf("empty IP")
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, 0, fmt.Errorf("invalid IP: %s", ip)
	}
	record, err := g.db.City(parsedIP)
	if err != nil {
		return 0, 0, fmt.Errorf("error querying GeoLite2 database: %v", err)
	}
	lat := record.Location.Latitude
	lon := record.Location.Longitude
	if lat == 0 && lon == 0 {
		return 0, 0, fmt.Errorf("no location found for IP: %s", ip)
	}
	return lat, lon, nil
}
