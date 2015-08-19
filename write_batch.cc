// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// WriteBatch holds a collection of updates to apply atomically to a DB.
//
// The updates are applied in the order in which they are added
// to the WriteBatch.  For example, the value of "key" will be "v3"
// after the following batch is written:
//
//    batch.Put("key", "v1");
//    batch.Delete("key");
//    batch.Put("key", "v2");
//    batch.Put("key", "v3");
//
// Multiple threads can invoke const methods on a WriteBatch without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same WriteBatch must use
// external synchronization.

#include <rocksdb/write_batch.h>
#include "db.h"
#include "write_batch.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(WriteBatch)
DEFINE_C_WRAP_DESTRUCTOR(WriteBatch)
DEFINE_C_WRAP_CONSTRUCTOR_COPY(WriteBatch)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(WriteBatch)

// Store the mapping "key->value" in the database.
void WriteBatchPutWithColumnFamily(const WriteBatch_t* write_batch,
                                   const ColumnFamilyHandle_t* column_family,
                                   const Slice_t* key, const Slice_t* value)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        GET_REP(write_batch, WriteBatch)->Put(GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice), GET_REP_REF(value, Slice));
    }
}

void WriteBatchPut(const WriteBatch_t* write_batch, 
                   const Slice_t* key, const Slice_t* value)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        GET_REP(write_batch, WriteBatch)->Put(GET_REP_REF(key, Slice), GET_REP_REF(value, Slice));
    }
}

// Merge "value" with the existing value of "key" in the database.
// "key->merge(existing, value)"
void WriteBatchMergeWithColumnFamily(const WriteBatch_t* write_batch,
                                     const ColumnFamilyHandle_t* column_family,
                                     const Slice_t* key, const Slice_t* value)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        GET_REP(write_batch, WriteBatch)->Merge(GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice), GET_REP_REF(value, Slice));
    }
}

void WriteBatchMerge(const WriteBatch_t* write_batch, 
                     const Slice_t* key, const Slice_t* value)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        GET_REP(write_batch, WriteBatch)->Merge(GET_REP_REF(key, Slice), GET_REP_REF(value, Slice));
    }
}

// If the database contains a mapping for "key", erase it.  Else do nothing.
void WriteBatchDeleteWithColumnFamily(const WriteBatch_t* write_batch,
                                      const ColumnFamilyHandle_t* column_family,
                                      const Slice_t* key)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        GET_REP(write_batch, WriteBatch)->Delete(GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice));
    }
}

void WriteBatchDelete(const WriteBatch_t* write_batch, const Slice_t* key)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        GET_REP(write_batch, WriteBatch)->Delete(GET_REP_REF(key, Slice));
    }
}

// Clear all updates buffered in this batch.
void WriteBatchClear(const WriteBatch_t* write_batch)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        GET_REP(write_batch, WriteBatch)->Clear();
    }
}

// Retrieve the serialized version of this batch.
String_t WriteBatchData(const WriteBatch_t* write_batch)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        const String &str = GET_REP(write_batch, WriteBatch)->Data();
        return NewStringTCopy(const_cast<String *>(&str));
    }
    
    return NewStringTDefault();
}

// Retrieve data size of the batch.
size_t WriteBatchGetDataSize(const WriteBatch_t* write_batch)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        GET_REP(write_batch, WriteBatch)->GetDataSize();
    }
}

// Returns the number of updates in the batch
int WriteBatchCount(const WriteBatch_t* write_batch)
{
    if (write_batch)
    {
        assert(GET_REP(write_batch, WriteBatch) != NULL);
        GET_REP(write_batch, WriteBatch)->Count();
    }
}

WriteBatch_t NewWriteBatchTRawArgs(const String_t* str)
{
    WriteBatch_t wrap;
    
    if (str)
    {
        assert(GET_REP(str, String) != NULL);
        wrap.rep = new WriteBatch(GET_REP_REF(str, String));
    }
    else
    {
        wrap.rep = new WriteBatch();
    }

    return wrap;
}
