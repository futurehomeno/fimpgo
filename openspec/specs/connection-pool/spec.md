# connection-pool Specification

## Purpose
A pool of MQTT connections for services that issue many concurrent requests, reusing idle
connections and shedding surplus ones.

## Requirements
### Requirement: Borrowing reuses an idle connection or opens one

Borrowing SHALL hand out the first idle connection it finds when that one is connected, and otherwise open a new
connection from the template with client ID `<prefix>_<n>`. It SHALL fail when the pool already holds
its maximum (20 by default) or the new connection cannot connect.

#### Scenario: Idle connection available

- **WHEN** a connection has been returned and is still connected
- **THEN** the next borrow gets that connection without opening a new one

#### Scenario: Pool full

- **WHEN** the pool holds its maximum number of connections, none idle
- **THEN** borrowing fails with `too many connections`

### Requirement: Returning a connection clears its subscriptions

Returning a connection SHALL unsubscribe it from every topic and mark it idle.

#### Scenario: Returned connection

- **WHEN** a borrowed connection subscribed to a response topic is returned
- **THEN** it is idle and no longer subscribed to that topic

### Requirement: Surplus idle connections are closed

Every 10 s, while the pool holds more than its normal size (2 by default), it SHALL close connections
idle for longer than the maximum idle age (30 s by default) until it is back at its normal size.
Stopping the pool SHALL close every connection.

#### Scenario: Burst over

- **WHEN** a burst grew the pool to 6 connections and all have been idle for over 30 s
- **THEN** within the next check the pool shrinks to 2 connections
