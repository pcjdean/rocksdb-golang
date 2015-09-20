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

#include "slice.h"

DEFINE_C_WRAP_CONSTRUCTOR(Slice)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(Slice, const char*, size_t)
DEFINE_C_WRAP_DESTRUCTOR(Slice)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(Slice)

DEFINE_C_WRAP_CONSTRUCTOR_DEC(SliceVector)
DEFINE_C_WRAP_DESTRUCTOR_DEC(SliceVector)

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

// Return the size of the SliceVector_t @slcv
size_t SliceVectorSize(SliceVector_t *slcv)
{
    return (slcv && GET_REP(slcv, SliceVector) ? GET_REP(slcv, SliceVector)->size() : 0);
}

// Return the Slice in the @index position of the SliceVector_t @slcv
Slice_t SliceVectorIndex(SliceVector_t *slcv, size_t index)
{
    Slice_t ret{nullptr};
    
    return (slcv && GET_REP(slcv, SliceVector) &&
            index < GET_REP(slcv, SliceVector)->size() ?
            Slice_t{&GET_REP_REF(slcv, SliceVector)[index]} : ret);
}

