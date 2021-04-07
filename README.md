# c4solver
Connect 4 Game solver in Go. It finds winning strategy in "Connect Four" game (also known as "Four in a row")

## Building & Running
To build, install Go and run:
```bash
go mod download
./build.sh
```

Run training mode to evaluate cached results:
```bash
./build.sh && ./c4s --train --size 7x6
```

Start a game in playing mode:
```bash
./build.sh && ./c4s --play --size 7x6
```

See help for usage:
```bash
./c4s -h
```

## Testing
```bash
go test ./...
```

## Run benchmarks
```bash
go test --bench=. ./...
```

## Profiling
```bash
go build && ./c4solver --profile
go tool pprof -http=:8080 cpuprof.prof
```
