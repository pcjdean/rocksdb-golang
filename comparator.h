// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_COMPARATOR_H_
#define GO_ROCKSDB_INCLUDE_COMPARATOR_H_

#include "types.h"
#include "slice.h"
#include "cstring.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(Comparator)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Comparator)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Comparator)

// Return a Comparator from a go Comparator interface
Comparator_t NewComparator(void* go_cmp);

Comparator_t GoBytewiseComparator();
Comparator_t GoReverseBytewiseComparator();

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_COMPARATOR_H_
