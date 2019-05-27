package discovery

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	"testing"
)

func SecondResponder() {
	mqt := fimpgo.NewMqttTransport("tcp://localhost:1883","fimpgotest-2","","",true,1,1)
	err := mqt.Start()
	if err != nil {
		//t.Error("Error connecting to broker ",err)
	}

	service := fimptype.Service{}

	resource := Resource{
		ResourceName:"test-app-2",
		ResourceType:ResourceTypeApp,
		Author:"aleks",
		IsInstanceConfigurable:false,
		InstanceId:"1",
		Version:"1",
		AppInfo: AppInfo{Services:[]fimptype.Service{service}},
	}


	responder := NewServiceDiscoveryResponder(mqt)
	responder.RegisterResource(resource)
	responder.Start()
}

func TestServiceDiscoveryResponder_Start(t *testing.T) {

	go SecondResponder()

	mqt := fimpgo.NewMqttTransport("tcp://localhost:1883","fimpgotest-1","","",true,1,1)
	err := mqt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	service := fimptype.Service{}

	resource := Resource{
		ResourceName:"test-app-1",
		ResourceType:ResourceTypeApp,
		Author:"aleks",
		IsInstanceConfigurable:false,
		InstanceId:"1",
		Version:"1",
		AppInfo: AppInfo{Services:[]fimptype.Service{service}},
	}


	responder := NewServiceDiscoveryResponder(mqt)
	responder.RegisterResource(resource)
	responder.Start()



	t.Log("Sending discovery request")
	discoveredResource := DiscoverResources(mqt,5)
	for _,r := range discoveredResource {
		t.Log("Discovered resource = "+r.ResourceName)
	}
	if len(discoveredResource) != 2 {
		t.Error("number of discovered resources doesn't match ")
	}

}
