# A stuck channel still delays every message by one receive timeout

PR #85 (review finding 2). After `share-channel-receive-timeout`, N stuck channels cost one receive timeout per message instead of N. One channel that is never read again still holds the dispatcher, and the channel registration lock, for the full timeout (10 s) on every message it is interested in.

## Reproduce

Register a channel that is never read and send more than one matching message per 10 s: the healthy channels get message k+1 only 10 s after message k, RegisterChannel and UnregisterChannel block for up to 10 s, the receive queue fills and on the 21st drop in a row the transport stops (`main msg queue stuck`).

## What would close it

A decision on the delivery rule, then a change:
- after a channel times out, send to it only without blocking (drop while full) until it drains once, or
- give each channel its own forwarding queue and goroutine, so one stuck consumer cannot slow the shared dispatcher.

Both drop messages for a stuck consumer sooner than today, which changes what consumers such as tpflow see.
