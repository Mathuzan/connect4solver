# c4solver
Connect 4 Game sovler in Go.

## Building & Running
```bash
go build && ./c4solver
```

## Testing
```bash
go test -v
```

## Profiling
```bash
./c4solver -profile
go tool pprof -http=:8080 cpuprof.prof
```
