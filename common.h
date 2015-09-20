// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_COMMON_H_
#define GO_ROCKSDB_INCLUDE_COMMON_H_

#ifdef __cplusplus

#include <vector>

typedef std::vector<bool> BoolVector;

#endif

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(BoolVector)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(BoolVector)
DEFINE_C_WRAP_DESTRUCTOR_DEC(BoolVector)

// Push the @val at the end of @bvc
void BoolVectorPushBack(BoolVector_t *bvc, bool val);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_COMMON_H_
