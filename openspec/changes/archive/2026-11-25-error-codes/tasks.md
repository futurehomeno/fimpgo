## 1. Error codes

- [x] 1.1 Failing test: each constant's wire value and `Str()` (one table test; it also fails to compile until the type exists)
- [x] 1.2 `fimptype/error_code.go`: `ErrorCodeT`, `Str()`, the ten constants
- [x] 1.3 `go vet ./...`, `golangci-lint run`, `make test` green
- [ ] 1.4 Release, so cliffhanger change `error-report-policy` can take it
