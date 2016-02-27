// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "db.h"
*/
import "C"

const (
	// Major version of go DB
	majorVersionGo int = 1
	// Minor version of go DB
	minorVersionGo int = 0;
)

var (
	// Major version of C++ DB - 4
	majorVersion int = int(C.DBGetMajorVersion());
	// Minor version of C++ DB - 5
	minorVersion int = int(C.DBGetMinorVersion());
)
