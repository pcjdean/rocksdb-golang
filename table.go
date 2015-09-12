// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// A Cache is an interface that maps keys to values.  It has internal
// synchronization and may be safely accessed concurrently from
// multiple threads.  It may automatically evict entries to make room
// for new entries.  Values have a specified charge against the cache
// capacity.  For example, a cache where the values are variable
// length strings, may use the length of the string as the charge for
// the string.
//
// A builtin cache implementation with a least-recently-used eviction
// policy is provided.  Clients may use their own implementations if
// they want something more sophisticated (like scan-resistance, a
// custom eviction policy, variable cache sizing, etc.)

package rocksdb

/*
#include "table.h"
*/
import "C"

import (
	"runtime"
)

// Wrap go TableFactory
type TableFactory struct {
	tbf C.PTableFactory_t
}

// Release resources
func (tbf *TableFactory) finalize() {
	var ctbf *C.PTableFactory_t= &tbf.tbf
	C.DeletePTableFactoryT(ctbf, toCBool(false))
}

// C TableFactory to go TableFactory
func (ctbf *C.PTableFactory_t) toTableFactory() (tbf *TableFactory) {
	tbf = &TableFactory{tbf: *ctbf}	
	runtime.SetFinalizer(tbf, finalize)
	return
}

// Wrap go BlockBasedTableOptions
type BlockBasedTableOptions struct {
	btop C.BlockBasedTableOptions_t
}

// Release resources
func (btop *BlockBasedTableOptions) finalize() {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	C.DeleteBlockBasedTableOptionsT(cbtop, toCBool(false))
}

// C BlockBasedTableOptions to go BlockBasedTableOptions
func (cbtop *C.BlockBasedTableOptions_t) toBlockBasedTableOptions() (btop *BlockBasedTableOptions) {
	btop = &BlockBasedTableOptions{btop: *cbtop}	
	runtime.SetFinalizer(btop, finalize)
	return
}

// Create default BlockBasedTableOptions
func NewBlockBasedTableOptions() *BlockBasedTableOptions {
	cbtop := C.NewBlockBasedTableOptionsTDefault()
	return cbtop.toBlockBasedTableOptions()
}

// If non-NULL use the specified cache for blocks.
// If NULL, rocksdb will automatically create and use an 8MB internal cache.
func (btop *BlockBasedTableOptions) SetBlockCache(cache *Cache) {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	C.BlockBasedTableOptions_set_block_cache(cbtop, &cache.cache)
}

// If non-nullptr, use the specified filter policy to reduce disk reads.
// Many applications will benefit from passing the result of
// NewBloomFilterPolicy() here.
func (btop *BlockBasedTableOptions) SetFilterPolicy(flp *FilterPolicy) {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	if nil == flp {
		def := C.NewPFilterPolicyTDefault()
		flp = def.toFilterPolicy()
	}
	C.BlockBasedTableOptions_set_filter_policy(cbtop, &flp.flp)
}

// Create default block based table factory.
func (btop *BlockBasedTableOptions) NewBlockBasedTableFactory() *TableFactory {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	ctbf := C.NewBlockBasedTableFactory(cbtop)
	return ctbf.toTableFactory()
}

// Wrap go PlainTableOptions
type PlainTableOptions struct {
	ptop C.PlainTableOptions_t
}

// Release resources
func (ptop *PlainTableOptions) finalize() {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.DeletePlainTableOptionsT(cptop, toCBool(false))
}

// C PlainTableOptions to go PlainTableOptions
func (cptop *C.PlainTableOptions_t) toPlainTableOptions() (ptop *PlainTableOptions) {
	ptop = &PlainTableOptions{ptop: *cptop}	
	runtime.SetFinalizer(ptop, finalize)
	return
}

// Create default PlainTableOptions
func NewPlainTableOptions() *PlainTableOptions {
	cptop := C.NewPlainTableOptionsTDefault()
	return cptop.toPlainTableOptions()
}

// Wrap go CuckooTableOptions
type CuckooTableOptions struct {
	ctop C.CuckooTableOptions_t
}

// Release resources
func (ctop *CuckooTableOptions) finalize() {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	C.DeleteCuckooTableOptionsT(cctop, toCBool(false))
}

// C CuckooTableOptions to go CuckooTableOptions
func (cctop *C.CuckooTableOptions_t) toCuckooTableOptions() (ctop *CuckooTableOptions) {
	ctop = &CuckooTableOptions{ctop: *cctop}	
	runtime.SetFinalizer(ctop, finalize)
	return
}

// Create default CuckooTableOptions
func NewCuckooTableOptions() *CuckooTableOptions {
	cctop := C.NewCuckooTableOptionsTDefault()
	return cctop.toCuckooTableOptions()
}
