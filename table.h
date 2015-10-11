// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_TABLE_H_
#define GO_ROCKSDB_INCLUDE_TABLE_H_

#include "types.h"
#include "cache.h"
#include "filterPolicy.h"

#ifdef __cplusplus
typedef shared_ptr<TableFactory> PTableFactory;
#endif

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(BlockBasedTableOptions)
DEFINE_C_WRAP_STRUCT(PlainTableOptions)
DEFINE_C_WRAP_STRUCT(CuckooTableOptions)
DEFINE_C_WRAP_STRUCT(PTableFactory)

DEFINE_C_WRAP_CONSTRUCTOR_DEC(BlockBasedTableOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(BlockBasedTableOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(BlockBasedTableOptions)
// Setter methods
DEFINE_C_WRAP_SETTER_WRAP_DEC(BlockBasedTableOptions, block_cache, PCache)
DEFINE_C_WRAP_SETTER_WRAP_DEC(BlockBasedTableOptions, filter_policy, PFilterPolicy)
PTableFactory_t NewBlockBasedTableFactory(const BlockBasedTableOptions_t* table_options);

DEFINE_C_WRAP_CONSTRUCTOR_DEC(PlainTableOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PlainTableOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(PlainTableOptions)
// Setter methods for PlainTableOptions
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, user_key_len, uint32_t)
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, bloom_bits_per_key, int)
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, hash_table_ratio, double)
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, index_sparseness, size_t)
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, huge_page_tlb_size, size_t)
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, encoding_type, char)
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, full_scan_mode, bool)
DEFINE_C_WRAP_SETTER_DEC(PlainTableOptions, store_index_in_file, bool)
PTableFactory_t NewPlainTableFactory(const PlainTableOptions_t* table_options);

DEFINE_C_WRAP_CONSTRUCTOR_DEC(CuckooTableOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CuckooTableOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(CuckooTableOptions)
// Setter methods for CuckooTableOptions
// Determines the utilization of hash tables. Smaller values
// result in larger hash tables with fewer collisions.
DEFINE_C_WRAP_SETTER_DEC(CuckooTableOptions, hash_table_ratio, double)
// A property used by builder to determine the depth to go to
// to search for a path to displace elements in case of
// collision. See Builder.MakeSpaceForKey method. Higher
// values result in more efficient hash tables with fewer
// lookups but take more time to build.
DEFINE_C_WRAP_SETTER_DEC(CuckooTableOptions, max_search_depth, uint32_t)
// In case of collision while inserting, the builder
// attempts to insert in the next cuckoo_block_size
// locations before skipping over to the next Cuckoo hash
// function. This makes lookups more cache friendly in case
// of collisions.
DEFINE_C_WRAP_SETTER_DEC(CuckooTableOptions, cuckoo_block_size, uint32_t)
// If this option is enabled, user key is treated as uint64_t and its value
// is used as hash value directly. This option changes builder's behavior.
// Reader ignore this option and behave according to what specified in table
// property.
DEFINE_C_WRAP_SETTER_DEC(CuckooTableOptions, identity_as_first_hash, bool)
// If this option is set to true, module is used during hash calculation.
// This often yields better space efficiency at the cost of performance.
// If this optino is set to false, # of entries in table is constrained to be
// power of two, and bit and is used to calculate hash, which is faster in
// general.
DEFINE_C_WRAP_SETTER_DEC(CuckooTableOptions, use_module_hash, bool)
// Create a CuckooTableFactory
PTableFactory_t NewCuckooTableFactory(const CuckooTableOptions_t* table_options);

DEFINE_C_WRAP_CONSTRUCTOR_DEC(PTableFactory)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PTableFactory)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(PTableFactory)
#ifdef __cplusplus
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(PTableFactory, TableFactory*)
#endif

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_TABLE_H_
