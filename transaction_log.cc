// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include "transaction_log.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(LogFile)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(LogFile)
DEFINE_C_WRAP_DESTRUCTOR(LogFile)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(LogFile)

DEFINE_C_WRAP_CONSTRUCTOR(TransactionLogIterator)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(TransactionLogIterator)
DEFINE_C_WRAP_DESTRUCTOR(TransactionLogIterator)

DEFINE_C_WRAP_CONSTRUCTOR(TransactionLogIterator_ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(TransactionLogIterator_ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(TransactionLogIterator_ReadOptions, bool)
DEFINE_C_WRAP_DESTRUCTOR(TransactionLogIterator_ReadOptions)
