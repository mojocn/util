package geoip

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
)

//返回值结构体
//需要满足以上要求
type Response struct {
	Country   string
	Province  string
	City      string
	ISP       string
	Latitude  float64
	Longitude float64
	TimeZone  string
}

var geoIpDb *geoip2.Reader

//json rpc 处理请求
//需要满足以上要求
func ParseGeo(ip string) (res *Response, err error) {
	res = &Response{}
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return nil, fmt.Errorf("%s is not a valid ip", ip)
	}
	//调用开源geoIp 数据库查询ip地址
	record, err := geoIpDb.City(netIp)
	if err != nil {
		return nil, err
	}
	res.City = record.City.Names["zh-CN"]
	res.Province = record.Subdivisions[0].Names["zh-CN"]
	res.Country = record.Country.Names["zh-CN"]
	res.Latitude = record.Location.Latitude
	res.Longitude = record.Location.Longitude
	res.TimeZone = record.Location.TimeZone

	return res, nil
}
func LoadIpMmdbFile(path string) error {
	mmdb, err := geoip2.Open(path)
	if err != nil {
		return err
	}
	geoIpDb = mmdb
	return nil
}
