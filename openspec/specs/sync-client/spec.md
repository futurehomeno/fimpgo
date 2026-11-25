# sync-client Specification

## Purpose
Request/response over the asynchronous bus: send a FIMP message and block until its answer arrives
or a timeout expires.

## Requirements
### Requirement: A request waits for its correlated response

A request SHALL be published to its topic and the caller SHALL block until a message arrives whose
correlation ID equals the request's UID, or, when a response topic, service and interface are given,
a message matching all three. It SHALL return a timeout error when nothing matches within the given
number of seconds.

#### Scenario: Correlated response

- **WHEN** a request with UID `abc` is sent and a message with `corid` `abc` arrives
- **THEN** that message is returned

#### Scenario: Unrelated traffic

- **WHEN** only messages with other correlation IDs arrive within the timeout
- **THEN** a timeout error is returned

### Requirement: Response subscriptions can be managed per request

With auto-subscribe, the client SHALL subscribe to the response topic before publishing the request
and unsubscribe after the response or timeout. Without it, the caller SHALL subscribe beforehand.
Subscribing or publishing through a client without a connection SHALL fail with `not connected`.

#### Scenario: Auto-subscribe

- **WHEN** a request is sent with auto-subscribe on response topic `pt:j1/mt:rsp/rt:app/rn:x/ad:1`
- **THEN** the client subscribes to it before publishing and is unsubscribed after the call returns

### Requirement: A global prefix targets one site

When a global prefix is set, the client SHALL send the request and subscribe for the response under
that prefix, so backend services can address one hub through the cloud broker.

#### Scenario: Backend request

- **WHEN** a backend service with global prefix `site-1` sends a request
- **THEN** it is published under `site-1/`

### Requirement: The client owns only the connection it opened

When the client opened its own broker connection it SHALL close it on Stop; a transport passed in by
the caller SHALL be left running.

#### Scenario: Shared transport

- **WHEN** a client built on an app's transport is stopped
- **THEN** the app's transport stays connected
