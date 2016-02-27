// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "compactionFilter.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// CompactionFilter allows an application to modify/delete a key-value at
// the time of compaction.
type ICompactionFilter interface {

	// The compaction process invokes this
	// method for kv that is being compacted. A return value
	// of false indicates that the kv should be preserved in the
	// output of this compaction run and a return value of true
	// indicates that this key-value should be removed from the
	// output of the compaction.  The application can inspect
	// the existing value of the key and make decision based on it.
	//
	// When the value is to be preserved, the application has the option
	// to modify the existing_value and pass it back through new_value.
	// value_changed needs to be set to true in this case.
	//
	// If multithreaded compaction is being used *and* a single CompactionFilter
	// instance was supplied via Options::compaction_filter, this method may be
	// called from different threads concurrently.  The application must ensure
	// that the call is thread-safe.
	//
	// If the CompactionFilter was created by a factory, then it will only ever
	// be used by a single thread that is doing the compaction run, and this
	// call does not need to be thread-safe.  However, multiple filters may be
	// in existence and operating concurrently.
	Filter(level int, key, exval []byte) (newval []byte, valchg bool, removed bool)

	// Returns a name that identifies this compaction filter.
	// The name will be printed to LOG file on start up for diagnosis.
	Name() string
}

// Wrap functions for ICompactionFilter

//export ICompactionFilterName
func ICompactionFilterName(ccpf unsafe.Pointer) *C.char {
	cpf := InterfacesGet(ccpf).(ICompactionFilter)
	return C.CString(cpf.Name())
}

//export ICompactionFilterFilter
func ICompactionFilterFilter(ccpf unsafe.Pointer, level C.int, key, exval *C.Slice_t, newval *C.String_t, valchg *C.bool) C.bool {
	cpf := InterfacesGet(ccpf).(ICompactionFilter)
	gnewval, gvalchg, removed := cpf.Filter(int(level), key.cToBytes(false), exval.cToBytes(false))
	*valchg = toCBool(gvalchg)
	if gvalchg {
		newval.setBytes(gnewval)
	}
	return toCBool(removed)
}

// Wrap go CompactionFilter
type CompactionFilter struct {
	cpf C.CompactionFilter_t
	// True if the CompactionFilter is closed
	closed bool
}

// Release resources
func (cpf *CompactionFilter) finalize() {
	if !cpf.closed {
		cpf.closed = true
		var ccpf *C.CompactionFilter_t= &cpf.cpf
		C.DeleteCompactionFilterT(ccpf, toCBool(false))
	}
}

// Close the @cpf
func (cpf *CompactionFilter) Close() {
	runtime.SetFinalizer(cpf, nil)
	cpf.finalize()
}

// C CompactionFilter to go CompactionFilter
func (ccpf *C.CompactionFilter_t) toCompactionFilter() (cpf *CompactionFilter) {
	cpf = &CompactionFilter{cpf: *ccpf}	
	runtime.SetFinalizer(cpf, finalize)
	return
}

// Return a new default CompactionFilter
func NewDefaultCompactionFilter() (cpf *CompactionFilter) {
	cpf = &CompactionFilter{cpf: C.CompactionFilter_t{nil}}	
	return
}

// Return a new CompactionFilter that uses ICompactionFilter
func NewCompactionFilter(itf ICompactionFilter) (cpf *CompactionFilter) {
	var iftp unsafe.Pointer = nil

	if nil != itf {
		iftp =InterfacesAddReference(itf)
	}
	ccpf := C.NewCompactionFilter(iftp)
	return ccpf.toCompactionFilter()
}

// Wrap go CompactionFilter_Context
type CompactionFilter_Context struct {
	cfc C.CompactionFilter_Context_t
}

// C CompactionFilter_Context to go CompactionFilter_Context
func (ccfc *C.CompactionFilter_Context_t) toCompactionFilter_Context() (cfc *CompactionFilter_Context) {
	cfc = &CompactionFilter_Context{cfc: *ccfc}	
	return
}

// Each compaction will create a new CompactionFilter allowing the
// application to know about different compactions
type ICompactionFilterFactory interface {

	// Create a ICompactionFilter
	CreateCompactionFilter(context *CompactionFilter_Context) ICompactionFilter

	// Returns a name that identifies this compaction filter factory.
	Name() string
}

// Wrap functions for ICompactionFilterFactory

//export ICompactionFilterFactoryName
func ICompactionFilterFactoryName(ccpf unsafe.Pointer) *C.char {
	cpf := InterfacesGet(ccpf).(ICompactionFilterFactory)
	return C.CString(cpf.Name())
}

//export ICompactionFilterFactoryCreateCompactionFilter
func ICompactionFilterFactoryCreateCompactionFilter(ccpf unsafe.Pointer, context *C.CompactionFilter_Context_t) (filter unsafe.Pointer) {
	cpf := InterfacesGet(ccpf).(ICompactionFilterFactory)
	filter = InterfacesAddReference(cpf.CreateCompactionFilter(context.toCompactionFilter_Context()))
	return
}

// Wrap go CompactionFilterFactory
type CompactionFilterFactory struct {
	cff C.PCompactionFilterFactory_t
}

// Release resources
func (cff *CompactionFilterFactory) finalize() {
	var ccff *C.PCompactionFilterFactory_t= &cff.cff
	C.DeletePCompactionFilterFactoryT(ccff, toCBool(false))
}

// C CompactionFilterFactory to go CompactionFilterFactory
func (ccff *C.PCompactionFilterFactory_t) toCompactionFilterFactory() (cff *CompactionFilterFactory) {
	cff = &CompactionFilterFactory{cff: *ccff}	
	runtime.SetFinalizer(cff, finalize)
	return
}

// Return a new default CompactionFilterFactory
func NewDefaultCompactionFilterFactory() (cff *CompactionFilterFactory) {
	cff = &CompactionFilterFactory{cff: C.NewPCompactionFilterFactoryTDefault()}	
	runtime.SetFinalizer(cff, finalize)
	return
}

// Return a new CompactionFilterFactory that uses ICompactionFilterFactory
func NewCompactionFilterFactory(itf ICompactionFilterFactory) (cff *CompactionFilterFactory) {
	var iftp unsafe.Pointer = nil

	if nil != itf {
		iftp =InterfacesAddReference(itf)
	}
	ccff := C.NewPCompactionFilterFactory(iftp)
	return ccff.toCompactionFilterFactory()
}

// Wrap go CompactionFilterContext
type CompactionFilterContext struct {
	cfc C.CompactionFilterContext_t
}

// C CompactionFilterContext_t to go CompactionFilterContext
func (ccfc *C.CompactionFilterContext_t) toCompactionFilterContext() (cfc *CompactionFilterContext) {
	cfc = &CompactionFilterContext{cfc: *ccfc}	
	return
}
