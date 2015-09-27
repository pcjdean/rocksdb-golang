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

package rocksdb

/*
#include <stdlib.h>
#include "env.h"
*/
import "C"

import (
	"runtime"
	"fmt"
	"unsafe"
)

const ( 
	DEBUG_LEVEL = iota
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
	FATAL_LEVEL
	NUM_INFO_LOG_LEVELS
)

// Wrap go Env
type Env struct {
	env C.Env_t
}

// Release resources
func (env *Env) finalize() {
	var cenv *C.Env_t = &env.env
	C.DeleteEnvT(cenv, toCBool(false))
}

// C env to go env
func (cenv *C.Env_t) toEnv() (env *Env) {
	env = &Env{env: *cenv}	
	runtime.SetFinalizer(env, finalize)
	return
}

// Wrap go Logger
type Logger struct {
	log C.Logger_t
}

// Release resources
func (log *Logger) finalize() {
	var clog *C.Logger_t = &log.log
	C.DeleteLoggerT(clog, toCBool(false))
}

// C logger to go logger. Delete the underlying c++ wrap object
// if del is true. 
func (clog *C.Logger_t) toLogger(del bool) (log *Logger) {
	log = &Logger{log: *clog}
	if del {	
		runtime.SetFinalizer(log, finalize)
	}
	return
}

// Flush to the OS buffers
func (log *Logger) Flush() {
	C.LoggerFlush(&log.log)
}

// Return the log level    
func (log *Logger) GetInfoLogLevel() {
	C.LoggerGetInfoLogLevel(&log.log)
}

// Set the log level. The level lower will not be logged.    
func (log *Logger) SetInfoLogLevel(level int) {
	C.LoggerSetInfoLogLevel(&log.log, C.int(level))
}

// a set of log functions with different log levels.
func (log *Logger) Header(format string, a ...interface{}) {
	str := C.CString(fmt.Sprintf(format, a...))
	defer C.free(unsafe.Pointer(str))
	C.LoggerHeader(&log.log, str)
}

// log functions with Debug log levels.
func (log *Logger) Debug(format string, a ...interface{}) {
	str := C.CString(fmt.Sprintf(format, a...))
	defer C.free(unsafe.Pointer(str))
	C.LoggerDebug(&log.log, str)
}

// log functions with Info log levels.
func (log *Logger) Info(format string, a ...interface{}) {
	str := C.CString(fmt.Sprintf(format, a...))
	defer C.free(unsafe.Pointer(str))
	C.LoggerInfo(&log.log, str)
}

// log functions with Error log levels.
func (log *Logger) Error(format string, a ...interface{}) {
	str := C.CString(fmt.Sprintf(format, a...))
	defer C.free(unsafe.Pointer(str))
	C.LoggerError(&log.log, str)
}

// log functions with Fatal log levels.
func (log *Logger) Fatal(format string, a ...interface{}) {
	str := C.CString(fmt.Sprintf(format, a...))
	defer C.free(unsafe.Pointer(str))
	C.LoggerFatal(&log.log, str)
}
