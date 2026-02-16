package transport

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

// BufferedStream - Implements in memory buffered stream. Underlying  buffer is flushed by reaching its max size or based on interval , depends what comes first.
// Content of the buffer is flushed either into file or into sink channel.
// -----------------------------------------------------------------------
// Source ---> buffer -> sink ---> File
//
//	|-----> Channel
//
// -----------------------------------------------------------------------
type BufferedStream struct {
	bufferMaxSize       int // number of messages
	bufferInterval      time.Duration
	buffer              []fimpgo.FimpMessage
	lock                sync.Mutex
	flushToFile         bool
	compressBeforeFlush bool
	filePrefix          string
	flushToSinkChannel  bool
	sinkChannel         chan []byte
	compressor          *fimpgo.MsgCompressor
	fileSinkDir         string
	close               chan struct{}
}

func (su *BufferedStream) SinkChannel() chan []byte {
	return su.sinkChannel
}

// returns nil on error
func NewBufferedStream(bufferSizeLimit int, bufferInterval time.Duration, compressBeforeFlush bool) *BufferedStream {
	if bufferInterval == 0 || bufferSizeLimit == 0 {
		log.Warn("[fimpgo] Invalid arguments")
		return nil
	}

	su := &BufferedStream{bufferMaxSize: bufferSizeLimit,
		bufferInterval:      bufferInterval,
		compressBeforeFlush: compressBeforeFlush,
		close:               make(chan struct{}, 1),
	}

	if su.compressBeforeFlush {
		su.compressor = fimpgo.NewMsgCompressor("", "")
	}

	go func() {
		ticker := time.NewTicker(su.bufferInterval)

		defer func() {
			ticker.Stop()
			ticker = nil
		}()

		for {
			select {
			case <-su.close:
				return
			case <-ticker.C:
				su.FlushBuffer()
			}
		}
	}()

	return su
}

// SetSourceStream - Configures source channel and starts message processing. Internal loop can be aborted by closing channel
func (su *BufferedStream) SetSourceStream(msgCh fimpgo.MessageCh) {
	go func() {
		for msg := range msgCh {
			su.EnqueueMessage(msg.Topic, msg.Payload)
		}
	}()
}

// EnqueueMessage - must be used to enqueue new message into the stream
func (su *BufferedStream) EnqueueMessage(topic string, msg *fimpgo.FimpMessage) {
	topic = strings.ReplaceAll(topic, "pt:j1/mt:evt", "")
	topic = strings.ReplaceAll(topic, "pt:j1/mt:cmd", "")
	msg.Topic = topic

	su.lock.Lock()
	su.buffer = append(su.buffer, *msg)
	shouldFlush := len(su.buffer) >= su.bufferMaxSize
	bufLen := len(su.buffer)
	su.lock.Unlock()

	if shouldFlush {
		su.FlushBuffer()
	}

	log.Tracef("[fimpgo] Msg queued len(buffer)=%d maxSize=%d", bufLen, su.bufferMaxSize)
}

func (su *BufferedStream) Size() int {
	su.lock.Lock()
	ret := len(su.buffer)
	su.lock.Unlock()
	return ret
}

// buffer is cleared even on serialization failure
func (su *BufferedStream) FlushBuffer() {
	su.lock.Lock()
	defer su.lock.Unlock()

	if len(su.buffer) == 0 {
		return
	}

	if err := su.serializeBuffer(); err != nil {
		log.Warnf("[fimpgo] Serialize buffer err: %v", err)
	}

	su.buffer = su.buffer[:0] // setting size to 0 without allocation
}

func (su *BufferedStream) ConfigureFileSink(filePrefix, path string) {
	su.flushToFile = true
	su.filePrefix = filePrefix
	su.fileSinkDir = path
}

func (su *BufferedStream) ConfigureChanelSink(size int) chan []byte {
	su.flushToSinkChannel = true
	su.sinkChannel = make(chan []byte, size)
	return su.sinkChannel
}

func (su *BufferedStream) serializeBuffer() error {
	for i := range su.buffer {
		if su.buffer[i].ValueType == fimpgo.VTypeObject {
			if err := su.buffer[i].GetObjectValue(&su.buffer[i].Value); err != nil {
				return err
			}
		}
	}

	bPayload, err := json.Marshal(su.buffer)
	if err != nil {
		return err
	}

	if su.compressBeforeFlush {
		bPayload, err = su.compressor.CompressBinMsg(bPayload)
		if err != nil {
			return err
		}
	}

	if su.flushToFile {
		fextension := "json"
		if su.compressBeforeFlush {
			fextension = "gz"
		}
		fname := fmt.Sprintf("%s/%s_%s.%s", su.fileSinkDir, su.filePrefix, time.Now().Format(time.RFC3339), fextension)
		err := os.WriteFile(fname, bPayload, 0644) //nolint:gosec
		if err != nil {
			return err
		}
	}

	if su.flushToSinkChannel {
		su.sinkChannel <- bPayload
	}

	return err
}

func (su *BufferedStream) Close() {
	select {
	case su.close <- struct{}{}:
	default:
	}
}
