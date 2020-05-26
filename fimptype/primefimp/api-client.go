package primefimp

import (
	"errors"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"sync"
	"time"
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
	notifChMux            sync.RWMutex
	isVincAppsSyncEnabled bool

}

func (mh *ApiClient) IsCacheEnabled() bool {
	return mh.isCacheEnabled
}

func (mh *ApiClient) SetIsCacheEnabled(isCacheEnabled bool) {
	mh.isCacheEnabled = isCacheEnabled
}
// if not enabled only rooms,areas,things and devices are synced
func (mh *ApiClient) EnableVincAppsSync(flag bool) {
	mh.isVincAppsSyncEnabled = flag
}

func (mh *ApiClient) IsCacheEmpty() bool {
	if len(mh.siteCache.Devices)==0 && len(mh.siteCache.Things)==0 && len(mh.siteCache.Rooms)==0 && len(mh.siteCache.Areas)==0 {
		return true
	}
	return false
}

// ValidateSiteCache validates cache , if empty it makes one reload attempt. The method can be used for cache lazy loading.
func (mh *ApiClient) ValidateAndReloadSiteCache()bool {
	if mh.IsCacheEmpty() {
		log.Debug("<PF-API> Empty site cache.Reloading...")
		mh.ReloadSiteToCache(1)
		if mh.IsCacheEmpty() {
			return false
		}
	}
	return true
}

// NewApiClient Creates a new api client. If isCacheEnabled it set to true , it will try to sync entire site on startup.
func NewApiClient(clientID string, mqttTransport *fimpgo.MqttTransport, loadSiteIntoCache bool) *ApiClient {
	api := &ApiClient{clientID: clientID, mqttTransport: mqttTransport}
	api.notifySubChannels = make(map[string]chan Notify)
	api.subFilters = make(map[string]NotifyFilter)
	api.sClient = fimpgo.NewSyncClient(mqttTransport)
	api.notifChMux = sync.RWMutex{}
	if loadSiteIntoCache {
		api.ReloadSiteToCache(3)
	}
	return api
}

//ReloadSiteToCache loads cache from vinculum and sets isCacheEnabled flag to true if operation was successful.
func (mh *ApiClient) ReloadSiteToCache(retry int) error {
	retry++
	var site *Site
	var err error
	for i:=1;i<retry;i++ {
		log.Debug("<PF-API> Reloading site into the cache.Attempt ",i)
		site, err = mh.GetSite(false)
		if err == nil {
			log.Debug("<PF-API> Site loaded successfully")
			break
		}else {
			log.Error("<PF-API> site sync error :",err.Error())
			time.Sleep(time.Second*time.Duration(5*i))
		}
	}
	if err != nil {
		mh.isCacheEnabled = false
		log.Errorf("<PF-API>: %s", err)
		return err
	} else {
		mh.isCacheEnabled = true
		mh.siteCache = *site
		log.Debug("<PF-API> Site info successfully loaded to cache")
	}
	return nil
}

// Loads site from file . File should be in exactly the same format as vinculum response
func (mh *ApiClient) LoadVincResponseFromFile(fileName string)error {
	bSite , err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	fimpMsg,err:= fimpgo.NewMessageFromBytes(bSite)
	if err != nil {
		return err
	}
	response, err := FimpToResponse(fimpMsg)
	if err != nil {
		return err
	}
	mh.siteCache = *SiteFromResponse(response)
	return nil
}

// RegisterChannel should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *ApiClient) RegisterChannel(channelId string, ch chan Notify) {
	mh.notifChMux.Lock()
	mh.notifySubChannels[channelId] = ch
	mh.notifChMux.Unlock()
}

// RegisterChannelWithFilter should be used if new message has to be sent to channel instead of callback.
// multiple channels can be registered , in that case a message bill be multicasted to all channels.
func (mh *ApiClient) RegisterChannelWithFilter(channelId string, ch chan Notify, filter NotifyFilter) {
	mh.notifChMux.Lock()
	mh.notifySubChannels[channelId] = ch
	mh.subFilters[channelId] = filter
	mh.notifChMux.Unlock()
}

