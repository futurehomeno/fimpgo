package fimpgo

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type (
	Option interface {
		apply(*MqttConnectionConfigs)
	}

	connectionLostHandler mqtt.ConnectionLostHandler
)

func applyDefaults(mqttConfigs *MqttConnectionConfigs) {
	// backwards compatibility
	mqttConfigs.connectionLostHandler = defaultConnectionLastHandler
}

func (clh connectionLostHandler) apply(connectionConfigs *MqttConnectionConfigs) {
	connectionConfigs.connectionLostHandler = mqtt.ConnectionLostHandler(clh)
}

func defaultConnectionLastHandler(_ mqtt.Client, err error) {
	log.Errorf("connection lost with MQTT broker . Error : %v", err)
}

func WithConnectionLostHandler(h mqtt.ConnectionLostHandler) Option {
	if h == nil {
		return connectionLostHandler(defaultConnectionLastHandler)
	}
	return connectionLostHandler(h)
}
