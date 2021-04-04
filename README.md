# c4solver
Connect 4 Game solver in Go. It finds winning strategy in "Connect Four" game (also known as "Four in a row")

## Building & Running
To build, run:
```bash
go build -o c4solver
```
or run `./build.sh`

Run training mode:
```bash
./build.sh && ./c4solver --train --size 7x6
```

Run playing a game mode:
```bash
./build.sh && ./c4solver --play --size 7x6
```

See help for usage:
```bash
./c4solver -h
```

## Testing
```bash
go test .
```

## Run benchmarks
```bash
go test --bench=.
```

## Profiling
```bash
go build && ./c4solver --profile
go tool pprof -http=:8080 cpuprof.prof
```
