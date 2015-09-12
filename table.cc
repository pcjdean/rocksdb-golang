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
