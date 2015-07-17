// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <rocksdb/metadata.h>
#include "metadata.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(ColumnFamilyMetaData)
DEFINE_C_WRAP_DESTRUCTOR(ColumnFamilyMetaData)

DEFINE_C_WRAP_CONSTRUCTOR(LiveFileMetaData)
DEFINE_C_WRAP_DESTRUCTOR(LiveFileMetaData)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(LiveFileMetaData)
