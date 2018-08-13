package fimpgo

import (
	log "github.com/Sirupsen/logrus"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"time"
)

type MessageCh chan *Message

type Message struct {
	Topic string
	Addr  *Address
	Payload *FimpMessage
}

// MqttAdapter , mqtt adapter .
type MqttTransport struct {
	client     MQTT.Client
	msgHandler MessageHandler
	subQos     byte
	pubQos     byte
	subs       map[string]byte
	subChannels map[string]MessageCh
	globalTopicPrefix string
	startFailRetryCount int

}

type MessageHandler func(topic string, addr *Address, iotMsg *FimpMessage , rawPayload []byte)

// NewMqttAdapter constructor
//serverUri="tcp://localhost:1883"
func NewMqttTransport(serverURI string, clientID string, username string, password string, cleanSession bool, subQos byte, pubQos byte) *MqttTransport {
	mh := MqttTransport{}
	opts := MQTT.NewClientOptions().AddBroker(serverURI)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(mh.onMessage)
	opts.SetCleanSession(cleanSession)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(mh.onConnect)
	//create and start a client using the above ClientOptions
	mh.client = MQTT.NewClient(opts)
	mh.pubQos = pubQos
	mh.subQos = subQos
	mh.subs = make(map[string]byte)
	mh.subChannels = make(map[string]MessageCh)
	mh.startFailRetryCount = 10
	return &mh
}

func(mh *MqttTransport) SetGlobalTopicPrefix(prefix string) {
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

// OnConnectionLost Handler, handling connection state outside of the lib
func (mh *MqttTransport) SetOnConnectionLostHandler(connHandler MQTT.ConnectionLostHandler) {
	mh.SetOnConnectionLostHandler(connHandler)
}

// RegisterChannel should be used if new message has to be send to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicated to all channels.
func (mh *MqttTransport) RegisterChannel(channelId string,messageCh MessageCh) {
	mh.subChannels[channelId] = messageCh
}
// UnregisterChannel shold be used to unregiter channel
func (mh *MqttTransport) UnregisterChannel(channelId string ) {
	delete(mh.subChannels,channelId)
}

// Start , starts adapter async.
func (mh *MqttTransport) Start() error {
	log.Info("<MqttAd> Connecting to MQTT broker ")
	var err error
	var delay time.Duration
	for i:=1;i<mh.startFailRetryCount;i++{
		if token := mh.client.Connect(); token.Wait() && token.Error() == nil {
			return nil
		}else {
			err = token.Error()
		}
		delay = time.Duration(i)*time.Duration(i)
		log.Infof("<MqttAd> Connection failed , retrying after %d sec.... ",delay)
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
	topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix,topic)
	log.Debug("<MqttAd> Subscribing to topic:", topic)
	if token := mh.client.Subscribe(topic, mh.subQos, nil); token.Wait() && token.Error() != nil {
		log.Error("<MqttAd> Can't subscribe. Error :", token.Error())
		return token.Error()
	}
	mh.subs[topic]=mh.subQos
	return nil
}

// Unsubscribe , unsubscribing from topic
func (mh *MqttTransport) Unsubscribe(topic string) error {
	topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix,topic)
	log.Debug("<MqttAd> Unsubscribing from topic:", topic)
	if token := mh.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	delete(mh.subs,topic)
	return nil
}

func (mh *MqttTransport) onConnect(client MQTT.Client) {
	log.Infof("<MqttAd> Connection established with MQTT broker .")
	if len(mh.subs) >0 {
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
		_,topic = DetachGlobalPrefixFromTopic(msg.Topic())
	}else {
		topic = msg.Topic()
	}
	// log.Debug("MSG: %s\n", msg.Payload())
	addr, err := NewAddressFromString(topic)
	if err != nil {
		log.Error("<MqttAd> Error processing address :" ,err)
		return
	}

	fimpMsg, err := NewMessageFromBytes(msg.Payload())
	if mh.msgHandler != nil {
		if err == nil {
			mh.msgHandler(topic, addr, fimpMsg , msg.Payload())
		} else {
			log.Debug(string(msg.Payload()))
			log.Error("<MqttAd> Error processing payload :" ,err)
		}
	}


	for i := range mh.subChannels {
		msg := Message{Topic:topic,Addr:addr,Payload:fimpMsg}
		select {
			case mh.subChannels[i] <- &msg:
				// send to channel
			default :
				log.Info("<MqttAd> Channel is not ready")
		}
	}
}

// Publish iotMsg
func (mh *MqttTransport) Publish(addr *Address, fimpMsg *FimpMessage) error {
	bytm, err := fimpMsg.SerializeToJson()
	topic := addr.Serialize()
	if mh.globalTopicPrefix != "" {
		topic = AddGlobalPrefixToTopic(mh.globalTopicPrefix,topic)
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
	var resultTopic ,globalPrefix string
	for i := range spt {
		if strings.Contains(spt[i], "pt:") {
			 //resultTopic= strings.Replace(topic, spt[0]+"/", "", 1)
			 resultTopic = strings.Join(spt[i:],"/")
			 globalPrefix = strings.Join(spt[:i],"/")
			 break
		}
	}

	// returns domain , topic
	return globalPrefix, resultTopic
}