## ADDED Requirements

### Requirement: All states include the last error

`Lifecycle.GetAllStates` SHALL return the text set with `SetLastError` as `LastErrorText`, and an empty text when none was set.

#### Scenario: Last error set

- **WHEN** `SetLastError("boom")` was called and `GetAllStates` is read
- **THEN** its `LastErrorText` is `boom`
