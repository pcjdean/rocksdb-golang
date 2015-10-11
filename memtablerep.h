// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_MEMTABLEREP_H_
#define GO_ROCKSDB_INCLUDE_MEMTABLEREP_H_

#include "types.h"

#ifdef __cplusplus
#include <memory>
typedef std::shared_ptr<MemTableRepFactory> PMemTableRepFactory;
#endif

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(PMemTableRepFactory)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(PMemTableRepFactory)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PMemTableRepFactory)
#ifdef __cplusplus
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(PMemTableRepFactory, MemTableRepFactory*)
#endif

// This uses a skip list to store keys. It is the default.
//
// Parameters:
//   lookahead: If non-zero, each iterator's seek operation will start the
//     search from the previously visited record (doing at most 'lookahead'
//     steps). This is an optimization for the access pattern including many
//     seeks with consecutive keys.
PMemTableRepFactory_t PMemTableRepFactoryNewSkipListFactory(size_t lookahead);

#ifndef ROCKSDB_LITE
// This creates MemTableReps that are backed by an std::vector. On iteration,
// the vector is sorted. This is useful for workloads where iteration is very
// rare and writes are generally not issued after reads begin.
//
// Parameters:
//   count: Passed to the constructor of the underlying std::vector of each
//     VectorRep. On initialization, the underlying array will be at least count
//     bytes reserved for usage.
PMemTableRepFactory_t PMemTableRepFactoryNewVectorRepFactory(size_t count);

// This class contains a fixed array of buckets, each
// pointing to a skiplist (null if the bucket is empty).
// bucket_count: number of fixed array buckets
// skiplist_height: the max height of the skiplist
// skiplist_branching_factor: probabilistic size ratio between adjacent
//                            link lists in the skiplist
PMemTableRepFactory_t PMemTableRepFactoryNewHashSkipListRepFactory(
    size_t bucket_count, int32_t skiplist_height,
    int32_t skiplist_branching_factor
);

// The factory is to create memtables based on a hash table:
// it contains a fixed array of buckets, each pointing to either a linked list
// or a skip list if number of entries inside the bucket exceeds
// threshold_use_skiplist.
// @bucket_count: number of fixed array buckets
// @huge_page_tlb_size: if <=0, allocate the hash table bytes from malloc.
//                      Otherwise from huge page TLB. The user needs to reserve
//                      huge pages for it to be allocated, like:
//                          sysctl -w vm.nr_hugepages=20
//                      See linux doc Documentation/vm/hugetlbpage.txt
// @bucket_entries_logging_threshold: if number of entries in one bucket
//                                    exceeds this number, log about it.
// @if_log_bucket_dist_when_flash: if true, log distribution of number of
//                                 entries when flushing.
// @threshold_use_skiplist: a bucket switches to skip list if number of
//                          entries exceed this parameter.
PMemTableRepFactory_t PMemTableRepFactoryNewHashLinkListRepFactory(
    size_t bucket_count, size_t huge_page_tlb_size,
    int bucket_entries_logging_threshold,
    bool if_log_bucket_dist_when_flash,
    uint32_t threshold_use_skiplist);

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
//   write_buffer_size: the write buffer size in bytes.
//   average_data_size: the average size of key + value in bytes.  This value
//     together with write_buffer_size will be used to compute the number
//     of buckets.
//   hash_function_count: the number of hash functions that will be used by
//     the cuckoo-hash.  The number also equals to the number of possible
//     buckets each key will have.
PMemTableRepFactory_t PMemTableRepFactoryNewHashCuckooRepFactory(
    size_t write_buffer_size, size_t average_data_size,
    unsigned int hash_function_count);
#endif  // ROCKSDB_LITE

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_MEMTABLEREP_H_
