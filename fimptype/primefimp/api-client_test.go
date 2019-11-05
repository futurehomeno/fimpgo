package primefimp

import (
	"strings"
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var brokerUrl = "tcp://dev-sdu-sm-beta.local:1884"
var brokerUser = "simsek"
var brokerPass = "SivErAmEtOnyRIDIci"

func TestPrimeFimp_ClientApi_Update(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	uuid := uuid.New().String()
	validClientID := strings.ReplaceAll(uuid, "-", "")[0:22]

	mqtt := fimpgo.NewMqttTransport(brokerUrl, validClientID, brokerUser, brokerPass, true, 1, 1)
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
	log.Infof("Site contains %d devices", len(site.Devices))
	client.Stop()
}

func TestPrimeFimp_ClientApi_Notify(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	uuid := uuid.New().String()
	validClientID := strings.ReplaceAll(uuid, "-", "")[0:22]

	mqtt := fimpgo.NewMqttTransport(brokerUrl, validClientID, brokerUser, brokerPass, true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	// Actual test
	notifyCh := make(chan Notify, 10)

	client := NewApiClient("test-1", mqtt, true) // (clientID string, mqttTransport *fimpgo.MqttTransport, isCacheEnabled bool)
	client.RegisterChannel("test-1-ch", notifyCh) // (channelId string, ch chan Notify)

	client.StartNotifyRouter()
	i := 0
	//TODO : FIX HERE
	for {
		select {
			case 
				log.Debug("<PF_API> New message received.")
			case <-time.After(time.Second * 10):
				log.Warn("<PF-API> Message is blocked, message is dropped")
		}
	}
	for {
		select {
		case msg := <-notifyCh:
			log.Infof("New notify message of cmd = %s,comp = %s", msg.Cmd, msg.Component)
			i++
			if i > 1 {
				client.Stop()
				break
			}
		}
	}
}
