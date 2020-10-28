package utils

import (
	"git.garena.com/shopee/loan-service/airpay_backend/public/common/utils/region_util"
	"github.com/shopspring/decimal"
	"os"
	"time"
)

const (
	StandardDateTimeFormat = "2006-01-02 15:04:05"
	StandardDateFormat     = "2006-01-02"
	StandardTimeFormat     = "15:04:05"
)

// 时间戳转换为string显示
func TimestampToString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func GetTimeZone() string {
	return os.Getenv("TZ")
}

func GetTimeUnix(timeStr string) uint32 {
	tz := GetTimeZone()
	loc, _ := time.LoadLocation(tz)
	tt, _ := time.ParseInLocation(StandardDateTimeFormat, timeStr, loc)
	return uint32(tt.Unix())
}

func CurrentRegionTime() time.Time {
	regionCode := region_util.ValueOf(GetRegion())
	return CurrentTimeByRegion(regionCode)
}

func CurrentTimeByRegion(region region_util.RegionCode) time.Time {
	loc := GetRegionTimeLocation(region)
	return time.Now().In(loc)
}

func GetLocalTimeLocation() *time.Location {
	return GetRegionTimeLocation(region_util.ValueOf(GetRegion()))
}

func GetRegionTimeLocation(region region_util.RegionCode) *time.Location {
	// time.LoadLocation not support load 'GMT+7' and 'GMT+8'.
	// Etc/GMT-7 = GMT+7
	// Etc/GMT-8 = GMT+8
	var zone string
	switch {
	case region.BeIncluded(region_util.RegionIndonesia, region_util.RegionThailand, region_util.RegionVietnam):
		zone = "Etc/GMT-7"
	case region.BeIncluded(region_util.RegionSingapore, region_util.RegionPhilippines, region_util.RegionMalaysia, region_util.RegionTaiWan):
		zone = "Etc/GMT-8"
	default:
		zone = "Local"
	}
	loc, _ := time.LoadLocation(zone)
	return loc
}

func TimeSubDays(startTime, endTime time.Time) int {
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, GetLocalTimeLocation())
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, GetLocalTimeLocation())

	return int(endTime.Sub(startTime).Hours()/24) + 1
}

func ConvertGoWeekdayToPythonWeekday(goWeekday int) int {
	switch goWeekday {
	case 0:
		return 6
	case 1:
		return 0
	case 2:
		return 1
	case 3:
		return 2
	case 4:
		return 3
	case 5:
		return 4
	case 6:
		return 5
	}
	return goWeekday
}

func CurrentRegionTimestamp() int {
	return int(CurrentRegionTime().Unix())
}

func GetOneDayStartTimestamp(day time.Time) int {
	newDayStart := int(time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, GetLocalTimeLocation()).Unix())
	return newDayStart
}

func GetRegionTimeStartOfMonth(timestamp int32) int32 {
	dateTime := time.Unix(int64(timestamp), 0)
	newDayStart := int32(time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, GetLocalTimeLocation()).Unix())
	return newDayStart
}

func GetRegionTimeStartOfWeek(timestamp int32) int32 {
	startOfDay := GetRegionTimeStartOfDay(timestamp)
	dateTime := time.Unix(int64(timestamp), 0)
	weekday := int(dateTime.Weekday())
	return startOfDay - int32(86400*ConvertGoWeekdayToPythonWeekday(weekday))
}

func GetRegionTimeStartOfDay(timestamp int32) int32 {
	dateTime := time.Unix(int64(timestamp), 0)
	newDayStart := int32(time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, GetLocalTimeLocation()).Unix())
	return newDayStart
}

func IsLeapYear(year int32) bool {
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}

func GetMonthLastDay(month int32, year int32) int32 {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	}
	if IsLeapYear(year) {
		return 29
	}
	return 28
}

func GetAddYearNum(month int32, addMonth int32) int32 {
	return int32(decimal.New(int64(month+addMonth), 0).
		Div(decimal.New(12, 0)).
		Sub(decimal.New(1, 0)).Ceil().IntPart())
}
