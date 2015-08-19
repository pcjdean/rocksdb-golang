// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <rocksdb/db.h>
#include "columnFamilyHandle.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(ColumnFamilyHandle)
DEFINE_C_WRAP_DESTRUCTOR(ColumnFamilyHandle)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(ColumnFamilyHandle)
String_t ColumnFamilyGetName(const ColumnFamilyHandle_t* column_family)
{
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    const std::string& name_str = GET_REP(column_family, ColumnFamilyHandle)->GetName();
    return NewStringT(const_cast<std::string*>(&name_str));
}
    
uint32_t ColumnFamilyGetID(const ColumnFamilyHandle_t* column_family)
{
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    return GET_REP(column_family, ColumnFamilyHandle)->GetID();
}
