package primefimp

import (
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

const VincEventTopic = "pt:j1/mt:evt/rt:app/rn:vinculum/ad:1"

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

// NewApiClient Creates a new api client.
func NewApiClient(clientID string, mqttTransport *fimpgo.MqttTransport, isCacheEnabled bool) *ApiClient {
	api := &ApiClient{clientID: clientID, mqttTransport: mqttTransport, isCacheEnabled: isCacheEnabled}
	api.notifySubChannels = make(map[string]chan Notify)
	api.subFilters = make(map[string]NotifyFilter)
	api.sClient = fimpgo.NewSyncClient(mqttTransport)
	if isCacheEnabled {
		site, err := api.GetSite(false)
		if err != nil {
			log.Errorf("Error: %s", err)
		} else {
			api.siteCache = *site
		}
	}
	return api
}

// RegisterChannel should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *ApiClient) RegisterChannel(channelId string, ch chan Notify) {
	mh.notifySubChannels[channelId] = ch
}

// RegisterChannelWithFilter should be used if new message has to be sent to channel instead of callback.
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

// Stop : This is destructor
func (mh *ApiClient) Stop() {
	mh.sClient.Stop()
	if mh.isNotifyRouterStarted {
		mh.stopFlag = true
		mh.inMsgChan <- &fimpgo.Message{}
	}

}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

// UpdateSite : Updates the site according to notification message.
func (mh *ApiClient) UpdateSite(notif *Notify) {
	log.Debugf("Command: %s & Component:%s", notif.Cmd, notif.Component)
	switch notif.Cmd {
	case CmdAdd:
		switch notif.Component {
		case ComponentArea:
			mh.siteCache.AddArea(notif.GetArea())
		case ComponentDevice:
			mh.siteCache.AddDevice(notif.GetDevice())
		case ComponentRoom:
			mh.siteCache.AddRoom(notif.GetRoom())
		case ComponentTimer:
			mh.siteCache.AddTimer(notif.GetTimer())
		case ComponentThing:
			mh.siteCache.AddThing(notif.GetThing())
		case ComponentShortcut:
			mh.siteCache.AddShortcut(notif.GetShortcut())
		default:
			log.Errorf("Unknown Component:%s cannot be added.", notif.Component)
		}
	case CmdDelete:
		err := mh.siteCache.RemoveWithID(notif.Component, int(notif.Id.(float64)))
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("%s with ID:%d is deleted", notif.Component, int(notif.Id.(float64)))
		}
	case CmdEdit:
		switch notif.Component {
		case ComponentArea:
			mh.siteCache.UpdateArea(notif.GetArea())
		case ComponentDevice:
			mh.siteCache.UpdateDevice(notif.GetDevice())
		case ComponentRoom:
			mh.siteCache.UpdateRoom(notif.GetRoom())
		case ComponentTimer:
			mh.siteCache.UpdateTimer(notif.GetTimer())
		//case ComponentHub:  //  TODO: Is this possible?
		//	mh.siteCache.UpdateHub(notif.GetHub())
		//case ComponentMode: //  TODO: Is this possible?
		//	mh.siteCache.UpdateMode(notif.GetMode())
		case ComponentThing:
			mh.siteCache.UpdateThing(notif.GetThing())
		case ComponentShortcut:
			mh.siteCache.UpdateShortcut(notif.GetShortcut())
		default:
			log.Error("Unknown component update occured. Report this as issue please")
		}
	case CmdSet:
		switch notif.Component {
		case ComponentRoom:
			roomIdx := mh.siteCache.FindIndex(ComponentRoom, int(notif.Id.(float64)))
			if roomIdx != -1 {
				log.Infof("Change in room id:%d", int(notif.Id.(float64)))
			} else {
				log.Errorf("Room with ID:%d not found. Adding", int(notif.Id.(float64)))
			}
		case ComponentHub:
			if notif.Id == "mode" {
				modeChange := notif.GetModeChange()
				if modeChange.Current != modeChange.Prev {
					log.Infof("Mode is changed from %s to %s", modeChange.Prev, modeChange.Current)
				} else {
					log.Infof("Mode is same again as %s", modeChange.Current)
				}
			}
		}
	}
	if mh.isNotifyRouterStarted { // make sure notify router is started
		for cid, nf := range mh.subFilters { // check all subfilters
			if nf.Cmd == notif.Cmd && nf.Component == notif.Component {
				mh.notifySubChannels[cid] <- *notif // send notification to corresponding subchannel if there is match
			}
		}
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
	mh.mqttTransport.Subscribe(VincEventTopic)

	for msg := range mh.inMsgChan {
		if mh.stopFlag {
			break
		}
		if msg.Topic != VincEventTopic {
			continue
		}
		notif, err := FimpToNotify(msg)
		if err != nil {
			log.Debug("<PF-API> Can't cast to Notify. Err:", err)
			continue
		} else {
			mh.UpdateSite(notif)
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

// GetDevices Gets the devices
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

// GetRooms Gets the rooms
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

// GetAreas Gets the areas
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

// GetThings Gets the things
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

// GetShortcuts Gets the shortcuts
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

// GetVincServices Gets vinculum services
func (mh *ApiClient) GetVincServices(fromCache bool) (VincServices, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentService})
		if err != nil {
			return VincServices{}, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return VincServices{}, err
		}
		return response.GetVincServices(), err
	}
	return VincServices{}, nil
}

// GetSite Gets the whole site information
func (mh *ApiClient) GetSite(fromCache bool) (*Site, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentThing, ComponentDevice, ComponentRoom, ComponentArea, ComponentShortcut, ComponentHouse, ComponentMode, ComponentService})
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
