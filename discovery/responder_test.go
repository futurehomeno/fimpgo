package discovery

import (
	"testing"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

func SecondResponder() {
	mqt := fimpgo.NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest-2", "", "", true, 1, 1, nil)
	err := mqt.Start()
	if err != nil {
		log.Error("Error connecting to broker ", err)
	}

	resource := Resource{
		ResourceName:           "test-app-2",
		ResourceType:           ResourceTypeApp,
		Author:                 "aleks",
		IsInstanceConfigurable: false,
		InstanceId:             "1",
		Version:                "1",
		AppInfo:                AppInfo{},
	}

	responder := NewServiceDiscoveryResponder(mqt)
	responder.RegisterResource(resource)
	responder.Start()
}

func TestServiceDiscoveryResponder_Start(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	go SecondResponder()

	mqt := fimpgo.NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest-1", "", "", true, 1, 1, nil)
	err := mqt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	resource := Resource{
		ResourceName:           "test-app-1",
		ResourceType:           ResourceTypeApp,
		Author:                 "aleks",
		IsInstanceConfigurable: false,
		InstanceId:             "1",
		Version:                "1",
		AppInfo:                AppInfo{},
	}

	responder := NewServiceDiscoveryResponder(mqt)
	responder.RegisterResource(resource)
	responder.Start()

	discoveredResource, _ := DiscoverResources(mqt, 2)

	if len(discoveredResource) != 2 {
		t.Fatal("number of discovered resources doesn't match ")
	}

	discoveredResource, _ = DiscoverResources(mqt, 2)

	if len(discoveredResource) != 2 {
		t.Fatal("number of discovered resources doesn't match ")
	}
}
