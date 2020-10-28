package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	DEV         = "dev"
	TEST        = "test"
	PRODUCT     = "live"
	UAT         = "uat"
	LIVESTAGING = "livestaging"
	TH          = "th"
	VN          = "vn"
)

const envKey = "ENV"

var (
	RuntimeEnv string
	APP_NAME   string
	APP_TYPE   string
	REGION     string //th or vn
)

//bool-true:有设置环境,false:没有设置该环境变量
func GetGeneralEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func GetAppName() string {
	return APP_NAME
}

func IsAppTransparentGateway(appName string) bool {
	return appName == "transparent-gateway"
}

func GetRegion() string {
	return REGION
}

func GetAppType() string {
	return APP_TYPE
}

func IsDev() bool {
	return strings.ToLower(GetEnv()) == DEV
}

func IsTest() bool {
	return strings.ToLower(GetEnv()) == TEST
}

func IsVN() bool {
	return strings.ToLower(GetRegion()) == VN
}

func IsTH() bool {
	return strings.ToLower(GetRegion()) == TH
}

func IsProduct() bool {
	return strings.ToLower(GetEnv()) == PRODUCT
}

func IsUat() bool {
	return strings.ToLower(GetEnv()) == UAT
}

func IsLivestaging() bool {
	return strings.ToLower(GetEnv()) == LIVESTAGING
}

func IsInMac() bool {
	return "darwin" == runtime.GOOS
}

func GetEnv() string {
	return RuntimeEnv
}

func initAppType() {
	val, _ := os.LookupEnv("APPTYPE")
	APP_TYPE = val
}

func initAppName() {
	val, _ := os.LookupEnv("APP_NAME")
	APP_NAME = val
}

func initRegion() {
	val, _ := os.LookupEnv("REGION")
	REGION = val
}

func init() {
	val, ok := os.LookupEnv(envKey)
	if !ok {
		RuntimeEnv = DEV
		fmt.Fprintf(os.Stderr, "ENV not set, force to env[%s]\n", RuntimeEnv)
		return
	}

	RuntimeEnv = val

	if RuntimeEnv == "" {
		RuntimeEnv = DEV
		fmt.Fprintf(os.Stderr, "invalid ENV[=%s], force to env[%s]\n", val, RuntimeEnv)
	}

	initAppName()
	initAppType()
	initRegion()
	/*
		switch RuntimeEnv {
		case DEV:
		case TEST:
		case PRODUCT:
		default:
		    old = RuntimeEnv
			RuntimeEnv = DEV
			fmt.Fprintf(os.Stderr, "invalid ENV[=%s], force to env[%s]\n", old, RuntimeEnv)
			return
		}*/

	fmt.Fprintf(os.Stderr, "set ENV=%s, appName=%s, appType=%s\n", RuntimeEnv, APP_NAME, APP_TYPE)
}
