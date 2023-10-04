package TinyCache

import (
	"TinyCache/util"
	"sync"
)

const (
	kNumShardBits = 4
	kNumShard     = 1 << kNumShardBits
)

type sharedCache struct {
	mu         []sync.Mutex
	caches     []*cache
	cacheBytes uint32
}

func HashKey(key string) uint32 {
	return util.Hash([]byte(key), 0)
}

func Shard(hash uint32) uint32 {
	return hash >> (32 - kNumShardBits)
}

func (sc *sharedCache) add(key string, value ByteView) {
	if sc.caches == nil {
		sc.caches = make([]*cache, kNumShard)
		perShard := (sc.cacheBytes + kNumShard - 1) / kNumShard
		for i := 0; i < kNumShard; i++ {
			sc.caches[i] = &cache{cacheBytes: perShard}
		}
	}
	hash := HashKey(key)
	shared := Shard(hash)
	sc.caches[shared].add(key, value)
}

func (sc *sharedCache) get(key string) (value ByteView, ok bool) {
	if sc.caches == nil {
		return
	}
	hash := HashKey(key)
	shared := Shard(hash)
	if v, ok := sc.caches[shared].get(key); ok {
		return v, true
	}
	return
}
