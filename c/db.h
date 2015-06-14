// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_DB_H_
#define GO_ROCKSDB_INCLUDE_DB_H_

#include "types.h"
#include "slice.h"
#include "options.h"
#include "status.h"
#include "string.h"
#include "write_batch.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(ColumnFamilyHandle)
extern String_t ColumnFamilyGetName(ColumnFamilyHandle_t* column_family);
extern uint32_t ColumnFamilyGetID(ColumnFamilyHandle_t* column_family);

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
                           bool error_if_log_file_exist);
Status_t DBOpenForReadOnlyWithColumnFamilies(const Options_t* options,
                                           const String_t* name,
                                           const ColumnFamilyDescriptor_t column_families[],
                                           const int size_col,
                                           ColumnFamilyHandle_t **handles,
                                           DB_t* dbptr, bool error_if_log_file_exist);
Status_t DBOpenWithColumnFamilies(const Options_t* options, const String_t* name,
                                const ColumnFamilyDescriptor_t column_families[], const int size_col,
                                ColumnFamilyHandle_t **handles, DB_t* dbptr);
Status_t DBListColumnFamilies(DBOptions_t* db_options,
                              const String_t* name,
                              const String_t **column_families, int* size_col);
Status_t DBCreateColumnFamily(DB_t* dbptr, const ColumnFamilyOptions_t* options,
                            const String_t* column_family_name,
                            ColumnFamilyHandle_t* handle);
Status_t DBDropColumnFamily(DB_t* dbptr, const ColumnFamilyHandle_t* column_family);
Status_t DBPutWithColumnFamily(DB_t* dbptr, const WriteOptions_t* options,
                           const ColumnFamilyHandle_t* column_family,
                           const Slice_t* key,
                           const Slice_t* value);
Status_t DBPut(DB_t* dbptr, const WriteOptions_t* optionss,
               const Slice_t* key,
               const Slice_t* value);
Status_t DBDeleteWithColumnFamily(DB_t* dbptr, const WriteOptions_t* options,
                                  const ColumnFamilyHandle_t* column_family,
                                  const Slice_t* key);
Status_t DBDelete(DB_t* dbptr, const WriteOptions_t* optionss, const Slice_t* key);
Status_t DBMergeWithColumnFamily(DB_t* dbptr, const WriteOptions_t* options,
                                 const ColumnFamilyHandle_t* column_family,
                                 const Slice_t* key,
                                 const Slice_t* value);
Status_t DBMerge(DB_t* dbptr, const WriteOptions_t* optionss,
               const Slice_t* key,
               const Slice_t* value);
Status_t DBWrite(DB_t* dbptr, const WriteOptions_t* optionss, WriteBatch_t* updates);
Status_t DBGetWithColumnFamily(DB_t* dbptr, const ReadOptions_t* options,
                               const ColumnFamilyHandle_t* column_family,
                               const Slice_t* key,
                               const String_t* value);
Status_t DBGet(DB_t* dbptr, const ReadOptions_t* options,
               const Slice_t* key,
               const String_t* value);

Status_t* DBMultiGetWithColumnFamily(DB_t* dbptr, const ReadOptions_t* options,
                                     const ColumnFamilyHandle_t column_families[],
                                     const int size_col,
                                     const Slice_t keys[],
                                     const int size_keys,
                                     String_t** values);
Status_t* DBMultiGet(DB_t* dbptr, const ReadOptions_t* options,
                     const Slice_t keys[],
                     const int size_keys,
                     String_t** values);
bool DBKeyMayExistWithColumnFamily(DB_t* dbptr, const ReadOptions_t* options
                                   ColumnFamilyHandle_t* column_family,
                                   const Slice_t* key,
                                   String_t* value,
                                   bool* value_found);
bool DBKeyMayExist(DB_t* dbptr, const ReadOptions_t* options
                   const Slice_t* key,
                   String_t* value,
                   bool* value_found);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_DB_H_
