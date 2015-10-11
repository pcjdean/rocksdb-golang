// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_PRIVATE_H_
#define GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_PRIVATE_H_

#ifdef __cplusplus
#include <memory>
#include <rocksdb/slice_transform.h>
using namespace rocksdb;

typedef std::shared_ptr<const SliceTransform> PConstSliceTransform;
#endif

#endif  // GO_ROCKSDB_INCLUDE_SLICE_TRANSFORM_PRIVATE_H_
