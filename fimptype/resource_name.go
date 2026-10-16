package fimptype

type ResourceNameT string

const (
	AdaxRn           ResourceNameT = "adax"
	AlarmServiceRn   ResourceNameT = "alarm_service"
	AngryDogRn       ResourceNameT = "angry_dog"
	AuthAPIRn        ResourceNameT = "auth-api"
	BackendServiceRn ResourceNameT = "backend-service"
	CloudAdapterRn   ResourceNameT = "cloud_adapter"
	CloudBridgeRn    ResourceNameT = "clbridge"
	DefaRn           ResourceNameT = "defa"
	DeviceMonitorRn  ResourceNameT = "device_monitor"
	EaseeRn          ResourceNameT = "easee"
	EcollectorRn     ResourceNameT = "ecollector"
	EnergyGuardRn    ResourceNameT = "energy_guard"
	EssRn            ResourceNameT = "ess"
	FhButlerRn       ResourceNameT = "fhbutler"
	GatewayRn        ResourceNameT = "gateway"
	GlamoxRn         ResourceNameT = "glamox"
	HoiaxRn          ResourceNameT = "hoiax"
	HubObserverRn    ResourceNameT = "hub_observer"
	HueRn            ResourceNameT = "hue"
	KindOwlRn        ResourceNameT = "kind_owl"
	MillRn           ResourceNameT = "mill"
	ModbusTCPRn      ResourceNameT = "modbus-tcp"
	NetatmoRn        ResourceNameT = "netatmo"
	OcppServerRn     ResourceNameT = "ocpp_server"
	RemoteClientRn   ResourceNameT = "remote-client"
	RobinRn          ResourceNameT = "robin"
	ScheduleRn       ResourceNameT = "schedule"
	SensiboRn        ResourceNameT = "sensibo"
	SimRn            ResourceNameT = "sim"
	SmarthomeAppRn   ResourceNameT = "smarthome-app"
	SonosRn          ResourceNameT = "sonos"
	TelemetryRn      ResourceNameT = "telemetry"
	TibberRn         ResourceNameT = "tibber"
	TimeOwlRn        ResourceNameT = "time_owl"
	TpFlowRn         ResourceNameT = "tpflow"
	TplexUIRn        ResourceNameT = "tplex-ui"
	VinculumRn       ResourceNameT = "vinculum"
	ZaptecRn         ResourceNameT = "zaptec"
	ZigbeeRn         ResourceNameT = "zigbee"
	ZwaveRn          ResourceNameT = "zw"

	FimpeeRn ResourceNameT = "fimpee"
)

func (s ResourceNameT) Str() string {
	return string(s)
}
