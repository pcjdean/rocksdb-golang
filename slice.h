// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_SLICE_H_
#define GO_ROCKSDB_INCLUDE_SLICE_H_

#ifdef __cplusplus

#include <vector>
#include <deque>
#include <rocksdb/slice.h>

using namespace rocksdb;

typedef std::vector<Slice> SliceVector;
typedef std::deque<Slice> SliceDeque;

#endif

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(Slice)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Slice)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(Slice, const char*, size_t)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Slice)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(Slice)

DEFINE_C_WRAP_STRUCT(SliceVector)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(SliceVector)
DEFINE_C_WRAP_DESTRUCTOR_DEC(SliceVector)

DEFINE_C_WRAP_STRUCT(SliceDeque)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(SliceDeque)
DEFINE_C_WRAP_DESTRUCTOR_DEC(SliceDeque)

const char* SliceData(Slice_t *slc);
size_t SliceSize(Slice_t *slc);

// Return the size of the SliceVector_t
size_t SliceVectorSize(SliceVector_t *slcv);

// Return the Slice in the @index position of the SliceVector_t
Slice_t SliceVectorIndex(SliceVector_t *slcv, size_t index);

// Return the size of the SliceDeque_t
size_t SliceDequeSize(SliceDeque_t *slcv);

// Return the Slice in the @index position of the SliceDeque_t
Slice_t SliceDequeIndex(SliceDeque_t *slcv, size_t index);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_SLICE_H_
