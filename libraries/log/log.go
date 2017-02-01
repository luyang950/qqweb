package log

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// err levels
const (
	EMERGENCY uint8 = iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var (
	levelMap = map[uint8]string{
		EMERGENCY: "emergency",
		ALERT:     "alert",
		CRITICAL:  "critical",
		ERROR:     "error",
		WARNING:   "warning",
		NOTICE:    "notice",
		INFO:      "info",
		DEBUG:     "debug",
	}
	LevelMapRev = map[string]uint8{
		"emergency": EMERGENCY,
		"alert":     ALERT,
		"critical":  CRITICAL,
		"error":     ERROR,
		"warning":   WARNING,
		"notice":    NOTICE,
		"info":      INFO,
		"debug":     DEBUG,
	}

	level2LogursLevel = map[uint8]logrus.Level{
		EMERGENCY: logrus.PanicLevel,
		ALERT:     logrus.FatalLevel,
		CRITICAL:  logrus.FatalLevel,
		ERROR:     logrus.ErrorLevel,
		WARNING:   logrus.WarnLevel,
		NOTICE:    logrus.InfoLevel,
		INFO:      logrus.InfoLevel,
		DEBUG:     logrus.DebugLevel,
	}
)

type Log struct {
	Dir       string // 默认记录日志所在文件夹
	Threshold uint8  // 默认日志记录过滤级别
	Name      string // 日志名
	Fields
	fileLock *sync.Mutex
	logger   *logrus.Logger
}

type Fields struct {
	Pos string
}

var DLog = New("logs", ERROR, "") // 默认日志实例

func New(dir string, threshold uint8, name string) *Log {
	return &Log{
		Dir:       dir,
		Threshold: threshold,
		Name:      name,
		fileLock:  new(sync.Mutex),
		logger:    logrus.New(),
	}
}

func (l *Log) log(errLevel uint8, isShow bool, args ...interface{}) (err error) {
	if errLevel > l.Threshold {
		return
	}

	l.fileLock.Lock()

	level := levelMap[errLevel]
	_, fileName, line, _ := runtime.Caller(2)
	l.Fields.Pos = fileName + " : " + strconv.Itoa(line)

	logName := l.Name + "_" + level + "_" + strings.Split(time.Now().String(), " ")[0] + ".log"
	if _, err := os.Stat(l.Dir); os.IsNotExist(err) {
		os.MkdirAll(l.Dir, 0744)
	}
	logFileName := l.Dir + "/" + logName

	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer func() {
		file.Close()
		l.fileLock.Unlock()
	}()
	if err != nil {
		fmt.Println("get file descriptor error, ", err)
		return err
	}

	l.logger.Formatter = &logrus.JSONFormatter{}
	l.logger.Out = file
	l.logger.Level = level2LogursLevel[l.Threshold]

	var entry *logrus.Entry
	if isShow == true {
		entry = l.logger.WithFields(logrus.Fields{
			"pos": l.Fields.Pos,
		})
	} else {
		entry = l.logger.WithFields(logrus.Fields{})
	}
	switch errLevel {
	case EMERGENCY:
		entry.Panic(args...)
	case ALERT, CRITICAL:
		entry.Fatal(args...)
	case ERROR:
		entry.Error(args...)
	case WARNING:
		entry.Warn(args...)
	case NOTICE, INFO:
		entry.Info(args...)
	case DEBUG:
		entry.Debug(args...)
	}

	return
}

/**
 * check error and write log
 * @return: true: error, false: not error
 */
func (l *Log) CheckErr(err error, errLevel uint8, sendEmail bool) (isErr bool) {
	if err != nil {
		isErr = true
		fmt.Println(err)
		//if sendEmail {
		//	email.WriteEmail(err.Error())
		//}
		l.log(errLevel, true, err.Error())
	}
	return
}

func (l *Log) Emergency(args ...interface{}) {
	l.log(EMERGENCY, true, args)
}

func (l *Log) Alert(args ...interface{}) {
	l.log(ALERT, true, args)
}

func (l *Log) Critical(args ...interface{}) {
	l.log(CRITICAL, true, args)
}

func (l *Log) Error(args ...interface{}) {
	l.log(ERROR, true, args)
}

func (l *Log) Warning(args ...interface{}) {
	l.log(WARNING, true, args)
}

func (l *Log) Notice(args ...interface{}) {
	l.log(NOTICE, true, args...)
}

func (l *Log) Info(args ...interface{}) {
	l.log(INFO, true, args...)
}

func (l *Log) Debug(args ...interface{}) {
	l.log(DEBUG, true, args...)
}

//func (l *Log) DebugData(args ...interface{}) {
//	dirTmp := l.Dir
//	nameTmp := l.Name
//	l.Dir = "logs/debug"
//	l.Name = args[0].(string)
//	info := args[1].(string)
//	l.log(ERROR, false, info+"============================================================="+fmt.Sprint(args[0])+":START=============================================================", util.Now())
//	l.log(ERROR, true, args)
//	l.log(ERROR, false, info+"=============================================================="+fmt.Sprint(args[0])+":END==============================================================", util.Now())
//
//	l.Dir = dirTmp
//	l.Name = nameTmp
//}
//
//func (l *Log) Trace(name interface{}, trace interface{}) {
//	l.log(INFO, false, "============================================================="+fmt.Sprint(name)+":"+fmt.Sprint(trace)+"=============================================================", util.Now())
//}