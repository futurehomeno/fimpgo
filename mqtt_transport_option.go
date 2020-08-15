package fimpgo

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type Option interface {
	apply(*MqttConnectionConfigs)
}

type connectionLostHandler mqtt.ConnectionLostHandler

func (clh connectionLostHandler) apply(connectionConfigs *MqttConnectionConfigs) {
	connectionConfigs.connectionLostHandler = mqtt.ConnectionLostHandler(clh)
}

func defaultConnectionLastHandler(client mqtt.Client, err error) {
	log.Errorf("connection lost with MQTT broker . Error : %v", err)
}

func WithConnectionLostHandler(h mqtt.ConnectionLostHandler) Option {
	if h == nil {
		return connectionLostHandler(defaultConnectionLastHandler)
	}
	return connectionLostHandler(h)
}
