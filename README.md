# c4solver
Connect 4 Game solver in Go. It finds winning strategy in "Connect Four" game (also known as "Four in a row")

## Building & Running
```bash
go build -o c4solver && ./c4solver --size 7x6
```

See help for usage:
```bash
./c4solver -h
```

## Testing
```bash
go test .
```

Run benchmarks
```bash
go test -bench=.
```

## Profiling
```bash
go build && ./c4solver -profile
go tool pprof -http=:8080 cpuprof.prof
```
