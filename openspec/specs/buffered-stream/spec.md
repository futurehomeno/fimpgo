# buffered-stream Specification

## Purpose
Batching FIMP messages in memory and flushing them as one JSON array to a file or a channel, for
apps that ship message history in bulk.

## Requirements
### Requirement: The buffer flushes on size or interval

A stream SHALL flush when it holds its maximum number of messages or when its interval elapses,
whichever comes first, and SHALL skip an empty flush. A zero size or interval SHALL not create a
stream. The buffer SHALL be emptied after each flush, also when serialization fails.

#### Scenario: Size reached

- **WHEN** a stream of size 10 receives its tenth message
- **THEN** all 10 are flushed at once

#### Scenario: Quiet period

- **WHEN** fewer messages than the size arrive within the interval
- **THEN** they are flushed when the interval elapses

### Requirement: Flushed batches are compact JSON

Each buffered message SHALL record its topic without the `pt:j1/mt:evt` or `pt:j1/mt:cmd` prefix,
object values SHALL be unpacked inline, and the batch SHALL be one JSON array, gzip-compressed when
the stream compresses. A file sink SHALL write `<prefix>_<RFC 3339 time>.json` (`.gz` when
compressed); a channel sink SHALL receive the bytes.

#### Scenario: Compressed file sink

- **WHEN** a compressing stream with file prefix `hist` flushes
- **THEN** a gzip file `hist_<time>.gz` holding the JSON array is written to the sink directory
