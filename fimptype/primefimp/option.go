package primefimp

import (
	"github.com/futurehomeno/fimpgo"
	"time"
)

type (
	connectionPoolConfig struct {
		clientIdPrefix          string
		initialSize             int
		minSize                 int
		maxSize                 int
		connectionConfiguration fimpgo.MqttConnectionConfigs
		lifetime                time.Duration
	}
	apiClientConfig struct {
		cloudService   string
		connectionPool *connectionPoolConfig
		globalPrefix   string
	}
)

type Option interface {
	apply(*apiClientConfig)
}

type cloudServiceOption string

func (cso cloudServiceOption) apply(config *apiClientConfig) {
	config.cloudService = string(cso)
}

type globalPrefixOption string

func (cso globalPrefixOption) apply(config *apiClientConfig) {
	config.globalPrefix = string(cso)
}


func WithCloudService(service string) Option {
	return cloudServiceOption(service)
}

func WithGlobalPrefix(prefix string) Option {
	return globalPrefixOption(prefix)
}

//type connectionPoolOption struct {
//	clientIdPrefix          string
//	initialSize             int
//	minSize                 int
//	maxSize                 int
//	connectionConfiguration fimpgo.MqttConnectionConfigs
//	lifetime                time.Duration
//}
//
//func (cpo connectionPoolOption) apply(config *apiClientConfig) {
//	config.connectionPool = &connectionPoolConfig{
//		clientIdPrefix:          cpo.clientIdPrefix,
//		initialSize:             cpo.initialSize,
//		maxSize:                 cpo.maxSize,
//		minSize:                 cpo.minSize,
//		connectionConfiguration: cpo.connectionConfiguration,
//		lifetime:                cpo.lifetime,
//	}
//
//}

//func WithConnectionPool(clientIdPrefix string, initialSize, minSize, maxSize int, lifetime time.Duration, connectionConfiguration fimpgo.MqttConnectionConfigs) Option {
//	if lifetime == 0 {
//		lifetime = 20 * time.Second
//	}
//	if initialSize < 0 {
//		initialSize = 0
//	}
//
//	if minSize < 1 {
//		minSize = 1
//	}
//
//	if maxSize > 100 {
//		maxSize = 100
//	}
//
//	return connectionPoolOption{clientIdPrefix, initialSize, minSize, maxSize, connectionConfiguration, lifetime}
//}
