// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// A Status encapsulates the result of an operation.  It may indicate success,
// or it may indicate an error with an associated error message.
//
// Multiple threads can invoke const methods on a Status without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Status must use
// external synchronization.

#include <rocksdb/status.h>
#include "status.h"

DEFINE_C_WRAP_CONSTRUCTOR(Status)
DEFINE_C_WRAP_DESTRUCTOR(Status)
DEFINE_C_WRAP_CONSTRUCTOR_COPY(Status)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(Status)

// Returns true iff the status indicates success.
inline bool StatusOk(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->ok() :
            false);
}

// Returns true iff the status indicates a NotFound error.
inline bool StatusIsNotFound(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsNotFound() :
            false);
}

// Returns true iff the status indicates a Corruption error.
inline bool StatusIsCorruption(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsCorruption() :
            false);
}

// Returns true iff the status indicates a NotSupported error.
inline bool StatusIsNotSupported(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsNotSupported() :
            false);
}

// Returns true iff the status indicates an IOError.
inline bool StatusIsInvalidArgument(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsIOError() :
            false);
}

// Returns true iff the status indicates an MergeInProgress.
inline bool StatusIsMergeInProgress(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsMergeInProgress() :
            false);
}

// Returns true iff the status indicates Incomplete
inline bool StatusIsIncomplete(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsIncomplete() :
            false);
}

// Returns true iff the status indicates Shutdown In progress
inline bool StatusIsShutdownInProgress(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsShutdownInProgress() :
            false);
}

inline bool StatusIsTimedOut(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsTimedOut() :
            false);
}

inline bool StatusIsAborted(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsAborted() :
            false);
}

// Returns true iff the status indicates that a resource is Busy and
// temporarily could not be acquired.
inline bool StatusIsBusy(Status_t *stat)
{
    return ((stat && GET_REP(stat)) ?
            GET_REP(stat)->IsBusy() :
            false);
}

// Return a string representation of this status suitable for printing.
// Returns the string "OK" for success.
inline String_t StatusToString(Status_t *stat) const
{
    return NewStringTmove(stat->ToString());
}

