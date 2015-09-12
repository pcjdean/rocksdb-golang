// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// A Cache is an interface that maps keys to values.  It has internal
// synchronization and may be safely accessed concurrently from
// multiple threads.  It may automatically evict entries to make room
// for new entries.  Values have a specified charge against the cache
// capacity.  For example, a cache where the values are variable
// length strings, may use the length of the string as the charge for
// the string.
//
// A builtin cache implementation with a least-recently-used eviction
// policy is provided.  Clients may use their own implementations if
// they want something more sophisticated (like scan-resistance, a
// custom eviction policy, variable cache sizing, etc.)

package rocksdb

/*
#include "cache.h"
*/
import "C"

import (
	"runtime"
)

// Wrap go cache
type Cache struct {
	cache C.PCache_t
}

// Release resources
func (cache *Cache) finalize() {
	var ccache *C.PCache_t= &cache.cache
	C.DeletePCacheT(ccache, toCBool(false))
}

// C cache to go cache
func (ccache *C.PCache_t) toCache() (cache *Cache) {
	cache = &Cache{cache: *ccache}	
	runtime.SetFinalizer(cache, finalize)
	return
}

// Create a new cache with a fixed size capacity. The cache is sharded
// to 2^numShardBits shards, by hash of the key. The total capacity
// is divided and evenly assigned to each shard.
//
// The functions without parameter numShardBits uses default value, which is 4
func NewLRUCache(capacity uint64, numShardBits ...int) *Cache {
	var cnumshbits C.int = 4

	if numShardBits != nil {
		cnumshbits = C.int(numShardBits[0])
	}

	ccache := C.NewPCacheTRawArgs(C.size_t(capacity), cnumshbits)
	return ccache.toCache()
}
