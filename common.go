// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "common.h"
*/
import "C"

// Set the C bool vector to go bool array
func (bvc *C.BoolVector_t) setBoolArray(vals []bool) {
	for _, val := range vals {
		C.BoolVectorPushBack(bvc, toCBool(val))
	}
	return
}
