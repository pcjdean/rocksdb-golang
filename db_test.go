// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

import (
	"os"
	"fmt"
	"bytes"
	"testing"
)

func (db *DB) checkGet(t *testing.T, ropts *ReadOptions, key []byte, valExp []byte) {
	val, stat := db.Get(ropts, key)
	if nil != valExp && !stat.Ok() {
		t.Fatalf("get err: stat = %s", stat)
	}

	if nil == valExp && !stat.IsNotFound() {
		t.Fatalf("get err: stat = %s", stat)
	} else if !bytes.Equal(val, valExp) {
		t.Fatalf("get err: expected = %v, got = %v", valExp, val)
	}
}

// Test from rocksdb's c_test.c.
func TestCMain(t *testing.T) {
	var (
		stat *Status
		db *DB
	)

	dbname := fmt.Sprintf("%s/rocksdb_go_test-%d", os.TempDir(), os.Geteuid);
	dbbackupname := fmt.Sprintf("%s/rocksdb_go_test-%d-backup", os.TempDir(), os.Geteuid);
	t.Log("create_objects")
	fmt.Printf("dbname = %s\n", dbname)
	fmt.Printf("dbbackupname = %s\n", dbbackupname)
	options := NewOptions()

	t.Log("Destroy")
	stat = DestroyDB(&dbname, options)
	if stat.Ok() {
		t.Error("DestroyDB")
	} else {
		t.Logf("DestroyDB: status = ", stat)
	}

	t.Log("open_error")
	_, stat, _ = Open(&dbname, options)
	if stat.Ok() {
		t.Error("open_error")
	} else {
		t.Logf("not open_error: ststus = ", stat)
	}

	t.Log("open")
	options.SetCreateIfMissing(true)
	db, stat, _ = Open(&dbname, options)
	if !stat.Ok() {
		t.Fatalf("open: stat = %s", stat)
	}
	t.Log("get")
	ropts := NewReadOptions()
	db.checkGet(t, ropts, []byte("foo"), nil)

	t.Log("put")
	woptions := NewWriteOptions()
	woptions.SetSync(true)
	db.Put(woptions, []byte("foo"), []byte("hello"))
	if !stat.Ok() {
		t.Fatalf("put err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))
}
