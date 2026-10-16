# Add the missing fimptype service and resource names

## Why

`fimptype.ServiceNameT` is the type of `Service.Name`, yet it holds only app identities: none of the
FIMP device services exist, so zigbee-ad and fimpee each keep their own typed copy of the whole list,
and adapters declare `ServiceNameT("...")` literals. Apps and cloud endpoints already on the bus
(telemetry, hub_observer, auth-api, robin, angry_dog, ...) are likewise redeclared or written
inline as topic strings in each consumer.

## What Changes

- `ServiceNameT`: every device service of the fimp-api spec (96, generated list in fimpee at
  fimp-api `9e69cfe`), plus the app services angrydog, auth-api, hub_observer, hub_update_registry,
  modbus-tcp, sim and telemetry.
- `ResourceNameT`: alarm_service, angry_dog, auth-api, cloud_adapter, hub_observer, modbus-tcp,
  ocpp_server, remote-client, robin, sim, telemetry and tplex-ui.
- Additive only; stacked on #65 (app-identity enums).

## Impact

`fimptype/service_name.go`, `fimptype/device_service_name.go`, `fimptype/resource_name.go`, new spec
`fimptype`. No consumer change is required; zigbee-ad, fimpee, telemetry-collector, kind-owl,
core-energy-guard and CloudBridge can drop their local copies in follow-ups.
