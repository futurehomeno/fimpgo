package fimpgo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/futurehomeno/fimpgo/transport"
	"github.com/futurehomeno/fimpgo/utils"
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

	connectionLostHandler MQTT.ConnectionLostHandler
}

type Message struct {
	Topic   string
	Addr    *Address
	Payload *FimpMessage
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

	globalTopicPrefixMux sync.RWMutex
	globalTopicPrefix    string
	defaultSourceLock    sync.RWMutex
	defaultSource        string
	startFailRetryCount  int
	certDir              string
	mqttOptions          *MQTT.ClientOptions
	receiveChTimeout     int
	syncPublishTimeout   time.Duration
	channelRegMux        sync.Mutex
	subMutex             sync.Mutex
	compressor          *transport.MsgCompressor
}

func (mh *MqttTransport) SetReceiveChTimeout(receiveChTimeout int) {
	mh.receiveChTimeout = receiveChTimeout
}

func (mh *MqttTransport) SetCertDir(certDir string) {
	mh.certDir = certDir
}

type MessageHandler func(topic string, addr *Address, iotMsg *FimpMessage, rawPayload []byte)

// NewMqttTransport constructor. serverUri="tcp://localhost:1883"
func NewMqttTransport(serverURI, clientID, username, password string, cleanSession bool, subQos byte, pubQos byte) *MqttTransport {
	mh := MqttTransport{}
	mh.mqttOptions = MQTT.NewClientOptions().AddBroker(serverURI)
	mh.mqttOptions.SetClientID(clientID)
	mh.mqttOptions.SetUsername(username)
	mh.mqttOptions.SetPassword(password)
	mh.mqttOptions.SetDefaultPublishHandler(mh.onMessage)
	mh.mqttOptions.SetCleanSession(cleanSession)
	mh.mqttOptions.SetAutoReconnect(true)
	mh.mqttOptions.SetConnectionLostHandler(mh.onConnectionLost)
	mh.mqttOptions.SetOnConnectHandler(mh.onConnect)
	mh.mqttOptions.SetWriteTimeout(time.Second*30)
	//create and start a client using the above ClientOptions
	mh.client = MQTT.NewClient(mh.mqttOptions)
	mh.pubQos = pubQos
	mh.subQos = subQos
	mh.subs = make(map[string]byte)
	mh.subChannels = make(map[string]MessageCh)
	mh.subFilters = make(map[string]FimpFilter)
	mh.subFilterFuncs = make(map[string]FilterFunc)
	mh.startFailRetryCount = 10
	mh.receiveChTimeout = 10
	mh.syncPublishTimeout = time.Second * 5
	return &mh
}

func NewMqttTransportFromConnection(client MQTT.Client, subQos byte, pubQos byte) *MqttTransport {
	mh := MqttTransport{}
	mh.client = client
	mh.pubQos = pubQos
	mh.subQos = subQos
	mh.subs = make(map[string]byte)
	mh.subChannels = make(map[string]MessageCh)
	mh.subFilters = make(map[string]FimpFilter)
	mh.subFilterFuncs = make(map[string]FilterFunc)
	mh.startFailRetryCount = 10
	mh.receiveChTimeout = 10
	mh.syncPublishTimeout = time.Second * 5
	return &mh
}

func NewMqttTransportFromConfigs(configs MqttConnectionConfigs, options ...Option) *MqttTransport {

	applyDefaults(&configs)

	// apply extra options
	for _, o := range options {
		o.apply(&configs)
	}

	mh := MqttTransport{}
	mh.mqttOptions = MQTT.NewClientOptions().AddBroker(configs.ServerURI)
	mh.mqttOptions.SetClientID(configs.ClientID)
	mh.mqttOptions.SetUsername(configs.Username)
	mh.mqttOptions.SetPassword(configs.Password)
	mh.mqttOptions.SetDefaultPublishHandler(mh.onMessage)
	mh.mqttOptions.SetCleanSession(configs.CleanSession)
	mh.mqttOptions.SetAutoReconnect(true)
	mh.mqttOptions.SetConnectionLostHandler(configs.connectionLostHandler)
	mh.mqttOptions.SetOnConnectHandler(mh.onConnect)

	//create and start a client using the above ClientOptions
	mh.client = MQTT.NewClient(mh.mqttOptions)
	mh.pubQos = configs.PubQos
	mh.subQos = configs.SubQos
	mh.subs = make(map[string]byte)
	mh.subChannels = make(map[string]MessageCh)
	mh.subFilters = make(map[string]FimpFilter)
	mh.subFilterFuncs = make(map[string]FilterFunc)
	mh.startFailRetryCount = 10
	mh.receiveChTimeout = 10
	mh.syncPublishTimeout = time.Second * 5
	mh.certDir = configs.CertDir
	mh.globalTopicPrefix = configs.GlobalTopicPrefix
	if configs.StartFailRetryCount == 0 {
		mh.startFailRetryCount = 10
	} else {
		mh.startFailRetryCount = configs.StartFailRetryCount
	}
	if configs.ReceiveChTimeout == 0 {
		mh.receiveChTimeout = 10
	} else {
		mh.receiveChTimeout = configs.ReceiveChTimeout
	}

	if configs.PrivateKeyFileName != "" && configs.CertFileName != "" {
		err := mh.ConfigureTls(configs.PrivateKeyFileName, configs.CertFileName, configs.CertDir, configs.IsAws)
		if err != nil {
			log.Error("Certificate loading error :", err.Error())
		}
	}

	return &mh
}

