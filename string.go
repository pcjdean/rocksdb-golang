// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include <stdlib.h>
#include "string.h"
*/
import "C"

type String struct {
	str C.String_t
}

func (rstr *String) GoString(del bool) string {
	var (
		cplustr *C.String_t = unsafe.Pointer(&rstr.str)
		cstr *C.char = C.StringGetCStr(cplustr)
		sz C.int = C.StringGetCStrLen(cplustr)
	)
	if del {
		defer C.DeleteStringT(cplustr, C.bool(false))
	}

	return C.GoStringN(cstr, sz);
}

func (rstr *String) SetGoString(str *string) {
	var (
		cstr *C.char = C.CString(string)
		cplustr *C.String_t = unsafe.Pointer(&rstr.str)
	)
	defer C.free(cstr)
	C.StringSetCStrN(cplustr, cstr, len(string))
}
