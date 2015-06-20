// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include "options.h"
*/
import "C"

type Options struct {
	opt C.Options_t
}

type DBOptions struct {
	dbopt C.DBOptions_t
}

