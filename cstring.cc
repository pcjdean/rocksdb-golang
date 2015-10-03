// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

#include <stdio.h>
#include "cstring.h"

DEFINE_C_WRAP_CONSTRUCTOR(String)
DEFINE_C_WRAP_CONSTRUCTOR_MOVE(String)
DEFINE_C_WRAP_CONSTRUCTOR_COPY(String)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(String, const char*, size_t)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(String)
DEFINE_C_WRAP_DESTRUCTOR(String)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(String)

DEFINE_C_WRAP_CONSTRUCTOR(StringVector)
DEFINE_C_WRAP_DESTRUCTOR(StringVector)

DEFINE_C_WRAP_CONSTRUCTOR(StringDeque)
DEFINE_C_WRAP_DESTRUCTOR(StringDeque)

// Return a c string of str
const char* StringGetCStr(String_t * str)
{
    return ((str && GET_REP(str, String)) ?
            GET_REP(str, String)->c_str() :
            nullptr);
}

// Return the length og str
size_t StringGetCStrLen(String_t *str)
{
    return ((str && GET_REP(str, String)) ?
            GET_REP(str, String)->length() :
            0);
}

//Set str to catr
void StringSetCStr(String_t * str, const char* cstr, size_t len)
{
    if (str && GET_REP(str, String))
    {
        GET_REP(str, String)->assign(cstr, len);
    }
}

// Push the @cstr at the end of @slcv
void StringVectorPushBack(StringVector_t *slcv, const char* cstr, size_t len)
{
    String s(cstr, len);
    
    if (slcv && GET_REP(slcv, StringVector))
    {
        GET_REP(slcv, StringVector)->push_back(s);
    }
    else
    {
        printf("StringVectorPushBack null pointer - slcv = %p\n", slcv);
    }
}

// Return the length of the StringDeque
size_t StringDequeSize(StringDeque_t * strdeq)
{
    size_t ret = 0;
    if (strdeq && GET_REP(strdeq, StringDeque))
    {
        ret = GET_REP(strdeq, StringDeque)->size();
    }
    return ret;
}

// Return the string at @index of the StringDeque
String_t StringDequeAt(StringDeque_t * strdeq, size_t index)
{
    String_t ret{nullptr};
    if (strdeq && GET_REP(strdeq, StringDeque))
    {
        ret.rep = &GET_REP_REF(strdeq, StringDeque)[index];
    }
    return ret;
}

