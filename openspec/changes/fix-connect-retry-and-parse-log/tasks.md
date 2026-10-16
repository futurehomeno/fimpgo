## 1. Connect retry

- [x] 1.1 Failing test: an unreachable broker leaves no "status can only transition" warning, and Start returns a timeout with the client disconnected
- [x] 1.2 Start reuses the pending connect attempt and disconnects it when giving up

## 2. Parse error log

- [x] 2.1 Failing test: an unparsable payload logs its length and top-level keys
- [x] 2.2 Add length and keys to the error line
