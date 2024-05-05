## Benchmark

### Single core

`go run cmd/bench/main.go -p 1 -n 4 -t 100`

| cpu       | gotruco | os           | go     | score   |
|-----------|---------|--------------|--------|---------|
| m2        | 0.2.1   | macOS 14.2   | 1.21   | 375,937 |
| i5-12600k | 0.2.1   | Ubuntu 22.04 | 1.18.3 | 279,112 |

### Multi core

`go run cmd/bench/main.go -p 4 -n 4 -t 100`

| cpu       | gotruco | os           | go     | score     |
|-----------|---------|--------------|--------|-----------|
| m2        | 0.2.1   | macOS 14.2   | 1.21   | 595,523   |
| i5-12600k | 0.2.1   | Ubuntu 22.04 | 1.18.3 | 1,137,748 |

### Max Multi core

Setting `-p` as the CPU's (max) number of (virtual) cores (i.e., "threads").

`go run cmd/bench/main.go -p ? -n 4 -t 100`

| cpu       | gotruco | os           | go     | `-p` | score     |
|-----------|---------|--------------|--------|------|-----------|
| m2        | 0.2.1   | macOS 14.2   | 1.21   | 8    | 741,297   |
| i5-12600k | 0.2.1   | Ubuntu 22.04 | 1.18.3 | 16   | 2,547,479 |