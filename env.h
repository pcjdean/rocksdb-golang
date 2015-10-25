// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_ENV_H_
#define GO_ROCKSDB_INCLUDE_ENV_H_

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(Env)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Env)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Env)
// Return a default environment suitable for the current operating
// system.  Sophisticated users may wish to provide their own Env
// implementation instead of relying on this default environment.
//
// The result of Default() belongs to rocksdb and must never be deleted.
Env_t NewEnvDefault();

DEFINE_C_WRAP_STRUCT(Logger)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Logger)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Logger)
// Flush to the OS buffers
void LoggerFlush(Logger_t* info_log);
// Return the log level    
int LoggerGetInfoLogLevel(Logger_t* info_log);
// Set the log level. The level lower will not be logged.    
void LoggerSetInfoLogLevel(Logger_t* info_log, int log_level);
    
// a set of log functions with different log levels.
void LoggerHeader(Logger_t* info_log, const char* msg);
void LoggerDebug(Logger_t* info_log, const char* msg);
void LoggerInfo(Logger_t* info_log, const char* msg);
void LoggerWarn(Logger_t* info_log, const char* msg);
void LoggerError(Logger_t* info_log, const char* msg);
void LoggerFatal(Logger_t* info_log, const char* msg);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_ENV_H_
