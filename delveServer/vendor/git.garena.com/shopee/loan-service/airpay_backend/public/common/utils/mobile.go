package utils

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"strings"
)

func GetCompleteMobile(phone, country string) string {
	num, err := phonenumbers.Parse(phone, strings.ToUpper(country))
	if err != nil {
		return phone
	}
	return fmt.Sprintf("%d%d", num.GetCountryCode(), num.GetNationalNumber())
}

func GetNationalMobile(phone, country string) string {
	num, err := phonenumbers.Parse(phone, strings.ToUpper(country))
	if err != nil {
		return phone
	}
	return fmt.Sprintf("0%d", num.GetNationalNumber())
}
