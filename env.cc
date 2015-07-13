// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// An Env is an interface used by the rocksdb implementation to access
// operating system functionality like the filesystem etc.  Callers
// may wish to provide a custom Env object when opening a database to
// get fine gain control; e.g., to rate limit file system operations.
//
// All Env implementations are safe for concurrent access from
// multiple threads without any external synchronization.

#include <rocksdb/env.h>
#include "env.h"

DEFINE_C_WRAP_CONSTRUCTOR(Env)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(Env)
DEFINE_C_WRAP_DESTRUCTOR(Env)
