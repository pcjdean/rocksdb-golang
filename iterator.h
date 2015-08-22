// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

#ifndef GO_ROCKSDB_INCLUDE_ITERATOR_H_
#define GO_ROCKSDB_INCLUDE_ITERATOR_H_

#include "types.h"
#include "slice.h"
#include "status.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(Iterator)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Iterator)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Iterator)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(Iterator)

bool IteratorValid(Iterator_t *it);
void IteratorSeekToFirst(Iterator_t *it);
void IteratorSeekToLast(Iterator_t *it);
void IteratorSeek(Iterator_t *it, const Slice_t* target);
void IteratorNext(Iterator_t *it);
void IteratorPrev(Iterator_t *it);
Slice_t IteratorKey(Iterator_t *it);
Slice_t IteratorValue(Iterator_t *it);
Status_t IteratorStatus(Iterator_t *it);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_ITERATOR_H_
