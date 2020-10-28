package log

import (
	"context"
	"fmt"
	"git.garena.com/shopee/loan-service/airpay_backend/public/common/metadata"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
)

const (
	ENUM_STACK_SKIP_TWO = 2
	OnceLogStrKey       = "OnceLogKey"
)

type OnceLogKey struct{}
type MonitorOnceLogHandlerType func(level LogLevel, grpcMethod string)

var (
	monitorOnceLogHandler MonitorOnceLogHandlerType
)

func SetMonitorOnceLogHandler(monitorOnceLogHandlerInput MonitorOnceLogHandlerType) {
	monitorOnceLogHandler = monitorOnceLogHandlerInput
}

type OnceLog struct {
	skipLevel  int
	uid        string
	grpcMethod string
	traceID    string
	//funcName   string
	//startTime time.Time
	customerPrefix  string
	prefixStr       string // "|method:%s|traceId:%s|uid:%s|"  前后有竖线
	blPrintFileLine bool   //是否打印文件名和行号,默认是 false 不打印, 高并发下比较耗性能，谨慎开启.参考: https://cloud.tencent.com/developer/article/1385947
}

// GetTraceID 返回jaeger traceid，ctx是grpc server method的context
func getTraceID(ctx context.Context) string {
	traceId := ""
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if spanCtx, ok := span.Context().(jaeger.SpanContext); ok {
			traceId = spanCtx.TraceID().String()
		}
	}
	return traceId
}

func GetTraceIDFromContext(ctx context.Context) string {
	return getTraceID(ctx)
}

func (c *OnceLog) SetCustomerPrefix(customerPrefix string) {
	c.customerPrefix = customerPrefix
}

func (c *OnceLog) SetUid(uid string) {
	c.uid = uid
}

func (c *OnceLog) SetGrpcMethod(grpcMethod string) {
	c.grpcMethod = grpcMethod
}

func (c *OnceLog) SetTraceID(traceID string) {
	c.traceID = traceID
}

func (c *OnceLog) SetSkipLevel(skipLevel int) {
	c.skipLevel = skipLevel
}

func (c *OnceLog) SetPrintFileLine(blPrintFileLine bool) {
	c.blPrintFileLine = blPrintFileLine
}

//实例化一个 OnceLog 对象
func NewOnceLogFromContext(ctx context.Context) *OnceLog {
	methodName, _ := grpc.Method(ctx)
	//fileName, line, funcName := findFileInfo(ENUM_STACK_SKIP_TWO)

	onceLog := new(OnceLog)
	if onceLog == nil {
		//打印错误日志
		Errorf("OnceLogFromContext faile.methodName:%s", methodName)
		return nil
	}

	onceLog.traceID = getTraceID(ctx)
	onceLog.uid = metadata.GetUidFromCtx(ctx)
	onceLog.grpcMethod = methodName
	onceLog.skipLevel = skipLevel
	onceLog.blPrintFileLine = true //默认先打开，量大了再关掉
	//onceLog.prefixStr = fmt.Sprintf("|grpcMethod:%s|funcName:%s|traceId:%s|uid:%s|file:%s|line:%d",
	//onceLog.grpcMethod, onceLog.funcName, onceLog.traceID, onceLog.uid, fileName, line)
	onceLog.prefixStr = fmt.Sprintf("|grpcMethod:%s|traceId:%s|uid:%s|",
		onceLog.grpcMethod, onceLog.traceID, onceLog.uid)

	return onceLog
}

//实例化一个 OnceLog 对象
func NewOnceLogFromContextGin(ctx *gin.Context) *OnceLog {
	methodName := ctx.Request.URL.RequestURI()
	//fileName, line, funcName := findFileInfo(ENUM_STACK_SKIP_TWO)

	onceLog := new(OnceLog)
	if onceLog == nil {
		//打印错误日志
		Errorf("OnceLogFromContext faile.methodName:%s", methodName)
		return nil
	}

	//onceLog.traceID = ctx.GetHeader("tracingId")
	//if onceLog.traceID == "" {
	//	onceLog.traceID = getTraceID(ctx.Request.Context())
	//}

	onceLog.traceID = getTraceID(ctx.Request.Context())
	if onceLog.traceID == "" {
		onceLog.traceID = ctx.GetHeader("tracingId")
	}

	onceLog.uid = ctx.GetHeader("uid")
	onceLog.skipLevel = skipLevel
	onceLog.grpcMethod = methodName
	onceLog.blPrintFileLine = true //默认先打开，量大了再关掉
	//onceLog.prefixStr = fmt.Sprintf("|grpcMethod:%s|funcName:%s|traceId:%s|uid:%s|file:%s|line:%d",
	//onceLog.grpcMethod, onceLog.funcName, onceLog.traceID, onceLog.uid, fileName, line)
	onceLog.prefixStr = fmt.Sprintf("|grpcMethod:%s|traceId:%s|uid:%s|",
		onceLog.grpcMethod, onceLog.traceID, onceLog.uid)
	ctx.Set(OnceLogStrKey, onceLog)
	return onceLog
}
func ExtractOnceLog(ctx context.Context) *OnceLog {
	onceLog, ok := ctx.Value(OnceLogKey{}).(*OnceLog)
	if !ok {
		return NewOnceLogFromContext(ctx)
	}
	return onceLog
}

