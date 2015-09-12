// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Slice is a simple structure containing a pointer into some external
// storage and a size.  The user of a Slice must ensure that the slice
// is not used after the corresponding external storage has been
// deallocated.
//
// Multiple threads can invoke const methods on a Slice without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Slice must use
// external synchronization.

#include <rocksdb/slice.h>
#include "slice.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(Slice)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(Slice, const char*, size_t)
DEFINE_C_WRAP_DESTRUCTOR(Slice)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(Slice)

// Return a pointer to the beginning of the referenced data
const char* SliceData(Slice_t *slc)
{
    const char* ret = nullptr;
    if (slc && GET_REP(slc, Slice))
    {
        ret = GET_REP(slc, Slice)->data();
    }
    return ret;
}

// Return the length (in bytes) of the referenced data
size_t SliceSize(Slice_t *slc)
{
    size_t ret = 0;
    if (slc && GET_REP(slc, Slice))
    {
        ret = GET_REP(slc, Slice)->size();
    }
    return ret;
}
