package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	log           *logrus.Logger
	AppLog        *logrus.Entry
	InitLog       *logrus.Entry
	CfgLog        *logrus.Entry
	ContextLog    *logrus.Entry
	HandlerLog    *logrus.Entry
	SessionQosLog *logrus.Entry
	ServParamLog  *logrus.Entry
	UtilLog       *logrus.Entry
	HttpLog       *logrus.Entry
	ConsumerLog   *logrus.Entry
	ProducerLog   *logrus.Entry
	GinLog        *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	AppLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "CFG"})
	ContextLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "Context"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "HDLR"})
	SessionQosLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "SessionQos"})
	ServParamLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "ServParam"})
	UtilLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "Util"})
	HttpLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "HTTP"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "Consumer"})
	ProducerLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "Producer"})
	GinLog = log.WithFields(logrus.Fields{"component": "NEF", "category": "GIN"})
}

func LogFileHook(logNfPath string, log5gcPath string) error {
	if fullPath, err := logger_util.CreateFree5gcLogFile(log5gcPath); err == nil {
		if fullPath != "" {
			free5gcLogHook, hookErr := logger_util.NewFileHook(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
			if hookErr != nil {
				return hookErr
			}
			log.Hooks.Add(free5gcLogHook)
		}
	} else {
		return err
	}

	if fullPath, err := logger_util.CreateNfLogFile(logNfPath, "nef.log"); err == nil {
		selfLogHook, hookErr := logger_util.NewFileHook(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
		if hookErr != nil {
			return hookErr
		}
		log.Hooks.Add(selfLogHook)
	} else {
		return err
	}

	return nil
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(enable bool) {
	log.SetReportCaller(enable)
}
