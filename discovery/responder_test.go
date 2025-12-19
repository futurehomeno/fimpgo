package discovery

import (
	"testing"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

func SecondResponder() {
	mqt := fimpgo.NewMqttTransport("tcp://localhost:1883", "fimpgotest-2", "", "", true, 1, 1)
	err := mqt.Start()
	if err != nil {
		//t.Error("Error connecting to broker ",err)
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

	mqt := fimpgo.NewMqttTransport("tcp://localhost:1883", "fimpgotest-1", "", "", true, 1, 1)
	err := mqt.Start()
	if err != nil {
		t.Error("Error connecting to broker ", err)
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

	t.Log("Sending discovery request 1 ")
	discoveredResource, _ := DiscoverResources(mqt, 2)
	for _, r := range discoveredResource {
		t.Log("Discovered resource = " + r.ResourceName)
	}
	if len(discoveredResource) != 2 {
		t.Fatal("number of discovered resources doesn't match ")
	}

	t.Log("Sending discovery request 2 ")
	discoveredResource, _ = DiscoverResources(mqt, 2)
	for _, r := range discoveredResource {
		t.Log("Discovered resource = " + r.ResourceName)
	}
	if len(discoveredResource) != 2 {
		t.Fatal("number of discovered resources doesn't match ")
	}

}
