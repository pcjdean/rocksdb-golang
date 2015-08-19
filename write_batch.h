// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_WRITE_BATCH_H_
#define GO_ROCKSDB_INCLUDE_WRITE_BATCH_H_
//
// WriteBatch::rep_ :=
//    sequence: fixed64
//    count: fixed32
//    data: record[count]
// record :=
//    kTypeValue varstring varstring
//    kTypeMerge varstring varstring
//    kTypeDeletion varstring
//    kTypeColumnFamilyValue varint32 varstring varstring
//    kTypeColumnFamilyMerge varint32 varstring varstring
//    kTypeColumnFamilyDeletion varint32 varstring varstring
// varstring :=
//    len: varint32
//    data: uint8[len]

#include "types.h"
#include "columnFamilyHandle.h"
#include "slice.h"
#include "cstring.h"


#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(WriteBatch)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(WriteBatch)
DEFINE_C_WRAP_DESTRUCTOR_DEC(WriteBatch)
DEFINE_C_WRAP_CONSTRUCTOR_COPY_DEC(WriteBatch)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(WriteBatch)

void WriteBatchPutWithColumnFamily(const WriteBatch_t* write_batch,
                                   const ColumnFamilyHandle_t* column_family,
                                   const Slice_t* key, const Slice_t* value);
void WriteBatchPut(const WriteBatch_t* write_batch, 
                   const Slice_t* key, const Slice_t* value);
    
void WriteBatchMergeWithColumnFamily(const WriteBatch_t* write_batch,
                                     const ColumnFamilyHandle_t* column_family,
                                     const Slice_t* key, const Slice_t* value);
void WriteBatchMerge(const WriteBatch_t* write_batch, 
                     const struct Slice_t* key, const Slice_t* value);

void WriteBatchDeleteWithColumnFamily(const WriteBatch_t* write_batch,
                                      const ColumnFamilyHandle_t* column_family,
                                      const Slice_t* key);
void WriteBatchDelete(const WriteBatch_t* write_batch, const Slice_t* key);
    
void WriteBatchClear(const WriteBatch_t* write_batch);

String_t WriteBatchData(const WriteBatch_t* write_batch);
size_t WriteBatchGetDataSize(const WriteBatch_t* write_batch);
int WriteBatchCount(const WriteBatch_t* write_batch);

WriteBatch_t NewWriteBatchTRawArgs(const String_t* str);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_WRITE_BATCH_H_
