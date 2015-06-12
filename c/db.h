// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_DB_H_
#define GO_ROCKSDB_INCLUDE_DB_H_

#include "types.h"
#include "slice.h"
#include "options.h"
#include "status.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(ColumnFamilyHandle)
extern String_t ColumnFamilyGetName(ColumnFamilyHandle_t* column_family);
extern uint32_t ColumnFamilyGetID(ColumnFamilyHandle_t* column_family);

DEFINE_C_WRAP_STRUCT(String)

DEFINE_C_WRAP_STRUCT(ColumnFamilyDescriptor)

DEFINE_C_WRAP_STRUCT(Snapshot)
extern SequenceNumber SnapshotGetSequenceNumber(Snapshot_t* snapshot);

DEFINE_C_WRAP_STRUCT(Range)
extern Range_t NewRangeTFromSlices(Slice_t* start, Slice_t* limit);

DEFINE_C_WRAP_STRUCT(DB)
Status_t DBOpen(const Options_t* options,
                const String_t* name,
                DB_t* dbptr)
Status_t DBOpenForReadOnly(const Options_t* options,
                           const String_t* name, DB_t* dbptr,
                           bool error_if_log_file_exist)

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_DB_H_
