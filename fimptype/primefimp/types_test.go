package primefimp

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMode(t *testing.T) {
	tb, _ := os.ReadFile("testdata/mode.json")

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
	tb, err := os.ReadFile("testdata/timer_with_actions.json")
	require.NoError(t, err)

	var timer Timer
	err = json.Unmarshal(tb, &timer)
	require.NoError(t, err)

	device37, ok := timer.Action.Device[37]
	require.True(t, ok, "Device 37 not found in timer action map=%v", timer.Action.Device)

	if device37["power"].(string) != "on" { //nolint:forcetypeassert
		t.Errorf("Wrong power value for device 37 exp=on act=%s", device37["power"].(string)) //nolint:forcetypeassert
	}
}

func TestTimerWithMode(t *testing.T) {
	tb, err := os.ReadFile("testdata/timer_with_mode.json")
	require.NoError(t, err)

	var timer Timer
	err = json.Unmarshal(tb, &timer)
	require.NoError(t, err)

	if timer.Mode != "vacation" {
		t.Errorf("Wrong action type exp=mode act=%s", timer.Mode)
	}
}

func TestTimerWithShortcut(t *testing.T) {
	tb, err := os.ReadFile("testdata/timer_with_shortcut.json")
	require.NoError(t, err)

	var timer Timer
	err = json.Unmarshal(tb, &timer)
	require.NoError(t, err)

	if timer.Shortcut != 1 {
		t.Errorf("Wrong action type exp=shortcut act=%d", timer.Shortcut)
	}
}

func TestPrimeFimpSendFimpWithTopicResponse(t *testing.T) {
	t.Skip()
	mqtt := fimpgo.NewMqttTransport(brokerUrl, "fimpgotest", brokerUser, brokerPass, true, 1, 1, nil)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	// Actual test
	syncClient := fimpgo.NewSyncClient(mqtt)

	reqAddr := fimpgo.Address{MsgType: fimptype.MsgTypeCmd, ResourceType: fimptype.ResourceTypeApp, ResourceName: fimptype.VinculumRn, ResourceAddress: "1"}
	respAddr := fimpgo.Address{MsgType: fimptype.MsgTypeRsp, ResourceType: fimptype.ResourceTypeApp, ResourceName: "fimpgo-test", ResourceAddress: "1"}
	if err := syncClient.AddSubscription(respAddr.Serialize()); err != nil {
		t.Error("Error adding subscription", err)
		t.Fail()
	}

	param := RequestParam{Components: []string{"device"}}
	req := Request{Cmd: "get", Param: &param}

	msg := fimpgo.NewMessage("cmd.pd7.request", "vinculum", fimptype.VTypeObject, req, nil, nil, nil)
	msg.ResponseToTopic = respAddr.Serialize()
	msg.Source = "fimpgo-test"
	response, err := syncClient.SendFimpWithTopicResponse(reqAddr.Serialize(), msg, respAddr.Serialize(), "temp_sensor", "", 5)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}
	resp := Response{}
	err = response.GetObjectValue(&resp)

	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}
	syncClient.Stop()
	if len(resp.GetDevices()) == 0 {
		t.Error("No rooms")
		t.Fail()
	}
}

func TestPrimeFimpClientApiGetDevices(t *testing.T) {
	t.Skip()

	mqtt := fimpgo.NewMqttTransport(brokerUrl, "some_name", brokerUser, brokerPass, true, 1, 1, nil)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		t.Fatal("Start MQTT err:", err)
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
	log.Infof("Site contains %d devices", len(devices))
	client.Stop()
}

func TestPrimeFimpClientApiGetShortcuts(t *testing.T) {
	t.Skip()

	mqtt := fimpgo.NewMqttTransport(brokerUrl, "some_name", brokerUser, brokerPass, true, 1, 1, nil)
	mqtt.SetMessageHandler(func(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawPayload []byte) {

	})
	client := NewApiClient("test-1", mqtt, false)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}
	devices, err := client.GetShortcuts(false)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}

	if len(devices) == 0 {
		t.Error("Site should have more then 0 devices ")
	}
	log.Infof("Site contains %d shortcuts", len(devices))
	client.Stop()
}

func TestPrimeFimpClientApiGetVincServices(t *testing.T) {
	t.Skip()

	mqtt := fimpgo.NewMqttTransport(brokerUrl, "some_name", brokerUser, brokerPass, true, 1, 1, nil)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	client := NewApiClient("test-1", mqtt, false)
	services, err := client.GetVincServices(false)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}

	if len(services.FireAlarm) == 0 {
		t.Error("Fire alarm service not found")
	}
	client.Stop()
}

func TestPrimeFimpClientApiGetSite(t *testing.T) {
	t.Skip()

	mqtt := fimpgo.NewMqttTransport(brokerUrl, "some_name", brokerUser, brokerPass, true, 1, 1, nil)
	err := mqtt.Start(10 * time.Second)
	if err != nil {
		t.Fatal("Start MQTT err:", err)
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
