package metadata

import (
	"context"
	"fmt"
	"git.garena.com/shopee/loan-service/airpay_backend/public/common/config"
	"git.garena.com/shopee/loan-service/airpay_backend/public/common/utils"
	"git.garena.com/shopee/loan-service/airpay_backend/public/gateway_proto/airpay_gateway"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"
	"strconv"
	"strings"
)

//APA-Js/APA-Native/APA-Python/APC-Js/APC-Native/APA-SDK
const ( //const 尽量私有，避免全局污染，外部需要访问再提供get函数
	from_apa_js           = "APA-Js"
	from_apa_native       = "APA-Native"
	from_apa_python       = "APA-Python"
	from_apc_js           = "APC-Js"
	from_apc_native       = "APC-Native"
	from_apa_sdk          = "APA-SDK"
	airpay_header_bin     = "airpay-header-bin"
	aripay_operator_email = "a-operator-email"
)

const (
	UserTypeFrom   = "A-UserType"
	UserType_User  = "APA-User"
	UserType_Staff = "APA-Staff"
)

func IsRequestFromStaff(ctx context.Context) bool {
	return strings.EqualFold(UserType_Staff, GetStringFromCtx(ctx, UserTypeFrom))
}

func GetAFromByCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, config.From)
}

func IsRequestFromApaJs(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apa_js)
}

func IsRequestFromApaNative(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apa_native)
}

func IsRequestFromApaPython(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apa_python)
}

func IsRequestFromApcJs(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apc_js)
}

func IsRequestFromApcNative(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apc_native)
}

func IsRequestFromApaSdk(ctx context.Context) bool {
	headerFrom := GetAFromByCtx(ctx)
	return strings.EqualFold(headerFrom, from_apa_sdk)
}

func GetUidFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, "uid")
}

func GetVersionFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, config.AppVer)
}

func GetSdkVersionFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, config.SdkVer)
}

//兼容 python过来的取值
func GetVersionCodeFromCtx(ctx context.Context) string {
	airpayHeader, err := GetAirpayHeader(ctx)
	if err == nil && airpayHeader != nil && airpayHeader.GetAppVersion() != 0 {
		return fmt.Sprint(airpayHeader.GetAppVersion())
	}
	return GetStringFromCtx(ctx, config.AppVerCode)
}

//兼容 python过来的取值
func GetLangFromCtx(ctx context.Context) string {
	airpayHeader, err := GetAirpayHeader(ctx)
	if err == nil && airpayHeader != nil && airpayHeader.Lang != "" {
		return airpayHeader.Lang
	}
	return GetStringFromCtx(ctx, config.Lang)
}

func GetOsFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, config.OS)
}

func GetOsVersionFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, config.OSVer)
}

func GetCallFromInfoFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, "call_from")
}

func GetUidInt64FromCtx(ctx context.Context) (int64, error) {
	strUid := GetUidFromCtx(ctx)
	return strconv.ParseInt(strUid, 10, 64)
}

func GetRemoteIP(ctx context.Context) string {
	return GetStringFromCtx(ctx, "remote_ip")
}

func GetIntFromCtx(ctx context.Context, key string) (int, error) {
	s := GetStringFromCtx(ctx, key)
	return strconv.Atoi(s)
}

func IsFromPython(ctx context.Context) bool {
	mData, blOK := metadata.FromIncomingContext(ctx) //从 ctx 提取元数据
	if !blOK {
		return false
	}

	strs := mData.Get(airpay_header_bin)
	if strs == nil || len(strs) == 0 {
		return false
	}
	return true
}

func GetStringFromCtx(ctx context.Context, key string) (s string) {
	mData, blOK := metadata.FromIncomingContext(ctx) //从 ctx 提取元数据
	if blOK {
		values := mData.Get(key)
		if values != nil && len(values) > 0 {
			s = values[0]
		}
	}
	return
}

func GetAirpayHeader(ctx context.Context) (*airpay_gateway.AirpayHeader, error) {
	str := GetStringFromCtx(ctx, airpay_header_bin)
	var data []byte = []byte(str)
	var airpayHeader airpay_gateway.AirpayHeader

	err := proto.Unmarshal(data, &airpayHeader)
	return &airpayHeader, err

}

func GetRemoteIPInt64(ctx context.Context) int64 {
	strIp := GetStringFromCtx(ctx, "remote_ip")
	return utils.InetAtoN(strIp)
}

func SetUidIntoContext(ctx context.Context, uid string) context.Context {
	metadatas := []string{}
	metadatas = append(metadatas, "uid", uid)
	ctx = metadata.AppendToOutgoingContext(ctx, metadatas...)
	return ctx
}

func SetRemoteipIntoCtx(ctx context.Context, remoteIp string) context.Context {
	metadatas := []string{}
	metadatas = append(metadatas, "remote_ip", remoteIp)
	ctx = metadata.AppendToOutgoingContext(ctx, metadatas...)
	return ctx
}

func GetOperatorEmailFromCtx(ctx context.Context) string {
	return GetStringFromCtx(ctx, aripay_operator_email)
}