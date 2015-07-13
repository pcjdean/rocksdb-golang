// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// An iterator yields a sequence of key/value pairs from a source.
// The following class defines the interface.  Multiple implementations
// are provided by this library.  In particular, iterators are provided
// to access the contents of a Table or a DB.
//
// Multiple threads can invoke const methods on an Iterator without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Iterator must use
// external synchronization.

#include <rocksdb/iterator.h>
#include "iterator.h"

DEFINE_C_WRAP_CONSTRUCTOR(Iterator)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(Iterator)
DEFINE_C_WRAP_DESTRUCTOR(Iterator)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(Iterator)
