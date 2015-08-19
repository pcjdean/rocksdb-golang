// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

#include <string>
#include "cstring.h"

DEFINE_C_WRAP_CONSTRUCTOR(String)
DEFINE_C_WRAP_CONSTRUCTOR_MOVE(String)
DEFINE_C_WRAP_CONSTRUCTOR_COPY(String)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(String, const char*, size_t)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(String)
DEFINE_C_WRAP_DESTRUCTOR(String)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(String)

const char* StringGetCStr(String_t * str)
{
    return ((str && GET_REP(str, String)) ?
            GET_REP(str, String)->c_str() :
            nullptr);
}

int StringGetCStrLen(String_t *str)
{
    return ((str && GET_REP(str, String)) ?
            GET_REP(str, String)->length() :
            0);
}

