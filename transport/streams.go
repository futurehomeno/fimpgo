package transport

import (
	"encoding/json"
	"fmt"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)
//BufferedStream - Implements in memory buffered stream. Underlying  buffer is flushed by reaching its max size or based on interval , depends what comes first.
// Content of the buffer is flushed either into file or into sink channel.
//-----------------------------------------------------------------------
// Source ---> buffer -> sink ---> File
//                         |-----> Channel
//-----------------------------------------------------------------------
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
	compressor          *MsgCompressor
	ticker              *time.Ticker
	fileSinkDir         string
}

func (su *BufferedStream) SinkChannel() chan []byte {
	return su.sinkChannel
}

func NewBufferedStream(bufferSizeLimit int, bufferInterval time.Duration,compressBeforeFlush bool) *BufferedStream {
	su := &BufferedStream{bufferMaxSize: bufferSizeLimit, bufferInterval: bufferInterval,compressBeforeFlush: compressBeforeFlush}
	if su.compressBeforeFlush {
		su.compressor = NewMsgCompressor("","")
	}
	if su.bufferInterval != 0 {
		su.ticker = time.NewTicker(time.Second * su.bufferInterval)
		go func() {
			for _ = range su.ticker.C {
				su.FlushBuffer()
			}
		}()
	}

	return su
}

//SetSourceStream - Configures source channel and starts message processing. Internal loop can be aborted by closing channel
func (su *BufferedStream) SetSourceStream(msgCh fimpgo.MessageCh) {
	go func() {
		for msg := range msgCh {
			su.EnqueueMessage(msg.Topic,msg.Payload)
		}
	}()
}


//EnqueueMessage - must be used to enqueue new message into the stream
func (su *BufferedStream) EnqueueMessage(topic string, msg *fimpgo.FimpMessage) {
	topic = strings.ReplaceAll(topic,"pt:j1/mt:evt","")
	topic = strings.ReplaceAll(topic,"pt:j1/mt:cmd","")
	msg.Topic = topic
	if len(su.buffer) >= su.bufferMaxSize {
		su.FlushBuffer()
	}
	su.lock.Lock()
	su.buffer = append(su.buffer, *msg)
	su.lock.Unlock()
	log.Debugf("Msg queued. Buffer size=%d , maxSize = %d",len(su.buffer),su.bufferMaxSize)
}

func (su *BufferedStream) Size()int {
	return len(su.buffer)
}

func (su *BufferedStream) FlushBuffer() {
	//var payload []byte
	su.lock.Lock()
	su.serializeBuffer()
	su.buffer = su.buffer[:0] // setting size to 0 without allocation
	su.lock.Unlock()
}

func (su *BufferedStream) ConfigureFileSink(filePrefix,path string) {
	su.flushToFile = true
	su.filePrefix = filePrefix
	su.fileSinkDir = path
}

func (su *BufferedStream) ConfigureChanelSink(size int) chan []byte {
	su.flushToSinkChannel = true
	su.sinkChannel = make( chan []byte,size)
	return su.sinkChannel
}

func (su *BufferedStream) serializeBuffer() error {
	log.Debugf("Serializing stream buffer.size=%d , maxSize = %d",len(su.buffer),su.bufferMaxSize)
	for i , _ := range su.buffer {
		if su.buffer[i].ValueType == fimpgo.VTypeObject {
			su.buffer[i].GetObjectValue(&su.buffer[i].Value)
		}
	}
	bPayload, err := json.Marshal(su.buffer)
	if err != nil {
		return err
	}
	if su.compressBeforeFlush {
		bPayload,err = su.compressor.CompressBinMsg(bPayload)
		if err != nil {
			log.Error("Compressor error : ",err.Error())
			return err
		}
		log.Debug("Compressed")
	}

	if su.flushToFile {
		fextension := "json"
		if su.compressBeforeFlush {
			fextension = "gz"
		}
		fname := fmt.Sprintf("%s/%s_%s.%s",su.fileSinkDir,su.filePrefix,time.Now().Format(time.RFC3339),fextension)
		err := ioutil.WriteFile(fname,bPayload,0777)
		if err != nil {
			return err
		}
		log.Debug("Serialized to file")
	}

	if su.flushToSinkChannel {
		su.sinkChannel <- bPayload
	}

	return err
}