func (mh *MqttTransport) SetGlobalTopicPrefix(prefix string) {
	mh.globalTopicPrefixMux.Lock()
	mh.globalTopicPrefix = prefix
	mh.globalTopicPrefixMux.Unlock()
}

func (mh *MqttTransport) getGlobalTopicPrefix() string {
	mh.globalTopicPrefixMux.RLock()
	defer mh.globalTopicPrefixMux.RUnlock()
	return mh.globalTopicPrefix
}

// SetDefaultSource safely sets default source name for all outgoing messages.
// Default source is used only if it was not set explicitly before.
func (mh *MqttTransport) SetDefaultSource(source string) {
	mh.defaultSourceLock.Lock()
	defer mh.defaultSourceLock.Unlock()

	mh.defaultSource = source
}

// ensureDefaultSource safely sets default source name for an outgoing message.
// Default source is used only if it was not set explicitly before.
func (mh *MqttTransport) ensureDefaultSource(message *FimpMessage) {
	if message.Source != "" {
		return
	}

	mh.defaultSourceLock.RLock()
	defer mh.defaultSourceLock.RUnlock()

	message.Source = mh.defaultSource
}

// SetStartAutoRetryCount Set number of retries transport will attempt on startup . Default value is 10
func (mh *MqttTransport) SetStartAutoRetryCount(count int) {
	mh.startFailRetryCount = count
}

// SetMessageHandler message handler setter
func (mh *MqttTransport) SetMessageHandler(msgHandler MessageHandler) {
	mh.msgHandler = msgHandler
}

// RegisterChannel should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *MqttTransport) RegisterChannel(channelId string, messageCh MessageCh) {
	mh.channelRegMux.Lock()
	mh.subChannels[channelId] = messageCh
	mh.channelRegMux.Unlock()
}

// UnregisterChannel should be used to unregister channel
func (mh *MqttTransport) UnregisterChannel(channelId string) {
	mh.channelRegMux.Lock()
	delete(mh.subChannels, channelId)
	delete(mh.subFilters, channelId)
	delete(mh.subFilterFuncs, channelId)
	mh.channelRegMux.Unlock()
}

// RegisterChannelWithFilter should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *MqttTransport) RegisterChannelWithFilter(channelId string, messageCh MessageCh, filter FimpFilter) {
	mh.channelRegMux.Lock()
	mh.subChannels[channelId] = messageCh
	mh.subFilters[channelId] = filter
	mh.channelRegMux.Unlock()
}

// RegisterChannelWithFilterFunc should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *MqttTransport) RegisterChannelWithFilterFunc(channelId string, messageCh MessageCh, filterFunc FilterFunc) {
	mh.channelRegMux.Lock()
	mh.subChannels[channelId] = messageCh
	mh.subFilterFuncs[channelId] = filterFunc
	mh.channelRegMux.Unlock()
}

func (mh *MqttTransport) Client() MQTT.Client {
	return mh.client
}

// Start , starts adapter async.
func (mh *MqttTransport) Start() error {
	log.Info("<MqttAd> Connecting to MQTT broker ")
	var err error
	var delay time.Duration
	for i := 1; i < mh.startFailRetryCount; i++ {
		if token := mh.client.Connect(); token.Wait() && token.Error() == nil {
			return nil
		} else {
			err = token.Error()
		}
		delay = time.Duration(i) * time.Duration(i)
		log.Infof("<MqttAd> Connection failed , retrying after %d sec.... ", delay)
		time.Sleep(delay * time.Second)
	}
	return err
}

