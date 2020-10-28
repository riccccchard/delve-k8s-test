package region_util

import (
	"strings"
)

type RegionCode string

// regions
const (
	RegionIndonesia   RegionCode = "id"
	RegionThailand    RegionCode = "th"
	RegionTaiWan      RegionCode = "tw"
	RegionPhilippines RegionCode = "ph"
	RegionVietnam     RegionCode = "vn"
	RegionSingapore   RegionCode = "sg"
	RegionMalaysia    RegionCode = "my"
)

func (rc RegionCode) Lower() string {
	return strings.ToLower(string(rc))
}

func (rc RegionCode) Upper() string {
	return strings.ToUpper(string(rc))
}

func (rc RegionCode) String() string {
	return rc.Upper()
}

func (rc RegionCode) BeIncluded(codes ...RegionCode) bool {
	if len(codes) == 0 {
		return false
	}
	lc := rc.Lower()
	for _, c := range codes {
		if c.Lower() == lc {
			return true
		}
	}
	return false
}

func (rc RegionCode) Equal(code RegionCode) bool {
	return rc.Lower() == code.Lower()
}

func ValueOf(region string) RegionCode {
	var regionCode RegionCode
	switch {
	case strings.ToLower(region) == "id":
		regionCode = RegionIndonesia
	case strings.ToLower(region) == "th":
		regionCode = RegionThailand
	case strings.ToLower(region) == "tw":
		regionCode = RegionTaiWan
	case strings.ToLower(region) == "ph":
		regionCode = RegionPhilippines
	case strings.ToLower(region) == "vn":
		regionCode = RegionVietnam
	case strings.ToLower(region) == "sg":
		regionCode = RegionSingapore
	case strings.ToLower(region) == "my":
		regionCode = RegionMalaysia
	default:
		regionCode = RegionSingapore
	}
	return regionCode
}
