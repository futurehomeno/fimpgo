package discovery

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/sirupsen/logrus"
	"time"
)

// DiscoverResources discovers resources around , timeout is set in seconds
func DiscoverResources(mqt *fimpgo.MqttTransport,timeout int) []Resource {
	msg := fimpgo.NewNullMessage("cmd.discovery.request", "system", nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeCmd, ResourceType: fimpgo.ResourceTypeDiscovery}
	resCh := make(fimpgo.MessageCh)
	mqt.Subscribe("pt:j1/mt:evt/rt:discovery")
	mqt.RegisterChannelWithFilter("resource-discovery-client",resCh, struct {
		Topic     string
		Service   string
		Interface string
	}{Topic: "pt:j1/mt:evt/rt:discovery",Service:"*",Interface:"*"})
	resultsCh := make(chan []Resource,20)
	// Response aggregator
	go func() {
		logrus.Info("Starting listener ")
		results := make([]Resource,0)
		stop :=false
		for {
			select {
			case msg :=<- resCh:
				logrus.Info("Discovery response from ",msg.Topic)
				res := Resource{}
				err := msg.Payload.GetObjectValue(&res)


				if err == nil {
					results = append(results,res)
				}else {
					logrus.Error("Error parsing object ",err)
				}

			case <-time.After(time.Duration(timeout)*time.Second):
				stop = true
				break
			}
			if stop {
				break
			}
		}
		resultsCh <- results
	}()
	//Sending request
	mqt.Publish(&adr,msg)
	result :=<- resultsCh
	logrus.Info("Func done  ",result)
	return result
}