// Stop stops adapter . Adapter can't be started again using Start . In order to start adapter it has to be re-initialized
func (mh *MqttTransport) Stop() {
	mh.client.Disconnect(250)
}

// Subscribe - subscribing for topic
func (mh *MqttTransport) Subscribe(topic string) error {
	if strings.TrimSpace(topic) == "" {
		return nil
	}

	mh.subMutex.Lock()
	defer mh.subMutex.Unlock()

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	topic = AddGlobalPrefixToTopic(mh.getGlobalTopicPrefix(), topic)
	log.Debug("<MqttAd> Subscribing to topic:", topic)
	token := mh.client.Subscribe(topic, mh.subQos, nil)
	isInTime := token.WaitTimeout(time.Second * 20)
	if token.Error() != nil {
		log.Error("<MqttAd> Can't subscribe. Error :", token.Error())
		return token.Error()
	} else if !isInTime {
		log.Error("<MqttAd> Subscribe operation timed out")
		return errors.New("subscribe timed out")
	}

	mh.subs[topic] = mh.subQos

	return nil
}

// Unsubscribe , unsubscribing from topic
func (mh *MqttTransport) Unsubscribe(topic string) error {
	mh.subMutex.Lock()
	defer mh.subMutex.Unlock()
	topic = AddGlobalPrefixToTopic(mh.getGlobalTopicPrefix(), topic)
	log.Debug("<MqttAd> Unsubscribing from topic:", topic)
	token := mh.client.Unsubscribe(topic)
	isInTime := token.WaitTimeout(time.Second * 20)
	if token.Error() != nil {
		return token.Error()
	} else if !isInTime {
		log.Error("<MqttAd> Unsubscribe operation timed out")
		return errors.New("unsubscribe timed out")
	}
	delete(mh.subs, topic)
	return nil
}

func (mh *MqttTransport) UnsubscribeAll() {
	var topics []string
	mh.subMutex.Lock()
	for i := range mh.subs {
		topics = append(topics, i)
	}
	mh.subMutex.Unlock()
	for _, t := range topics {
		if err := mh.Unsubscribe(t); err != nil {
			log.Error(errors.Wrap(err, "unsubscribing from topic"))
		}
	}
}

func (mh *MqttTransport) onConnectionLost(_ MQTT.Client, err error) {
	log.Errorf("<MqttAd> Connection lost with MQTT broker . Error : %v", err)
}

func (mh *MqttTransport) onConnect(_ MQTT.Client) {
	mh.subMutex.Lock()
	defer mh.subMutex.Unlock()

	log.Infof("<MqttAd> Connection established with MQTT broker .")
	if len(mh.subs) > 0 {
		if token := mh.client.SubscribeMultiple(mh.subs, nil); token.Wait() && token.Error() != nil {
			log.Error("Can't subscribe. Error :", token.Error())
		}
	}
}

//define a function for the default message handler
func (mh *MqttTransport) onMessage(_ MQTT.Client, msg MQTT.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("<MqttAd> onMessage CRASHED with error :", r)
		}
	}()
	log.Tracef("<MqttAd> New msg from TOPIC: %s", msg.Topic())
	var topic string
	if strings.TrimSpace(mh.globalTopicPrefix) != "" {
		_, topic = DetachGlobalPrefixFromTopic(msg.Topic())
	} else {
		topic = msg.Topic()
	}

	addr, err := NewAddressFromString(topic)
	if err != nil {
		log.Error("<MqttAd> Error processing address :", err)
		return
	}
	var fimpMsg *FimpMessage

	switch addr.PayloadType {
	case DefaultPayload:
		fimpMsg, err = NewMessageFromBytes(msg.Payload())
	case CompressedJsonPayload:
		fimpMsg,err = mh.compressor.DecompressFimpMsg(msg.Payload())
	default:
		// This means binary payload , for instance compressed message
	}

	if mh.msgHandler != nil {
		if err == nil  {
			mh.msgHandler(topic, addr, fimpMsg, msg.Payload())
		} else {
			log.Trace(string(msg.Payload()))
			log.Error("<MqttAd> Error processing payload :", err)
			return
		}
	}

	mh.channelRegMux.Lock()
	defer mh.channelRegMux.Unlock()

	for i := range mh.subChannels {
		if !mh.isChannelInterested(i, topic, addr, fimpMsg) {
			continue
		}
		var fmsg Message
		if addr.PayloadType == DefaultPayload {
			fmsg = Message{Topic: topic, Addr: addr, Payload: fimpMsg}
		}else {
			// message receiver should do decompressions
			fmsg = Message{Topic: topic, Addr: addr, RawPayload: msg.Payload()}
		}
		timer := time.NewTimer(time.Second * time.Duration(mh.receiveChTimeout))
		select {
		case mh.subChannels[i] <- &fmsg:
			timer.Stop()
			// send to channel
		case <-timer.C:
			log.Info("<MqttAd> Channel is not read for ", mh.receiveChTimeout)
		}
	}

}

