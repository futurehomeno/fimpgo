package discovery

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/sirupsen/logrus"
)

const (
	AppCurrentStateNotConfigured = "NOT_CONFIGURED"
	AppCurrentStateRunning       = "RUNNING"
	AppCurrentStateERROR         = "ERROR"

	discoverChanName = "discovery-responder"
)

type Resource struct {
	ResourceName     string                 `json:"resource_name"`      // zigbee , fimpui
	ResourceType     fimptype.ResourceTypeT `json:"resource_type"`      // ad - adapter , app -  application
	ResourceFullName string                 `json:"resource_full_name"` // full name is a name for app store or another registry
	Description      string                 `json:"description"`
	Author           string                 `json:"author"`
	Version          string                 `json:"version"`
	PackageName      string                 `json:"package_name"`    // in some cases package may have different name from service/resource name
	State            string                 `json:"state"`           // Current application state
	ConfigRequired   bool                   `json:"config_required"` // if true , the adapter should be configured before it can be used
	Configs          map[string]string      `json:"configs"`         // configurations params
	Props            map[string]string      `json:"props"`
	DocUrl           string                 `json:"doc_url"` // Url for
	// if true , the instance of adapter/app has to be configured before it can be used . false - adapter/app can be used without instance configuration
	IsInstanceConfigurable bool `json:"is_instance_configurable"`
	// Some system configurations can allow to run multiple instances of the same app or adapter , for instance multiple hubs under the same site and with radio module every hub //nolint:lll
	InstanceId string `json:"instance_id"`
}

type ServiceDiscoveryResponder struct {
	mqtt                  *fimpgo.MqttTransport
	resource              Resource
	discoveryRequestTopic string
	responderTopic        string
	requestsCh            fimpgo.MessageCh
	stopSignal            chan bool
}

func NewServiceDiscoveryResponder(mqtt *fimpgo.MqttTransport) *ServiceDiscoveryResponder {
	inst := &ServiceDiscoveryResponder{mqtt: mqtt, discoveryRequestTopic: "pt:j1/mt:cmd/rt:discovery", responderTopic: "pt:j1/mt:evt/rt:discovery"}
	inst.stopSignal = make(chan bool, 1)
	inst.requestsCh = make(fimpgo.MessageCh)
	return inst
}

// Start responder service listener
func (sr *ServiceDiscoveryResponder) Start() {
	if err := sr.mqtt.Subscribe(sr.discoveryRequestTopic); err != nil {
		logrus.Error("[fimpgo] Discovery responder subscribe err:", err)
		return
	}

	sr.mqtt.RegisterChannelWithFilter(discoverChanName, sr.requestsCh, struct {
		Topic     string
		Service   fimptype.ServiceNameT
		Interface string
	}{Topic: sr.discoveryRequestTopic, Service: "*", Interface: "*"})

	go sr.responder()
}

// Stop responder service listener
func (sr *ServiceDiscoveryResponder) Stop() {
	sr.stopSignal <- true

	if err := sr.mqtt.Unsubscribe(sr.discoveryRequestTopic); err != nil {
		logrus.Errorf("[fimpgo] Discovery responder unsubscribe err: %v", err)
	}

	sr.mqtt.UnregisterChannel(discoverChanName)
}

// RegisterResource should be invoked to register resource
func (sr *ServiceDiscoveryResponder) RegisterResource(res Resource) {
	sr.resource = res
}

func (sr *ServiceDiscoveryResponder) responder() {
	for {
		select {
		case <-sr.requestsCh:
			msg := fimpgo.NewMessage("evt.discovery.report", "system", fimptype.VTypeObject, sr.resource, nil, nil, nil)
			adr := fimpgo.Address{MsgType: fimptype.MsgTypeEvt, ResourceType: fimptype.ResourceTypeDiscovery}
			if err := sr.mqtt.Publish(&adr, msg); err != nil {
				logrus.Error("[fimpgo] Discovery responder publish err: ", err)
			}
		case <-sr.stopSignal:
			return
		}
	}
}
