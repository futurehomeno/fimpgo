package fimptype

type ServiceNameT string

const (
	BalanceGuardService    ServiceNameT = "balance_guard"
	CloudBridgeService     ServiceNameT = "clbridge"
	DefaService            ServiceNameT = "defa"
	EcollectorSrevice      ServiceNameT = "ecollector"
	EnergyGuardService     ServiceNameT = "energy_guard"
	EaseeService           ServiceNameT = "easee"
	FhButlerService        ServiceNameT = "fhbutler"
	GatewayService         ServiceNameT = "gateway"
	KindOwlService         ServiceNameT = "kind_owl"
	MaxGuardService        ServiceNameT = "max_guard"
	PriceGuardService      ServiceNameT = "price_guard"
	ScheduleService        ServiceNameT = "schedule"
	TibberService          ServiceNameT = "tibber"
	TimeOwlService         ServiceNameT = "time_owl"
	TpFlowService          ServiceNameT = "tpflow"
	VinculumService        ServiceNameT = "vinculum"
	ZaptecService          ServiceNameT = "zaptec"
	ZigbeeService          ServiceNameT = "zigbee"
	ZwaveDeprecatedService ServiceNameT = "zw"
	ZwaveService           ServiceNameT = "zwave-ad"
)

func (s ServiceNameT) Str() string {
	return string(s)
}
