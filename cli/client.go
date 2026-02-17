package main

import (
	"flag"
	"time"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	log "github.com/sirupsen/logrus"
)

var (
	mqtt *fimpgo.MqttTransport
	done = make(chan struct{})
)

func onMsg(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	log.Infof("[fimpgo] New msg %s", topic)
}

func onMqttError(err error) {
	log.Errorf("[fimpgo] Mqtt err: %s", err.Error())
	mqtt.Stop()
}

func main() {
	mqttHost := flag.String("host", "127.0.0.1:1883", "MQTT broker URL, for instance cube.local:1883")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	log.Infof("[fimpgo] Broker url %s", *mqttHost)
	mqtt = fimpgo.NewMqttTransport("tcp://"+*mqttHost, "", "", "", true, 1, 1, onMqttError)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		log.Error("[fimpgo] Error connecting to broker ", err)
		return
	}

	log.Infof("[fimpgo] '%s' connected to the broker", *mqttHost)

	mqtt.SetMessageHandler(onMsg)

	if err := mqtt.Subscribe("#"); err != nil {
		log.Errorf("[fimpgo] Subscribe # err: %v", err)
		return
	}

	log.Info("[fimpgo] Publish message")

	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimptype.ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	if err := mqtt.Publish(&adr, msg); err != nil {
		log.Errorf("[fimpgo] Publish err: %v", err)
	}

	<-done
}
