// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_DB_H_
#define GO_ROCKSDB_INCLUDE_DB_H_

#include "types.h"
#include "slice.h"
#include "options.h"
#include "status.h"
#include "cstring.h"
#include "write_batch.h"
#include "iterator.h"
#include "env.h"
#include "metadata.h"
#include "transaction_log.h"
#include "snapshot.h"
#include "columnFamilyHandle.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(TablePropertiesCollection)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(TablePropertiesCollection)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(TablePropertiesCollection)
DEFINE_C_WRAP_DESTRUCTOR_DEC(TablePropertiesCollection)

DEFINE_C_WRAP_STRUCT(ColumnFamilyDescriptor)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(ColumnFamilyDescriptor)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(ColumnFamilyDescriptor, String, ColumnFamilyOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(ColumnFamilyDescriptor)

DEFINE_C_WRAP_STRUCT(Range)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Range)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(Range, Slice, Slice)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Range)
Range_t NewRangeTFromSlices(Slice_t* start, Slice_t* limit);

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
                              String_t **column_families, int* size_col);
Status_t DBCreateColumnFamily(const DB_t* dbptr, const ColumnFamilyOptions_t* options,
                            const String_t* column_family_name,
                            ColumnFamilyHandle_t* handle);
Status_t DBDropColumnFamily(const DB_t* dbptr, const ColumnFamilyHandle_t* column_family);
Status_t DBPutWithColumnFamily(const DB_t* dbptr, const WriteOptions_t* options,
                           const ColumnFamilyHandle_t* column_family,
                           const Slice_t* key,
                           const Slice_t* value);
Status_t DBPut(const DB_t* dbptr, const WriteOptions_t* optionss,
               const Slice_t* key,
               const Slice_t* value);
Status_t DBDeleteWithColumnFamily(const DB_t* dbptr, const WriteOptions_t* options,
                                  const ColumnFamilyHandle_t* column_family,
                                  const Slice_t* key);
Status_t DBDelete(const DB_t* dbptr, const WriteOptions_t* optionss, const Slice_t* key);
Status_t DBMergeWithColumnFamily(const DB_t* dbptr, const WriteOptions_t* options,
                                 const ColumnFamilyHandle_t* column_family,
                                 const Slice_t* key,
                                 const Slice_t* value);
Status_t DBMerge(const DB_t* dbptr, const WriteOptions_t* optionss,
               const Slice_t* key,
               const Slice_t* value);
Status_t DBWrite(const DB_t* dbptr, const WriteOptions_t* optionss, WriteBatch_t* updates);
Status_t DBGetWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                               const ColumnFamilyHandle_t* column_family,
                               const Slice_t* key,
                               const String_t* value);
Status_t DBGet(const DB_t* dbptr, const ReadOptions_t* options,
               const Slice_t* key,
               const String_t* value);

Status_t* DBMultiGetWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                                     const ColumnFamilyHandle_t column_families[],
                                     const int size_col,
                                     const Slice_t keys[],
                                     const int size_keys,
                                     String_t** values);
Status_t* DBMultiGet(const DB_t* dbptr, const ReadOptions_t* options,
                     const Slice_t keys[],
                     const int size_keys,
                     String_t** values);
bool DBKeyMayExistWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                                   const ColumnFamilyHandle_t* column_family,
                                   const Slice_t* key,
                                   String_t* value,
                                   bool* value_found);
bool DBKeyMayExist(const DB_t* dbptr, const ReadOptions_t* options,
                   const Slice_t* key,
                   String_t* value,
                   bool* value_found);
Iterator_t DBNewIteratorWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                                         const ColumnFamilyHandle_t* column_family);
Iterator_t DBNewIterator(const DB_t* dbptr, const ReadOptions_t* options);
Status_t DBNewIterators(const DB_t* dbptr, const ReadOptions_t* options,
                        const ColumnFamilyHandle_t column_families[],
                        const int size_col,
                        Iterator_t** values,
                        int *val_sz);
Snapshot_t DBGetSnapshot(const DB_t* dbptr);
void DBReleaseSnapshot(const DB_t* dbptr, const Snapshot_t* snapshot);
bool DBGetPropertyWithColumnFamily(const DB_t* dbptr, const ColumnFamilyHandle_t* column_family,
                                   const Slice_t* property, String_t* value);
bool DBGetProperty(const DB_t* dbptr, const Slice_t* property, String_t* value);
bool DBGetIntPropertyWithColumnFamily(const DB_t* dbptr, 
                                      const ColumnFamilyHandle_t* column_family,
                                      const Slice_t* property, uint64_t* value);
bool DBGetIntProperty(const DB_t* dbptr, 
                      const Slice_t* property, uint64_t* value);
void DBGetApproximateSizesWithColumnFamily(const DB_t* dbptr, 
                                           const ColumnFamilyHandle_t* column_family,
                                           const Range_t* range, int n,
                                           uint64_t* sizes);
