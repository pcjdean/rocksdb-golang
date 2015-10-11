// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_H_
#define GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_H_

#include "types.h"
#include "slice.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(SliceTransform)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(SliceTransform)
DEFINE_C_WRAP_DESTRUCTOR_DEC(SliceTransform)

DEFINE_C_WRAP_STRUCT(PConstSliceTransform)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(PConstSliceTransform)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PConstSliceTransform)

// Return a SliceTransform from a go SliceTransform interface
SliceTransform_t NewSliceTransform(void* go_stf);

// Return a SharedSliceTransform from a go SliceTransform interface
PConstSliceTransform_t NewSharedSliceTransform(void* go_stf);

PConstSliceTransform_t GoNewFixedPrefixTransform(size_t prefix_len);
PConstSliceTransform_t GoNewCappedPrefixTransform(size_t cap_len);
PConstSliceTransform_t GoNewNoopTransform();

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_H_
