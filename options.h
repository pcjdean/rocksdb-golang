// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_OPTIONS_H_
#define GO_ROCKSDB_INCLUDE_OPTIONS_H_

#ifdef __cplusplus
#include <rocksdb/options.h>
using namespace rocksdb;
#endif

#include "types.h"
#include "snapshot.h"
#include "table.h"
#include "comparator.h"
#include "compactionfilter.h"
#include "mergeOperator.h"
#include "sliceTransform.h"
#include "memtablerep.h"
#include "env.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(ColumnFamilyOptions)
DEFINE_C_WRAP_STRUCT(DBOptions)
DEFINE_C_WRAP_STRUCT(Options)
DEFINE_C_WRAP_STRUCT(ReadOptions)
DEFINE_C_WRAP_STRUCT(WriteOptions)
DEFINE_C_WRAP_STRUCT(FlushOptions)
DEFINE_C_WRAP_STRUCT(CompactionOptions)
DEFINE_C_WRAP_STRUCT(CompactRangeOptions)

// Cast Options* to DBOptions*
DEFINE_C_WRAP_STATIC_CAST_DEC(Options, DBOptions)
// Cast Options* to ColumnFamilyOptions*
DEFINE_C_WRAP_STATIC_CAST_DEC(Options, ColumnFamilyOptions)


DEFINE_C_WRAP_CONSTRUCTOR_DEC(ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(ColumnFamilyOptions, Options)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(ColumnFamilyOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(ColumnFamilyOptions)
// Get/Set methods
DEFINE_C_WRAP_GETTER_DEC(ColumnFamilyOptions, compression, int)
DEFINE_C_WRAP_SETTER_DEC(ColumnFamilyOptions, compression, int)
DEFINE_C_WRAP_GETTER_DEC(ColumnFamilyOptions, write_buffer_size, size_t)
DEFINE_C_WRAP_SETTER_DEC(ColumnFamilyOptions, write_buffer_size, size_t)
// Set method for memtable factory
DEFINE_C_WRAP_SETTER_WRAP_DEC(ColumnFamilyOptions, memtable_factory, PMemTableRepFactory)
// Set method for table factory
DEFINE_C_WRAP_SETTER_WRAP_DEC(ColumnFamilyOptions, table_factory, PTableFactory)
// Set method for merge operator.
DEFINE_C_WRAP_SETTER_WRAP_DEC(ColumnFamilyOptions, merge_operator, PMergeOperator)
// Set method for prefix extractor.
DEFINE_C_WRAP_SETTER_WRAP_DEC(ColumnFamilyOptions, prefix_extractor, PConstSliceTransform);
// Get/Set methods for comparator
DEFINE_C_WRAP_SETTER_WRAP_DEC(ColumnFamilyOptions, comparator, Comparator)
// Get/Set methods for compaction filter
DEFINE_C_WRAP_SETTER_WRAP_DEC(ColumnFamilyOptions, compaction_filter, CompactionFilter)
DEFINE_C_WRAP_SETTER_WRAP_DEC(ColumnFamilyOptions, compaction_filter_factory, PCompactionFilterFactory)
void ColumnFamilyOptions_set_compression_per_level(ColumnFamilyOptions_t* opt,
                                                   int* level_values,
                                                   size_t num_levels);
void ColumnFamilyOptions_set_compression_options(
    ColumnFamilyOptions_t* opt, int w_bits, int level, int strategy);


DEFINE_C_WRAP_CONSTRUCTOR_DEC(DBOptions)
#ifdef __cplusplus
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(DBOptions, const Options&)
#endif
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(DBOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(DBOptions)
// Get/Set methods
DEFINE_C_WRAP_GETTER_DEC(DBOptions, create_if_missing, bool)
DEFINE_C_WRAP_SETTER_DEC(DBOptions, create_if_missing, bool)
DEFINE_C_WRAP_GETTER_DEC(DBOptions, error_if_exists, bool)
DEFINE_C_WRAP_SETTER_DEC(DBOptions, error_if_exists, bool)
// Setter method for mmap reads
DEFINE_C_WRAP_SETTER_DEC(DBOptions, allow_mmap_reads, bool)
// Get/Set methods for @env
DEFINE_C_WRAP_SETTER_WRAP_DEC(DBOptions, env, Env)
DEFINE_C_WRAP_GETTER_WRAP_DEC(DBOptions, env, Env)
// Get/Set methods for @info_log
DEFINE_C_WRAP_SETTER_WRAP_DEC(DBOptions, info_log, PLogger)
// Get/Set methods for @paranoid_checks
DEFINE_C_WRAP_SETTER_DEC(DBOptions, paranoid_checks, bool)
// Get/Set methods for @max_open_files
DEFINE_C_WRAP_SETTER_DEC(DBOptions, max_open_files, int)


DEFINE_C_WRAP_CONSTRUCTOR_DEC(Options)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(Options, DBOptions, ColumnFamilyOptions)
#ifdef __cplusplus
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(Options, const DBOptions&, const ColumnFamilyOptions&)
#endif
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(Options)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Options)


DEFINE_C_WRAP_CONSTRUCTOR_DEC(ReadOptions)
#ifdef __cplusplus
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(ReadOptions, bool, bool)
#endif
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(ReadOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(ReadOptions)
void ReadOptions_set_snapshot(ReadOptions_t* opt, const Snapshot_t* snap);
// Get/Set methods
// Set method for @iterate_upper_bound
DEFINE_C_WRAP_SETTER_WRAP_DEC(ReadOptions, iterate_upper_bound, Slice)
// Get/Set methods for @verify_checksums
DEFINE_C_WRAP_SETTER_DEC(ReadOptions, verify_checksums, bool)
// Get/Set methods for @fill_cache
DEFINE_C_WRAP_SETTER_DEC(ReadOptions, fill_cache, bool)



DEFINE_C_WRAP_CONSTRUCTOR_DEC(WriteOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(WriteOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(WriteOptions)
// Get/Set methods
DEFINE_C_WRAP_GETTER_DEC(WriteOptions, sync, bool)
DEFINE_C_WRAP_SETTER_DEC(WriteOptions, sync, bool)


DEFINE_C_WRAP_CONSTRUCTOR_DEC(FlushOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(FlushOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(FlushOptions)


DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(CompactionOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionOptions)


DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactRangeOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(CompactRangeOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactRangeOptions)


#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_OPTIONS_H_
