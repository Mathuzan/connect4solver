package main

import (
	"io/ioutil"

	pb "github.com/igrek51/connect4solver/proto"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func SaveCache(cache *EndingCache) error {
	dephtCaches := cacheToProto(cache)

	outBytes, err := proto.Marshal(dephtCaches)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cache to proto")
	}
	if err := ioutil.WriteFile("cache.bin", outBytes, 0644); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	return nil
}

func cacheToProto(cache *EndingCache) *pb.DepthCaches {
	dephtCaches := &pb.DepthCaches{
		DepthCaches: make([]*pb.DepthCache, len(cache.depthCaches)),
	}
	for d, depthCache := range cache.depthCaches {
		entriesMap := map[uint64]uint32{}
		for k, v := range depthCache {
			entriesMap[k] = uint32(v)
		}
		dephtCaches.DepthCaches[d] = &pb.DepthCache{
			Entries: entriesMap,
		}
	}
	return dephtCaches
}
