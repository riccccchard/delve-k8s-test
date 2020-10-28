package metadata

import (
	"context"
	"strconv"
	"strings"
)

const (
	from_apa_ms_js     = "APA-MS-Js"
	from_apa_ms_native = "APA-MS-Native"
)

func GetStaffIdFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, "staff_id")
}

func GetStaffIdInt64FromCtx(ctx context.Context) (uint64, error) {
	staffId := GetStaffIdFromCtx(ctx)
	return strconv.ParseUint(staffId, 10, 64)
}

func GetStaffDeviceIdFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, "staff_device_id")
}

func GetStaffDeviceIdInt64FromCtx(ctx context.Context) (uint64, error) {
	staffDeviceId := GetStringFromCtx(ctx, "staff_device_id")
	return strconv.ParseUint(staffDeviceId, 10, 64)
}

func IsRequestFromApaMsJs(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apa_ms_js)
}

func IsRequestFromApaMsNative(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apa_ms_native)
}
