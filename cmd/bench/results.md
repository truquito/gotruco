## Benchmark

### Single core

`go run cmd/bench/main.go -p 1 -n 4 -t 100`

| cpu | os    | version | total   |
|-----|-------|---------|---------|
| m2  | macOS | 14.2    | 375,937 |

### Multi core

`go run cmd/bench/main.go -p 4 -n 4 -t 100`

| cpu | os    | version | total   |
|-----|-------|---------|---------|
| m2  | macOS | 14.2    | 595,523 |