// UnregisterChannel shold be used to unregiter channel
func (mh *ApiClient) UnregisterChannel(channelId string) {
	mh.notifChMux.Lock()
	delete(mh.notifySubChannels, channelId)
	delete(mh.subFilters, channelId)
	mh.notifChMux.Unlock()
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

// UpdateSite : Updates the internal cache according to notification message.
func (mh *ApiClient) UpdateSite(notif *Notify) {
	log.Tracef("Command: %s & Component:%s", notif.Cmd, notif.Component)
	if !mh.isVincAppsSyncEnabled {
		if  notif.Component != ComponentArea && notif.Component != ComponentDevice && notif.Component != ComponentThing && notif.Component != ComponentRoom {
			log.Debugf("Component skipped")
			return
		}
	}
	switch notif.Cmd {
	case CmdAdd:
		switch notif.Component {
		case ComponentArea:
			mh.siteCache.AddArea(notif.GetArea())
		case ComponentDevice:
			mh.siteCache.AddDevice(notif.GetDevice())
		case ComponentRoom:
			mh.siteCache.AddRoom(notif.GetRoom())
		case ComponentThing:
			mh.siteCache.AddThing(notif.GetThing())
		case ComponentShortcut:
			if mh.isVincAppsSyncEnabled {
				mh.siteCache.AddShortcut(notif.GetShortcut())
			}
		case ComponentTimer:
			mh.siteCache.AddTimer(notif.GetTimer())
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
		case ComponentThing:
			mh.siteCache.UpdateThing(notif.GetThing())
		case ComponentTimer:
			mh.siteCache.UpdateTimer(notif.GetTimer())
		case ComponentShortcut:
			mh.siteCache.UpdateShortcut(notif.GetShortcut())
		default:
			log.Error("Unknown component update occured. Report this as issue please")
		}
	case CmdSet:
		switch notif.Component {
		case ComponentRoom:
			//roomIdx := mh.siteCache.FindIndex(ComponentRoom, int(notif.Id.(float64)))
			//if roomIdx != -1 {
			//	log.Infof("Change in room id:%d", int(notif.Id.(float64)))
			//} else {
			//	log.Errorf("Room with ID:%d not found. Adding", int(notif.Id.(float64)))
			//}
		case ComponentHub:
			//if notif.Id == "mode" {
			//	modeChange := notif.GetModeChange()
			//	if modeChange.Current != modeChange.Prev {
			//		log.Infof("Mode is changed from %s to %s", modeChange.Prev, modeChange.Current)
			//	} else {
			//		log.Infof("Mode is same again as %s", modeChange.Current)
			//	}
			//}
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
			if mh.isNotifyRouterStarted { // make sure notify router is started
				mh.notifChMux.RLock()
				for cid, nfCh := range mh.notifySubChannels { // check all subfilters
				    nfFilter , ok := mh.subFilters[cid]
				    var send bool
				    if ok {
						if nfFilter.Cmd == notif.Cmd && nfFilter.Component == notif.Component {
							send = true
						}
					}else {
						send = true
					}
					if send {
						select {
						case nfCh <- *notif: // send notification to corresponding subchannel if there is match
						default:
							log.Warnf("<PF-API> Send channel %s is blocked ",cid)
						}
					}
				}
				mh.notifChMux.RUnlock()
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
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Devices,nil
		}
	}

	return nil, errors.New("cache is empty")
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
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Rooms,nil
		}

	}
	return nil, errors.New("cache is empty")
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
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Areas,nil
		}
	}
	return nil, errors.New("cache is empty")
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
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Things,nil
		}
	}
	return nil, errors.New("cache is empty")
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
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Shortcuts,nil
		}
	}
	return nil, errors.New("cache is empty")
}

// GetShortcuts Gets the modes
func (mh *ApiClient) GetModes(fromCache bool) ([]Mode, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentMode})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetModes(), err
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Modes,nil
		}
	}
	return nil, errors.New("cache is empty")
}

// GetShortcuts Gets the modes
func (mh *ApiClient) GetTimers(fromCache bool) ([]Timer, error) {
	if !fromCache {
		fimpResponse, err := mh.sendGetRequest([]string{ComponentTimer})
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		return response.GetTimers(), err
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Timers,nil
		}
	}
	return nil, errors.New("cache is empty")
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
	}else if mh.isCacheEnabled {
		if mh.ValidateAndReloadSiteCache() {
			return mh.siteCache.Services,nil
		}
	}
	return VincServices{}, errors.New("cache is empty")
}

// GetSite Gets the whole site information
func (mh *ApiClient) GetSite(fromCache bool) (*Site, error) {
	if !fromCache {
		var components []string
		if mh.isVincAppsSyncEnabled {
			components = []string {ComponentThing, ComponentDevice, ComponentRoom, ComponentArea, ComponentShortcut, ComponentHouse, ComponentMode, ComponentService}
		}else {
			components = []string {ComponentThing, ComponentDevice, ComponentRoom, ComponentArea}
		}

		fimpResponse, err := mh.sendGetRequest(components)
		if err != nil {
			return nil, err
		}
		response, err := FimpToResponse(fimpResponse)
		if err != nil {
			return nil, err
		}
		if mh.isCacheEnabled {
			// Sync cache if cache is enabled
			mh.siteCache = *SiteFromResponse(response)
			return &mh.siteCache, err
		} else {
			return SiteFromResponse(response), err
		}
	}else {
		if mh.ValidateAndReloadSiteCache() {
			return &mh.siteCache, nil
		}

	}
	return nil, errors.New("cache is empty")
}
