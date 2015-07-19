// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_TRANSACTION_LOG_ITERATOR_H_
#define GO_ROCKSDB_INCLUDE_TRANSACTION_LOG_ITERATOR_H_

#ifdef __cplusplus
#include <rocksdb/transaction_log.h>
#endif

#include "types.h"

#ifdef __cplusplus
typedef rocksdb::TransactionLogIterator::ReadOptions TransactionLogIterator_ReadOptions;
#endif

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(LogFile)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(LogFile)
DEFINE_C_WRAP_DESTRUCTOR_DEC(LogFile)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(LogFile)

DEFINE_C_WRAP_STRUCT(TransactionLogIterator)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(TransactionLogIterator)
DEFINE_C_WRAP_DESTRUCTOR_DEC(TransactionLogIterator)

DEFINE_C_WRAP_STRUCT(TransactionLogIterator_ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(TransactionLogIterator_ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(TransactionLogIterator_ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(TransactionLogIterator_ReadOptions, bool)
DEFINE_C_WRAP_DESTRUCTOR_DEC(TransactionLogIterator_ReadOptions)

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_TRANSACTION_LOG_ITERATOR_H_
