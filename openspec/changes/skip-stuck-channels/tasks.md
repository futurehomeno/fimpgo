## 1. Skip a stuck channel while full

- [x] 1.1 Failing test: after a miss the next message returns at once, the healthy channel gets both; after one read the channel gets the next message, and the one after waits the timeout again
- [x] 1.2 Mark a channel stuck on timeout, skip it while full, clear the mark on the next delivery and on unregister
- [x] 1.3 Failing test, then clear the mark on register too: a channel replacing a stuck one under the same name is not skipped (review of PR 86)
