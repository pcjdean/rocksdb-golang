// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Abstract handle to particular state of a DB.
// A Snapshot is an immutable object and can therefore be safely
// accessed from multiple threads without any external synchronization.

package rocksdb

/*
#include "snapshot.h"
*/
import "C"

import (
	"runtime"
)

type Snapshot struct {
	snp C.Snapshot_t
	db *DB
}

func (snp *Snapshot) finalize() {
	if snp.db != nil {
		snp.db.ReleaseSnapshot(snp)
	}
}

func (csnp *C.Snapshot_t) toSnapshot(db *DB) (snp *Snapshot) {
	snp = &Snapshot{snp: *csnp, db: db}	
	runtime.SetFinalizer(snp, finalize)
	return
}

func (snp *Snapshot) GetSequenceNumber() SequenceNumber {
	var csnp *C.Snapshot_t = &snp.snp
	return SequenceNumber(C.SnapshotGetSequenceNumber(csnp))
}
