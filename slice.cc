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
