// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// An iterator yields a sequence of key/value pairs from a source.
// The following class defines the interface.  Multiple implementations
// are provided by this library.  In particular, iterators are provided
// to access the contents of a Table or a DB.
//
// Multiple threads can invoke const methods on an Iterator without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Iterator must use
// external synchronization.

#include <rocksdb/iterator.h>
#include <rocksdb/slice.h>
#include "iterator.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(Iterator)
DEFINE_C_WRAP_DESTRUCTOR(Iterator)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(Iterator)

// An iterator is either positioned at a key/value pair, or
// not valid.  This method returns true iff the iterator is valid.
bool IteratorValid(Iterator_t *it)
{
    return ((it && GET_REP(it, Iterator)) ?
            GET_REP(it, Iterator)->Valid() :
            false);
}

// Position at the first key in the source.  The iterator is Valid()
// after this call iff the source is not empty.
void IteratorSeekToFirst(Iterator_t *it)
{
    if (it && GET_REP(it, Iterator))
    {
        GET_REP(it, Iterator)->SeekToFirst();
    }
}

// Position at the last key in the source.  The iterator is
// Valid() after this call iff the source is not empty.
void IteratorSeekToLast(Iterator_t *it)
{
    if (it && GET_REP(it, Iterator))
    {
        GET_REP(it, Iterator)->SeekToLast();
    }
}

// Position at the first key in the source that at or past target
// The iterator is Valid() after this call iff the source contains
// an entry that comes at or past target.
void IteratorSeek(Iterator_t *it, const Slice_t* target)
{
    if (it && GET_REP(it, Iterator) &&
        target && GET_REP(target, Slice))
    {
        GET_REP(it, Iterator)->Seek(GET_REP_REF(target, Slice));
    }
}

// Moves to the next entry in the source.  After this call, Valid() is
// true iff the iterator was not positioned at the last entry in the source.
// REQUIRES: Valid()
void IteratorNext(Iterator_t *it)
{
    if (it && GET_REP(it, Iterator))
    {
        GET_REP(it, Iterator)->Next();
    }
}

// Moves to the previous entry in the source.  After this call, Valid() is
// true iff the iterator was not positioned at the first entry in source.
// REQUIRES: Valid()
void IteratorPrev(Iterator_t *it)
{
    if (it && GET_REP(it, Iterator))
    {
        GET_REP(it, Iterator)->Prev();
    }
}

// Return the key for the current entry.  The underlying storage for
// the returned slice is valid only until the next modification of
// the iterator.
// REQUIRES: Valid()
Slice_t IteratorKey(Iterator_t *it)
{
    Slice slc;
    if (it && GET_REP(it, Iterator))
    {
        slc = GET_REP(it, Iterator)->key();
    }
    return NewSliceTRawArgs(slc.data(), slc.size());
}

// Return the value for the current entry.  The underlying storage for
// the returned slice is valid only until the next modification of
// the iterator.
// REQUIRES: !AtEnd() && !AtStart()
Slice_t IteratorValue(Iterator_t *it)
{
    Slice slc;
    if (it && GET_REP(it, Iterator))
    {
        slc = GET_REP(it, Iterator)->value();
    }
    return NewSliceTRawArgs(slc.data(), slc.size());
}

// If an error has occurred, return it.  Else return an ok status.
// If non-blocking IO is requested and this operation cannot be
// satisfied without doing some IO, then this returns Status::Incomplete().
Status_t IteratorStatus(Iterator_t *it)
{
    Status ret;
    if (it && GET_REP(it, Iterator))
    {
        ret = GET_REP(it, Iterator)->status();
    }
    else
        ret = invalid_status;
    return NewStatusTCopy(&ret);
}
