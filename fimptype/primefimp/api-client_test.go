package primefimp

import (
	"fmt"
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

	client := NewApiClient("test-1", mqtt, false)
	client.RegisterChannel("test-1-ch", notifyCh)

	client.StartNotifyRouter()
	i := 0
	for msg := range notifyCh {
		if msg.Component == ComponentDevice {
			log.Infof("New notify from device %s", *msg.GetDevice().Client.Name)
			fmt.Printf("New notify from device %s", *msg.GetDevice().Client.Name)
		}
		if msg.Component == ComponentArea {
			log.Infof("New notify from area %s", msg.GetArea().Name)
			fmt.Printf("New notify from area %s", msg.GetArea().Name)
		}
		log.Infof("New notify message of cmd = %s,comp = %s", msg.Cmd, msg.Component)
		i++
		if i > 1 {
			break
		}
	}
	client.Stop()
}
