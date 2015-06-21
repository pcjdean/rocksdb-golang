// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

type Finalizer interface {
	func Finalize()
}

func Finalize(obj *Finalizer) {
	obj.Finalize()
}