void DBGetApproximateSizes(const DB_t* dbptr, 
                           const Range_t* range, int n,
                           uint64_t* sizes);
Status_t DBCompactRangeWithColumnFamily(const DB_t* dbptr, 
                                        const CompactRangeOptions_t* compact_range_options,
                                        const ColumnFamilyHandle_t* column_family,
                                        const Slice_t* begin, const Slice_t* end);
Status_t DBCompactRange(const DB_t* dbptr, 
                        const CompactRangeOptions_t* compact_range_options,
                        const Slice_t* begin, const Slice_t* end);
Status_t DBSetOptionsWithColumnFamily(const DB_t* dbptr, 
                                      const ColumnFamilyHandle_t* column_family,
                                      const String_t new_options[],
                                      int n);
Status_t DBSetOptions(const DB_t* dbptr, 
                      const String_t new_options[],
                      const int n);
Status_t DBCompactFilesWithColumnFamily(const DB_t* dbptr, 
                                        const CompactionOptions_t* compact_options,
                                        const ColumnFamilyHandle_t* column_family,
                                        const String_t input_file_names[],
                                        const int n,
                                        const int output_level, const int output_path_id);
Status_t DBCompactFiles(const DB_t* dbptr, 
                        const CompactionOptions_t* compact_options,
                        const String_t input_file_names[],
                        const int n,
                        const int output_level, const int output_path_id);
int DBNumberLevelsWithColumnFamily(const DB_t* dbptr, 
                                   const ColumnFamilyHandle_t* column_family);
int DBNumberLevels(const DB_t* dbptr);
int DBMaxMemCompactionLevelWithColumnFamily(const DB_t* dbptr, 
                                            const ColumnFamilyHandle_t* column_family);
int DBMaxMemCompactionLevel(const DB_t* dbptr);
int DBLevel0StopWriteTriggerWithColumnFamily(const DB_t* dbptr, 
                                             const ColumnFamilyHandle_t* column_family);
int DBLevel0StopWriteTrigger(const DB_t* dbptr);
String_t DBGetName(const DB_t* dbptr);
Env_t DBGetEnv(const DB_t* dbptr);
Options_t DBGetOptionsWithColumnFamily(const DB_t* dbptr, 
                                       const ColumnFamilyHandle_t* column_family);
Options_t DBGetOptions(const DB_t* dbptr);
DBOptions_t DBGetDBOptions(const DB_t* dbptr);
Status_t DBFlushWithColumnFamily(const DB_t* dbptr, 
                                 const FlushOptions_t* options,
                                 const ColumnFamilyHandle_t* column_family);
Status_t DBFlush(const DB_t* dbptr, 
                 const FlushOptions_t* options);
SequenceNumber DBGetLatestSequenceNumber(const DB_t* dbptr);
Status_t DBDisableFileDeletions(const DB_t* dbptr);
Status_t DBEnableFileDeletions(const DB_t* dbptr, bool force);
Status_t DBGetLiveFiles(const DB_t* dbptr,
                        String_t **live_files,
                        int* n,
                        uint64_t* manifest_file_size,
                        bool flush_memtable);
Status_t DBGetSortedWalFiles(const DB_t* dbptr, LogFile_t **files, int* n);
Status_t DBGetUpdatesSince(const DB_t* dbptr, SequenceNumber seq_number,
                           TransactionLogIterator_t* iter,
                           const TransactionLogIterator_ReadOptions_t* read_options);
Status_t DBDeleteFile(const DB_t* dbptr, String_t* name);
void DBGetLiveFilesMetaData(const DB_t* dbptr, LiveFileMetaData_t **metadata, int* n);
void DBGetColumnFamilyMetaDataWithColumnFamily(const DB_t* dbptr, 
                                               const ColumnFamilyHandle_t* column_family,
                                               ColumnFamilyMetaData_t* metadata);
void DBGetColumnFamilyMetaData(const DB_t* dbptr, 
                               ColumnFamilyMetaData_t* metadata);
Status_t DBGetDbIdentity(const DB_t* dbptr, String_t* identity);
ColumnFamilyHandle_t DBDefaultColumnFamily(const DB_t* dbptr);
Status_t DBGetPropertiesOfAllTablesWithColumnFamily(const DB_t* dbptr, 
                                                    const ColumnFamilyHandle_t* column_family,
                                                    TablePropertiesCollection_t* props);
Status_t DBGetPropertiesOfAllTables(const DB_t* dbptr, 
                                    TablePropertiesCollection_t* props);
Status_t DBDestroyDB(const String_t* name, const Options_t* options);
Status_t DBRepairDB(const String_t* dbname, const Options_t* options);


// Return the major version of DB.
int DBGetMajorVersion();
// Return the minor version of DB.
int DBGetMinorVersion();

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_DB_H_
