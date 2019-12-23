package geoip

import (
	"testing"
)

func TestParseGeo(t *testing.T) {
	err := LoadIpMmdbFile("../_data/GeoLite2-City.mmdb")
	if err != nil {
		t.Error(err)
		return
	}
	res, err := ParseGeo("39.106.87.48")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
