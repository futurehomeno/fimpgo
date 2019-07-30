package primefimp

import (
	"encoding/json"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestMode(t *testing.T) {
	tb, _ := ioutil.ReadFile("testdata/mode.json")

	var mode []Mode

	err := json.Unmarshal(tb, &mode)
	if err != nil {
		t.Error(err.Error())
	}

	if mode[0].Id == "" {
		t.Errorf("Error unmarshling mode")
	}
}

func TestTimerWithActions(t *testing.T) {
	tb, _ := ioutil.ReadFile("testdata/timer_with_actions.json")

	var timer Timer

	json.Unmarshal(tb, &timer)

	device37 := timer.Action.Action.Device["37"].(map[string]interface{})

	if device37["power"].(string) != "on" {
		t.Errorf("Wrong power value for device 37. Expecting: on, Got: %s", device37["power"].(string))
	}
}

func TestTimerWithMode(t *testing.T) {
	tb, _ := ioutil.ReadFile("testdata/timer_with_mode.json")

	var timer Timer

	json.Unmarshal(tb, &timer)

	if timer.Action.Type != "mode" {
		t.Errorf("Wrong action type. Expection: mode, Got: %s", timer.Action.Mode)
	}
}

func TestTimerWithShortcut(t *testing.T) {
	tb, _ := ioutil.ReadFile("testdata/timer_with_shortcut.json")

	var timer Timer

	json.Unmarshal(tb, &timer)

	if timer.Action.Type != "shortcut" {
		t.Errorf("Wrong action type. Expection: mode, Got: %s", timer.Action.Mode)
	}
}

func TestPrimeFimp_SendFimpWithTopicResponse(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://cube.local:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	// Actual test
	syncClient := fimpgo.NewSyncClient(mqtt)

	reqAddr := fimpgo.Address{MsgType: fimpgo.MsgTypeCmd, ResourceType: fimpgo.ResourceTypeApp, ResourceName: "vinculum", ResourceAddress: "1"}
	respAddr := fimpgo.Address{MsgType: fimpgo.MsgTypeRsp, ResourceType: fimpgo.ResourceTypeApp, ResourceName: "fimpgo-test", ResourceAddress: "1"}
	syncClient.AddSubscription(respAddr.Serialize())

	param := RequestParam{Components: []string{"device"}}
	req := Request{Cmd: "get", Param: param}

	msg := fimpgo.NewMessage("cmd.pd7.request", "vinculum", fimpgo.VTypeObject, req, nil, nil, nil)
	msg.ResponseToTopic = respAddr.Serialize()
	msg.Source = "fimpgo-test"
	response, err := syncClient.SendFimpWithTopicResponse(reqAddr.Serialize(), msg, respAddr.Serialize(), "temp_sensor", "", 5)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}
	resp := Response{}
	err = response.GetObjectValue(&resp)

	t.Log(resp.Success)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}
	syncClient.Stop()
	if len(resp.GetDevices()) == 0 {
		t.Error("No rooms")
		t.Fail()
	}
	t.Log("Response test - OK , total number of devices = ", len(resp.GetDevices()))

}

func TestPrimeFimp_ClientApi_GetDevices(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://cube.local:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	client := NewApiClient("test-1", mqtt, false)
	devices, err := client.GetDevices(false)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}

	if len(devices) == 0 {
		t.Error("Site should have more then 0 devices ")
	}
	log.Infof("SIte contains %d devices", len(devices))
	client.Stop()
}

func TestPrimeFimp_ClientApi_GetSite(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://cube.local:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	client := NewApiClient("test-1", mqtt, false)
	site, err := client.GetSite(false)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}

	if len(site.Devices) == 0 {
		t.Error("Site should have more then 0 devices ")
	}
	log.Infof("SIte contains %d devices", len(site.Devices))
	client.Stop()
}

func TestPrimeFimp_ClientApi_Notify(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://cube.local:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	// Actual test
	notifyCh := make(chan Notify, 10)

	client := NewApiClient("test-1", mqtt, false)
	client.RegisterChannel("test-1-ch", notifyCh)

	client.StartNotifyRouter()
	i := 0
	for msg := range notifyCh {
		if msg.Component == ComponentDevice {
			log.Infof("New notify from device %s", msg.GetDevice().Client.Name)
		}
		log.Infof("New notify message of cmd = %s,comp = %s", msg.Cmd, msg.Component)
		i++
		if i > 3 {
			break
		}
	}
	client.Stop()

}