// isChannelInterested validates if channel is interested in message. Filtering is executed against either static filters or filter function
func (mh *MqttTransport) isChannelInterested(chanName string, topic string, addr *Address, msg *FimpMessage) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Error("<MqttAd> Filter CRASHED with error :", r)
		}
	}()

	filterFunc, ok := mh.subFilterFuncs[chanName]
	if ok {
		return filterFunc(topic, addr, msg)
	}
	filter, ok := mh.subFilters[chanName]
	if !ok {
		// no filters has been set
		return true
	}
	if msg != nil {
		if utils.RouteIncludesTopic(filter.Topic, topic) &&
			(msg.Service == filter.Service || filter.Service == "*") &&
			(msg.Type == filter.Interface || filter.Interface == "*") {
			return true
		}
	}else {
		// It means binary payload , and message can't be parsed
		if utils.RouteIncludesTopic(filter.Topic, topic) {
			return true
		}
	}

	return false
}

// Publish publishes message to FIMP address
func (mh *MqttTransport) Publish(addr *Address, fimpMsg *FimpMessage) error {
	mh.ensureDefaultSource(fimpMsg)

	var bytm []byte
	var err error
	switch addr.PayloadType {
	case DefaultPayload:
		bytm, err = fimpMsg.SerializeToJson()
	case CompressedJsonPayload:
		bytm, err = mh.compressor.CompressFimpMsg(fimpMsg)

	}
	topic := addr.Serialize()
	if strings.TrimSpace(mh.globalTopicPrefix) != "" {
		topic = AddGlobalPrefixToTopic(mh.getGlobalTopicPrefix(), topic)
	}
	if err == nil {
		log.Trace("<MqttAd> Publishing msg to topic:", topic)
		mh.client.Publish(topic, mh.pubQos, false, bytm)
		return nil
	}
	return err
}

// PublishToTopic publishes iotMsg to string topic
func (mh *MqttTransport) PublishToTopic(topic string, fimpMsg *FimpMessage) error {
	mh.ensureDefaultSource(fimpMsg)

	byteMessage, err := fimpMsg.SerializeToJson()
	if err != nil {
		return err
	}

	if strings.TrimSpace(mh.globalTopicPrefix) != "" {
		topic = AddGlobalPrefixToTopic(mh.getGlobalTopicPrefix(), topic)
	}

	log.Trace("<MqttAd> Publishing msg to topic:", topic)
	return mh.client.Publish(topic, mh.pubQos, false, byteMessage).Error()
}

// RespondToRequest should be used by a service to respond to request
func (mh *MqttTransport) RespondToRequest(requestMsg *FimpMessage, responseMsg *FimpMessage) error {
	if requestMsg.ResponseToTopic == "" {
		return errors.New("empty response topic")
	}
	return mh.PublishToTopic(requestMsg.ResponseToTopic, responseMsg)
}

func (mh *MqttTransport) PublishSync(addr *Address, fimpMsg *FimpMessage) error {
	mh.ensureDefaultSource(fimpMsg)

	var bytm []byte
	var err error
	switch addr.PayloadType {
	case DefaultPayload:
		bytm, err = fimpMsg.SerializeToJson()
	case CompressedJsonPayload:
		bytm, err = mh.compressor.CompressFimpMsg(fimpMsg)

	}
	topic := addr.Serialize()
	if strings.TrimSpace(mh.globalTopicPrefix) != "" {
		topic = AddGlobalPrefixToTopic(mh.getGlobalTopicPrefix(), topic)
	}
	if err == nil {
		log.Trace("<MqttAd> Publishing msg to topic:", topic)
		token := mh.client.Publish(topic, mh.pubQos, false, bytm)
		if token.WaitTimeout(mh.syncPublishTimeout) && token.Error() == nil {
			return nil
		} else {
			return token.Error()
		}
	}
	return err
}

