package fimptype

type ServiceNameT string

const (
	AdaxService              ServiceNameT = "adax"
	AngryDogService          ServiceNameT = "angrydog"
	AuthAPIService           ServiceNameT = "auth-api"
	BalanceGuardService      ServiceNameT = "balance_guard"
	CloudBridgeService       ServiceNameT = "clbridge"
	DefaService              ServiceNameT = "defa"
	DeviceMonitorService     ServiceNameT = "device_monitor"
	EaseeService             ServiceNameT = "easee"
	EcollectorService        ServiceNameT = "ecollector"
	EnergyGuardService       ServiceNameT = "energy_guard"
	EssService               ServiceNameT = "ess"
	FhButlerService          ServiceNameT = "fhbutler"
	GatewayService           ServiceNameT = "gateway"
	GlamoxService            ServiceNameT = "glamox"
	HoiaxService             ServiceNameT = "hoiax"
	HubObserverService       ServiceNameT = "hub_observer"
	HubUpdateRegistryService ServiceNameT = "hub_update_registry"
	HueService               ServiceNameT = "hue"
	KindOwlService           ServiceNameT = "kind_owl"
	MaxGuardService          ServiceNameT = "max_guard"
	MillService              ServiceNameT = "mill"
	ModbusTCPService         ServiceNameT = "modbus-tcp"
	NetatmoService           ServiceNameT = "netatmo"
	PriceGuardService        ServiceNameT = "price_guard"
	ScheduleService          ServiceNameT = "schedule"
	SensiboService           ServiceNameT = "sensibo"
	SimService               ServiceNameT = "sim"
	SolarGuardService        ServiceNameT = "solar_guard"
	SonosService             ServiceNameT = "sonos"
	TelemetryService         ServiceNameT = "telemetry"
	TibberService            ServiceNameT = "tibber"
	TimeOwlService           ServiceNameT = "time_owl"
	TpFlowService            ServiceNameT = "tpflow"
	VinculumService          ServiceNameT = "vinculum"
	ZaptecService            ServiceNameT = "zaptec"
	ZigbeeService            ServiceNameT = "zigbee"
	ZwaveDeprecatedService   ServiceNameT = "zw"
	ZwaveService             ServiceNameT = "zwave-ad"
)

func (s ServiceNameT) Str() string {
	return string(s)
}
