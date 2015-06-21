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

type cString struct {
	str C.String_t
}

func (rstr *cString) GoString(del bool) string {
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

func NewCStringFromString(str *string) (str *cString) {
	var cstr *C.char = C.CString(string)
	defer C.free(cstr)
	slc = &cString{str: C.NewStringTRawArgs(unsafe.Pointer(cstr), len(string))}
	return
}
