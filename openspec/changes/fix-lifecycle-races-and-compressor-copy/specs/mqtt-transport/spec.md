## ADDED Requirements

### Requirement: Compressed payloads stay intact

A compressed payload SHALL stay unchanged after it is returned, whatever the compressor encodes next.

#### Scenario: Back-to-back compressed publishes

- **WHEN** two messages are compressed one after the other with the same compressor
- **THEN** the first payload still decompresses to the first message
