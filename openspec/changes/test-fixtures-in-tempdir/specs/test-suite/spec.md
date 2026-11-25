## ADDED Requirements

### Requirement: Test runs leave the working tree unchanged

The test suite SHALL write any files it creates to a temporary directory, so a test run never modifies or adds files in the repository.

#### Scenario: Config and key store round-trip

- **WHEN** the config save/load test and the key store save/load test run on a clean checkout
- **THEN** both tests pass
- **AND** the working tree has no modified or untracked files afterwards
