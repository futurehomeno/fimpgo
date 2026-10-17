## ADDED Requirements

### Requirement: Lifecycle state is safe for concurrent use

Lifecycle state changes, state reads, subscribing and unsubscribing SHALL be safe to call from any goroutine at the same time.

#### Scenario: Subscribers come and go while states change

- **WHEN** components subscribe and unsubscribe while the config, auth and connection states change
- **THEN** the app keeps running and every state change takes effect

### Requirement: Waiting for a state cannot miss it

Waiting for a state SHALL return once the state holds, including when it changes while the wait is starting.

#### Scenario: State reached while the wait starts

- **WHEN** the app state becomes `RUNNING` after a wait for `RUNNING` has read the old state but before it has subscribed
- **THEN** the wait returns

#### Scenario: State event dropped for a busy wait

- **WHEN** the state reaches `RUNNING` but its event is dropped because the wait's buffer is full
- **THEN** the wait still returns, after rechecking the state

#### Scenario: Wait shares an id with a subscriber

- **WHEN** a wait uses the same id as an existing subscriber
- **THEN** that subscriber keeps receiving events after the wait returns
