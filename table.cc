// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Currently we support two types of tables: plain table and block-based table.
//   1. Block-based table: this is the default table type that we inherited from
//      LevelDB, which was designed for storing data in hard disk or flash
//      device.
//   2. Plain table: it is one of RocksDB's SST file format optimized
//      for low query latency on pure-memory or really low-latency media.
//
// A tutorial of rocksdb table formats is available here:
//   https://github.com/facebook/rocksdb/wiki/A-Tutorial-of-RocksDB-SST-formats
//
// Example code is also available
//   https://github.com/facebook/rocksdb/wiki/A-Tutorial-of-RocksDB-SST-formats#wiki-examples

#include <rocksdb/table.h>
#include "filterPolicyPrivate.h"
#include "table.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(BlockBasedTableOptions)
DEFINE_C_WRAP_DESTRUCTOR(BlockBasedTableOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(BlockBasedTableOptions)

// If non-NULL use the specified cache for blocks.
// If NULL, rocksdb will automatically create and use an 8MB internal cache.
DEFINE_C_WRAP_SETTER_WRAP(BlockBasedTableOptions, block_cache, PCache)

// If non-nullptr, use the specified filter policy to reduce disk reads.
// Many applications will benefit from passing the result of
// NewBloomFilterPolicy() here.
DEFINE_C_WRAP_SETTER_WRAP(BlockBasedTableOptions, filter_policy, PFilterPolicy)

// Create default block based table factory.
PTableFactory_t NewBlockBasedTableFactory(const BlockBasedTableOptions_t* table_options)
{
    return NewPTableFactoryTRawArgs(NewBlockBasedTableFactory((table_options && GET_REP(table_options, BlockBasedTableOptions)) ?
                                                      GET_REP_REF(table_options, BlockBasedTableOptions) :
                                                      BlockBasedTableOptions()));
}


DEFINE_C_WRAP_CONSTRUCTOR(PlainTableOptions)
DEFINE_C_WRAP_DESTRUCTOR(PlainTableOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(PlainTableOptions)

// @user_key_len: plain table has optimization for fix-sized keys, which can
//                be specified via user_key_len.  Alternatively, you can pass
//                `kPlainTableVariableLength` if your keys have variable
//                lengths.
DEFINE_C_WRAP_SETTER(PlainTableOptions, user_key_len, uint32_t)

// @bloom_bits_per_key: the number of bits used for bloom filer per prefix.
//                      You may disable it by passing a zero.
DEFINE_C_WRAP_SETTER(PlainTableOptions, bloom_bits_per_key, int)

// @hash_table_ratio: the desired utilization of the hash table used for
//                    prefix hashing.
//                    hash_table_ratio = number of prefixes / #buckets in the
//                    hash table
DEFINE_C_WRAP_SETTER(PlainTableOptions, hash_table_ratio, double)

// @index_sparseness: inside each prefix, need to build one index record for
//                    how many keys for binary search inside each hash bucket.
//                    For encoding type kPrefix, the value will be used when
//                    writing to determine an interval to rewrite the full
//                    key. It will also be used as a suggestion and satisfied
//                    when possible.
DEFINE_C_WRAP_SETTER(PlainTableOptions, index_sparseness, size_t)

// @huge_page_tlb_size: if <=0, allocate hash indexes and blooms from malloc.
//                      Otherwise from huge page TLB. The user needs to
//                      reserve huge pages for it to be allocated, like:
//                          sysctl -w vm.nr_hugepages=20
//                      See linux doc Documentation/vm/hugetlbpage.txt
DEFINE_C_WRAP_SETTER(PlainTableOptions, huge_page_tlb_size, size_t)

// @encoding_type: how to encode the keys. See enum EncodingType above for
//                 the choices. The value will determine how to encode keys
//                 when writing to a new SST file. This value will be stored
//                 inside the SST file which will be used when reading from
//                 the file, which makes it possible for users to choose
//                 different encoding type when reopening a DB. Files with
//                 different encoding types can co-exist in the same DB and
//                 can be read.
DEFINE_C_WRAP_SETTER_CAST(PlainTableOptions, encoding_type, char, EncodingType)

// @full_scan_mode: mode for reading the whole file one record by one without
//                  using the index.
DEFINE_C_WRAP_SETTER(PlainTableOptions, full_scan_mode, bool)

// @store_index_in_file: compute plain table index and bloom filter during
//                       file building and store it in file. When reading
//                       file, index will be mmaped instead of recomputation.
DEFINE_C_WRAP_SETTER(PlainTableOptions, store_index_in_file, bool)

// -- Plain Table with prefix-only seek
// For this factory, you need to set Options.prefix_extrator properly to make it
// work. Look-up will starts with prefix hash lookup for key prefix. Inside the
// hash bucket found, a binary search is executed for hash conflicts. Finally,
// a linear search is used.
PTableFactory_t NewPlainTableFactory(const PlainTableOptions_t* table_options)
{
    return NewPTableFactoryTRawArgs(NewPlainTableFactory((table_options && GET_REP(table_options, PlainTableOptions)) ?
                                                      GET_REP_REF(table_options, PlainTableOptions) :
                                                      PlainTableOptions()));
}

DEFINE_C_WRAP_CONSTRUCTOR(CuckooTableOptions)
DEFINE_C_WRAP_DESTRUCTOR(CuckooTableOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(CuckooTableOptions)

// Setter methods for CuckooTableOptions
// Determines the utilization of hash tables. Smaller values
// result in larger hash tables with fewer collisions.
DEFINE_C_WRAP_SETTER(CuckooTableOptions, hash_table_ratio, double)

// A property used by builder to determine the depth to go to
// to search for a path to displace elements in case of
// collision. See Builder.MakeSpaceForKey method. Higher
// values result in more efficient hash tables with fewer
// lookups but take more time to build.
DEFINE_C_WRAP_SETTER(CuckooTableOptions, max_search_depth, uint32_t)

// In case of collision while inserting, the builder
// attempts to insert in the next cuckoo_block_size
// locations before skipping over to the next Cuckoo hash
// function. This makes lookups more cache friendly in case
// of collisions.
DEFINE_C_WRAP_SETTER(CuckooTableOptions, cuckoo_block_size, uint32_t)

// If this option is enabled, user key is treated as uint64_t and its value
// is used as hash value directly. This option changes builder's behavior.
// Reader ignore this option and behave according to what specified in table
// property.
DEFINE_C_WRAP_SETTER(CuckooTableOptions, identity_as_first_hash, bool)

// If this option is set to true, module is used during hash calculation.
// This often yields better space efficiency at the cost of performance.
// If this optino is set to false, # of entries in table is constrained to be
// power of two, and bit and is used to calculate hash, which is faster in
// general.
DEFINE_C_WRAP_SETTER(CuckooTableOptions, use_module_hash, bool)

// Cuckoo Table Factory for SST table format using Cache Friendly Cuckoo Hashing
PTableFactory_t NewCuckooTableFactory(const CuckooTableOptions_t* table_options)
{
    return NewPTableFactoryTRawArgs(NewCuckooTableFactory((table_options && GET_REP(table_options, CuckooTableOptions)) ?
                                                      GET_REP_REF(table_options, CuckooTableOptions) :
                                                      CuckooTableOptions()));
}

DEFINE_C_WRAP_CONSTRUCTOR(PTableFactory)
DEFINE_C_WRAP_DESTRUCTOR(PTableFactory)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(PTableFactory, TableFactory*)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(PTableFactory)
