## 1. Serialize Start and Stop

- [x] 1.1 Failing test: concurrent `Start` calls on a fresh transport all succeed and a later `Stop` returns
- [x] 1.2 Failing test: a `Start` issued while `Stop` waits for a busy reader leaves `Stop` and a later `Stop` returning
- [x] 1.3 One mutex serializes `Start` and `Stop`; `Stop` closes done before taking it; `Start` waits for the previous reader

## 2. Stop returns once disconnected

- [x] 2.1 Failing test: `Start` right after `Stop` with a retry count of 1 succeeds, 20 times in a row
- [x] 2.2 `Stop` disconnects with a 250 ms quiesce, so paho finishes the disconnect before `Stop` returns
