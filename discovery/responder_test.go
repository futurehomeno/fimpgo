package discovery

import (
	"testing"
	"time"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

func SecondResponder() {
	mqtt := fimpgo.NewMqttTransport("tcp://127.0.0.1:11883", "fimpgotest-2", "", "", true, 1, 1, nil)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		log.Error("Error connecting to broker ", err)
	}

	resource := Resource{
		ResourceName:           "test-app-2",
		ResourceType:           fimpgo.ResourceTypeApp,
		Author:                 "aleks",
		IsInstanceConfigurable: false,
		InstanceId:             "1",
		Version:                "1",
	}

	responder := NewServiceDiscoveryResponder(mqtt)
	responder.RegisterResource(resource)
	responder.Start()
}

func TestServiceDiscoveryResponder_Start(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	go SecondResponder()

	mqtt := fimpgo.NewMqttTransport("tcp://127.0.0.1:11883", "fimpgotest-1", "", "", true, 1, 1, nil)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	resource := Resource{
		ResourceName:           "test-app-1",
		ResourceType:           fimpgo.ResourceTypeApp,
		Author:                 "aleks",
		IsInstanceConfigurable: false,
		InstanceId:             "1",
		Version:                "1",
	}

	responder := NewServiceDiscoveryResponder(mqtt)
	responder.RegisterResource(resource)
	responder.Start()

	discoveredResource, _ := DiscoverResources(mqtt, 2)

	if len(discoveredResource) != 2 {
		t.Fatalf("Number of discovered resources doesn't match act=%d exp=%d", len(discoveredResource), 2)
	}

	discoveredResource, _ = DiscoverResources(mqtt, 2)

	if len(discoveredResource) != 2 {
		t.Fatalf("Number of discovered resources doesn't match act=%d exp=%d", len(discoveredResource), 2)
	}
}
