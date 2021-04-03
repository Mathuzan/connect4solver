package main

import (
	"io/ioutil"
	"os"

	pb "github.com/igrek51/connect4solver/proto"
	log "github.com/igrek51/log15"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

const outCacheFile = "cache.bin"

func SaveCache(cache *EndingCache) error {
	maxDepth := int(cache.maxCacheDepth / 2)

	log.Debug("Saving cache...", log.Ctx{
		"filename": outCacheFile,
		"entries":  cache.Size(),
		"maxDepth": maxDepth,
	})

	dephtCaches := cacheToProto(cache, maxDepth)

	outBytes, err := proto.Marshal(dephtCaches)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cache to proto")
	}
	if err := ioutil.WriteFile(outCacheFile, outBytes, 0644); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	log.Info("Cache saved", log.Ctx{
		"filename": outCacheFile,
		"entries":  cache.Size(),
		"maxDepth": maxDepth,
	})
	return nil
}

func LoadCache(board *Board) (*EndingCache, error) {
	in, err := ioutil.ReadFile(outCacheFile)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}
	dephtCaches := &pb.DepthCaches{}
	if err := proto.Unmarshal(in, dephtCaches); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal protobuf")
	}

	cache := protoToCache(dephtCaches, board.w, board.h)

	log.Debug("Cache loaded", log.Ctx{
		"filename": outCacheFile,
		"entries":  cache.Size(),
		"maxDepth": cache.HighestDepth(),
	})

	return cache, nil
}

func CacheFileExists() bool {
	_, err := os.Stat(outCacheFile)
	return err == nil
}

func cacheToProto(cache *EndingCache, maxDepth int) *pb.DepthCaches {
	dephtCaches := &pb.DepthCaches{
		DepthCaches: make([]*pb.DepthCache, len(cache.depthCaches)),
	}
	for d, depthCache := range cache.depthCaches {
		if d <= maxDepth {
			entriesMap := map[uint64]uint32{}
			for k, v := range depthCache {
				entriesMap[k] = uint32(v)
			}
			dephtCaches.DepthCaches[d] = &pb.DepthCache{
				Entries: entriesMap,
			}
		}
	}
	return dephtCaches
}

func protoToCache(dephtCaches *pb.DepthCaches, boardW int, boardH int) *EndingCache {
	cache := NewEndingCache(boardW, boardH)
	for d, depthCache := range dephtCaches.DepthCaches {
		for k, v := range depthCache.Entries {
			cache.depthCaches[d][k] = Player(v)
		}
		cache.cachedEntries += uint64(len(depthCache.Entries))
	}
	return cache
}
