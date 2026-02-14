package fimpgo

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/futurehomeno/fimpgo/security"
	"github.com/futurehomeno/fimpgo/utils"
)

type MessageHandler func(topic string, addr *Address, iotMsg *FimpMessage, rawPayload []byte)

func NewMqttTransport(serverURI, clientID, username, password string, cleanSession bool, subQos byte, pubQos byte, errHandler func(error)) *MqttTransport {
	mh := MqttTransport{}
	clientOptions := defaultClientOptions(serverURI, clientID, username, password, cleanSession)
	clientOptions.SetDefaultPublishHandler(mh.onMessage)
	clientOptions.SetOnConnectHandler(mh.onConnect)

	mh.client = MQTT.NewClient(clientOptions)
	mh.pubQos = pubQos
	mh.subQos = subQos
	mh.subs = make(map[string]byte)
	mh.subChannels = make(map[string]MessageCh)
	mh.subFilters = make(map[string]FimpFilter)
	mh.subFilterFuncs = make(map[string]FilterFunc)
	mh.mainQueue = make(chan MQTT.Message, defaultMainQueueSize)
	mh.startFailRetryCount = 10
	mh.receiveChTimeout.Store(10)
	mh.syncPublishTimeout = time.Second * 5
	mh.compressor = NewMsgCompressor("", "")
	mh.errorHandler = errHandler
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
	mh.mainQueue = make(chan MQTT.Message, defaultMainQueueSize)
	mh.startFailRetryCount = 10
	mh.receiveChTimeout.Store(10)
	mh.syncPublishTimeout = time.Second * 5
	mh.compressor = NewMsgCompressor("", "")
	return &mh
}

func NewMqttTransportFromConfigs(cfg MqttConnectionConfigs, errHandler func(error), options ...Option) *MqttTransport {
	applyDefaults(&cfg)

	for _, o := range options {
		o.apply(&cfg)
	}

	mh := &MqttTransport{}

	if cfg.PrivateKeyFileName != "" && cfg.CertFileName != "" {
		mh = NewMqttTransportTLS(cfg.ServerURI, cfg.ClientID, cfg.Username, cfg.Password, cfg.CleanSession, cfg.SubQos, cfg.PubQos, cfg.errorHandler,
			cfg.PrivateKeyFileName, cfg.CertFileName, cfg.CertDir, cfg.IsAws)
	} else {
		mh = NewMqttTransport(cfg.ServerURI, cfg.ClientID, cfg.Username, cfg.Password, cfg.CleanSession, cfg.SubQos, cfg.PubQos, cfg.errorHandler)
	}

	if mh == nil {
		return nil
	}

	if cfg.StartFailRetryCount > 0 {
		mh.startFailRetryCount = cfg.StartFailRetryCount
	}

	if cfg.ReceiveChTimeout > 0 {
		mh.receiveChTimeout.Store(cfg.ReceiveChTimeout)
	}

	if cfg.MainQueueSize > 0 {
		mh.mainQueue = make(chan MQTT.Message, cfg.MainQueueSize)
	}

	return mh
}

func (mh *MqttTransport) Start(timeout time.Duration) error {
	mh.connState.Init()

	// try to connect with retries
	err := func() (ret error) {
		for i := 1; i <= mh.startFailRetryCount; i++ {
			token := mh.client.Connect()

			if !token.WaitTimeout(timeout) {
				ret = errors.New("timeout")
				continue
			}

			ret = token.Error()

			if ret == nil {
				return nil
			}

			log.Warnf("[fimpgo] MQTT connect failed %d/%d err: %v", i, mh.startFailRetryCount, ret)
			delay := time.Duration(i) * time.Duration(i)
			time.Sleep(delay * time.Second)
		}

		return ret
	}()

	if err != nil {
		return err
	}

	mh.incMsgsWg = sync.WaitGroup{}
	mh.incMsgsWg.Add(1)

	ret := mh.connState.WaitConnected(timeout)

	if ret == nil {
		go mh.handleIncomingMessages()
	}

	return ret
}

func (mh *MqttTransport) IsConnected() bool {
	return mh.connState.IsConnected()
}

func (mh *MqttTransport) Stop() {
	mh.connState.OnDone()

	if !mh.connState.IsConnected() {
		return
	}

	mh.incMsgsWg.Wait()
	time.Sleep(100 * time.Millisecond)
}

