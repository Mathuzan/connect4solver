# c4solver
Connect 4 Game solver in Go. It finds winning strategy in "Connect Four" game (also known as "Four in a row")

## Building & Running
To build, install [Go](https://golang.org/doc/install) and run:
```bash
go mod download
./build.sh
```

Run training mode to evaluate cached results, precalculating every possible scenario:
```bash
./build.sh && ./c4solver --train --size 7x6
```
![](docs/train5x5.gif)

Start a game in playing mode:
```bash
./build.sh && ./c4solver --play --size 7x6
```

You can play against computer AI or analyze each player moves, showing best game endings for moves:
![](docs/play-6x5.gif)

See help for usage and possible options:
```console
$ ./c4solver --help
Usage of ./c4solver:
  -autoattack-a
    	Make player A move automatically
  -autoattack-b
    	Make player B move automatically
  -height int
    	board height (default 6)
  -hide-a
    	Hide endings hints for player A
  -hide-b
    	Hide endings hints for player B
  -nocache
    	Load cached endings from file
  -play
    	Playing mode
  -profile
    	Enable pprof CPU profiling
  -scores
    	Show scores of each move
  -size string
    	board size (eg. 7x6)
  -train
    	Training mode
  -width int
    	board width (default 7)
  -win int
    	win streak (default 4)
```

## Performance Tricks
Tricks used to improve performance:
- Representing board as binary number,
- Checking against winning condition using bitwise operators,
- Caching best game endings for later boards,
- Short-circuit when finding winning result,
- Winning strategy heuristics,
- Disregarding mirrored boards.

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
This produces the following CPU profiling graph, showing the places where CPU spends most of the time:  
![](docs/cpu-profiling-graph.png)
