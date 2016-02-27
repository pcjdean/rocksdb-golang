// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_H_
#define GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_H_

#include "types.h"
#include "common.h"
#include "sliceTransform.h"
#include "cstring.h"
#include "slice.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(CompactionFilterContext)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilterContext)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilterContext)

DEFINE_C_WRAP_STRUCT(CompactionFilter)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilter)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilter)
// Return a CompactionFilter from a go ICompactionFilter
CompactionFilter_t NewCompactionFilter(void* go_cpf);

// Definitions for CompactionFilter::Context
DEFINE_C_WRAP_STRUCT(CompactionFilter_Context)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilter_Context)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilter_Context)

DEFINE_C_WRAP_STRUCT(PCompactionFilterFactory)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(PCompactionFilterFactory)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(PCompactionFilterFactory)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PCompactionFilterFactory)
// Return a CompactionFilterFactory from a go ICompactionFilterFactory
PCompactionFilterFactory_t NewPCompactionFilterFactory(void* go_cpflt);


#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_H_
