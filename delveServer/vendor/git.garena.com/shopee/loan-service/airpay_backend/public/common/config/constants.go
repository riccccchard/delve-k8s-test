package config

import "google.golang.org/grpc/codes"

const ErrorMicroServPanic codes.Code = 10012
const ErrorShouldTrans2Uat codes.Code = 10014 //应该把请求转发到uat域名

const (
	TimeStamp  string = "A-TimeStamp"
	Token      string = "A-Token"
	OS         string = "A-OS"
	OSVer      string = "A-OSVer"
	AppVer     string = "A-AppVer"     //字符串版本号
	SdkVer     string = "A-Sdk-Ver" //sdk字符串版本号
	AppVerCode string = "A-AppVerCode" //整数型版本号
	Sign       string = "A-Sign"
	Lang       string = "A-Lang"
	Channel    string = "A-Channel"
	From       string = "A-From" // APA-Js/APA-Native/APA-Python/APC-Js/APC-Native
)
