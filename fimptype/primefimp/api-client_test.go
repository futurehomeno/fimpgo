package primefimp

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/futurehomeno/fimpgo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var brokerUrl = "tcp://cube.local:1883"
var brokerUser = ""
var brokerPass = ""
var testSiteGuid = ""
var awsIotEndpoint = "ssl://xxxxxxxxxx.iot.xxxxxxx.amazonaws.com:443"

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

	client := NewApiClient("test-1", mqtt, true)
	client.StartNotifyRouter()
	site, err := client.GetSite(false)
	if err != nil {
		t.Error("Error", err)
		t.Fail()
	}

	for _, r := range site.Rooms {
		log.Infof("Room %s , area = %d ", r.Alias, r.Area)

	}

	if len(site.Devices) == 0 {
		t.Error("Site should have more then 0 devices ")
	}

	notifyCh := make(chan Notify, 10)
	client.RegisterChannel("test-run-1", notifyCh)
	go func() {
		for {
			newMsg := <-notifyCh
			if newMsg.Component != "device" {
				continue
			}
			log.Infof("Update from component : %s , command : %s ", newMsg.Component, newMsg.Cmd)
			for _, r := range site.Devices {
				var name string
				if r.Client.Name != nil {
					name = *r.Client.Name
				}
				log.Infof("Device id = %d , name = %s ", r.ID, name)
			}
		}
	}()
	log.Infof("Site contains %d devices", len(site.Devices))
	time.Sleep(20 * time.Minute)
	client.Stop()
}

