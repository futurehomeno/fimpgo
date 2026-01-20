package main

import (
	"flag"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

var (
	mqtt *fimpgo.MqttTransport
	done = make(chan struct{})
)

func onMsg(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	log.Infof("New message from topic %s", topic)
	log.Infof("Message received with type: %s", iotMsg.Type)
}

func onMqttError(err error) {
	log.Errorf("Mqtt err: %s", err.Error())
	mqtt.Stop()
	close(done)
}

func main() {

	mqttHost := flag.String("host", "localhost:1883", "MQTT broker URL , for instance cube.local:1883")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	log.Infof("Broker url %s", *mqttHost)
	mqtt = fimpgo.NewMqttTransport("tcp://"+*mqttHost, "", "", "", true, 1, 1, onMqttError)
	err := mqtt.Start()
	log.Infof("Connected to broker %s", *mqttHost)
	if err != nil {
		log.Error("Error connecting to broker ", err)
	}

	mqtt.SetMessageHandler(onMsg)
	//time.Sleep(time.Second*1)
	if err := mqtt.Subscribe("#"); err != nil {
		log.Error(err)
	}
	log.Info("Publishing message")

	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	if err := mqtt.Publish(&adr, msg); err != nil {
		log.Error(err)
	}

	<-done
}
