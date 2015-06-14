// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// A Status encapsulates the result of an operation.  It may indicate success,
// or it may indicate an error with an associated error message.
//
// Multiple threads can invoke const methods on a Status without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Status must use
// external synchronization.

#include <status.h>
#include "status.h"

DEFINE_C_WRAP_CONSTRUCTOR(Status)
DEFINE_C_WRAP_DESTRUCTOR(Status)
DEFINE_C_WRAP_CONSTRUCTOR_COPY(Status)
