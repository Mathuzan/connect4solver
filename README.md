# c4solver
`c4solver` is "Connect 4" Game solver written in Go. It finds winning strategy in "Connect Four" game (also known as "Four in a row"). The solver proves that on `7x6` board, first player has a winning strategy (can always win).

# Game Rules
> Connect Four is a two-player board game, in which the players take turns dropping colored discs into a seven-column, six-row vertically suspended grid. The pieces fall straight down, occupying the lowest available space within the column. The objective of the game is to be the first to form a horizontal, vertical, or diagonal line of four of one's own discs.

## Building
To build, install [Go](https://golang.org/doc/install) and run:
```bash
go mod download
go build -o c4solver
```

## Running
### Training mode
Run training mode to evaluate cached results, precalculating every possible scenario:
![](docs/train5x5.gif)

The following proves that on 7x6 board, first player has winning strategy:
```bash
./build.sh && ./c4solver --train --size 7x6
```
Precalculating every possible scenario and traversing the decision tree might take a long time on large boards for the first time. Cached results are stored in protobuf format and will be used when playing a game.

### Playing mode
Start a game in playing mode:
```bash
./build.sh && ./c4solver --play --size 7x6
```

You can play against computer AI or analyze each player moves, showing best game endings for moves (**W** - Win, **T** - Tie, **L** - Lose):
![](docs/play-6x5.gif)

If you want to challenge yourself versus Unbeatable "C4" AI, you can hide the move hints for yourself and enable automatic moves for computer player (Autoattack):
```bash
./c4solver --play --size 7x6 --autoattack-a --hide-b
```

## Help / Usage
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
- Representing whole board as binary number,
- Checking against winning condition using bitwise operators,
- Caching best game endings for later boards - different moves sequences lead to the same board,
- Disregarding mirrored boards - reflected boards can be treated as the same,
- Short-circuit when finding winning result,
- Winning strategy heuristics - start from middle moves,
- Checking only current player's move local neighbourhood when checking winning condition - don't need to check all rows & columns each time.

The optimized solver algorithm is able to consider over 4 millions boards per second, running on a regular laptop.
Still, it takes around week to solve `7x6` board since number of possible combinations is enormous.

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
./build.sh && ./c4solver --profile
go tool pprof -http=:8080 cpuprof.prof
```
This produces the following CPU profiling graph, showing the places where CPU spends most of the time:  
![](docs/cpu-profiling-graph.png)
