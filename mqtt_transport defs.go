package fimpgo

import (
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
	ReceiveChTimeout    int
	IsAws               bool // Should be set to true if cloud broker is AwS IoT platform .
	MainQueueSize       int

	connectionLostHandler MQTT.ConnectionLostHandler
	errorHandler          func(err error)
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

	//doneSignal           chan struct{}
	mainQueue            chan MQTT.Message
	mainQueueOverflowCnt atomic.Uint32
	errorHandler         func(err error)

	globalTopicPrefixLock sync.RWMutex
	_globalTopicPrefix    string
	defaultSourceLock     sync.RWMutex
	defaultSource         string
	startFailRetryCount   int
	certDir               string
	//mqttOptions          *MQTT.ClientOptions
	receiveChTimeout   int
	syncPublishTimeout time.Duration
	channelRegLock     sync.Mutex // channel registration
	subscribeLock      sync.Mutex // subscribe
	compressor         *MsgCompressor
}
