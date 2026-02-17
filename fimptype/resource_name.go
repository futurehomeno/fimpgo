package fimptype

type ResourceNameT string

const (
	BackendServiceRn ResourceNameT = "backend-service"
	CloudBridgeRn    ResourceNameT = "clbridge"
	DefaRn           ResourceNameT = "defa"
	EcollectorRn     ResourceNameT = "ecollector"
	EnergyGuardRn    ResourceNameT = "energy_guard"
	EaseeRn          ResourceNameT = "easee"
	FhButlerRn       ResourceNameT = "fhbutler"
	GatewayRn        ResourceNameT = "gateway"
	KindOwlRn        ResourceNameT = "kind_owl"
	ScheduleRn       ResourceNameT = "schedule"
	TibberRn         ResourceNameT = "tibber"
	TimeOwlRn        ResourceNameT = "time_owl"
	TpFlowRn         ResourceNameT = "tpflow"
	VinculumRn       ResourceNameT = "vinculum"
	ZaptecRn         ResourceNameT = "zaptec"
	ZigbeeRn         ResourceNameT = "zigbee"
	ZwaveRn          ResourceNameT = "zw"

	FimpeeRn ResourceNameT = "fimpee"
)

func (s ResourceNameT) Str() string {
	return string(s)
}
