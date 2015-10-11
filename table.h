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
