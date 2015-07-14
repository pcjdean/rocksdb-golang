// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_SLICE_H_
#define GO_ROCKSDB_INCLUDE_SLICE_H_

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(Slice)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Slice)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(Slice, const char*, size_t)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Slice)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(Slice)

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_SLICE_H_
