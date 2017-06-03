package fimpgo

import (
	log "github.com/Sirupsen/logrus"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttAdapter , mqtt adapter .
type MqttTransport struct {
	client     MQTT.Client
	msgHandler MessageHandler
	subQos     byte
	pubQos     byte
}

type MessageHandler func(topic string, addr *Address, iotMsg *FimpMessage)

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
	opts.SetConnectionLostHandler(mh.onConnectionLost)
	opts.SetOnConnectHandler(mh.onConnect)
	//create and start a client using the above ClientOptions
	mh.client = MQTT.NewClient(opts)
	mh.pubQos = pubQos
	mh.subQos = subQos
	return &mh
}

// SetMessageHandler message handler setter
func (mh *MqttTransport) SetMessageHandler(msgHandler MessageHandler) {
	mh.msgHandler = msgHandler
}

// Start , starts adapter async.
func (mh *MqttTransport) Start() error {
	log.Info("Connecting to MQTT broker ")
	if token := mh.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Stop , stops adapter.
func (mh *MqttTransport) Stop() {
	mh.client.Disconnect(250)
}

// Subscribe - subscribing for topic
func (mh *MqttTransport) Subscribe(topic string) error {
	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	log.Debug("Subscribing to topic:", topic)
	if token := mh.client.Subscribe(topic, mh.subQos, nil); token.Wait() && token.Error() != nil {
		log.Error("Can't subscribe. Error :", token.Error())
		return token.Error()
	}
	return nil
}

// Unsubscribe , unsubscribing from topic
func (mh *MqttTransport) Unsubscribe(topic string) error {
	log.Debug("Unsubscribing from topic:", topic)
	if token := mh.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mh *MqttTransport) onConnectionLost(client MQTT.Client, err error) {
	log.Errorf("Connection lost with MQTT broker . Error : %v", err)
}

func (mh *MqttTransport) onConnect(client MQTT.Client) {
	log.Infof("Connection established with MQTT broker .")
}

//define a function for the default message handler
func (mh *MqttTransport) onMessage(client MQTT.Client, msg MQTT.Message) {
	log.Debugf(" New msg from TOPIC: %s", msg.Topic())
	// log.Debug("MSG: %s\n", msg.Payload())
	addr, err := NewAddressFromString(msg.Topic())
	if err != nil {
		log.Error("Error processing address :" ,err)
		return
	}
	fimpMsg, err := NewMessageFromBytes(msg.Payload())
	if err == nil {
		mh.msgHandler(msg.Topic(), addr, fimpMsg)
	} else {
		log.Debug(string(msg.Payload()))
		log.Error("Error processing payload :" ,err)

	}
}

// Publish iotMsg
func (mh *MqttTransport) Publish(addr *Address, fimpMsg *FimpMessage) error {
	bytm, err := fimpMsg.SerializeToJson()
	topic := addr.Serialize()
	if err == nil {
		log.Debug("Publishing msg to topic:", topic)
		mh.client.Publish(topic, mh.pubQos, false, bytm)
		return nil
	}
	return err
}

func (mh *MqttTransport) PublishRaw(topic string, bytem []byte) {
	log.Debug("Publishing msg to topic:", topic)
	mh.client.Publish(topic, mh.pubQos, false, bytem)
}
