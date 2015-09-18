// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_H_
#define GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_H_

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(SliceTransform)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(SliceTransform)
DEFINE_C_WRAP_DESTRUCTOR_DEC(SliceTransform)


// Return a SliceTransform from a go SliceTransform interface
SliceTransform_t SliceTransformNewSliceTransform(void* go_stf);

SliceTransform_t SliceTransformNewFixedPrefixTransform(size_t prefix_len);
SliceTransform_t SliceTransformNewCappedPrefixTransform(size_t cap_len);
SliceTransform_t SliceTransformNewNoopTransform();

#ifdef __cplusplus
}  /* end extern "C" */
#endif

}

#endif  // GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_H_
