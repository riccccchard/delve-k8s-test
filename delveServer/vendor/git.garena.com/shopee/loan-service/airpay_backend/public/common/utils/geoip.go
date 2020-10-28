package utils

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

var (
	geoipFile *geoip2.Reader
)

func InitGeoIp(geoPath string) error {
	if geoipFile != nil {
		return nil
	}
	var err error
	geoipFile, err = geoip2.Open(geoPath)
	if err != nil {
		return err
	}
	return nil
}

func CloseGeoIp() {
	if geoipFile != nil {
		defer geoipFile.Close()
	}
}

func GeoIPToCountry(ip int64) string {
	if geoipFile == nil {
		return ""
	}

	ipStr := InetNtoA(ip)

	record, err := geoipFile.City(net.ParseIP(ipStr))
	if err != nil {
		return ""
	}

	if record != nil {
		return record.Country.IsoCode
	}

	return ""
}
