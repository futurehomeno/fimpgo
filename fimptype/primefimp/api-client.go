package primefimp

import (
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"time"
)

type NotifyFilter struct {
	Cmd       string
	Component string
}

type ApiClient struct {
	clientID              string
	mqttTransport         *fimpgo.MqttTransport
	sClient               *fimpgo.SyncClient
	siteCache             Site
	isCacheEnabled        bool
	notifySubChannels     map[string]chan Notify
	subFilters            map[string]NotifyFilter
	inMsgChan             fimpgo.MessageCh
	stopFlag              bool
	isNotifyRouterStarted bool
}

func NewApiClient(clientID string, mqttTransport *fimpgo.MqttTransport, isCacheEnabled bool) *ApiClient {
	api := &ApiClient{clientID: clientID, mqttTransport: mqttTransport, isCacheEnabled: isCacheEnabled}
	api.notifySubChannels = make(map[string]chan Notify)
	api.subFilters = make(map[string]NotifyFilter)
	api.sClient = fimpgo.NewSyncClient(mqttTransport)
	return api
}

// RegisterChannel should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *ApiClient) RegisterChannel(channelId string, ch chan Notify) {
	mh.notifySubChannels[channelId] = ch
}

// RegisterChannel should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *ApiClient) RegisterChannelWithFilter(channelId string, ch chan Notify, filter NotifyFilter) {
	mh.notifySubChannels[channelId] = ch
	mh.subFilters[channelId] = filter
}

// UnregisterChannel shold be used to unregiter channel
func (mh *ApiClient) UnregisterChannel(channelId string) {
	delete(mh.notifySubChannels, channelId)
	delete(mh.subFilters, channelId)
}

func (mh *ApiClient) StartNotifyRouter() {
	go func() {
		mh.isNotifyRouterStarted = true
		for {
			if mh.stopFlag {
				break
			}
			mh.notifyRouter()
			log.Info("<PF-API> Restarting notify router")
		}
		log.Info("<PF-API> Notify router stopped ")
	}()
}

// This is destructor
func (mh *ApiClient) Stop() {
	mh.sClient.Stop()
	if mh.isNotifyRouterStarted {
		mh.stopFlag = true
		mh.inMsgChan <- &fimpgo.Message{}
	}

}

// Receives notify messages and forwards them using filtering
func (mh *ApiClient) notifyRouter() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("<PF-API> notify router CRASHED with error :", r)
		}
	}()

	mh.inMsgChan = make(fimpgo.MessageCh, 10)
	mh.mqttTransport.RegisterChannel(mh.clientID, mh.inMsgChan)
	mh.mqttTransport.Subscribe("pt:j1/mt:evt/rt:app/rn:vinculum/ad:1")

	for msg := range mh.inMsgChan {
		if mh.stopFlag {
			break
		}
		notif, err := FimpToNotify(msg)
		if err != nil {
			log.Debug("<PF-API> Can't cast to Notify . Err:", err)
			continue
		}
		for _, subCh := range mh.notifySubChannels {
			select {
			case subCh <- *notif:

			case <-time.After(time.Second * 10):
				log.Warn("<PF-API> Message is blocked , message is dropped ")
			}

		}
	}
}

func (mh *ApiClient) sendGetRequest(components []string) (*fimpgo.FimpMessage, error) {
	reqAddr := fimpgo.Address{MsgType: fimpgo.MsgTypeCmd, ResourceType: fimpgo.ResourceTypeApp, ResourceName: "vinculum", ResourceAddress: "1"}
	respAddr := fimpgo.Address{MsgType: fimpgo.MsgTypeRsp, ResourceType: fimpgo.ResourceTypeApp, ResourceName: mh.clientID, ResourceAddress: "1"}
	mh.sClient.AddSubscription(respAddr.Serialize())

	param := RequestParam{Components: components}
	req := Request{Cmd: CmdGet, Param: param}

	msg := fimpgo.NewMessage("cmd.pd7.request", "vinculum", fimpgo.VTypeObject, req, nil, nil, nil)
	msg.ResponseToTopic = respAddr.Serialize()
	msg.Source = mh.clientID
	return mh.sClient.SendFimpWithTopicResponse(reqAddr.Serialize(), msg, respAddr.Serialize(), "", "", 5)
}

func (mh *ApiClient) GetDevices(fromCache bool) ([]Device, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentDevice})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetDevices(), err
	}
	return nil, nil
}

func (mh *ApiClient) GetRooms(fromCache bool) ([]Room, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentRoom})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetRooms(), err
	}
	return nil, nil
}

func (mh *ApiClient) GetAreas(fromCache bool) ([]Area, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentArea})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetAreas(), err
	}
	return nil, nil
}

func (mh *ApiClient) GetThings(fromCache bool) ([]Thing, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentThing})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetThings(), err
	}
	return nil, nil
}

func (mh *ApiClient) GetShortcuts(fromCache bool) ([]Shortcut, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentShortcut})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetShortcuts(), err
	}
	return nil, nil
}

func (mh *ApiClient) GetVincServices(fromCache bool) (map[string]interface{}, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentService})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetVincServices(), err
	}
	return nil, nil
}


func (mh *ApiClient) GetSite(fromCache bool) (*Site, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentThing, ComponentDevice, ComponentRoom, ComponentArea, ComponentShortcut, ComponentHouse, ComponentMode})
		if err != nil {
			return nil, err
		}

		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		if mh.isCacheEnabled {
			mh.siteCache = *SiteFromResponse(response)
			return &mh.siteCache, err
		} else {
			return SiteFromResponse(response), err
		}
	}
	return &mh.siteCache, nil
}
