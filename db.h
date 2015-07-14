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
#include "iterator.h"
#include "env.h"
#include "metadata.h"
#include "transaction_log.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(ColumnFamilyHandle)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(ColumnFamilyHandle)
DEFINE_C_WRAP_DESTRUCTOR_DEC(ColumnFamilyHandle)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(ColumnFamilyHandle)
extern String_t ColumnFamilyGetName(ColumnFamilyHandle_t* column_family);
extern uint32_t ColumnFamilyGetID(ColumnFamilyHandle_t* column_family);

DEFINE_C_WRAP_STRUCT(TablePropertiesCollection)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(TablePropertiesCollection)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(TablePropertiesCollection)
DEFINE_C_WRAP_DESTRUCTOR_DEC(TablePropertiesCollection)

DEFINE_C_WRAP_STRUCT(ColumnFamilyDescriptor)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(ColumnFamilyDescriptor, String, ColumnFamilyOptions)
// DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(ColumnFamilyDescriptor, String, ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(ColumnFamilyDescriptor, kDefaultColumnFamilyName, ColumnFamilyOptions())
DEFINE_C_WRAP_DESTRUCTOR_DEC(ColumnFamilyDescriptor)

DEFINE_C_WRAP_STRUCT(Snapshot)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Snapshot)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Snapshot)
extern SequenceNumber SnapshotGetSequenceNumber(Snapshot_t* snapshot);

DEFINE_C_WRAP_STRUCT(Range)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Range)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(Range, Slice, Slice)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Range)
extern Range_t NewRangeTFromSlices(Slice_t* start, Slice_t* limit);

DEFINE_C_WRAP_STRUCT(DB)
DEFINE_C_WRAP_DESTRUCTOR_DEC(DB)
Status_t DBOpen(const Options_t* options,
                const String_t* name,
                DB_t* dbptr);
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
bool DBKeyMayExistWithColumnFamily(DB_t* dbptr, const ReadOptions_t* options,
                                   ColumnFamilyHandle_t* column_family,
                                   const Slice_t* key,
                                   String_t* value,
                                   bool* value_found);
bool DBKeyMayExist(DB_t* dbptr, const ReadOptions_t* options,
                   const Slice_t* key,
                   String_t* value,
                   bool* value_found);
Iterator_t DBNewIteratorWithColumnFamily(DB_t* dbptr, const ReadOptions_t* options,
                                         ColumnFamilyHandle_t* column_family);
Iterator_t DBNewIterator(DB_t* dbptr, const ReadOptions_t* options);
Status_t DBNewIterators(DB_t* dbptr, const ReadOptions_t* options,
                        const ColumnFamilyHandle_t column_families[],
                        const int size_col,
                        Iterator_t** values,
                        int *val_sz);
Snapshot_t DBGetSnapshot(DB_t* dbptr);
void DBReleaseSnapshot(DB_t* dbptr, const Snapshot_t* snapshot);
bool DBGetPropertyWithColumnFamily(DB_t* dbptr, const ReadOptions_t* options,
                                   ColumnFamilyHandle_t* column_family,
                                   const Slice_t* property, String_t* value);
bool DBGetProperty(DB_t* dbptr, const ReadOptions_t* options,
                   const Slice_t* property, String_t* value);
bool DBGetIntPropertyWithColumnFamily(DB_t* dbptr, 
                                      ColumnFamilyHandle_t* column_family,
                                      const Slice_t* property, uint64_t* value);
bool DBGetIntProperty(DB_t* dbptr, 
                      const Slice_t* property, uint64_t* value);
void DBGetApproximateSizesWithColumnFamily(DB_t* dbptr, 
                                           ColumnFamilyHandle_t* column_family,
                                           const Range_t* range, int n,
                                           uint64_t* sizes);
void DBGetApproximateSizes(DB_t* dbptr, 
                           const Range_t* range, int n,
                           uint64_t* sizes);
Status_t DBCompactRangeWithColumnFamily(DB_t* dbptr, 
                                        ColumnFamilyHandle_t* column_family,
                                        const Slice_t* begin, const Slice_t* end,
                                        bool reduce_level, int target_level,
                                        uint32_t target_path_id);
