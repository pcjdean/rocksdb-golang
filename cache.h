// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_CACHE_H_
#define GO_ROCKSDB_INCLUDE_CACHE_H_

#ifdef __cplusplus
#include <rocksdb/cache.h>
using namespace rocksdb;
#endif

#include "types.h"

#ifdef __cplusplus
typedef std::shared_ptr<Cache> PCache;
#endif

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(PCache)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(PCache)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(PCache, size_t, int)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PCache)

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_CACHE_H_
