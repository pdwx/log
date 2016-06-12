package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"path"

	"math"

	"encoding/json"

	"github.com/Sirupsen/logrus"
)

var log *logrus.Logger
var enableCallFunc bool = true

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

type FileLogConfig struct {
	Filename string       `json:"filename"`
	MaxLines int          `json:"maxlines"`
	MaxSize  int          `json:"maxsize"`
	Daily    bool         `json:"daily"`
	MaxDays  int64        `json:"maxdays"`
	Rotate   bool         `json:"rotate"`
	Level    logrus.Level `json:"level"`
}

func init() {
	log = &logrus.Logger{}
	log.Out = os.Stdout
	log.Hooks = make(logrus.LevelHooks)
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{
		DisableColors: true,
	}
}

func DefaultFileLogConfig() *FileLogConfig {
	flc := &FileLogConfig{}
	os.Mkdir("./log", 0777)
	pname := path.Base(os.Args[0])
	flc.Filename = fmt.Sprintf("./log/%s-%d-%s.log", pname, os.Getpid(), time.Now().Format("2013-01-02"))
	flc.MaxLines = math.MaxInt32
	flc.MaxSize = 1024 * 1024 * 512 //512MB
	flc.Daily = true
	flc.MaxDays = 90
	flc.Rotate = true
	flc.Level = logrus.DebugLevel
	return flc
}

func DefaultFileHook(fileConfig *FileLogConfig) (logrus.Hook, error) {
	if fileConfig == nil {
		fileConfig = DefaultFileLogConfig()
	}
	b, err := json.Marshal(fileConfig)
	if err != nil {
		return nil, err
	}
	hook, err := NewFileHook(string(b))
	if err != nil {
		return nil, err
	}
	return hook, err
}

func AddHook(hook logrus.Hook) {
	log.Hooks.Add(hook)
}

func SetFormatter(formatter logrus.Formatter) {
	log.Formatter = formatter
}

func SetLogLevel(level logrus.Level) {
	log.Level = level
}

func DisableCallFunc() {
	enableCallFunc = false
}

func Debugf(format string, args ...interface{}) {
	if log.Level >= logrus.DebugLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Debugf(format, args...)
		} else {
			logrus.NewEntry(log).Debugf(format, args...)
		}
	}
}

func Infof(format string, args ...interface{}) {
	if log.Level >= logrus.InfoLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Infof(format, args...)
		} else {
			logrus.NewEntry(log).Infof(format, args...)
		}
	}
}

func Printf(format string, args ...interface{}) {
	if enableCallFunc {
		logrus.NewEntry(log).WithField("func", getFuncCall()).Printf(format, args...)
	} else {
		logrus.NewEntry(log).Printf(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Warnf(format, args...)
		} else {
			logrus.NewEntry(log).Warnf(format, args...)
		}
	}
}

func Warningf(format string, args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Warnf(format, args...)
		} else {
			logrus.NewEntry(log).Warnf(format, args...)
		}
	}
}

func Errorf(format string, args ...interface{}) {
	if log.Level >= logrus.ErrorLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Errorf(format, args...)
		} else {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Errorf(format, args...)
		}
	}
}

func Fatalf(format string, args ...interface{}) {
	if log.Level >= logrus.FatalLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Fatalf(format, args...)
		} else {
			logrus.NewEntry(log).Fatalf(format, args...)
		}
	}
	os.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	if log.Level >= logrus.PanicLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Panicf(format, args...)
		} else {
			logrus.NewEntry(log).Panicf(format, args...)
		}
	}
}

type PackageLog struct {
	Level logrus.Level
}

func (plog *PackageLog) Debugf(format string, args ...interface{}) {
	if plog.Level >= logrus.DebugLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Infof(format, args...)
		} else {
			logrus.NewEntry(log).Infof(format, args...)
		}
	}
}

func (plog *PackageLog) Infof(format string, args ...interface{}) {
	if plog.Level >= logrus.InfoLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Infof(format, args...)
		} else {
			logrus.NewEntry(log).Infof(format, args...)
		}
	}
}

func (plog *PackageLog) Printf(format string, args ...interface{}) {
	if enableCallFunc {
		logrus.NewEntry(log).WithField("func", getFuncCall()).Printf(format, args...)
	} else {
		logrus.NewEntry(log).Printf(format, args...)
	}
}

func (plog *PackageLog) Warnf(format string, args ...interface{}) {
	if plog.Level >= logrus.WarnLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Warnf(format, args...)
		} else {
			logrus.NewEntry(log).Warnf(format, args...)
		}
	}
}

func (plog *PackageLog) Warningf(format string, args ...interface{}) {
	if plog.Level >= logrus.WarnLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Warnf(format, args...)
		} else {
			logrus.NewEntry(log).Warnf(format, args...)
		}
	}
}

func (plog *PackageLog) Errorf(format string, args ...interface{}) {
	if plog.Level >= logrus.ErrorLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Errorf(format, args...)
		} else {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Errorf(format, args...)
		}
	}
}

func (plog *PackageLog) Fatalf(format string, args ...interface{}) {
	if plog.Level >= logrus.FatalLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Fatalf(format, args...)
		} else {
			logrus.NewEntry(log).Fatalf(format, args...)
		}
	}
	os.Exit(1)
}

func (plog *PackageLog) Panicf(format string, args ...interface{}) {
	if plog.Level >= logrus.PanicLevel {
		if enableCallFunc {
			logrus.NewEntry(log).WithField("func", getFuncCall()).Panicf(format, args...)
		} else {
			logrus.NewEntry(log).Panicf(format, args...)
		}
	}
}

func getFuncCall() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	return fmt.Sprintf("%s:%d", filename, line)
}
