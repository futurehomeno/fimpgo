# vinculum-client Specification

## Purpose
Reading and changing the hub's site model (things, devices, rooms, areas, shortcuts, modes, timers,
state) through Vinculum's `pd7` API, with an optional cache kept current by Vinculum's notifications.

## Requirements
### Requirement: Requests go to Vinculum and wait for its response

A get or set SHALL be sent as `cmd.pd7.request` with an object value to
`pt:j1/mt:cmd/rt:app/rn:vinculum/ad:1`, with source and response topic
`pt:j1/mt:rsp/rt:app/rn:<client id>/ad:1`, and SHALL wait up to 5 s for the response. A client for a
cloud service SHALL receive on `rt:cloud/rn:backend-service/ad:<service>` instead, and a client with a
global prefix SHALL address that site.

#### Scenario: Get devices

- **WHEN** a client with ID `kind-owl` asks Vinculum for its devices
- **THEN** the request carries `resp_to` `pt:j1/mt:rsp/rt:app/rn:kind-owl/ad:1` and the devices from the response are returned

#### Scenario: Vinculum silent

- **WHEN** Vinculum does not answer within 5 s
- **THEN** the call returns a timeout error

### Requirement: The site cache loads on demand

A cached read SHALL serve from the cache, reloading it once when it holds no devices, things, rooms
or areas, and SHALL return an empty-cache error when it is still empty. Loading the site SHALL fetch
things, devices, rooms and areas, plus shortcuts, house, modes and services when app sync is enabled,
retrying with a delay of 5 s × attempt.

#### Scenario: Cold cache

- **WHEN** a cached device list is asked before the site was loaded
- **THEN** the site is loaded once and its devices are returned

### Requirement: Notifications keep the cache current

The notify router SHALL subscribe to `pt:j1/mt:evt/rt:app/rn:vinculum/ad:1`, apply each add, edit and
delete to the cache, and forward it to every registered channel whose command and component filter
matches, or to every unfiltered channel. Without app sync only rooms, areas, things and devices SHALL
update the cache. A full channel SHALL miss the notification with a warning, and the router SHALL
restart after a panic until it is stopped.

#### Scenario: Device removed

- **WHEN** Vinculum notifies `delete` for device 12
- **THEN** device 12 is no longer in the cache and a channel filtered on `delete`/`device` receives the notification

#### Scenario: Busy channel

- **WHEN** a registered channel's buffer is full
- **THEN** that channel misses the notification and other channels still receive it
