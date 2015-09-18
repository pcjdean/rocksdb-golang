// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_H_
#define GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_H_

#ifdef __cplusplus
#include <rocksdb/compaction_filter.h>
using namespace rocksdb;
#endif

#include "types.h"

#ifdef __cplusplus
typedef rocksdb::CompactionFilter::Context CompactionFilter_Context;

typedef std::shared_ptr<CompactionFilterFactory> PCompactionFilterFactory;
typedef std::shared_ptr<CompactionFilterFactoryV2> PCompactionFilterFactoryV2;
#endif

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(CompactionFilterContext)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilterContext)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilterContext)

DEFINE_C_WRAP_STRUCT(CompactionFilter)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilter)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilter)

// Definitions for CompactionFilter::Context
DEFINE_C_WRAP_STRUCT(CompactionFilter_Context)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilter_Context)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilter_Context)

DEFINE_C_WRAP_STRUCT(CompactionFilterV2)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilterV2)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilterV2)

DEFINE_C_WRAP_STRUCT(PCompactionFilterFactory)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(PCompactionFilterFactory)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PCompactionFilterFactory)

DEFINE_C_WRAP_STRUCT(PCompactionFilterFactoryV2)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(PCompactionFilterFactoryV2)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PCompactionFilterFactoryV2)


#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_H_
