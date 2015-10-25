// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// An Env is an interface used by the rocksdb implementation to access
// operating system functionality like the filesystem etc.  Callers
// may wish to provide a custom Env object when opening a database to
// get fine gain control; e.g., to rate limit file system operations.
//
// All Env implementations are safe for concurrent access from
// multiple threads without any external synchronization.

#include <rocksdb/env.h>
#include "env.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(Env)
DEFINE_C_WRAP_DESTRUCTOR(Env)
// Return a default environment suitable for the current operating
// system.  Sophisticated users may wish to provide their own Env
// implementation instead of relying on this default environment.
//
// The result of Default() belongs to rocksdb and must never be deleted.
Env_t NewEnvDefault()
{
    Env_t ret;
    ret.rep = Env::Default();
    return ret;
}


DEFINE_C_WRAP_CONSTRUCTOR(Logger)
DEFINE_C_WRAP_DESTRUCTOR(Logger)

// Flush to the OS buffers
void LoggerFlush(Logger_t* info_log)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        GET_REP(info_log, Logger)->Flush();
    }
}

// Return the log level    
int LoggerGetInfoLogLevel(Logger_t* info_log)
{
    int ret = INFO_LEVEL;
    
    if (info_log && GET_REP(info_log, Logger))
    {
        ret = GET_REP(info_log, Logger)->GetInfoLogLevel();
    }

    return ret;
}

// Set the log level. The level lower will not be logged.    
void LoggerSetInfoLogLevel(Logger_t* info_log, int log_level)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        GET_REP(info_log, Logger)->SetInfoLogLevel(InfoLogLevel(log_level));
    }
}

// a set of log functions with different log levels.
void LoggerHeader(Logger_t* info_log, const char* msg)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        Header(GET_REP(info_log, Logger), msg);
    }
}

// log functions with Debug log levels.
void LoggerDebug(Logger_t* info_log, const char* msg)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        Debug(GET_REP(info_log, Logger), msg);
    }
}

// log functions with Info log levels.
void LoggerInfo(Logger_t* info_log, const char* msg)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        Info(GET_REP(info_log, Logger), msg);
    }
}

// log functions with Warn log levels.
void LoggerWarn(Logger_t* info_log, const char* msg)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        Warn(GET_REP(info_log, Logger), msg);
    }
}

// log functions with Error log levels.
void LoggerError(Logger_t* info_log, const char* msg)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        Error(GET_REP(info_log, Logger), msg);
    }
}

// log functions with Fatal log levels.
void LoggerFatal(Logger_t* info_log, const char* msg)
{
    if (info_log && GET_REP(info_log, Logger))
    {
        Fatal(GET_REP(info_log, Logger), msg);
    }
}

