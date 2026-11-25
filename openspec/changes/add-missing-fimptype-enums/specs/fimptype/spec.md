## ADDED Requirements

### Requirement: Device service names
`fimptype` SHALL declare a `ServiceNameT` constant for every device service defined by the fimp-api
specification at commit `9e69cfe`, named `<CamelCaseName>Service` and valued with the exact wire name.

#### Scenario: device service constant matches the wire name
- **WHEN** a consumer builds a service definition with `fimptype.OutBinSwitchService`
- **THEN** the service name on the wire is `out_bin_switch`

#### Scenario: spelling follows the spec
- **WHEN** a consumer uses `fimptype.AlarmWaterValveService`
- **THEN** its value is `alarm_water_valve`

### Requirement: App and cloud identities
Each `ServiceNameT` and `ResourceNameT` constant `fimptype` declares for an app, adapter or cloud
endpoint SHALL be valued with the exact name that endpoint uses in the `serv` field and in the `rn:`
topic segment respectively.

#### Scenario: app whose service and resource names differ
- **WHEN** a consumer addresses Angry Dog
- **THEN** `fimptype.AngryDogRn` is `angry_dog` and `fimptype.AngryDogService` is `angrydog`

#### Scenario: cloud endpoint resource
- **WHEN** a consumer builds the cloud telemetry config topic
- **THEN** `fimptype.TelemetryRn` is `telemetry`
