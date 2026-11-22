## 1. Skip a stuck channel while full

- [x] 1.1 Failing test: after a miss the next message returns at once, the healthy channel gets both; after one read the channel gets the next message, and the one after waits the timeout again
- [x] 1.2 Mark a channel stuck on timeout, skip it while full, clear the mark on the next delivery and on unregister
