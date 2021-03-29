# c4solver
Connect 4 Game solver in Go.

## Building & Running
```bash
go build && ./c4solver
```

See help for usage:
```bash
./c4solver -h
```

## Testing
```bash
go test -v
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
