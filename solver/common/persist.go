package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	pb "github.com/igrek51/connect4solver/proto"
)

func SaveCache(cache ICache, boardW, boardH int) error {
	maxDepth := int(cache.MaxCachedDepth() / 2)
	filename := cacheFilename(boardW, boardH)

	log.Debug("Converting to protobuf struct...", log.Ctx{
		"filename": filename,
	})
	startTime := time.Now()
	protoCache, entriesLen := cacheToProto(cache, maxDepth, boardW, boardH)
	log.Debug("Marshalling protobuf...", log.Ctx{
		"splitTime": time.Since(startTime),
		"entries":   entriesLen,
	})
	outBytes, err := proto.Marshal(protoCache)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cache to proto")
	}
	log.Debug("Saving cache file...", log.Ctx{
		"splitTime": time.Since(startTime),
		"bytes":     len(outBytes),
	})
	if err := ioutil.WriteFile(filename, outBytes, 0644); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	log.Debug("Cache saved", log.Ctx{
		"filename":     filename,
		"savedEntries": entriesLen,
		"allEntries":   cache.Size(),
		"maxDepth":     maxDepth,
		"duration":     time.Since(startTime),
	})
	return nil
}

func LoadCache(cache ICache, boardW, boardH int) error {
	filename := cacheFilename(boardW, boardH)
	log.Debug("Loading cache file...", log.Ctx{
		"filename": filename,
	})
	startTime := time.Now()
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "error reading file")
	}
	log.Debug("Unmarshalling protobuf...", log.Ctx{
		"splitTime": time.Since(startTime),
		"bytes":     len(in),
	})
	dephtCaches := &pb.DepthCaches{}
	if err := proto.Unmarshal(in, dephtCaches); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf")
	}
	log.Debug("Converting protobuf struct...", log.Ctx{
		"splitTime": time.Since(startTime),
	})

	protoToCache(dephtCaches, cache, boardW, boardH)

	log.Debug("Cache loaded", log.Ctx{
		"filename": filename,
		"entries":  cache.Size(),
		"duration": time.Since(startTime),
	})
	return nil
}

func MustSaveCache(cache ICache, boardW, boardH int) {
	err := SaveCache(cache, boardW, boardH)
	if err != nil {
		panic(errors.Wrap(err, "saving cache"))
	}
}

func MustLoadCache(cache ICache, boardW, boardH int) {
	err := LoadCache(cache, boardW, boardH)
	if err != nil {
		panic(errors.Wrap(err, "loading cache"))
	}
}

func CacheFileExists(board *Board) bool {
	filename := cacheFilename(board.W, board.H)
	_, err := os.Stat(filename)
	return err == nil
}

func cacheToProto(cache ICache, maxDepth int, boardW int, boardH int) (*pb.DepthCaches, uint64) {
	dephtCaches := &pb.DepthCaches{
		DepthCaches: make([]*pb.DepthCache, boardW*boardH),
	}
	entriesLen := uint64(0)
	for d, depthCache := range cache.DepthCaches() {
		if d <= maxDepth {
			boardsPlayerA := []uint64{}
			boardsPlayerB := []uint64{}
			boardsTie := []uint64{}
			for k, v := range depthCache {
				if v == PlayerA {
					boardsPlayerA = append(boardsPlayerA, k)
				} else if v == PlayerB {
					boardsPlayerB = append(boardsPlayerB, k)
				} else if v == Empty {
					boardsTie = append(boardsTie, k)
				}
			}
			dephtCaches.DepthCaches[d] = &pb.DepthCache{
				BoardsPlayerA: boardsPlayerA,
				BoardsPlayerB: boardsPlayerB,
				BoardsTie:     boardsTie,
			}
			entriesLen += uint64(len(depthCache))
		}
	}
	return dephtCaches, entriesLen
}

func protoToCache(dephtCaches *pb.DepthCaches, cache ICache, boardW int, boardH int) {
	for d, depthCache := range dephtCaches.DepthCaches {
		for _, k := range depthCache.BoardsPlayerA {
			cache.SetEntry(d, k, PlayerA)
		}
		for _, k := range depthCache.BoardsPlayerB {
			cache.SetEntry(d, k, PlayerB)
		}
		for _, k := range depthCache.BoardsTie {
			cache.SetEntry(d, k, Empty)
		}
	}
}

func cacheFilename(boardW, boardH int) string {
	return fmt.Sprintf("cache/cache_%dx%d.protobuf", boardW, boardH)
}
