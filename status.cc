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
#include "cstring.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(Status)
DEFINE_C_WRAP_DESTRUCTOR(Status)
DEFINE_C_WRAP_CONSTRUCTOR_COPY(Status)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY(Status)

// Returns true iff the status indicates success.
bool StatusOk(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->ok() :
            false);
}

// Returns true iff the status indicates a NotFound error.
bool StatusIsNotFound(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsNotFound() :
            false);
}

// Returns true iff the status indicates a Corruption error.
bool StatusIsCorruption(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsCorruption() :
            false);
}

// Returns true iff the status indicates a NotSupported error.
bool StatusIsNotSupported(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsNotSupported() :
            false);
}

// Returns true iff the status indicates an IOError.
bool StatusIsInvalidArgument(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsIOError() :
            false);
}

// Returns true iff the status indicates an MergeInProgress.
bool StatusIsMergeInProgress(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsMergeInProgress() :
            false);
}

// Returns true iff the status indicates Incomplete
bool StatusIsIncomplete(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsIncomplete() :
            false);
}

// Returns true iff the status indicates Shutdown In progress
bool StatusIsShutdownInProgress(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsShutdownInProgress() :
            false);
}

bool StatusIsTimedOut(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsTimedOut() :
            false);
}

bool StatusIsAborted(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsAborted() :
            false);
}

// Returns true iff the status indicates that a resource is Busy and
// temporarily could not be acquired.
bool StatusIsBusy(Status_t *stat)
{
    return ((stat && GET_REP(stat, Status)) ?
            GET_REP(stat, Status)->IsBusy() :
            false);
}

// Return a string representation of this status suitable for printing.
// Returns the string "OK" for success.
String_t StatusToString(Status_t *stat)
{
    String str = GET_REP(stat, Status)->ToString();
    return NewStringTMove(&str);
}

