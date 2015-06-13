// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

#include <string>
#include "string.h"

typedef std::string String;

DEFINE_C_WRAP_CONSTRUCTOR(String)
DEFINE_C_WRAP_DESTRUCTOR(String)
