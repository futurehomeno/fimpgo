## ADDED Requirements

### Requirement: Library packages are covered by unit tests

Every library package except `cli` SHALL have unit tests exercising its exported behaviour, and the tests SHALL reach no network host other than loopback (an `httptest` server or the local test broker).

#### Scenario: Coverage run

- **WHEN** `go test -p 1 -cover ./...` runs on a clean checkout with the local test broker up
- **THEN** all tests pass
- **AND** total statement coverage is at least 60%
- **AND** the working tree has no modified or untracked files afterwards
