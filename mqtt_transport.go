package fimpgo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/futurehomeno/fimpgo/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type MessageCh chan *Message

type Message struct {
	Topic   string
	Addr    *Address
	Payload *FimpMessage
	//RawPayload []byte
}

type FimpFilter struct {
	Topic     string
	Service   string
	Interface string
}

type FilterFunc func(topic string, addr *Address, iotMsg *FimpMessage) bool

// MqttAdapter , mqtt adapter .
type MqttTransport struct {
	client              MQTT.Client
	msgHandler          MessageHandler
	subQos              byte
	pubQos              byte
	subs                map[string]byte
	subChannels         map[string]MessageCh
	subFilters          map[string]FimpFilter
	subFilterFuncs      map[string]FilterFunc
	globalTopicPrefix   string
	startFailRetryCount int
	certDir             string
	mqttOptions         *MQTT.ClientOptions
}

func (mh *MqttTransport) SetCertDir(certDir string) {
	mh.certDir = certDir
}

type MessageHandler func(topic string, addr *Address, iotMsg *FimpMessage, rawPayload []byte)

// NewMqttAdapter constructor
//serverUri="tcp://localhost:1883"
func NewMqttTransport(serverURI string, clientID string, username string, password string, cleanSession bool, subQos byte, pubQos byte) *MqttTransport {
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
	//create and start a client using the above ClientOptions
	mh.client = MQTT.NewClient(mh.mqttOptions)
	mh.pubQos = pubQos
	mh.subQos = subQos
	mh.subs = make(map[string]byte)
	mh.subChannels = make(map[string]MessageCh)
	mh.subFilters = make(map[string]FimpFilter)
	mh.subFilterFuncs = make(map[string]FilterFunc)
	mh.startFailRetryCount = 10
	return &mh
}

func (mh *MqttTransport) SetGlobalTopicPrefix(prefix string) {
	mh.globalTopicPrefix = prefix
}

// Set number of retries transport will attempt on startup . Default value is 10
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
	mh.subChannels[channelId] = messageCh
}

// UnregisterChannel shold be used to unregiter channel
func (mh *MqttTransport) UnregisterChannel(channelId string) {
	delete(mh.subChannels, channelId)
	delete(mh.subFilters, channelId)
	delete(mh.subFilterFuncs, channelId)
}

// RegisterChannel should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *MqttTransport) RegisterChannelWithFilter(channelId string, messageCh MessageCh, filter FimpFilter) {
	mh.subChannels[channelId] = messageCh
	mh.subFilters[channelId] = filter
}

// RegisterChannel should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *MqttTransport) RegisterChannelWithFilterFunc(channelId string, messageCh MessageCh, filterFunc FilterFunc) {
	mh.subChannels[channelId] = messageCh
	mh.subFilterFuncs[channelId] = filterFunc
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

// Stop , stops adapter.
func (mh *MqttTransport) Stop() {
	mh.client.Disconnect(250)
}

// Subscribe - subscribing for topic
func (mh *MqttTransport) Subscribe(topic string) error {
	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix, topic)
	log.Debug("<MqttAd> Subscribing to topic:", topic)
	if token := mh.client.Subscribe(topic, mh.subQos, nil); token.Wait() && token.Error() != nil {
		log.Error("<MqttAd> Can't subscribe. Error :", token.Error())
		return token.Error()
	}
	mh.subs[topic] = mh.subQos
	return nil
}

// Unsubscribe , unsubscribing from topic
func (mh *MqttTransport) Unsubscribe(topic string) error {
	topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix, topic)
	log.Debug("<MqttAd> Unsubscribing from topic:", topic)
	if token := mh.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	delete(mh.subs, topic)
	return nil
}

func (mh *MqttTransport) onConnectionLost(client MQTT.Client, err error) {
	log.Errorf("<MqttAd> Connection lost with MQTT broker . Error : %v", err)
}

func (mh *MqttTransport) onConnect(client MQTT.Client) {
	log.Infof("<MqttAd> Connection established with MQTT broker .")
	if len(mh.subs) > 0 {
		if token := mh.client.SubscribeMultiple(mh.subs, nil); token.Wait() && token.Error() != nil {
			log.Error("Can't subscribe. Error :", token.Error())
		}
	}
}

//define a function for the default message handler
func (mh *MqttTransport) onMessage(client MQTT.Client, msg MQTT.Message) {
	log.Debugf("<MqttAd> New msg from TOPIC: %s", msg.Topic())
	var topic string
	if mh.globalTopicPrefix != "" {
		_, topic = DetachGlobalPrefixFromTopic(msg.Topic())
	} else {
		topic = msg.Topic()
	}

	// log.Debug("MSG: %s\n", msg.Payload())
	addr, err := NewAddressFromString(topic)
	if err != nil {
		log.Error("<MqttAd> Error processing address :", err)
		return
	}

	fimpMsg, err := NewMessageFromBytes(msg.Payload())
	if mh.msgHandler != nil {
		if err == nil {
			mh.msgHandler(topic, addr, fimpMsg, msg.Payload())
		} else {
			log.Debug(string(msg.Payload()))
			log.Error("<MqttAd> Error processing payload :", err)
		}
	}

	for i := range mh.subChannels {

		if !mh.isChannelInterested(i, topic, addr, fimpMsg) {
			continue
		}

		msg := Message{Topic: topic, Addr: addr, Payload: fimpMsg}
		select {
		case mh.subChannels[i] <- &msg:
			// send to channel
		default:
			log.Info("<MqttAd> Channel is not ready")
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

	if utils.RouteIncludesTopic(filter.Topic, topic) &&
		(msg.Service == filter.Service || filter.Service == "*") &&
		(msg.Type == filter.Interface || filter.Interface == "*") {
		return true

	}
	return false
}

// Publish iotMsg
func (mh *MqttTransport) Publish(addr *Address, fimpMsg *FimpMessage) error {
	bytm, err := fimpMsg.SerializeToJson()
	topic := addr.Serialize()
	if mh.globalTopicPrefix != "" {
		topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix, topic)
	}
	if err == nil {
		log.Debug("<MqttAd> Publishing msg to topic:", topic)
		mh.client.Publish(topic, mh.pubQos, false, bytm)
		return nil
	}
	return err
}

func (mh *MqttTransport) PublishRaw(topic string, bytem []byte) {
	log.Debug("<MqttAd> Publishing msg to topic:", topic)
	mh.client.Publish(topic, mh.pubQos, false, bytem)
}

// AddGlobalPrefixToTopic , adds prefix to topic .
func AddGlobalPrefixToTopic(domain string, topic string) string {
	// Check if topic is already prefixed with  "/" if yes then concat without adding "/"
	// 47 is code of "/"
	if topic[0] == 47 {
		return domain + topic
	}
	if domain == "" {
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
// The method should be used to configure mutual TLS , like AwS IoT core is using . Also it configures TLS protocol switch .
// Cert dir should contains all CA root certificates .
// isAws flag controls AWS specific TLS protocol switch.
func (mh *MqttTransport) ConfigureTls(privateKeyFileName, certFileName, certDir string, isAws bool) error {
	mh.certDir = certDir
	privateKeyFileName = filepath.Join(certDir,privateKeyFileName)
	certFileName = filepath.Join(certDir,certFileName)
	TLSConfig := &tls.Config{InsecureSkipVerify: false}
	if isAws {
		TLSConfig.NextProtos = []string{"x-amzn-mqtt-ca"}
	}

	certPool, err := mh.getCACertPool()
	if err != nil {
		return err
	}
	TLSConfig.RootCAs = certPool

	if certFileName != "" {
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
