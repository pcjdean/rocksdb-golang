// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// +build !lite

package rocksdb

/*
#include "memtablerep.h"
*/
import "C"

// This creates MemTableReps that are backed by an std::vector. On iteration,
// the vector is sorted. This is useful for workloads where iteration is very
// rare and writes are generally not issued after reads begin.
//
// Parameters:
//   count uint64: Passed to the constructor of the underlying std::vector of each
//     VectorRep. On initialization, the underlying array will be at least count
//     bytes reserved for usage.
func NewVectorRepFactory (count uint64) (mtf *MemTableRepFactory) {
	var (
		ccnt C.size_t = C.size_t(count)
	)
	
	cmtf := C.PMemTableRepFactoryNewVectorRepFactory(ccnt)
	mtf = cmtf.toMemTableRepFactory()
	return
}

// This class contains a fixed array of buckets, each
// pointing to a skiplist (null if the bucket is empty).
// bucket_count uint64: number of fixed array buckets
// skiplist_height int32: the max height of the skiplist
// skiplist_branching_factor int32: probabilistic size ratio between adjacent
//                            link lists in the skiplist
func NewHashSkipListRepFactory (param ...interface{}) (mtf *MemTableRepFactory) {
	var (
		bcnt C.size_t = 1000000
		skh C.int32_t = 4
		skbf C.int32_t = 4
	)

	const (
		BUCKET_COUNT = iota
		SKIPLIST_HEIGHT
		SKIPLIST_BRANCHING_FACTOR
	)

	lp := len(param)

	if BUCKET_COUNT < lp {
		bcnt = C.size_t(param[BUCKET_COUNT].(uint64))
	} 

	if SKIPLIST_HEIGHT < lp {
		skh = C.int32_t(param[SKIPLIST_HEIGHT].(int32))
	} 

	if SKIPLIST_BRANCHING_FACTOR < lp {
		skbf = C.int32_t(param[SKIPLIST_BRANCHING_FACTOR].(int32))
	} 
	
	cmtf := C.PMemTableRepFactoryNewHashSkipListRepFactory(bcnt, skh, skbf)
	mtf = cmtf.toMemTableRepFactory()
	return
}

// The factory is to create memtables based on a hash table:
// it contains a fixed array of buckets, each pointing to either a linked list
// or a skip list if number of entries inside the bucket exceeds
// threshold_use_skiplist.
// @bucket_count uint64: number of fixed array buckets
// @huge_page_tlb_size uint64: if <=0, allocate the hash table bytes from malloc.
//                      Otherwise from huge page TLB. The user needs to reserve
//                      huge pages for it to be allocated, like:
//                          sysctl -w vm.nr_hugepages=20
//                      See linux doc Documentation/vm/hugetlbpage.txt
// @bucket_entries_logging_threshold int: if number of entries in one bucket
//                                    exceeds this number, log about it.
// @if_log_bucket_dist_when_flash bool: if true, log distribution of number of
//                                 entries when flushing.
// @threshold_use_skiplist uint32: a bucket switches to skip list if number of
//                          entries exceed this parameter.
func NewHashLinkListRepFactory (param ...interface{}) (mtf *MemTableRepFactory) {
	var (
		bcnt C.size_t = 50000
		hpts C.size_t = 0
		belt C.int = 4096
		ilbdwf C.bool = C.true
		tus C.uint32_t = 256
	)

	const (
		BUCKET_COUNT = iota
		HUGE_PAGE_TLB_SIZE
		BUCKET_ENTRIES_LOGGING_THRESHOLD
		IF_LOG_BUCKET_DIST_WHEN_FLASH
		THRESHOLD_USE_SKIPLIST
	)

	lp := len(param)

	if BUCKET_COUNT < lp {
		bcnt = C.size_t(param[BUCKET_COUNT].(uint64))
	} 

	if HUGE_PAGE_TLB_SIZE < lp {
		hpts = C.size_t(param[HUGE_PAGE_TLB_SIZE].(uint64))
	} 

	if BUCKET_ENTRIES_LOGGING_THRESHOLD < lp {
		belt = C.int(param[BUCKET_ENTRIES_LOGGING_THRESHOLD].(int))
	} 

	if IF_LOG_BUCKET_DIST_WHEN_FLASH < lp {
		ilbdwf = toCBool(param[IF_LOG_BUCKET_DIST_WHEN_FLASH].(bool))
	} 

	if THRESHOLD_USE_SKIPLIST < lp {
		tus = C.uint32_t(param[THRESHOLD_USE_SKIPLIST].(uint32))
	} 
	
	cmtf := C.PMemTableRepFactoryNewHashLinkListRepFactory(bcnt, hpts, belt, ilbdwf, tus)
	mtf = cmtf.toMemTableRepFactory()
	return
}

// This factory creates a cuckoo-hashing based mem-table representation.
// Cuckoo-hash is a closed-hash strategy, in which all key/value pairs
// are stored in the bucket array itself intead of in some data structures
// external to the bucket array.  In addition, each key in cuckoo hash
// has a constant number of possible buckets in the bucket array.  These
// two properties together makes cuckoo hash more memory efficient and
// a constant worst-case read time.  Cuckoo hash is best suitable for
// point-lookup workload.
//
// When inserting a key / value, it first checks whether one of its possible
// buckets is empty.  If so, the key / value will be inserted to that vacant
// bucket.  Otherwise, one of the keys originally stored in one of these
// possible buckets will be "kicked out" and move to one of its possible
// buckets (and possibly kicks out another victim.)  In the current
// implementation, such "kick-out" path is bounded.  If it cannot find a
// "kick-out" path for a specific key, this key will be stored in a backup
// structure, and the current memtable to be forced to immutable.
//
// Note that currently this mem-table representation does not support
// snapshot (i.e., it only queries latest state) and iterators.  In addition,
// MultiGet operation might also lose its atomicity due to the lack of
// snapshot support.
//
// Parameters:
//   write_buffer_size uint64: the write buffer size in bytes.
//   average_data_size uint64: the average size of key + value in bytes.  This value
//     together with write_buffer_size will be used to compute the number
//     of buckets.
//   hash_function_count uint: the number of hash functions that will be used by
//     the cuckoo-hash.  The number also equals to the number of possible
//     buckets each key will have.
func NewHashCuckooRepFactory (wbufsz uint64, param ...interface{}) (mtf *MemTableRepFactory) {
	var (
		ads C.size_t = 64
		hfc C.uint = 4
	)

	const (
		AVERAGE_DATA_SIZE = iota
		HASH_FUNCTION_COUNT
	)

	lp := len(param)

	if AVERAGE_DATA_SIZE < lp {
		ads = C.size_t(param[AVERAGE_DATA_SIZE].(uint64))
	} 

	if HASH_FUNCTION_COUNT < lp {
		hfc = C.uint(param[HASH_FUNCTION_COUNT].(uint))
	} 
	
	cmtf := C.PMemTableRepFactoryNewHashCuckooRepFactory(C.size_t(wbufsz), ads, hfc)
	mtf = cmtf.toMemTableRepFactory()
	return
}
