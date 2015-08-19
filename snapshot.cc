// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Abstract handle to particular state of a DB.
// A Snapshot is an immutable object and can therefore be safely
// accessed from multiple threads without any external synchronization.

#include <rocksdb/db.h>
#include "snapshot.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(Snapshot)

SequenceNumber SnapshotGetSequenceNumber(Snapshot_t* snapshot)
{
    assert(GET_REP(snapshot, Snapshot) != NULL);
    return GET_REP(snapshot, Snapshot)->GetSequenceNumber();
}