// Subscribe - subscribing for topic
func (mh *MqttTransport) Subscribe(topic string) error {
	if strings.TrimSpace(topic) == "" {
		return nil
	}

	topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix(), topic)

	mh.subscribeLock.Lock()
	defer mh.subscribeLock.Unlock()

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	token := mh.client.Subscribe(topic, mh.subQos, nil)
	isInTime := token.WaitTimeout(time.Second * 20)
	if token.Error() != nil {
		return token.Error()
	} else if !isInTime {
		return errors.New("subscribe timed out")
	}

	mh.subs[topic] = mh.subQos
	return nil
}

// Unsubscribe , unsubscribing from topic
func (mh *MqttTransport) Unsubscribe(topic string) error {
	topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix(), topic)

	mh.subscribeLock.Lock()
	defer mh.subscribeLock.Unlock()

	token := mh.client.Unsubscribe(topic)
	isInTime := token.WaitTimeout(time.Second * 20)
	if token.Error() != nil {
		return token.Error()
	} else if !isInTime {
		return errors.New("unsubscribe timed out")
	}
	delete(mh.subs, topic)
	return nil
}

func (mh *MqttTransport) UnsubscribeAll() error {
	var ret string
	var topics []string
	mh.subscribeLock.Lock()
	for i := range mh.subs {
		topics = append(topics, i)
	}
	mh.subscribeLock.Unlock()
	for _, t := range topics {
		if err := mh.Unsubscribe(t); err != nil {
			ret += fmt.Sprintf("Unsubscribe from topic %s err: %s\n", t, err.Error())
		}
	}

	if ret != "" {
		return errors.New(ret)
	}

	return nil
}

func (mh *MqttTransport) SetGlobalTopicPrefix(prefix string) {
	mh.globalTopicPrefixLock.Lock()
	mh._globalTopicPrefix = strings.TrimSpace(prefix)
	mh.globalTopicPrefixLock.Unlock()
}

func (mh *MqttTransport) globalTopicPrefix() string {
	mh.globalTopicPrefixLock.RLock()
	defer mh.globalTopicPrefixLock.RUnlock()
	return mh._globalTopicPrefix
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
	mh.channelRegLock.Lock()
	mh.subChannels[channelId] = messageCh
	mh.channelRegLock.Unlock()
}

// UnregisterChannel should be used to unregister channel
func (mh *MqttTransport) UnregisterChannel(channelId string) {
	mh.channelRegLock.Lock()
	delete(mh.subChannels, channelId)
	delete(mh.subFilters, channelId)
	delete(mh.subFilterFuncs, channelId)
	mh.channelRegLock.Unlock()
}

// RegisterChannelWithFilter should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *MqttTransport) RegisterChannelWithFilter(channelId string, messageCh MessageCh, filter FimpFilter) {
	mh.channelRegLock.Lock()
	mh.subChannels[channelId] = messageCh
	mh.subFilters[channelId] = filter
	mh.channelRegLock.Unlock()
}

// RegisterChannelWithFilterFunc should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *MqttTransport) RegisterChannelWithFilterFunc(channelId string, messageCh MessageCh, filterFunc FilterFunc) {
	mh.channelRegLock.Lock()
	mh.subChannels[channelId] = messageCh
	mh.subFilterFuncs[channelId] = filterFunc
	mh.channelRegLock.Unlock()
}

func onConnectionLost(client MQTT.Client, err error) {
	options := client.OptionsReader()
	log.Warnf("[fimpgo] Client=%s lost connection with the broker err: %v", options.ClientID(), err)
}

func onConnectionNotifEvt(client MQTT.Client, _type MQTT.ConnectionNotification) {
	options := client.OptionsReader()
	log.Errorf("[fimpgo] Client=%s notification %s", options.ClientID(), connectionNotifStr(_type.Type()))
}

func (mh *MqttTransport) onConnect(client MQTT.Client) {
	mh.subscribeLock.Lock()
	defer mh.subscribeLock.Unlock()

	options := client.OptionsReader()
	log.Infof("[fimpgo] '%s' connected to the broker", options.ClientID())

	if len(mh.subs) > 0 {
		if token := mh.client.SubscribeMultiple(mh.subs, nil); token.Wait() && token.Error() != nil {
			log.Error("[fimpgo] Subscribe error:", token.Error())
		}
	}

	mh.connState.OnConnect()
}