Status_t DBCompactRange(DB_t* dbptr, 
                        const Slice_t* begin, const Slice_t* end,
                        bool reduce_level, int target_level,
                        uint32_t target_path_id);
Status_t DBSetOptionsWithColumnFamily(DB_t* dbptr, 
                                      ColumnFamilyHandle_t* column_family,
                                      const String_t new_options[],
                                      int n);
Status_t DBSetOptions(DB_t* dbptr, 
                      const String_t new_options[],
                      const int n);
Status_t DBCompactFilesWithColumnFamily(DB_t* dbptr, 
                                        const CompactionOptions_t compact_options,
                                        ColumnFamilyHandle_t* column_family,
                                        const String_t input_file_names[],
                                        const int n,
                                        const int output_level, const int output_path_id);
Status_t DBCompactFiles(DB_t* dbptr, 
                        const CompactionOptions_t compact_options,
                        const String_t input_file_names[],
                        const int n,
                        const int output_level, const int output_path_id);
int DBNumberLevelsWithColumnFamily(DB_t* dbptr, 
                                   ColumnFamilyHandle_t* column_family);
int DBNumberLevels(DB_t* dbptr);
int DBMaxMemCompactionLevelWithColumnFamily(DB_t* dbptr, 
                                            ColumnFamilyHandle_t* column_family);
int DBMaxMemCompactionLevel(DB_t* dbptr);
int DBLevel0StopWriteTriggerWithColumnFamily(DB_t* dbptr, 
                                             ColumnFamilyHandle_t* column_family);
int DBLevel0StopWriteTrigger(DB_t* dbptr);
String_t DBGetName(DB_t* dbptr);
Env_t DBGetEnv(DB_t* dbptr);
Options_t DBGetOptionsWithColumnFamily(DB_t* dbptr, 
                                       ColumnFamilyHandle_t* column_family);
Options_t DBGetOptions(DB_t* dbptr);
DBOptions_t DBGetDBOptions(DB_t* dbptr);
Status_t DBFlushWithColumnFamily(DB_t* dbptr, 
                                 const FlushOptions_t* options,
                                 ColumnFamilyHandle_t* column_family);
Status_t DBFlush(DB_t* dbptr, 
                 const FlushOptions_t* options);
SequenceNumber DBGetLatestSequenceNumber(DB_t* dbptr);
Status_t DBDisableFileDeletions(DB_t* dbptr);
Status_t DBEnableFileDeletions(DB_t* dbptr, bool force);
Status_t DBGetLiveFiles(DB_t* dbptr,
                        const String_t **live_files,
                        int* n,
                        uint64_t* manifest_file_size,
                        bool flush_memtable);
Status_t DBGetSortedWalFiles(DB_t* dbptr, LogFile_t **files, int* n);
Status_t DBGetUpdatesSince(DB_t* dbptr, SequenceNumber seq_number,
                           TransactionLogIterator_t* iter,
                           const TransactionLogIterator_ReadOptions_t* read_options);
Status_t DBDeleteFile(DB_t* dbptr, String_t* name);
void DBGetLiveFilesMetaData(DB_t* dbptr, LiveFileMetaData_t **metadata, int* n);
void DBGetColumnFamilyMetaDataWithColumnFamily(DB_t* dbptr, 
                                               ColumnFamilyHandle_t* column_family,
                                               ColumnFamilyMetaData_t* metadata);
void DBGetColumnFamilyMetaData(DB_t* dbptr, 
                               ColumnFamilyMetaData_t* metadata);
Status_t DBGetDbIdentity(DB_t* dbptr, String_t* identity);
ColumnFamilyHandle_t DBDefaultColumnFamily(DB_t* dbptr);
Status_t DBGetPropertiesOfAllTablesWithColumnFamily(DB_t* dbptr, 
                                                    ColumnFamilyHandle_t* column_family,
                                                    TablePropertiesCollection_t* props);
Status_t DBGetPropertiesOfAllTables(DB_t* dbptr, 
                                    TablePropertiesCollection_t* props);
Status_t DBDestroyDB(const String_t* name, const Options_t* options);
Status_t DBRepairDB(const String_t* dbname, const Options_t* options);


#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_DB_H_