func TestPrimeFimp_ClientApi_Notify(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	validClientID := strings.ReplaceAll(uuid.New().String(), "-", "")[0:22]

	mqtt := fimpgo.NewMqttTransport(brokerUrl, validClientID, brokerUser, brokerPass, true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	// Actual test
	notifyCh := make(chan Notify, 10)
	apiclientid := uuid.New().String()[0:12]
	client := NewApiClient(apiclientid, mqtt, true) // (clientID string, mqttTransport *fimpgo.MqttTransport, isCacheEnabled bool)
	channelID := uuid.New().String()[0:12]
	// Using "RegisterChannel" will send a message to our notify channel for all messages
	// If you want to use filters, check "RegisterChannelWithFilter"
	client.RegisterChannel(channelID, notifyCh) // (channelId string, ch chan Notify)
	client.StartNotifyRouter()
	// Notify router is started. Now please, make 3 "add", "edit" or "delete" actions to finalize the test.
	i := 0
	limit := 3
	for {
		select {
		case msg := <-notifyCh:
			log.Infof("Check %d/%d: New notify message of cmd = %s,comp = %s", i, limit, msg.Cmd, msg.Component)
			i++
			if i > limit {
				client.Stop()
				break
			}
		}
	}
}

func TestPrimeFimp_SiteLazyLoading(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	validClientID := strings.ReplaceAll(uuid.New().String(), "-", "")[0:22]

	mqtt := fimpgo.NewMqttTransport(brokerUrl, validClientID, brokerUser, brokerPass, true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	// Actual test
	apiclientid := uuid.New().String()[0:12]
	client := NewApiClient(apiclientid, mqtt, false)
	if !client.IsCacheEmpty() {
		t.Error("Cache is not empty.Must be empty")
	}
	_, err = client.GetSite(true)
	if err != nil || client.IsCacheEmpty() {
		t.Error("Cache is empty. Cache must contain data.")
	}

}

func TestPrimeFimp_LoadStates(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	validClientID := strings.ReplaceAll(uuid.New().String(), "-", "")[0:22]
	mqtt := fimpgo.NewMqttTransport(awsIotEndpoint, validClientID, brokerUser, brokerPass, true, 1, 1)
	mqtt.ConfigureTls("awsiot.private.key","awsiot.crt","./datatools/certs",true)
	mqtt.SetGlobalTopicPrefix(testSiteGuid)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	// Actual test
	apiclientid := uuid.New().String()[0:12]
	client := NewApiClient(apiclientid, mqtt, false,WithCloudService("test-proc-1"))
	client.SetResponsePayloadType(fimpgo.CompressedJsonPayload)
	state, err := client.GetState()
	if err != nil || len(state.Devices)==0 {
		t.Error("Cache is empty. Cache must contain data.")
	}else {
		t.Log("STATES - All Good .Number of states = ",len(state.Devices))
	}
	shortcuts, err := client.GetShortcuts(false)
	if err != nil || len(shortcuts)==0 {
		t.Error("Cache is empty. Cache must contain data.")
	}else {
		t.Log("SHORTCUTS - All Good . Number of shortcuts = ",len(shortcuts))
	}
}

func TestPrimeFimp_LoadStatesWithConnPool(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	validClientID := strings.ReplaceAll(uuid.New().String(), "-", "")[0:22]
	transportConfigs := fimpgo.MqttConnectionConfigs{
		ServerURI:           awsIotEndpoint,
		ClientID:            validClientID,
		CleanSession:        true,
		SubQos:              1,
		PubQos:              1,
		GlobalTopicPrefix:   "",
		CertDir:             "./datatools/certs",
		PrivateKeyFileName:  "awsiot.private.key",
		CertFileName:        "awsiot.crt",
		IsAws:               true,
		StartFailRetryCount: 5,
	}

	connPool := fimpgo.NewMqttConnectionPool(3,5,20,time.Minute*3,transportConfigs,"lib_code_test_pool")
	connPool.Start()

	var successCounter int

	go func() {
		for i := 0; i < 3; i++ {
			connId, conn,err := connPool.BorrowConnection()
			if err != nil {
				t.Fatal("Connection pool error , Err:",err.Error())
			}
			client := NewApiClient(validClientID, conn, false, WithCloudService("test-proc-1"),WithGlobalPrefix(testSiteGuid))
			client.SetResponsePayloadType(fimpgo.CompressedJsonPayload)
			state, err := client.GetState()
			if err != nil || len(state.Devices)==0 {
				t.Fatal("Cache is empty. Cache must contain data.")
			}else {
				t.Log("STATES - All Good .Number of states = ",len(state.Devices))
				successCounter++
			}
			connPool.ReturnConnection(connId)
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			connId, conn, err := connPool.BorrowConnection()
			if err != nil {
				t.Fatal("Connection pool error , Err:", err.Error())
			}
			client := NewApiClient(validClientID, conn, false, WithCloudService("test-proc-1"), WithGlobalPrefix(testSiteGuid))
			client.SetResponsePayloadType(fimpgo.CompressedJsonPayload)
			shortcuts, err := client.GetShortcuts(false)
			if err != nil || len(shortcuts) == 0 {
				t.Fatal("Cache is empty. Cache must contain data.")
			} else {
				t.Log("SHORTCUTS - All Good . Number of shortcuts = ", len(shortcuts))
				successCounter++
			}
			connPool.ReturnConnection(connId)
		}
	}()

	time.Sleep(time.Second*10)

	if successCounter != 6 {
		t.Fatal("something went wrong")
	}else {
		t.Log("______ALL____GOOOD_______")
	}

}

func TestPrimeFimp_ClientApi_Notify_With_Filter(t *testing.T) {
	log.SetLevel(log.TraceLevel)

	validClientID := strings.ReplaceAll(uuid.New().String(), "-", "")[0:22]

	mqtt := fimpgo.NewMqttTransport(brokerUrl, validClientID, brokerUser, brokerPass, true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	// Actual test
	channelIDAdd := uuid.New().String()[0:12]
	notifyAreaAdd := make(chan Notify, 10)
	notifyfilterAreaAdd := NotifyFilter{Cmd: CmdAdd, Component: ComponentArea}

	channelIDDelete := uuid.New().String()[0:12]
	notifyAreaDelete := make(chan Notify, 10)
	notifyfilterAreaDelete := NotifyFilter{Cmd: CmdDelete, Component: ComponentArea}

	channelIDEdit := uuid.New().String()[0:12]
	notifyAreaEdit := make(chan Notify, 10)
	notifyfilterAreaEdit := NotifyFilter{Cmd: CmdEdit, Component: ComponentArea}

	apiclientid := uuid.New().String()[0:12]
	client := NewApiClient(apiclientid, mqtt, true)                                             // (clientID string, mqttTransport *fimpgo.MqttTransport, isCacheEnabled bool)
	client.RegisterChannelWithFilter(channelIDAdd, notifyAreaAdd, notifyfilterAreaAdd)          // (channelId string, ch chan Notify, filter NotifyFilter)
	client.RegisterChannelWithFilter(channelIDDelete, notifyAreaDelete, notifyfilterAreaDelete) // (channelId string, ch chan Notify, filter NotifyFilter)
	client.RegisterChannelWithFilter(channelIDEdit, notifyAreaEdit, notifyfilterAreaEdit)       // (channelId string, ch chan Notify, filter NotifyFilter)
	client.StartNotifyRouter()
	// We started the channel with filter now let's add an area, edit the name and then delete it to finalize the test.
	addarea := 0
	deletearea := 0
	editarea := 0

	closeChan := make(chan string)
	go func() {
		for {
			select {
			case msg := <-notifyAreaAdd:
				addarea++
				log.Infof("Check %s: New notify message of cmd = %s,comp = %s", msg.Cmd, msg.Cmd, msg.Component)
				if addarea > 0 && deletearea > 0 && editarea > 0 {
					client.Stop()
					closeChan <- "shit"
					break
				}
			case msg := <-notifyAreaDelete:
				deletearea++
				log.Infof("Check %s: New notify message of cmd = %s,comp = %s", msg.Cmd, msg.Cmd, msg.Component)
				if addarea > 0 && deletearea > 0 && editarea > 0 {
					client.Stop()
					closeChan <- "shit"
					break
				}
			case msg := <-notifyAreaEdit:
				editarea++
				log.Infof("Check %s: New notify message of cmd = %s,comp = %s", msg.Cmd, msg.Cmd, msg.Component)
				if addarea > 0 && deletearea > 0 && editarea > 0 {
					client.Stop()
					closeChan <- "shit"
					break
				}
			}
		}
	}()

	<-closeChan
	t.Log("Tadaaa")
}

func TestPrimeFimp_LoadSiteFromFile(t *testing.T) {
	fApi := NewApiClient("pf-test", nil, false)
	err := fApi.LoadVincResponseFromFile("testdata/site-info-response.json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	site, err := fApi.GetSite(true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	thing := site.GetThingById(9)
	if thing.Name != "YR temperature report" {
		t.Error("name doesn't match")
		t.FailNow()
	}
	t.Log(thing.Name)

	device := site.GetDeviceByServiceAddress("/rt:dev/rn:flow/ad:1/sv:out_bin_switch/ad:7zfeSQx3Q8")
	if device.ID != 12 {
		t.Error("device id doesn't match")
		t.FailNow()
	}
	t.Log(*device.Client.Name)

	room := site.GetRoomById(4)
	if room.ID != 4 {
		t.Error("room id doesn't match")
		t.FailNow()
	}
	t.Log(room.Alias)
}

func TestPrimefimp_LoadStateFromFile(t *testing.T) {
	const deviceCount = 18
	bSite, err := ioutil.ReadFile("testdata/state.json")
	if err != nil {
		t.Fatal(err)
	}
	fimpMsg, err := fimpgo.NewMessageFromBytes(bSite)
	if err != nil {
		t.Fatal(err)
	}
	response, err := FimpToResponse(fimpMsg)
	if err != nil {
		t.Fatal(err)
	}

	state, err := response.GetState()
	if err != nil {
		t.Fatal(err)
	}

	if deviceCount != len(state.Devices) {
		t.Fatal("device counts do not match")
	}

	// the current state.json file has 7 devices with the "meter_elec" service
	const meterElecDevices = 7
	filteredDevices := state.Devices.FilterDevicesByService("meter_elec")
	if len(filteredDevices) != meterElecDevices {
		t.Fatal(fmt.Sprintf("meter_elec devices count does not match. expected %d, got %d", meterElecDevices, len(filteredDevices)))
	}

	// tue current state.json file has 6 attributes with the "meter" name
	const meterAttributes = 7
	filteredDevices = state.Devices.FilterDevicesByAttribute("meter")
	if len(filteredDevices) != meterAttributes {
		t.Fatal(fmt.Sprintf("meter_elec devices count does not match. expected %d, got %d", meterAttributes, len(filteredDevices)))
	}
}
