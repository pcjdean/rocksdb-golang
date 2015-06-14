// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <options.h>
#include "options.h"

DEFINE_C_WRAP_CONSTRUCTOR(ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(ColumnFamilyOptions, Options)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(ColumnFamilyOptions)
DEFINE_C_WRAP_DESTRUCTOR(ColumnFamilyOptions)

DEFINE_C_WRAP_CONSTRUCTOR(DBOptions)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(DBOptions, Options)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(DBOptions)
DEFINE_C_WRAP_DESTRUCTOR(DBOptions)

DEFINE_C_WRAP_CONSTRUCTOR(Options)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(Options, DBOptions, ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(Options)
DEFINE_C_WRAP_DESTRUCTOR(Options)

DEFINE_C_WRAP_CONSTRUCTOR(ReadOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(ReadOptions, bool, bool)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(ReadOptions)
DEFINE_C_WRAP_DESTRUCTOR(ReadOptions)

DEFINE_C_WRAP_CONSTRUCTOR(WriteOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(WriteOptions)
DEFINE_C_WRAP_DESTRUCTOR(WriteOptions)