func ExtractOnceLogGin(ctx *gin.Context) *OnceLog {
	onceLog, ok := ctx.Get(OnceLogStrKey)
	if !ok {
		return NewOnceLogFromContextGin(ctx)
	}
	return onceLog.(*OnceLog)
}

func (c *OnceLog) Trace(fields ...interface{}) {
	head := c.getFullPrefix(TraceLvl, c.skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Trace(fields...)
}

func (c *OnceLog) Debug(fields ...interface{}) {
	head := c.getFullPrefix(DebugLvl, c.skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Debug(fields...)
}

func (c *OnceLog) Info(fields ...interface{}) {
	head := c.getFullPrefix(InfoLvl, c.skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Info(fields...)
}

func (c *OnceLog) Warn(fields ...interface{}) {
	head := c.getFullPrefix(WarnLvl, c.skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Warn(fields...)
}

func (c *OnceLog) Error(fields ...interface{}) {
	head := c.getFullPrefix(ErrorLvl, c.skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Error(fields...)
}

func (c *OnceLog) Critical(fields ...interface{}) {
	head := c.getFullPrefix(CriticalLvl, c.skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Critical(fields...)
}

// 格式化接口
func (c *OnceLog) Tracef(format string, fields ...interface{}) {
	head := c.getFullPrefix(TraceLvl, c.skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Trace(head, formatStr)
}

func (c *OnceLog) Debugf(format string, fields ...interface{}) {
	head := c.getFullPrefix(DebugLvl, c.skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Debug(head, formatStr)
}

func (c *OnceLog) Infof(format string, fields ...interface{}) {
	head := c.getFullPrefix(InfoLvl, c.skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Info(head, formatStr)
}

func (c *OnceLog) Warnf(format string, fields ...interface{}) {
	head := c.getFullPrefix(WarnLvl, c.skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Warn(head, formatStr)
}

func (c *OnceLog) Errorf(format string, fields ...interface{}) {
	head := c.getFullPrefix(ErrorLvl, c.skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Error(head, formatStr)
}

func (c *OnceLog) Criticalf(format string, fields ...interface{}) {
	head := c.getFullPrefix(CriticalLvl, c.skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Critical(head, formatStr)
}

//todo del
func (c *OnceLog) formatLog2(level LogLevel, skip int) string {
	//获取日志打印文件、行数、函数名
	fileName, line, funcName := findFileInfo(skip)
	return fmt.Sprintf("[%s] [%v] [%v/%v]", seelog.LogLevel(level).String(), funcName, fileName, line)
}

func (c *OnceLog) getFullPrefix(level LogLevel, skip int) string {
	//获取日志打印文件、行数、函数名
	fileName, line, funcName := "not_print", 0, "not_print"
	if c.blPrintFileLine {
		fileName, line, funcName = findFileInfo(skip)
	}
	codeInfo := fmt.Sprintf("file:%s/%v|func:%s|", fileName, line, funcName)
	var customerPrefix = ""
	if c.customerPrefix != "" {
		customerPrefix = fmt.Sprintf("%s|", c.customerPrefix)
	}
	c.prefixStr = fmt.Sprintf("|grpcMethod:%s|traceId:%s|uid:%s|",
		c.grpcMethod, c.traceID, c.uid)
	fullPrefix := c.prefixStr + codeInfo + customerPrefix

	if monitorOnceLogHandler != nil { //上报日志条数
		monitorOnceLogHandler(level, c.grpcMethod)
	}

	return fmt.Sprintf("%s%s", seelog.LogLevel(level).String(), fullPrefix)
}

// 日志格式按参数顺序组装key1:value1|key2:value2
func (c *OnceLog) Tracew(fields ...interface{}) {
	head := c.getFullPrefix(TraceLvl, c.skipLevel)
	formatStr := formatw(fields)
	seelog.Trace(head, formatStr)
}

func (c *OnceLog) Debugw(fields ...interface{}) {
	head := c.getFullPrefix(DebugLvl, c.skipLevel)
	formatStr := formatw(fields)
	seelog.Debug(head, formatStr)
}

func (c *OnceLog) Infow(fields ...interface{}) {
	head := c.getFullPrefix(InfoLvl, c.skipLevel)
	formatStr := formatw(fields)
	seelog.Info(head, formatStr)
}

func (c *OnceLog) Warnw(fields ...interface{}) {
	head := c.getFullPrefix(WarnLvl, c.skipLevel)
	formatStr := formatw(fields)
	seelog.Warn(head, formatStr)
}

func (c *OnceLog) Errorw(fields ...interface{}) {
	head := c.getFullPrefix(ErrorLvl, c.skipLevel)
	formatStr := formatw(fields)
	seelog.Error(head, formatStr)
}

func (c *OnceLog) Criticalw(fields ...interface{}) {
	head := c.getFullPrefix(CriticalLvl, c.skipLevel)
	formatStr := formatw(fields)
	seelog.Critical(head, formatStr)
}

func (c *OnceLog) Flush() {
	seelog.Flush()
}

func Flush() {
	FlushFlowLog()
	seelog.Flush()
}
