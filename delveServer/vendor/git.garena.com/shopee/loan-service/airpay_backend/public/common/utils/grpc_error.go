package utils

import (
	pb "git.garena.com/shopee/loan-service/airpay_backend/public/common/utils/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCError(code uint32, errMsg, showMsg string) error {
	msg := &pb.ExtraMsg{
		ShowMessage: showMsg,
	}
	s := status.New(codes.Code(code), errMsg)
	ns, err := s.WithDetails(msg)
	if err != nil {
		return s.Err()
	}
	return ns.Err()
}
func NewGRPCErrorReply(code uint32, errMsg, showMsg string, reply []byte) error {
	msg := &pb.ExtraMsg{
		ShowMessage: showMsg,
		Reply:       reply,
	}
	s := status.New(codes.Code(code), errMsg)
	ns, err := s.WithDetails(msg)
	if err != nil {
		return s.Err()
	}
	return ns.Err()
}

//return: code, errMsg, showMsg
func ExtractGRPCError(e error) (codes.Code, string, string) {
	var code codes.Code
	var errMsg string
	var showMsg string

	st, _ := status.FromError(e)
	code = st.Code()
	errMsg = st.Message()
	showMsg = errMsg //showMsg 默认等于 errMsg
	details := st.Details()
	if len(details) > 0 {
		sp, ok := details[0].(*pb.ExtraMsg)
		if ok {
			showMsg = sp.GetShowMessage()
		}
	}
	return code, errMsg, showMsg
}

//return: code, errMsg, showMsg
func ExtractGRPCErrorReply(e error) (codes.Code, string, string, []byte) {
	var code codes.Code
	var errMsg string
	var showMsg string
	var reply []byte
	st, _ := status.FromError(e)
	code = st.Code()
	errMsg = st.Message()
	showMsg = errMsg //showMsg 默认等于 errMsg
	details := st.Details()
	if len(details) > 0 {
		sp, ok := details[0].(*pb.ExtraMsg)
		if ok {
			showMsg = sp.GetShowMessage()
			reply = sp.GetReply()
		}
	}
	return code, errMsg, showMsg, reply
}

/**
 * GRPC ERROR CODE:1-16 [0:OK,17:_maxCode]
 * code in [1-16] and showMessage=="" return true else retrun false
 * param:code,showMessage
 * return:bool
 */
func IsGRPCError(code codes.Code, showMessage string) bool {
	if code > 0 && code < 17 && showMessage == "" {
		return true
	}
	return false
}

//key 去掉双引号做匹配
var strToCode = map[string]codes.Code{
	"OK": codes.OK,
	"CANCELLED":/* [sic] */ codes.Canceled,
	"UNKNOWN":             codes.Unknown,
	"INVALID_ARGUMENT":    codes.InvalidArgument,
	"DEADLINE_EXCEEDED":   codes.DeadlineExceeded,
	"NOT_FOUND":           codes.NotFound,
	"ALREADY_EXISTS":      codes.AlreadyExists,
	"PERMISSION_DENIED":   codes.PermissionDenied,
	"RESOURCE_EXHAUSTED":  codes.ResourceExhausted,
	"FAILED_PRECONDITION": codes.FailedPrecondition,
	"ABORTED":             codes.Aborted,
	"OUT_OF_RANGE":        codes.OutOfRange,
	"UNIMPLEMENTED":       codes.Unimplemented,
	"INTERNAL":            codes.Internal,
	"UNAVAILABLE":         codes.Unavailable,
	"DATA_LOSS":           codes.DataLoss,
	"UNAUTHENTICATED":     codes.Unauthenticated,
}
