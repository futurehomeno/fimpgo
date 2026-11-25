package formatters

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudzikFormatter_Format(t *testing.T) {
	ts := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		name  string
		entry *log.Entry
		want  string
	}{
		{"info", &log.Entry{Time: ts, Level: log.InfoLevel, Message: "hi"}, "01-02 03:04:05 I hi\n"},
		{"error with field", &log.Entry{Time: ts, Level: log.ErrorLevel, Message: "bad", Data: log.Fields{"k": 1}}, "01-02 03:04:05 E bad k=1\n"},
		{"trace", &log.Entry{Time: ts, Level: log.TraceLevel, Message: "t"}, "01-02 03:04:05 T t\n"},
		{"unknown level", &log.Entry{Time: ts, Level: log.Level(42), Message: "x"}, "01-02 03:04:05 D x\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBudzikFormatter().Format(tt.entry)
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(b))
		})
	}
}
