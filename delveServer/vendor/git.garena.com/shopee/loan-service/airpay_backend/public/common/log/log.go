package log

import (
	"fmt"
	"git.garena.com/shopee/loan-service/airpay_backend/public/common/utils"
	"github.com/cihub/seelog"
	"github.com/go-xorm/core"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"
)

type LogLevel seelog.LogLevel

const (
	TraceLvl    = seelog.TraceLvl
	DebugLvl    = seelog.DebugLvl
	InfoLvl     = seelog.InfoLvl
	WarnLvl     = seelog.WarnLvl
	ErrorLvl    = seelog.ErrorLvl
	CriticalLvl = seelog.CriticalLvl
	Off         = seelog.Off
)

const skipLevel = 3

const rlogConfig = `
<seelog type="asynctimer" asyncinterval="5000000" minlevel="{{.Level}}" maxlevel="critical">
    <outputs formatid="common">
        <rollingfile formatid="common" type="date" filename="./log/{{.RandDir}}/rlog.log" datepattern="2006-01-02"/>
    </outputs>
    <formats>
        <format id="common" format="%Date(2006-01-02 15:04:05.000000Z07:00) %Msg%n" />
    </formats>
</seelog>
`

var rLog *rLogger

func InitLog(level LogLevel) (string, error) {
	randDir := ""

	podName, ok := os.LookupEnv("POD_NAME")
	if ok && podName != "" {
		randDir = podName
	}

	if randDir == "" && utils.IsInsideDocker() {
		rand.Seed(time.Now().UnixNano() ^ int64(os.Getpid()))
		randDir = fmt.Sprintf("%s-%d", time.Now().Format("20060102150405"), rand.Int()%1000)
	}

	initLog(level, randDir)

	InitFlowLog(randDir) //初始化流水日志

	return randDir, nil
}

func initLog(level LogLevel, randDir string) error {
	configMap := map[string]string{
		"Level":   fmt.Sprintf("%s", seelog.LogLevel(level).String()),
		"RandDir": randDir,
	}
	tmpl, err := template.New("config").Parse(rlogConfig)
	sb := &strings.Builder{}
	tmpl.Execute(sb, configMap)
	logger, err := seelog.LoggerFromConfigAsString(sb.String())
	if err != nil {
		return err
	}
	rLog = &rLogger{logger, false}
	err = seelog.ReplaceLogger(rLog.logger)

	if err != nil {
		return err
	}
	return nil
}

type rLogger struct {
	logger  seelog.LoggerInterface
	showSQL bool
}

func (log *rLogger) Debug(v ...interface{}) {
	log.logger.Debug(v)
}

func (log *rLogger) Debugf(format string, v ...interface{}) {
	log.logger.Debugf(format, v)
}

func (log *rLogger) Info(v ...interface{}) {
	log.logger.Info(v)
}

func (log *rLogger) Infof(format string, v ...interface{}) {
	log.logger.Infof(format, v)
}

func (log *rLogger) Warn(v ...interface{}) {
	log.logger.Warn(v)
}

func (log *rLogger) Warnf(format string, v ...interface{}) {
	log.logger.Warnf(format, v)
}

func (log *rLogger) Error(v ...interface{}) {
	log.logger.Error(v)
}

func (log *rLogger) Errorf(format string, params ...interface{}) {
	log.logger.Errorf(format, params)
}

func (log *rLogger) Level() core.LogLevel {
	return 0
}

func (log *rLogger) SetLevel(core.LogLevel) {}

func (log *rLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		log.showSQL = true
	} else {
		log.showSQL = show[0]
	}
}

func (log *rLogger) IsShowSQL() bool {
	return log.showSQL
}

func GetRLog() *rLogger {
	return rLog
}

func Trace(fields ...interface{}) {
	head := formatLog(TraceLvl, skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Trace(fields...)
}

func Debug(fields ...interface{}) {
	head := formatLog(DebugLvl, skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Debug(fields...)
}

func Info(fields ...interface{}) {
	head := formatLog(InfoLvl, skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Info(fields...)
}

func Warn(fields ...interface{}) {
	head := formatLog(WarnLvl, skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Warn(fields...)
}

func Error(fields ...interface{}) {
	head := formatLog(ErrorLvl, skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Error(fields...)
}

func Critical(fields ...interface{}) {
	head := formatLog(CriticalLvl, skipLevel)
	fields = append([]interface{}{head}, fields...)
	seelog.Critical(fields...)
}

// 格式化接口
func Tracef(format string, fields ...interface{}) {
	head := formatLog(TraceLvl, skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Trace(head, formatStr)
}

func Debugf(format string, fields ...interface{}) {
	head := formatLog(DebugLvl, skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Debug(head, formatStr)
}

func Infof(format string, fields ...interface{}) {
	head := formatLog(InfoLvl, skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Info(head, formatStr)
}

func Warnf(format string, fields ...interface{}) {
	head := formatLog(WarnLvl, skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Warn(head, formatStr)
}

func Errorf(format string, fields ...interface{}) {
	head := formatLog(ErrorLvl, skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Error(head, formatStr)
}

func Criticalf(format string, fields ...interface{}) {
	head := formatLog(CriticalLvl, skipLevel)
	formatStr := fmt.Sprintf(format, fields...)
	seelog.Critical(head, formatStr)
}

//标准格式化日志
func formatLog(level LogLevel, skip int) string {
	//获取日志打印文件、行数、函数名
	fileName, line, funcName := findFileInfo(skip)

	return fmt.Sprintf("[%s] [%v] [%v/%v]", seelog.LogLevel(level).String(), funcName, fileName, line)
}

//查找文件信息
func findFileInfo(skip int) (string, int, string) {
	fileName, line, funcName := "???", 0, "???"

	// 这个获取函数名、文件名的操作有性能损耗，日志量大的话要关掉
	//https://cloud.tencent.com/developer/article/1385947
	pc, fileName, line, ok := runtime.Caller(skip)

	if ok {
		funcName = runtime.FuncForPC(pc).Name()
		funcName = filepath.Ext(funcName)
		funcName = strings.TrimPrefix(funcName, ".")

		fileName = filepath.Base(fileName)
	}
	return fileName, line, funcName
}
