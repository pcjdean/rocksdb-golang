// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <stdio.h>
#include "common.h"

DEFINE_C_WRAP_CONSTRUCTOR(BoolVector)
DEFINE_C_WRAP_DESTRUCTOR(BoolVector)

// Push the @val at the end of @bvc
void BoolVectorPushBack(BoolVector_t *bvc, bool val)
{
    if (bvc && GET_REP(bvc, BoolVector))
    {
        GET_REP(bvc, BoolVector)->push_back(val);
    }
    else
    {
        printf("BoolVectorPushBack null pointer - bvc = %p\n", bvc);
    }
}
