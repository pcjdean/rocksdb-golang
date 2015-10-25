// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include "compactionfilterPrivate.h"
#include "mergeOperatorPrivate.h"
#include "sliceTransformPrivate.h"
#include "options.h"

DEFINE_C_WRAP_STATIC_CAST(Options, DBOptions)
DEFINE_C_WRAP_STATIC_CAST(Options, ColumnFamilyOptions)

DEFINE_C_WRAP_CONSTRUCTOR(ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(ColumnFamilyOptions, Options)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(ColumnFamilyOptions)
DEFINE_C_WRAP_DESTRUCTOR(ColumnFamilyOptions)
// Get/Set methods
DEFINE_C_WRAP_GETTER(ColumnFamilyOptions, compression, int)
DEFINE_C_WRAP_SETTER_CAST(ColumnFamilyOptions, compression, int, CompressionType)
DEFINE_C_WRAP_GETTER(ColumnFamilyOptions, write_buffer_size, size_t)
DEFINE_C_WRAP_SETTER(ColumnFamilyOptions, write_buffer_size, size_t)

// This is a factory that provides MemTableRep objects.
// Default: a factory that provides a skip-list-based implementation of
// MemTableRep.
DEFINE_C_WRAP_SETTER_WRAP(ColumnFamilyOptions, memtable_factory, PMemTableRepFactory)

// This is a factory that provides TableFactory objects.
// Default: a block-based table factory that provides a default
// implementation of TableBuilder and TableReader with default
// BlockBasedTableOptions.
DEFINE_C_WRAP_SETTER_WRAP(ColumnFamilyOptions, table_factory, PTableFactory)

// REQUIRES: The client must provide a merge operator if Merge operation
// needs to be accessed. Calling Merge on a DB without a merge operator
// would result in Status::NotSupported. The client must ensure that the
// merge operator supplied here has the same name and *exactly* the same
// semantics as the merge operator provided to previous open calls on
// the same DB. The only exception is reserved for upgrade, where a DB
// previously without a merge operator is introduced to Merge operation
// for the first time. It's necessary to specify a merge operator when
// openning the DB in this case.
// Default: nullptr
DEFINE_C_WRAP_SETTER_WRAP(ColumnFamilyOptions, merge_operator, PMergeOperator)

// If non-nullptr, use the specified function to determine the
// prefixes for keys.  These prefixes will be placed in the filter.
// Depending on the workload, this can reduce the number of read-IOP
// cost for scans when a prefix is passed via ReadOptions to
// db.NewIterator().  For prefix filtering to work properly,
// "prefix_extractor" and "comparator" must be such that the following
// properties hold:
//
// 1) key.starts_with(prefix(key))
// 2) Compare(prefix(key), key) <= 0.
// 3) If Compare(k1, k2) <= 0, then Compare(prefix(k1), prefix(k2)) <= 0
// 4) prefix(prefix(key)) == prefix(key)
//
// Default: nullptr
DEFINE_C_WRAP_SETTER_WRAP(ColumnFamilyOptions, prefix_extractor, PConstSliceTransform)

void ColumnFamilyOptions_set_compression_per_level(ColumnFamilyOptions_t* opt,
                                                   int* level_values,
                                                   size_t num_levels)
{
    assert(opt != NULL);
    assert(GET_REP(opt, ColumnFamilyOptions) != NULL);
    GET_REP(opt, ColumnFamilyOptions)->compression_per_level.resize(num_levels);
    for (uint64_t i = 0; i < num_levels; ++i) {
        GET_REP(opt, ColumnFamilyOptions)->compression_per_level[i] =
            static_cast<CompressionType>(level_values[i]);
    }
}

void ColumnFamilyOptions_set_compression_options(
    ColumnFamilyOptions_t* opt, int w_bits, int level, int strategy)
{
    assert(opt != NULL);
    assert(GET_REP(opt, ColumnFamilyOptions) != NULL);
    GET_REP(opt, ColumnFamilyOptions)->compression_opts.window_bits = w_bits;
    GET_REP(opt, ColumnFamilyOptions)->compression_opts.level = level;
    GET_REP(opt, ColumnFamilyOptions)->compression_opts.strategy = strategy;
}

// -------------------
// Parameters that affect behavior

// Comparator used to define the order of keys in the table.
// Default: a comparator that uses lexicographic byte-wise ordering
//
// REQUIRES: The client must ensure that the comparator supplied
// here has the same name and orders keys *exactly* the same as the
// comparator provided to previous open calls on the same DB.
DEFINE_C_WRAP_SETTER_PTR_WRAP(ColumnFamilyOptions, comparator, Comparator)

// A single CompactionFilter instance to call into during compaction.
// Allows an application to modify/delete a key-value during background
// compaction.
//
// If the client requires a new compaction filter to be used for different
// compaction runs, it can specify compaction_filter_factory instead of this
// option.  The client should specify only one of the two.
// compaction_filter takes precedence over compaction_filter_factory if
// client specifies both.
//
// If multithreaded compaction is being used, the supplied CompactionFilter
// instance may be used from different threads concurrently and so should be
// thread-safe.
//
// Default: nullptr
DEFINE_C_WRAP_SETTER_PTR_WRAP(ColumnFamilyOptions, compaction_filter, CompactionFilter)

// This is a factory that provides compaction filter objects which allow
// an application to modify/delete a key-value during background compaction.
//
// A new filter will be created on each compaction run.  If multithreaded
// compaction is being used, each created CompactionFilter will only be used
// from a single thread and so does not need to be thread-safe.
//
// Default: a factory that doesn't provide any object
DEFINE_C_WRAP_SETTER_WRAP(ColumnFamilyOptions, compaction_filter_factory, PCompactionFilterFactory)

// Version TWO of the compaction_filter_factory
// It supports rolling compaction
//
// Default: a factory that doesn't provide any object
DEFINE_C_WRAP_SETTER_WRAP(ColumnFamilyOptions, compaction_filter_factory_v2, PCompactionFilterFactoryV2)

DEFINE_C_WRAP_CONSTRUCTOR(DBOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(DBOptions, const Options&)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(DBOptions)
DEFINE_C_WRAP_DESTRUCTOR(DBOptions)
DEFINE_C_WRAP_GETTER(DBOptions, create_if_missing, bool)
DEFINE_C_WRAP_SETTER(DBOptions, create_if_missing, bool)
DEFINE_C_WRAP_GETTER(DBOptions, error_if_exists, bool)
DEFINE_C_WRAP_SETTER(DBOptions, error_if_exists, bool)
// Allow the OS to mmap file for reading sst tables. Default: false
DEFINE_C_WRAP_SETTER(DBOptions, allow_mmap_reads, bool)
// Use the specified object to interact with the environment,
// e.g. to read/write files, schedule background work, etc.
// Default: Env::Default()
DEFINE_C_WRAP_SETTER_PTR_WRAP(DBOptions, env, Env)
DEFINE_C_WRAP_GETTER_WRAP(DBOptions, env, Env)


DEFINE_C_WRAP_CONSTRUCTOR(Options)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(Options, DBOptions, ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(Options, const DBOptions&, const ColumnFamilyOptions&)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(Options)
DEFINE_C_WRAP_DESTRUCTOR(Options)


DEFINE_C_WRAP_CONSTRUCTOR(ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(ReadOptions, bool, bool)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(ReadOptions)
DEFINE_C_WRAP_DESTRUCTOR(ReadOptions)

// If "snapshot" is non-nullptr, read as of the supplied snapshot
// (which must belong to the DB that is being read and which must
// not have been released).  If "snapshot" is nullptr, use an impliicit
// snapshot of the state at the beginning of this read operation.
// Default: nullptr
void ReadOptions_set_snapshot(ReadOptions_t* opt, const Snapshot_t* snap)
{
    assert(opt != NULL);
    assert(GET_REP(opt, ReadOptions) != NULL);
    GET_REP(opt, ReadOptions)->snapshot = (snap ? GET_REP(snap, Snapshot) : nullptr);
}

// "iterate_upper_bound" defines the extent upto which the forward iterator
// can returns entries. Once the bound is reached, Valid() will be false.
// "iterate_upper_bound" is exclusive ie the bound value is
// not a valid entry.  If iterator_extractor is not null, the Seek target
// and iterator_upper_bound need to have the same prefix.
// This is because ordering is not guaranteed outside of prefix domain.
// There is no lower bound on the iterator. If needed, that can be easily
// implemented
//
// Default: nullptr
DEFINE_C_WRAP_SETTER_PTR_WRAP(ReadOptions, iterate_upper_bound, Slice)

DEFINE_C_WRAP_CONSTRUCTOR(WriteOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(WriteOptions)
DEFINE_C_WRAP_DESTRUCTOR(WriteOptions)
// Get/Set methods
DEFINE_C_WRAP_GETTER(WriteOptions, sync, bool)
DEFINE_C_WRAP_SETTER(WriteOptions, sync, bool)


DEFINE_C_WRAP_CONSTRUCTOR(FlushOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(FlushOptions)
DEFINE_C_WRAP_DESTRUCTOR(FlushOptions)


DEFINE_C_WRAP_CONSTRUCTOR(CompactionOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(CompactionOptions)
DEFINE_C_WRAP_DESTRUCTOR(CompactionOptions)
