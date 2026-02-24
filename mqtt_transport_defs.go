package fimpgo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultMainQueueSize = 100
)

type MessageCh chan *Message

type MqttConnectionConfigs struct {
	ServerURI           string
	ClientID            string
	Username            string
	Password            string
	CleanSession        bool
	SubQos              byte
	PubQos              byte
	GlobalTopicPrefix   string // Should be set for communicating one single hub via cloud
	StartFailRetryCount int
	CertDir             string // full path to directory where all certificates are stored. Cert dir should contains all CA root certificates .
	PrivateKeyFileName  string //
	CertFileName        string //
	ReceiveChTimeout    uint32
	IsAws               bool // Should be set to true if cloud broker is AwS IoT platform .
	MainQueueSize       int
}

type Message struct {
	Topic      string
	Addr       *Address
	Payload    *FimpMessage
	RawPayload []byte
}

type FimpFilter struct {
	Topic     string
	Service   string
	Interface string
}

func connectionNotifStr(evt MQTT.ConnectionNotificationType) string {
	switch evt {
	case MQTT.ConnectionNotificationTypeConnected:
		return "connected"
	case MQTT.ConnectionNotificationTypeConnecting:
		return "connecting"
	case MQTT.ConnectionNotificationTypeFailed:
		return "connection_failed"
	case MQTT.ConnectionNotificationTypeLost:
		return "connection_lost"
	case MQTT.ConnectionNotificationTypeBroker:
		return "broker"
	case MQTT.ConnectionNotificationTypeBrokerFailed:
		return "broker_failed"
	}

	return fmt.Sprintf("unknown(%d)", evt)
}

type FilterFunc func(topic string, addr *Address, iotMsg *FimpMessage) bool

type MqttTransport struct {
	client         MQTT.Client
	msgHandler     MessageHandler
	subQos         byte
	pubQos         byte
	subs           map[string]byte
	subChannels    map[string]MessageCh
	subFilters     map[string]FimpFilter
	subFilterFuncs map[string]FilterFunc

	connState ConnStateT
	incMsgsWg sync.WaitGroup // WaitGroup for incoming message handler

	mainQueue                   chan MQTT.Message
	mainQueueOverflowCnt        atomic.Uint32
	errorHandler                func(err error)
	connectionLostCustomHandler func(client MQTT.Client, err error)

	globalTopicPrefixLock sync.RWMutex
	_globalTopicPrefix    string
	defaultSourceLock     sync.RWMutex
	defaultSource         string
	startFailRetryCount   int
	certDir               string
	receiveChTimeout      atomic.Uint32
	syncPublishTimeout    time.Duration
	channelRegLock        sync.Mutex // channel registration
	subscribeLock         sync.Mutex // subscribe
	compressor            *MsgCompressor
}
