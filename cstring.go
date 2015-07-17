// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include <stdlib.h>
#include "cstring.h"
*/
import "C"

type cString struct {
	str C.String_t
}

func (rstr *cString) goString(del bool) string {
	var (
		cplustr *C.String_t = unsafe.Pointer(&rstr.str)
		cstr *C.char = C.StringGetCStr(cplustr)
		sz C.int = C.StringGetCStrLen(cplustr)
	)
	if del {
		defer C.DeleteStringT(cplustr, C.bool(false))
	}

	if cstr && sz > 0 {
		return C.GoStringN(cstr, sz);
	} else {
		return nil
	}
}

func (str *cString) del()  {
	C.DeleteStringT(unsafe.Pointers(&str.str), false)
}

func (ccstr *C.String_t) cToString() (str string) {
	cstr := cString{str: *ccstr}
	str = cstr.goString(true)
	return
}

func newCStringFromString(str *string) (str *cString) {
	var cstr *C.char = C.CString(string)
	defer C.free(cstr)
	slc = &cString{str: C.NewStringTRawArgs(unsafe.Pointer(cstr), len(string))}
	return
}

func newCString() (str *cString) {
	str = &cString{str: C.NewStringTDefault()}
	return
}

func newStringArrayFromCArray(ccstr *C.String_t, sz uint) (strs []String) {
	defer C.DeleteStringTArray(cstr)
	strs = make([]string, sz)
	for i := 0; i < sz; i++ {
		cstr := cString{str: (*[sz]C.String_t)(unsafe.Pointer(ccstr))[i]}
		strs[i] = cstr.goString(true)
	}
	return
}

func newcStringsFromStringArray(strs []String, sz uint) (cstrs []*cString) {
	cstrs = make([]*cString, len(strs))
	for i, str := range strs {
		cstrs[i] = newCStringFromString(str)
	}
	return
}

func (cstrs []*cString) del() {
	for _, cstr := range cstrs {
		cstr.del()
	}
}

func (cstrs []*cString) toCArray (ccstrs []C.String_t) {
	ccstrs = make([]C.String_t, len(cstrs))
	for i, cstr := range cstrs {
		ccstrs[i] = cstr.str
	}
	return
}