func (mh *MqttTransport) PublishRaw(topic string, bytem []byte) {
	log.Trace("<MqttAd> Publishing msg to topic:", topic)
	mh.client.Publish(topic, mh.pubQos, false, bytem)
}

func (mh *MqttTransport) PublishRawSync(topic string, bytem []byte) error {
	log.Trace("<MqttAd> Publishing msg to topic:", topic)
	token := mh.client.Publish(topic, mh.pubQos, false, bytem)
	if token.WaitTimeout(mh.syncPublishTimeout) && token.Error() == nil {
		return nil
	} else {
		return token.Error()
	}

}

// AddGlobalPrefixToTopic , adds prefix to topic .
func AddGlobalPrefixToTopic(domain string, topic string) string {
	// Check if topic is already prefixed with  "/" if yes then concat without adding "/"
	// 47 is code of "/"
	if topic[0] == 47 {
		return domain + topic
	}

	if strings.TrimSpace(domain) == "" {
		return topic
	}
	return domain + "/" + topic
}

// DetachGlobalPrefixFromTopic detaches domain from topic
func DetachGlobalPrefixFromTopic(topic string) (string, string) {
	spt := strings.Split(topic, "/")
	var resultTopic, globalPrefix string
	for i := range spt {
		if strings.Contains(spt[i], "pt:") {
			//resultTopic= strings.Replace(topic, spt[0]+"/", "", 1)
			resultTopic = strings.Join(spt[i:], "/")
			globalPrefix = strings.Join(spt[:i], "/")
			break
		}
	}

	// returns domain , topic
	return globalPrefix, resultTopic
}

// ConfigureTls The method should be used to configure mutual TLS , like AwS IoT core is using . Also it configures TLS protocol switch .
// Cert dir should contains all CA root certificates .
// IsAws flag controls AWS specific TLS protocol switch.
func (mh *MqttTransport) ConfigureTls(privateKeyFileName, certFileName, certDir string, isAws bool) error {
	mh.certDir = certDir
	privateKeyFileName = filepath.Join(certDir, privateKeyFileName)
	certFileName = filepath.Join(certDir, certFileName)
	TLSConfig := &tls.Config{InsecureSkipVerify: false}
	if isAws {
		TLSConfig.NextProtos = []string{"x-amzn-mqtt-ca"}
	}

	certPool, err := mh.getCACertPool()
	if err != nil {
		return err
	}
	TLSConfig.RootCAs = certPool

	if strings.TrimSpace(certFileName) != "" {
		certPool, err := mh.getCertPool(certFileName)
		if err != nil {
			return err
		}
		TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		TLSConfig.ClientCAs = certPool
	}
	if privateKeyFileName != "" {
		if certFileName == "" {
			return fmt.Errorf("key specified but cert is not specified")
		}
		cert, err := tls.LoadX509KeyPair(certFileName, privateKeyFileName)
		if err != nil {
			return err
		}
		TLSConfig.Certificates = []tls.Certificate{cert}
	}
	mh.mqttOptions.SetTLSConfig(TLSConfig)
	mh.client = MQTT.NewClient(mh.mqttOptions)
	return nil

}

// configuring CA certificate pool
func (mh *MqttTransport) getCACertPool() (*x509.CertPool, error) {
	certs := x509.NewCertPool()
	cafile := filepath.Join(mh.certDir, "root-ca-1.pem")
	pemData, err := ioutil.ReadFile(cafile)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)

	cafile = filepath.Join(mh.certDir, "root-ca-2.pem")
	pemData, err = ioutil.ReadFile(cafile)
	certs.AppendCertsFromPEM(pemData)

	cafile = filepath.Join(mh.certDir, "root-ca-3.pem")
	pemData, err = ioutil.ReadFile(cafile)
	certs.AppendCertsFromPEM(pemData)
	log.Infof("CA certificates are loaded.")
	return certs, nil
}

// configuring certificate pool
func (mh *MqttTransport) getCertPool(certFile string) (*x509.CertPool, error) {
	certs := x509.NewCertPool()
	pemData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)
	log.Infof("Certificate is loaded.")
	return certs, nil
}
