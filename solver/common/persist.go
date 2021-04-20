package common

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	pb "github.com/igrek51/connect4solver/proto"
)

func SaveCache(cache ICache, boardW, boardH int) error {
	maxDepth := int(cache.MaxCachedDepth() / 2)
	filename := cacheFilename(boardW, boardH)

	log.Debug("Saving cache...", log.Ctx{
		"filename": filename,
	})
	protoCache, entriesLen := cacheToProto(cache, maxDepth, boardW, boardH)
	outBytes, err := proto.Marshal(protoCache)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cache to proto")
	}
	if err := ioutil.WriteFile(filename, outBytes, 0644); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	log.Debug("Cache saved", log.Ctx{
		"filename":     filename,
		"savedEntries": entriesLen,
		"allEntries":   cache.Size(),
		"maxDepth":     maxDepth,
	})

	return nil
}

func LoadCache(cache ICache, boardW, boardH int) error {
	filename := cacheFilename(boardW, boardH)
	log.Debug("Loading cache...", log.Ctx{
		"filename": filename,
	})
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "error reading file")
	}
	dephtCaches := &pb.DepthCaches{}
	if err := proto.Unmarshal(in, dephtCaches); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf")
	}

	protoToCache(dephtCaches, cache, boardW, boardH)

	log.Debug("Cache loaded", log.Ctx{
		"filename": filename,
		"entries":  cache.Size(),
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
			entriesMap := map[uint64]uint32{}
			for k, v := range depthCache {
				entriesMap[k] = uint32(v)
			}
			dephtCaches.DepthCaches[d] = &pb.DepthCache{
				Entries: entriesMap,
			}
			entriesLen += uint64(len(depthCache))
		}
	}
	return dephtCaches, entriesLen
}

func protoToCache(dephtCaches *pb.DepthCaches, cache ICache, boardW int, boardH int) {
	for d, depthCache := range dephtCaches.DepthCaches {
		for k, v := range depthCache.Entries {
			cache.SetEntry(d, k, Player(v))
		}
	}
}

func cacheFilename(boardW, boardH int) string {
	return fmt.Sprintf("cache/cache_%dx%d.protobuf", boardW, boardH)
}
