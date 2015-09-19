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

import (
	"unsafe"
)

// Go wrap C string
type cString struct {
	str C.String_t
}

// Go wrap C string array
type cStringPtrAry []*cString

// Go wrap C string to go string
func (rstr *cString) goString(del bool) (str string) {
	var (
		cplustr *C.String_t = &rstr.str
		cstr *C.char = C.StringGetCStr(cplustr)
		sz C.int = C.StringGetCStrLen(cplustr)
	)
	if del {
		defer C.DeleteStringT(cplustr, toCBool(false))
	}

	if unsafe.Pointer(cstr) != nil && sz > 0 {
		str = C.GoStringN(cstr, sz);
	}
	return
}

// Go wrap C string to go bytes
func (rstr *cString) goBytes(del bool) (str []byte) {
	var (
		cplustr *C.String_t = &rstr.str
		cstr *C.char = C.StringGetCStr(cplustr)
		sz C.int = C.StringGetCStrLen(cplustr)
	)
	if del {
		defer C.DeleteStringT(cplustr, toCBool(false))
	}

	if unsafe.Pointer(cstr) != nil && sz > 0 {
		str = C.GoBytes(unsafe.Pointer(cstr), sz);
	}
	return
}

// Delete the go wrap C string
func (str *cString) del()  {
	C.DeleteStringT(&str.str, toCBool(false))
}

// C string to go string
func (ccstr *C.String_t) cToString() (str string) {
	cstr := cString{str: *ccstr}
	str = cstr.goString(true)
	return
}

// C string to go bytes
func (ccstr *C.String_t) cToBytes() (str []byte) {
	cstr := cString{str: *ccstr}
	str = cstr.goBytes(true)
	return
}

// Set C string to go bytes
func (ccstr *C.String_t) setBytes(str []byte) {
	cstr := C.CString(str)
	defer C.free(cstr)
	C.StringSetCStr(ccstr, cstr, C.uint64ToSizeT(C.uint64_t(len(str))))
}

// The caller is responsible for delete the returned cstr
func newCStringFromString(str *string) (cstr *cString) {
	var ccstr *C.char = C.CString(*str)
	defer C.free(unsafe.Pointer(ccstr))
	cstr = &cString{str: C.NewStringTRawArgs(ccstr, C.uint64ToSizeT(C.uint64_t(len(*str))))}
	return
}

// The caller is responsible for delete the returned cstr
func newCString() (str *cString) {
	str = &cString{str: C.NewStringTDefault()}
	return
}

// C string array to go string array
func newStringArrayFromCArray(ccstr *C.String_t, sz uint) (strs []string) {
	defer C.DeleteStringTArray(ccstr)
	strs = make([]string, sz)
	for i := uint(0); i < sz; i++ {
		cstr := cString{str: (*[arrayDimenMax]C.String_t)(unsafe.Pointer(ccstr))[i]}
		strs[i] = cstr.goString(true)
	}
	return
}

// C string array to go bytes array
func newBytesFromCArray(ccstr *C.String_t, sz uint) (strs [][]byte) {
	defer C.DeleteStringTArray(ccstr)
	strs = make([][]byte, sz)
	for i := uint(0); i < sz; i++ {
		cstr := cString{str: (*[arrayDimenMax]C.String_t)(unsafe.Pointer(ccstr))[i]}
		strs[i] = cstr.goBytes(true)
	}
	return
}

// The caller is responsible for delete the returned cstrs
func newcStringsFromStringArray(strs []string) (cstrs []*cString) {
	cstrs = make([]*cString, len(strs))
	for i, str := range strs {
		cstrs[i] = newCStringFromString(&str)
	}
	return
}

// Delete go wrap C string array
func (cstrs *cStringPtrAry) del() {
	for _, cstr := range *cstrs {
		cstr.del()
	}
}

// Go wrap C string array to C string array
func (cstrs *cStringPtrAry) toCArray() (ccstrs []C.String_t) {
	ccstrs = make([]C.String_t, len(*cstrs))
	for i, cstr := range *cstrs {
		ccstrs[i] = cstr.str
	}
	return
}
