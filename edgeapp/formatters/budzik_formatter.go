package formatters

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type BudzikFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *BudzikFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	return fmt.Appendf(nil, "%s %s %s\n", timestamp, f.LevelDesc[entry.Level], entry.Message), nil
}

func NewBudzikFormatter() *BudzikFormatter {
	lvlDesc := []string{"PANIC", "FATAL", "E", "W", "I", "D", "T", "?"}
	return &BudzikFormatter{TimestampFormat: "01-02 15:04:05", LevelDesc: lvlDesc}
}
