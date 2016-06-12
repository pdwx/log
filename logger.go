package log

import (
	"fmt"
	"os"
	"time"

	"path"

	"math"

	"encoding/json"

	"github.com/Sirupsen/logrus"
)

var log *logrus.Logger

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
	log.Level = logrus.InfoLevel
	log.Formatter = &TextFormatter{
		FullTimestamp: true,
		FuncCallDepth: 4,
	}
}

func DefaultFileLogConfig() *FileLogConfig {
	flc := &FileLogConfig{}
	os.Mkdir("log", 0777)
	pname := path.Base(os.Args[0])
	flc.Filename = fmt.Sprintf("/log/%s-%d-%s.log", pname, os.Getpid(), time.Now().Format("2013-01-02"))
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

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.NewEntry(log).WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.NewEntry(log).WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return logrus.NewEntry(log).WithError(err)
}

func Debugf(format string, args ...interface{}) {
	if log.Level >= logrus.DebugLevel {
		logrus.NewEntry(log).Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if log.Level >= logrus.InfoLevel {
		logrus.NewEntry(log).Infof(format, args...)
	}
}

func Printf(format string, args ...interface{}) {
	logrus.NewEntry(log).Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		logrus.NewEntry(log).Warnf(format, args...)
	}
}

func Warningf(format string, args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		logrus.NewEntry(log).Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if log.Level >= logrus.ErrorLevel {
		logrus.NewEntry(log).Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if log.Level >= logrus.FatalLevel {
		logrus.NewEntry(log).Fatalf(format, args...)
	}
	os.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	if log.Level >= logrus.PanicLevel {
		logrus.NewEntry(log).Panicf(format, args...)
	}
}

func Debug(args ...interface{}) {
	if log.Level >= logrus.DebugLevel {
		logrus.NewEntry(log).Debug(args...)
	}
}

func Info(args ...interface{}) {
	if log.Level >= logrus.InfoLevel {
		logrus.NewEntry(log).Info(args...)
	}
}

func Print(args ...interface{}) {
	logrus.NewEntry(log).Info(args...)
}

func Warn(args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		logrus.NewEntry(log).Warn(args...)
	}
}

func Warning(args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		logrus.NewEntry(log).Warn(args...)
	}
}

func Error(args ...interface{}) {
	if log.Level >= logrus.ErrorLevel {
		logrus.NewEntry(log).Error(args...)
	}
}

func Fatal(args ...interface{}) {
	if log.Level >= logrus.FatalLevel {
		logrus.NewEntry(log).Fatal(args...)
	}
	os.Exit(1)
}

func Panic(args ...interface{}) {
	if log.Level >= logrus.PanicLevel {
		logrus.NewEntry(log).Panic(args...)
	}
}

func Debugln(args ...interface{}) {
	if log.Level >= logrus.DebugLevel {
		logrus.NewEntry(log).Debugln(args...)
	}
}

func Infoln(args ...interface{}) {
	if log.Level >= logrus.InfoLevel {
		logrus.NewEntry(log).Infoln(args...)
	}
}

func Println(args ...interface{}) {
	logrus.NewEntry(log).Println(args...)
}

func Warnln(args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		logrus.NewEntry(log).Warnln(args...)
	}
}

func Warningln(args ...interface{}) {
	if log.Level >= logrus.WarnLevel {
		logrus.NewEntry(log).Warnln(args...)
	}
}

func Errorln(args ...interface{}) {
	if log.Level >= logrus.ErrorLevel {
		logrus.NewEntry(log).Errorln(args...)
	}
}

func Fatalln(args ...interface{}) {
	if log.Level >= logrus.FatalLevel {
		logrus.NewEntry(log).Fatalln(args...)
	}
	os.Exit(1)
}

func Panicln(args ...interface{}) {
	if log.Level >= logrus.PanicLevel {
		logrus.NewEntry(log).Panicln(args...)
	}
}
