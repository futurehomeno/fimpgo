package formatters

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type BudzikFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *BudzikFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)

	level := "D"

	if int(entry.Level) >= 0 && int(entry.Level) < len(f.LevelDesc) {
		level = f.LevelDesc[int(entry.Level)]
	}

	ret := fmt.Appendf(nil, "%s %s %s", timestamp, level, entry.Message)
	for k, v := range entry.Data {
		ret = fmt.Appendf(ret, " %s=%v", k, v)
	}

	ret = fmt.Appendf(ret, "\n")

	return ret, nil
}

func NewBudzikFormatter() *BudzikFormatter {
	lvlDesc := []string{"PANIC", "FATAL", "E", "W", "I", "D", "T", "?"}
	return &BudzikFormatter{TimestampFormat: "01-02 15:04:05", LevelDesc: lvlDesc}
}
