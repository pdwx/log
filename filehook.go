package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"path/filepath"

	"strings"

	"github.com/Sirupsen/logrus"
)

// FileHook to write logs to a log file
type FileHook struct {
	sync.Mutex

	// File name
	Filename   string `json:"filename"`
	fileWriter *os.File

	// Rotate at line
	MaxLines         int `json:"maxlines"`
	maxLinesCurLines int

	// Rotate at size
	MaxSize        int `json:"maxsize"`
	maxSizeCurSize int
	// Rotate daily
	Daily         bool  `json:"daily"`
	MaxDays       int64 `json:"maxdays"`
	dailyOpenDate int

	Rotate bool `json:"rotate"`

	Level logrus.Level `json:"level"`

	Perm os.FileMode `json:"perm"`

	fileNameOnly, suffix string // like "project.log", project is fileNameOnly and .log is suffix
}

// NewFileHook create a hook with json config
// josnCfg like:
func NewFileHook(jsonCfg string) (*FileHook, error) {
	hook := &FileHook{
		Filename: "",
		MaxLines: 1000000,
		MaxSize:  1 << 28, //256 MB
		Daily:    true,
		MaxDays:  7,
		Rotate:   true,
		Level:    logrus.InfoLevel,
		Perm:     0660,
	}
	err := json.Unmarshal([]byte(jsonCfg), hook)
	if err != nil {
		return nil, err
	}
	hook.suffix = filepath.Ext(hook.Filename)
	hook.fileNameOnly = strings.TrimSuffix(hook.Filename, hook.suffix)
	if hook.suffix == "" {
		hook.suffix = ".log"
	}
	err = hook.startLogger()
	return hook, err
}

// Levels return FileHook supports log levels
func (hook *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}

}

// Fire is write log to file
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	if entry.Level > hook.Level {
		return nil
	}
	if hook.Rotate {
		if hook.needRotate() {
			hook.Lock()
			if hook.needRotate() {
				if err = hook.doRotate(time.Now()); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q):%s\n", hook.Filename, err)
				}
			}
			hook.Unlock()
		}
	}
	hook.Lock()
	_, err = hook.fileWriter.Write([]byte(line))
	if err == nil {
		hook.maxLinesCurLines++
		hook.maxSizeCurSize += len(line)
	}
	hook.Unlock()
	return err
}

// start file logger. create log file and set to locker-inside file writer.
func (hook *FileHook) startLogger() error {
	file, err := hook.createLogFile()
	if err != nil {
		return err
	}
	if hook.fileWriter != nil {
		hook.fileWriter.Close()
	}
	hook.fileWriter = file
	return hook.initFd()
}

func (hook *FileHook) needRotate() bool {
	return (hook.MaxLines > 0 && hook.maxLinesCurLines >= hook.MaxLines) ||
		(hook.MaxSize > 0 && hook.maxSizeCurSize >= hook.MaxSize) ||
		(hook.Daily && time.Now().Day() != hook.dailyOpenDate)
}

func (hook *FileHook) createLogFile() (*os.File, error) {
	// Open the log file
	fd, err := os.OpenFile(hook.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, hook.Perm)
	return fd, err
}

func (hook *FileHook) initFd() error {
	fd := hook.fileWriter
	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	hook.maxSizeCurSize = int(fInfo.Size())
	hook.dailyOpenDate = time.Now().Day()
	hook.maxLinesCurLines = 0
	if fInfo.Size() > 0 {
		count, err := hook.lines()
		if err != nil {
			return err
		}
		hook.maxLinesCurLines = count
	}
	return nil
}

func (hook *FileHook) lines() (int, error) {
	fd, err := os.Open(hook.Filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func (hook *FileHook) doRotate(logTime time.Time) error {
	_, err := os.Lstat(hook.Filename)
	if err != nil {
		return err
	}
	// file exists
	// Find the next available number
	num := 1
	fName := ""
	if hook.MaxLines > 0 || hook.MaxSize > 0 {
		for ; err == nil && num <= 999; num++ {
			fName = hook.fileNameOnly + fmt.Sprintf(".%s.%03d%s", logTime.Format("2006-01-02"), num, hook.suffix)
			_, err = os.Lstat(fName)
		}
	} else {
		fName = fmt.Sprintf("%s.%s%s", hook.fileNameOnly, logTime.Format("2006-01-02"), hook.suffix)
		_, err = os.Lstat(fName)
	}
	// return error if the last file checked still existed
	if err == nil {
		return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", hook.Filename)
	}

	// close fileWriter before rename
	hook.fileWriter.Close()

	// Rename the file to its new found name
	// even if occurs error,we MUST guarantee to  restart new logger
	renameErr := os.Rename(hook.Filename, fName)
	// re-start logger
	startLoggerErr := hook.startLogger()
	go hook.deleteOldLog()

	if startLoggerErr != nil {
		return fmt.Errorf("Rotate StartLogger: %s\n", startLoggerErr)
	}
	if renameErr != nil {
		return fmt.Errorf("Rotate: %s\n", renameErr)
	}
	return nil
}

func (hook *FileHook) deleteOldLog() {
	dir := filepath.Dir(hook.Filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*hook.MaxDays) {
			if strings.HasPrefix(filepath.Base(path), hook.fileNameOnly) &&
				strings.HasSuffix(filepath.Base(path), hook.suffix) {
				os.Remove(path)
			}
		}
		return
	})
}
