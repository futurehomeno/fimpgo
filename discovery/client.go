package discovery

import (
	"time"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

// DiscoverResources discovers resources around , timeout is set in seconds
func DiscoverResources(mqt *fimpgo.MqttTransport, timeout int) ([]Resource, error) {
	msg := fimpgo.NewNullMessage("cmd.discovery.request", "system", nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeCmd, ResourceType: fimpgo.ResourceTypeDiscovery}
	resCh := make(fimpgo.MessageCh)
	channel := "resource-discovery-client"
	if err := mqt.Subscribe("pt:j1/mt:evt/rt:discovery"); err != nil {
		return nil, err
	}
	mqt.RegisterChannelWithFilter(channel, resCh, struct {
		Topic     string
		Service   string
		Interface string
	}{Topic: "pt:j1/mt:evt/rt:discovery", Service: "*", Interface: "*"})

	defer func() {
		if err := mqt.Unsubscribe("pt:j1/mt:evt/rt:discovery"); err != nil {
			log.Error("[fimpgo] Unsubscribe err:", err)
		}
		mqt.UnregisterChannel("resource-discovery-client")
	}()

	resultsCh := make(chan []Resource, 20)
	// Response aggregator
	go func() {
		results := make([]Resource, 0)
		stop := false
		for !stop {
			select {
			case msg := <-resCh:
				res := Resource{}
				err := msg.Payload.GetObjectValue(&res)

				if err == nil {
					results = append(results, res)
				} else {
					log.Error("[fimpgo] Parsing object err:", err)
				}

			case <-time.After(time.Duration(timeout) * time.Second):
				stop = true
			}
		}

		resultsCh <- results
	}()

	//Sending request
	if err := mqt.Publish(&adr, msg); err != nil {
		return nil, err
	}

	result := <-resultsCh
	return result, nil
}