// onMessage is a message handler registered with MQTT client.
// It enqueues incoming messages to an intermediate queue or drops them if the queue is full.
// The intermediate queue is required, because the handler cannot be blocking.
func (mh *MqttTransport) onMessage(_ MQTT.Client, msg MQTT.Message) {
	select {
	case mh.mainQueue <- msg:
		mh.mainQueueOverflowCnt.Store(0)
		return
	default:
		mh.mainQueueOverflowCnt.Add(1)

		// stop MQTT and inform higher layer when unrecoverable situation occurs
		if mh.mainQueueOverflowCnt.Load() > 20 {
			mh.Stop()

			if mh.errorHandler != nil {
				mh.errorHandler(errors.New("main msg queue stuck"))
			}
		} else {
			log.Error("[fimpgo] Main msg queue overflow")
		}
	}
}

func (mh *MqttTransport) handleIncomingMessages() {
	defer mh.incMsgsWg.Done()

	for {
		select {
		case <-mh.connState.DoneC():
			mh.client.Disconnect(250)
			return
		case msg := <-mh.mainQueue:
			mh.handleIncomingMessage(msg)
		}
	}
}

func (mh *MqttTransport) handleIncomingMessage(msg MQTT.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("[fimpgo] handleIncomingMessage crash %v", r)
			log.Info(string(debug.Stack()))
		}
	}()

	var topic string
	if mh._globalTopicPrefix != "" {
		_, topic = DetachGlobalPrefixFromTopic(msg.Topic())
	} else {
		topic = msg.Topic()
	}

	addr, err := NewAddressFromString(topic)
	if err != nil {
		log.Errorf("[fimpgo] Processing topic=%v err:%v", topic, err)
		return
	}

	var fimpMsg *FimpMessage

	switch addr.PayloadType {
	case DefaultPayload:
		fimpMsg, err = NewMessageFromBytes(msg.Payload())
	case CompressedJsonPayload:
		if mh.compressor == nil {
			log.Warnf("[fimpgo] Compressor is not initialized for topic=%s", topic)
			return
		}

		fimpMsg, err = mh.compressor.DecompressFimpMsg(msg.Payload())
	default:
		// This means unknown binary payload, for instance compressed message
		log.Warnf("[fimpgo] Unknown payload type=%s topic=%s", addr.PayloadType, topic)
		return
	}

	if err != nil {
		log.Errorf("[fimpgo] Processing payload from topic=%s err: %v", topic, err)
		log.Tracef("[fimpgo] Payload preview len=%d: %.100s", len(msg.Payload()), msg.Payload())
		return
	}

	if mh.msgHandler != nil {
		mh.msgHandler(topic, addr, fimpMsg, msg.Payload())
	}

	var msgChs []MessageCh
	var chNames []string

	mh.channelRegLock.Lock()
	for i := range mh.subChannels {
		if mh.isChannelInterested(i, topic, addr, fimpMsg) {
			msgChs = append(msgChs, mh.subChannels[i])
			chNames = append(chNames, i)
		}
	}
	mh.channelRegLock.Unlock()

	for i, c := range msgChs {
		timer := time.NewTimer(time.Second * time.Duration(mh.receiveChTimeout.Load()))

		select {
		case c <- &Message{Topic: topic, Addr: addr, Payload: fimpMsg}:
		case <-timer.C:
			log.Warnf("[fimpgo] Channel %s not read for %d sec", chNames[i], mh.receiveChTimeout.Load())
		}

		timer.Stop()
	}
}

// isChannelInterested validates if channel is interested in message. Filtering is executed against either static filters or filter function
func (mh *MqttTransport) isChannelInterested(chanName string, topic string, addr *Address, msg *FimpMessage) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("[fimpgo] isChannelInterested crash %v", r)
			log.Info(string(debug.Stack()))
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
	} else {
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
	if addr.PayloadType == "" {
		addr.PayloadType = DefaultPayload
	}
	switch addr.PayloadType {
	case DefaultPayload:
		bytm, err = fimpMsg.SerializeToJson()
	case CompressedJsonPayload:
		bytm, err = mh.compressor.CompressFimpMsg(fimpMsg)
	default:
		// This means unknown binary payload , for instance compressed message
		log.Warnf("[fimpgo] Publish - unknown binary PayloadType=%v", addr.PayloadType)
	}
	if err != nil {
		return err
	}
	topic := addr.Serialize()
	if mh._globalTopicPrefix != "" {
		topic = AddGlobalPrefixToTopic(mh._globalTopicPrefix, topic)
	}

	log.Trace("[fimpgo] Publishing msg to topic:", topic)
	mh.client.Publish(topic, mh.pubQos, false, bytm)
	return nil
}

