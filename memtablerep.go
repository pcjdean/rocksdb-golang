// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// This file contains the interface that must be implemented by any collection
// to be used as the backing store for a MemTable. Such a collection must
// satisfy the following properties:
//  (1) It does not store duplicate items.
//  (2) It uses MemTableRep::KeyComparator to compare items for iteration and
//     equality.
//  (3) It can be accessed concurrently by multiple readers and can support
//     during reads. However, it needn't support multiple concurrent writes.
//  (4) Items are never deleted.
// The liberal use of assertions is encouraged to enforce (1).
//
// The factory will be passed an MemTableAllocator object when a new MemTableRep
// is requested.
//
// Users can implement their own memtable representations. We include three
// types built in:
//  - SkipListRep: This is the default; it is backed by a skip list.
//  - HashSkipListRep: The memtable rep that is best used for keys that are
//  structured like "prefix:suffix" where iteration within a prefix is
//  common and iteration across different prefixes is rare. It is backed by
//  a hash map where each bucket is a skip list.
//  - VectorRep: This is backed by an unordered std::vector. On iteration, the
// vector is sorted. It is intelligent about sorting; once the MarkReadOnly()
// has been called, the vector will only be sorted once. It is optimized for
// random-write-heavy workloads.
//
// The last four implementations are designed for situations in which
// iteration over the entire collection is rare since doing so requires all the
// keys to be copied into a sorted data structure.

package rocksdb

/*
#include "memtablerep.h"
*/
import "C"

import (
	"runtime"
)

// Wrap go MemTableRepFactory
type MemTableRepFactory struct {
	mtf C.PMemTableRepFactory_t
}

// Release resources
func (mtf *MemTableRepFactory) finalize() {
	var cmtf *C.PMemTableRepFactory_t= &mtf.mtf
	C.DeletePMemTableRepFactoryT(cmtf, toCBool(false))
}

// C MemTableRepFactory to go MemTableRepFactory
func (cmtf *C.PMemTableRepFactory_t) toMemTableRepFactory() (mtf *MemTableRepFactory) {
	mtf = &MemTableRepFactory{mtf: *cmtf}	
	runtime.SetFinalizer(mtf, finalize)
	return
}

// This uses a skip list to store keys. It is the default.
//
// Parameters:
//   lookahead: If non-zero, each iterator's seek operation will start the
//     search from the previously visited record (doing at most 'lookahead'
//     steps). This is an optimization for the access pattern including many
//     seeks with consecutive keys.
func NewSkipListFactory (lookahead uint64) (mtf *MemTableRepFactory) {
	var (
		clookahead C.size_t =  C.size_t(lookahead)
	)
	
	cmtf := C.PMemTableRepFactoryNewSkipListFactory(clookahead)
	mtf = cmtf.toMemTableRepFactory()
	return
}

// TODO Add memtable interface
