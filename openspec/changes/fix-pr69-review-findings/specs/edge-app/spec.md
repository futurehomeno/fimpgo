## ADDED Requirements

### Requirement: The hub token can be requested again

`FhOAuth2Client.LoadHubTokenFromCB` SHALL release its MQTT connection after each request so that a later call connects again.

#### Scenario: Request after CloudBridge did not answer

- **WHEN** `LoadHubTokenFromCB` got no answer from CloudBridge
- **AND** `LoadHubTokenFromCB` is called again
- **THEN** the second call connects, subscribes and sends the request again
