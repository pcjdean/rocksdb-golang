// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_COLUMNFAMILYHANDLE_H_
#define GO_ROCKSDB_INCLUDE_COLUMNFAMILYHANDLE_H_

#include "types.h"
#include "cstring.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(ColumnFamilyHandle)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(ColumnFamilyHandle)
DEFINE_C_WRAP_DESTRUCTOR_DEC(ColumnFamilyHandle)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(ColumnFamilyHandle)
String_t ColumnFamilyGetName(const ColumnFamilyHandle_t* column_family);
uint32_t ColumnFamilyGetID(const ColumnFamilyHandle_t* column_family);


#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_COLUMNFAMILYHANDLE_H_
