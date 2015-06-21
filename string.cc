// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

#include <string>
#include "string.h"

typedef std::string String;

DEFINE_C_WRAP_CONSTRUCTOR(String)
DEFINE_C_WRAP_CONSTRUCTOR_MOVE(String)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(String, const char*, size_t)
DEFINE_C_WRAP_DESTRUCTOR(String)

inline const char* StringGetCStr(String_t * str)
{
    return ((str && GET_REP(str)) ?
            GET_REP(str)->c_str() :
            nullptr);
}

inline int StringGetCStrLen(String_t *str)
{
    return ((str && GET_REP(str)) ?
            GET_REP(str)->length() :
            0);
}

