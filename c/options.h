// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_OPTIONS_H_
#define GO_ROCKSDB_INCLUDE_OPTIONS_H_

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(ColumnFamilyOptions)
DEFINE_C_WRAP_STRUCT(DBOptions)
DEFINE_C_WRAP_STRUCT(Options)

DEFINE_C_WRAP_STRUCT(ReadOptions)
DEFINE_C_WRAP_STRUCT(WriteOptions)

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_OPTIONS_H_
