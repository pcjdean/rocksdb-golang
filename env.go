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
#include "env.h"
*/
import "C"

import (
	"runtime"
)

type Env struct {
	env C.Env_t
}

func (env *Env) finalize() {
	var cenv *C.Env_t = &env.env
	C.DeleteEnvT(cenv, toCBool(false))
}

func (cenv *C.Env_t) toEnv() (env *Env) {
	env = &Env{env: *cenv}	
	runtime.SetFinalizer(env, finalize)
	return
}
