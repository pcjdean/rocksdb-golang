// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_STRING_H_
#define GO_ROCKSDB_INCLUDE_STRING_H_

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(String)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(String)
DEFINE_C_WRAP_CONSTRUCTOR_MOVE_DEC(String)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(String, const char*, size_t)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(String)
DEFINE_C_WRAP_DESTRUCTOR_DEC(String)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(String)

extern const char* StringGetCStr(String_t * str);
extern int StringGetCStrLen(String_t * str);
extern void StringSetCStrN(String_t *str, char *cstr, int sz);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_STRING_H_
