package log

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"

	"github.com/Sirupsen/logrus"
)

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	// Disable call funcCallDepth
	DisableFuncCallDepth bool
	FuncCallDepth        int
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = logrus.DefaultTimestampFormat
	}
	if !f.DisableFuncCallDepth {
		_, file, line, ok := runtime.Caller(f.FuncCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		data["func"] = fmt.Sprintf("%s:%d", filename, line)
	}
	data["time"] = entry.Time.Format(timestampFormat)
	data["msg"] = entry.Message
	data["level"] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
