## 1. One receive timeout per message

- [x] 1.1 Failing test: three stuck channels, one healthy and one late reader with a 2 s timeout; the message is delivered in under 3.5 s, the healthy and the late reader get it
- [x] 1.2 Failing test: a full channel closed during the wait does not crash the process and the registration lock is free afterwards
- [x] 1.3 Non-blocking send first, then one select over the full channels and one timer in the receiving goroutine
- [x] 1.4 Failing test: two channels getting one message receive distinct `Message` values; one is allocated per channel
- [x] 1.5 Failing test: a timeout changed during the wait does not change the timeout named in the warning
