package log

import (
	"fmt"
	log "github.com/cihub/seelog"
	"strconv"
	"sync"
	"time"
)

type MonitorFLogHandlerType func(grpcMethod string)

var (
	flowLog            log.LoggerInterface
	once               sync.Once
	_randDir           string
	monitorFLogHandler MonitorFLogHandlerType
)

const maxLogLen = 200 //一开始不要太大,只要把resp的commonHeader打印出来就好
const configTemplate = `
<seelog type="asynctimer" asyncinterval="5000000" minlevel="trace" maxlevel="critical">

    <outputs formatid="common">
        <rollingfile formatid="common" type="date" filename="./log/%s/flog.log" datepattern="2006-01-02" maxrolls="30" />
    </outputs>

    <formats>
        <format id="common" format="%s" />
    </formats>

</seelog>
`

func doInitFlowLog() {
	logFormat := "%Date(2006-01-02 15:04:05.000000Z07:00)%Msg%n" //Z07:00 表示把时区也打印出来
	config := fmt.Sprintf(configTemplate, _randDir, logFormat)
	flowLog, _ = log.LoggerFromConfigAsString(config)
}

func SetMonitorFLogHandler(monitorFLogHandlerTypeInput MonitorFLogHandlerType) {
	monitorFLogHandler = monitorFLogHandlerTypeInput
}

func InitFlowLog(randDir string) {
	_randDir = randDir
	once.Do(doInitFlowLog)
}

func FlowLogInfo(traceId string, retCode int, timeCost time.Duration, uid string, fullMethod string, custStr string, req interface{}, resp interface{}) {
	// 2019-03-13 16:07:20.069468|traceid|code|耗时(ms)|/hello.Greeter/SayHello|【自定义】|name:"you",id:123|message:"hello you",code:0|
	//flowLog.Infof("|%s|%d|%d|%s|%s|%s|%v|%v|",
	//traceId, retCode, timeCost, uid, fullMethod, custStr, req, resp)
	flowLog.Infof("|%s|%d|%d|%s|%s|%s|%v|%v|",
		traceId, retCode, timeCost, uid, fullMethod, custStr, req, "<rsp not print>")

}

/***
 * param printResp(bool) default false
 * false:output <rsp not print>
 * true:output len>2000?resp[0:2000]+"...leave( "+(len-2000)+")":resp
 * eg:false <rsp not print>
 *    true msggfdklgfidugidfhgkjf...(leave 500)
 */
func FlowLogInfoWithFlag(traceId string, retCode int, timeCost time.Duration, uid string, fullMethod string, custStr string, req interface{}, resp interface{}, printResp bool) {
	var respStr = "<rsp not print>" //default not print
	if printResp {
		respStr = fmt.Sprint(resp)
		strLen := len(respStr)
		if strLen > maxLogLen {
			respStr = respStr[0:maxLogLen-1] + "...(leave " + strconv.Itoa(strLen-maxLogLen) + ")"
		}
	}

	if monitorFLogHandler != nil { //上报flog的调用次数
		monitorFLogHandler(fullMethod)
	}

	flowLog.Infof("|%s|%d|%d|%s|%s|%s|%v|%v|",
		traceId, retCode, timeCost, uid, fullMethod, custStr, req, respStr)
}

/***
 * param printResp(bool) default false,param printRequest(bool) default true
 *
 */
func FlowLogInfoWithConfig(traceId string, retCode int, timeCost time.Duration, uid string, fullMethod string, custStr string, req interface{}, resp interface{}, printResp bool, printRequest bool, callFrom string) {
	var respStr = "<rsp not print>" //default not print
	var requestStr = fmt.Sprint(req)
	if printResp {
		respStr = fmt.Sprint(resp)
		strLen := len(respStr)
		if strLen > maxLogLen {
			respStr = respStr[0:maxLogLen-1] + "...(leave " + strconv.Itoa(strLen-maxLogLen) + ")"
		}
	}
	if !printRequest {
		requestStr = "<request not print>" // not print
	}

	if monitorFLogHandler != nil { //上报flog的调用次数
		monitorFLogHandler(fullMethod)
	}

	flowLog.Infof("|%s|%d|%d|%s|%s|%s|%s|%s|%s|",
		traceId, retCode, timeCost, uid, fullMethod, callFrom, requestStr, respStr, custStr)
}

func FlushFlowLog() {
	if flowLog != nil {
		flowLog.Flush()
	}
}
