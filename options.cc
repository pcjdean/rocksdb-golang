// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

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

DEFINE_C_WRAP_CONSTRUCTOR(DBOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(DBOptions, const Options&)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(DBOptions)
DEFINE_C_WRAP_DESTRUCTOR(DBOptions)
DEFINE_C_WRAP_GETTER(DBOptions, create_if_missing, bool)
DEFINE_C_WRAP_SETTER(DBOptions, create_if_missing, bool)
DEFINE_C_WRAP_GETTER(DBOptions, error_if_exists, bool)
DEFINE_C_WRAP_SETTER(DBOptions, error_if_exists, bool)

DEFINE_C_WRAP_CONSTRUCTOR(Options)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(Options, DBOptions, ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(Options, const DBOptions&, const ColumnFamilyOptions&)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(Options)
DEFINE_C_WRAP_DESTRUCTOR(Options)

DEFINE_C_WRAP_CONSTRUCTOR(ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(ReadOptions, bool, bool)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(ReadOptions)
DEFINE_C_WRAP_DESTRUCTOR(ReadOptions)

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