// PublishToTopic publishes iotMsg to string topic
func (mh *MqttTransport) PublishToTopic(topic string, fimpMsg *FimpMessage) error {
	mh.ensureDefaultSource(fimpMsg)

	byteMessage, err := fimpMsg.SerializeToJson()
	if err != nil {
		return err
	}
	addr, err := NewAddressFromString(topic)
	if err == nil {
		if addr.PayloadType == CompressedJsonPayload {
			byteMessage, err = mh.compressor.CompressBinMsg(byteMessage)
			if err != nil {
				return err
			}
		}
	}

	if mh._globalTopicPrefix != "" {
		topic = AddGlobalPrefixToTopic(mh._globalTopicPrefix, topic)
	}

	log.Trace("[fimpgo] Publishing msg to topic:", topic)
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
	if addr.PayloadType == "" {
		addr.PayloadType = DefaultPayload
	}
	switch addr.PayloadType {
	case DefaultPayload:
		bytm, err = fimpMsg.SerializeToJson()
	case CompressedJsonPayload:
		bytm, err = mh.compressor.CompressFimpMsg(fimpMsg)

	}
	topic := addr.Serialize()
	if mh._globalTopicPrefix != "" {
		topic = AddGlobalPrefixToTopic(mh._globalTopicPrefix, topic)
	}

	if err == nil {
		log.Trace("[fimpgo] Publishing msg to topic:", topic)
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
	log.Trace("[fimpgo] Publishing msg to topic:", topic)
	mh.client.Publish(topic, mh.pubQos, false, bytem)
}

func (mh *MqttTransport) PublishRawSync(topic string, bytem []byte) error {
	log.Trace("[fimpgo] Publishing msg to topic:", topic)
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
	if len(topic) > 0 && topic[0] == '/' {
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
		payloadTypeHdr := "pt:"
		if len(spt[i]) >= len(payloadTypeHdr) && strings.Contains(spt[i], payloadTypeHdr) {
			resultTopic = strings.Join(spt[i:], "/")
			globalPrefix = strings.Join(spt[:i], "/")
			break
		}
	}

	// returns domain , topic
	return globalPrefix, resultTopic
}

func NewMqttTransportTLS(serverURI, clientID, username, password string, cleanSession bool, subQos byte, pubQos byte, errHandler func(error),
	privKeyFileName, certFileName, certDir string, isAWS bool) *MqttTransport {
	mh := &MqttTransport{}
	clientOptions := defaultClientOptions(serverURI, clientID, username, password, cleanSession)
	clientOptions.SetDefaultPublishHandler(mh.onMessage)
	clientOptions.SetOnConnectHandler(mh.onConnect)

	mh.pubQos = pubQos
	mh.subQos = subQos
	mh.subs = make(map[string]byte)
	mh.subChannels = make(map[string]MessageCh)
	mh.subFilters = make(map[string]FimpFilter)
	mh.subFilterFuncs = make(map[string]FilterFunc)
	mh.mainQueue = make(chan MQTT.Message, defaultMainQueueSize)
	mh.startFailRetryCount = 10
	mh.receiveChTimeout.Store(10)
	mh.syncPublishTimeout = time.Second * 5
	mh.compressor = NewMsgCompressor("", "")
	mh.errorHandler = errHandler

	mh.certDir = certDir
	configTLS, err := security.TLSConfig(privKeyFileName, certFileName, certDir, isAWS)

	if err != nil {
		log.Errorf("[fimpgo] TLS config err: %v", err)
		return nil
	}

	clientOptions.SetTLSConfig(configTLS)
	mh.client = MQTT.NewClient(clientOptions)

	return mh
}

func (mh *MqttTransport) SetReceiveChTimeout(receiveChTimeout uint32) {
	mh.receiveChTimeout.Store(receiveChTimeout)
}

func (mh *MqttTransport) SetCertDir(certDir string) {
	mh.certDir = certDir
}

func defaultClientOptions(serverURI, clientID, username, password string, cleanSession bool) *MQTT.ClientOptions {
	clientOptions := MQTT.NewClientOptions().AddBroker(serverURI)
	clientOptions.SetClientID(clientID)
	clientOptions.SetUsername(username)
	clientOptions.SetPassword(password)
	clientOptions.SetCleanSession(cleanSession)
	clientOptions.SetAutoReconnect(true)
	clientOptions.SetConnectRetry(true)
	clientOptions.SetWriteTimeout(time.Second * 30)
	clientOptions.SetConnectionLostHandler(onConnectionLost)
	clientOptions.SetConnectionNotificationHandler(onConnectionNotifEvt)

	return clientOptions
}
