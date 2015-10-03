// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_STRING_H_
#define GO_ROCKSDB_INCLUDE_STRING_H_

#ifdef __cplusplus

#include <string>
#include <vector>

typedef std::string String;
typedef std::vector<String> StringVector;
typedef std::deque<String> StringDeque;

#endif

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(String)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(String)
DEFINE_C_WRAP_CONSTRUCTOR_MOVE_DEC(String)
DEFINE_C_WRAP_CONSTRUCTOR_COPY_DEC(String)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(String, const char*, size_t)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(String)
DEFINE_C_WRAP_DESTRUCTOR_DEC(String)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(String)

DEFINE_C_WRAP_STRUCT(StringVector)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(StringVector)
DEFINE_C_WRAP_DESTRUCTOR_DEC(StringVector)

DEFINE_C_WRAP_STRUCT(StringDeque)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(StringDeque)
DEFINE_C_WRAP_DESTRUCTOR_DEC(StringDeque)

const char* StringGetCStr(String_t * str);
size_t StringGetCStrLen(String_t * str);

// Return the string at @index of the StringDeque
String_t StringDequeAt(StringDeque_t * strdeq, size_t index);

// Return the length of the StringDeque
size_t StringDequeSize(StringDeque_t * strdeq);

//Set str to cstr
void StringSetCStr(String_t * str, const char* cstr, size_t len);

// Push the @cstr at the end of @slcv
void StringVectorPushBack(StringVector_t *slcv, const char * cstr, size_t len);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_STRING_H_
