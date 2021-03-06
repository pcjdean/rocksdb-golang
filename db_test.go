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

func (db *DB) checkGet(t *testing.T, ropts *ReadOptions, key []byte, valExp []byte, cfh ...*ColumnFamilyHandle) {
	val, stat := db.Get(ropts, key, cfh...)
	if nil != valExp && !stat.Ok() {
		t.Fatalf("err: get err: stat = %s", stat)
	}

	if nil == valExp && !stat.IsNotFound() {
		t.Fatalf("err: get nil err: stat = %s", stat)
	} else if !bytes.Equal(val, valExp) {
		t.Fatalf("err: get err: expected = %s, got = %s", valExp, val)
	}
}

func checkCompaction(t *testing.T, dbname *string, options *Options, roptions *ReadOptions, woptions *WriteOptions) (db *DB) {
	var stat *Status
	db, stat, _ = Open(options, dbname)
	if !stat.Ok() {
		t.Fatalf("checkCompaction err: open: stat = %s", stat)
	}

	stat = db.Put(woptions, []byte("foo"), []byte("foovalue"))
	if !stat.Ok() {
		t.Fatalf("checkCompaction err: put err: stat = %s", stat)
	}
	db.checkGet(t, roptions, []byte("foo"), []byte("foovalue"))

	stat = db.Put(woptions, []byte("bar"), []byte("barvalue"))
	if !stat.Ok() {
		t.Fatalf("checkCompaction err: put err: stat = %s", stat)
	}
	db.checkGet(t, roptions, []byte("bar"), []byte("barvalue"))

	stat = db.Put(woptions, []byte("baz"), []byte("bazvalue"))
	if !stat.Ok() {
		t.Fatalf("checkCompaction err: put err: stat = %s", stat)
	}
	db.checkGet(t, roptions, []byte("baz"), []byte("bazvalue"))

	// Force compaction
	cropt := NewCompactRangeOptions()
	db.CompactRange(cropt, nil, nil)

	// should have filtered bar, but not foo
	db.checkGet(t, roptions, []byte("foo"), []byte("foovalue"))
	db.checkGet(t, roptions, []byte("bar"), nil)
	db.checkGet(t, roptions, []byte("baz"), []byte("newbazvalue"))
	return
}

func checkCondition(t *testing.T, cond bool) {
	if !cond {
		t.Fatal("err: checkCondition:")
	}
}

func checkIter(t *testing.T, iter *Iterator, key, val string) {
	str := iter.Key();
	checkCondition(t, bytes.Equal(str, []byte(key)));
	str = iter.Value();
	checkCondition(t, bytes.Equal(str, []byte(val)));
}

// Custom Comparator filter
type testComparator struct {
	t *testing.T
}

func (tcf *testComparator) Name() string {
	return "foo"
}
	
func (tcf *testComparator) Compare(a, b []byte) int {
	ret := bytes.Compare(a, b)
	// tcf.t.Logf("testComparator a = %s, b = %s, ret = %v", a, b, ret)
	return ret
}
	
func (tcf *testComparator) FindShortestSeparator(start, limit []byte) []byte {
	return nil
}
	
func (tcf *testComparator) FindShortSuccessor(key []byte) []byte {
	return nil
}

// Custom Merge Operator
type testMergeOperator struct {
	t *testing.T
}

func (tcf *testMergeOperator) Name() string {
	return "TestMergeOperator"
}
	
func (tcf *testMergeOperator) FullMerge(key []byte, exval []byte, opdlist [][]byte, logger *Logger) (suc bool, newval []byte) {
	// tcf.t.Logf("testMergeOperator::FullMerge key = %s, exval = %s, opdlist = %s", key, exval, opdlist)
	newval = []byte("fake")
	suc = true
	return
}
	
func (tcf *testMergeOperator) PartialMerge(key []byte, leftopd []byte, rightopd []byte, logger *Logger) (suc bool, newval []byte) {
	// tcf.t.Logf("testMergeOperator::PartialMerge key = %s, leftopd = %s, rightopd = %s", key, leftopd, rightopd)
	newval = []byte("fake")
	suc = true
	return
}
	
func (tcf *testMergeOperator) PartialMergeMulti(key []byte, opdlist [][]byte, logger *Logger) (suc bool, newval []byte) {
	// tcf.t.Logf("testMergeOperator::PartialMergeMulti key = %s, opdlist = %s", key, opdlist)
	suc = false
	return
}

type testFilterPolicy struct {
	t *testing.T
	fakeResult bool
}

// Custom FilterPolicy
func (tfp testFilterPolicy) Name() string {
	return "testFilterPolicy"
}

func (tfp testFilterPolicy) CreateFilter(keys [][]byte) []byte {
	// tfp.t.Logf("CreateFilter keys = %s", keys)
	return []byte("fake")
}

func (tfp testFilterPolicy) KeyMayMatch(key, filter []byte) bool {
	// tfp.t.Logf("KeyMayMatch keys = %s, filter = %s, tfp.fakeResult = %v", key, filter, tfp.fakeResult)
	checkCondition(tfp.t, 4 == len(filter))
	checkCondition(tfp.t, bytes.Equal(filter, []byte("fake")));
	return tfp.fakeResult
}

func (tfp testFilterPolicy) GetFilterBitsBuilder() IFilterBitsBuilder {
	return nil
}

func (tfp testFilterPolicy) GetFilterBitsReader() IFilterBitsReader {
	return nil
}

// Custom compaction filter
type testCompactionFilter struct {
	t *testing.T
}

func (tcf *testCompactionFilter) Name() string {
	return "foo"
}

func (tcf *testCompactionFilter) Filter(level int, key, exval []byte) (newval []byte, valchg bool, removed bool) {
	if 3 == len(key) {
		if  bytes.Equal([]byte("bar"), key) {
			removed = true
		} else if bytes.Equal([]byte("baz"), key) {
			valchg = true;
			newval = []byte("newbazvalue");
			removed = false
		}
	}
	// tcf.t.Logf("testCompactionFilter level = %v, key = %s, exval = %s, newval = %s, valchg = %v, removed = %v", level, string(key), string(exval), string(newval), valchg, removed)
	return
}

// Custom CompactionFilterFactory filter
type testCompactionFilterFactory struct {
	t *testing.T
}

func (tcff *testCompactionFilterFactory) Name() string {
	return "foo"
}

func (tcff *testCompactionFilterFactory) CreateCompactionFilter(context *CompactionFilter_Context) ICompactionFilter {
	// tcff.t.Logf("testCompactionFilterFactory context = %v", context)
	cff := &testCompactionFilter{t: tcff.t}
	return cff
}

// Custom SliceTransform filter
type testSliceTransform struct {
	t *testing.T
}

func (tstf *testSliceTransform) Name() string {
	return "TestCFV2PrefixExtractor"
}

func (tstf *testSliceTransform) Transform(src []byte) (offset, sz uint64) {
	// Verify keys are maximum length 4; this verifies fix for a
	// prior bug which was passing the RocksDB-encoded key with
	// logical timestamp suffix instead of parsed user key.
	l := uint64(len(src))
	if 4 < l {
		tstf.t.Fatalf("testSliceTransform::Transform - key %v is not user key\n", string(src))
	}
	if 3 > l {
		sz = l
	} else {
		sz = 3
	}
	offset = 0
	return 
}

func (tstf *testSliceTransform) InDomain(src []byte) bool {
	return true
}

func (tstf *testSliceTransform) InRange(dst []byte) bool {
	return true
}

func (tstf *testSliceTransform) SameResultWhenAppended(prefix []byte) bool {
	return false
}

// Test from rocksdb's c_test.c.
func TestCMain(t *testing.T) {
	var (
		stat *Status
		db *DB
	)

	dbname := fmt.Sprintf("%s/rocksdb_go_test-%d", os.TempDir(), os.Geteuid);
	dbbackupname := fmt.Sprintf("%s/rocksdb_go_test-%d-backup", os.TempDir(), os.Geteuid);
	fmt.Printf("rocksdbgo version = %d.%d\n", majorVersionGo, minorVersionGo)
	fmt.Printf("rocksdb version = %d.%d\n", majorVersion, minorVersion)
	t.Log("phase: create_objects")
	fmt.Printf("dbname = %s\n", dbname)
	fmt.Printf("dbbackupname = %s\n", dbbackupname)
	
	icmp := &testComparator{t: t}
	cmp := NewComparator(icmp)
	env := NewEnvDefault();
	cache := NewLRUCache(100000)

	options := NewOptions()

	options.SetComparator(cmp)
	options.SetErrorIfExists(true)
	options.SetEnv(env)
	options.SetInfoLog(NewPLoggerDefault())
	options.SetWriteBufferSize(100000)
	options.SetParanoidChecks(true)
	options.SetMaxOpenFiles(10)
	table_options := NewBlockBasedTableOptions()
	table_options.SetBlockCache(cache)
	options.SetTableFactory(table_options.NewBlockBasedTableFactory())

	options.SetCompression(NoCompression)
	options.SetCompressionOptions(-14, -1, 0)
	compressionLevels := []int{
		NoCompression, NoCompression, NoCompression, NoCompression}
	options.SetCompressionPerLevel(compressionLevels)

	t.Log("phase: Destroy")
	stat = DestroyDB(options, &dbname)
	t.Logf("DestroyDB: status = %s", stat)

	t.Log("phase: open_error")
	_, stat, _ = Open(options, &dbname)
	if stat.Ok() {
		t.Error("err: open_error")
	} else {
		t.Logf("open_error: status = %s", stat)
	}

	t.Log("phase: open")
	options.SetCreateIfMissing(true)
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("err: open: stat = %s", stat)
	}
	t.Log("phase: get")
	ropts := NewReadOptions()

	ropts.SetVerifyChecksums(true)
	ropts.SetFillCache(false)

	db.checkGet(t, ropts, []byte("foo"), nil)

	t.Log("phase: put")
	woptions := NewWriteOptions()
	woptions.SetSync(true)
	stat = db.Put(woptions, []byte("foo"), []byte("hello"))
	if !stat.Ok() {
		t.Fatalf("err: put err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))

	t.Log("phase: backup_and_restore")
	stat = DestroyDB(options, &dbbackupname)
	if !stat.Ok() {
		t.Fatalf("backup_and_restore: DestroyDB: status = %s", stat)
	}
	var be *BackupEngine
	be, stat = BackupEngineOpen(options.Env(), NewBackupableDBOptions(&dbbackupname))
	if !stat.Ok() {
		t.Fatalf("backup_and_restore: BackupEngineOpen: status = %s", stat)
	}
	stat = be.CreateNewBackup(db)
	if !stat.Ok() {
		t.Fatalf("backup_and_restore: CreateNewBackup: status = %s", stat)
	}
	stat = db.Delete(woptions, []byte("foo"))
	if !stat.Ok() {
		t.Fatalf("backup_and_restore: Delete: status = %s", stat)
	}
	db.Close()
	stat = DestroyDB(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("backup_and_restore: dbname - DestroyDB: status = %s", stat)
	}
	restore_options := NewRestoreOptions()
	restore_options.SetKeepLogFile(false)
	stat = be.RestoreDBFromLatestBackup(&dbname, &dbname, restore_options)
	if !stat.Ok() {
		t.Fatalf("backup_and_restore: RestoreDBFromLatestBackup: status = %s", stat)
	}
	restore_options.Close()
	options.SetErrorIfExists(false);
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("err: backup_and_restore: open: stat = %s", stat)
	}
	options.SetErrorIfExists(true);
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))
	be.Close()

	t.Log("phase: compactall")
	cropt := NewCompactRangeOptions()
	stat = db.CompactRange(cropt, nil, nil)
	if !stat.Ok() {
		t.Fatalf("err: compactall: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))

	t.Log("phase: compactrange")
	stat = db.CompactRange(cropt, []byte("a"), []byte("z"))
	if !stat.Ok() {
		t.Fatalf("err: compactrange: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))

	t.Log("phase: writebatch")
	wb := NewWriteBatch()
	wb.Put([]byte("foo"), []byte("a"))
	wb.Clear()
	wb.Put([]byte("bar"), []byte("b"))
	wb.Put([]byte("box"), []byte("c"))
	wb.Delete([]byte("bar"))
	stat = db.Write(woptions, wb)
	if !stat.Ok() {
		t.Fatalf("err: writebatch Write: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))
	db.checkGet(t, ropts, []byte("bar"), nil)
	db.checkGet(t, ropts, []byte("box"), []byte("c"))
	// StartPhase("writebatch");
	// {
	// 	int pos = 0;
	// 	rocksdb_writebatch_iterate(wb, &pos, CheckPut, CheckDel);
	// 	CheckCondition(pos == 3);
	// }
	wb.Close()

	t.Log("phase: writebatch_rep")
	wb1 := NewWriteBatch()
	wb1.Put([]byte("baz"), []byte("d"))
	wb1.Put([]byte("quux"), []byte("e"))
	wb1.Delete([]byte("quux"))
	wb2 := NewWriteBatchFromBytes(wb1.Data())
	checkCondition(t, wb1.Count() == wb2.Count())
	checkCondition(t, wb1.GetDataSize() == wb2.GetDataSize())
	checkCondition(t, bytes.Equal(wb1.Data(), wb2.Data()))
	wb1.Close()
	wb2.Close()

	t.Log("phase: iter")
	iter := db.NewIterator(ropts)
	checkCondition(t, !iter.Valid())
	iter.SeekToFirst()
	checkCondition(t, iter.Valid())
	checkIter(t, iter, "box", "c")
	iter.Next()
	checkIter(t, iter, "foo", "hello")
	iter.Prev()
	checkIter(t, iter, "box", "c")
	iter.Prev()
	checkCondition(t, !iter.Valid())
	iter.SeekToLast()
	checkIter(t, iter, "foo", "hello")
	iter.Seek([]byte("b"))
	checkIter(t, iter, "box", "c")
	iter.Close()

	t.Log("phase: approximate_sizes")
	rngs := []*Range{NewRange([]byte("a"), []byte("k00000000000000010000")), NewRange([]byte("k00000000000000010000"), []byte("z"))}
	n := 20000
	woptions.SetSync(false)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("k%020d", i)
		val := fmt.Sprintf("v%020d", i)
		stat = db.Put(woptions, []byte(key), []byte(val))
		if !stat.Ok() {
			t.Fatalf("err: approximate_sizes put: stat = %s", stat)
		}
	}
	szs := db.GetApproximateSizes(rngs)
	checkCondition(t, szs[0] > 0)
	checkCondition(t, szs[1] > 0)

	t.Log("phase: property")
	val, res := db.GetProperty([]byte("nosuchprop"))
	checkCondition(t, !res)
	checkCondition(t, len(val) == 0)
	val, res = db.GetProperty([]byte("rocksdb.stats"))
	checkCondition(t, res)
	checkCondition(t, len(val) > 0)

	t.Log("phase: snapshot")
	snap := db.GetSnapshot()
	stat = db.Delete(woptions, []byte("foo"))
	if !stat.Ok() {
		t.Fatalf("err: snapshot Delete: stat = %s", stat)
	}
	ropts.SetSnapshot(snap)
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))
	ropts.SetSnapshot(nil)
	db.checkGet(t, ropts, []byte("foo"), nil)
	db.ReleaseSnapshot(snap)

	t.Log("phase: repair")
	// If we do not compact here, then the lazy deletion of
	// files (https://reviews.facebook.net/D6123) would leave
	// around deleted files and the repair process will find
	// those files and put them back into the database.
	stat = db.CompactRange(cropt, nil, nil)
	db.Close()
	options.SetErrorIfExists(false);
	options.SetCreateIfMissing(false)
	stat = RepairDB(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("err: repair RepairDB: stat = %s", stat)
	}
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Logf("err: repair Open: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), nil)
	db.checkGet(t, ropts, []byte("bar"), nil)
	db.checkGet(t, ropts, []byte("box"), []byte("c"))
	options.SetErrorIfExists(true);
	options.SetCreateIfMissing(true)

	t.Log("phase: filter")
	var policy *FilterPolicy
	tfp := &testFilterPolicy{t: t, fakeResult: true}
	for run := 0; run < 2; run++ {
		if 0 == run {
			policy = NewFilterPolicy(tfp)
		} else {
			policy = NewBloomFilterPolicy(10)
		}
		table_options.SetFilterPolicy(policy)
		db.Close()
		stat = DestroyDB(options, &dbname)
		t.Logf("filter: DestroyDB: status = %s", stat)
		options.SetTableFactory(table_options.NewBlockBasedTableFactory())
		db, stat, _ = Open(options, &dbname)
		if !stat.Ok() {
			t.Fatalf("err: open: stat = %s", stat)
		}
		stat = db.Put(woptions, []byte("foo"), []byte("foovalue"))
		if !stat.Ok() {
			t.Fatalf("err: put err: stat = %s", stat)
		}
		stat = db.Put(woptions, []byte("bar"), []byte("barvalue"))
		if !stat.Ok() {
			t.Fatalf("err: put err: stat = %s", stat)
		}
		db.CompactRange(cropt, nil, nil)

		tfp.fakeResult = true
		db.checkGet(t, ropts, []byte("foo"), []byte("foovalue"))
		db.checkGet(t, ropts, []byte("bar"), []byte("barvalue"))

		if 0 == run {
			// Must not find value when custom filter returns false
			tfp.fakeResult = false
			db.checkGet(t, ropts, []byte("foo"), nil)
			db.checkGet(t, ropts, []byte("bar"), nil)

			tfp.fakeResult = true
			db.checkGet(t, ropts, []byte("foo"), []byte("foovalue"))
			db.checkGet(t, ropts, []byte("bar"), []byte("barvalue"))
		}

		table_options.SetFilterPolicy(nil)
		options.SetTableFactory(table_options.NewBlockBasedTableFactory())
	}

	t.Log("phase: compaction_filter")
	options_with_filter := NewOptions()
	options_with_filter.SetCreateIfMissing(true)
	cpf := &testCompactionFilter{t: t}
	db.Close()
	stat = DestroyDB(options_with_filter, &dbname)
	t.Logf("compaction_filter: DestroyDB: status = %s", stat)
	cfilter := NewCompactionFilter(cpf)
	options_with_filter.SetCompactionFilter(cfilter)
	db = checkCompaction(t, &dbname, options_with_filter, ropts, woptions)
	options_with_filter.SetCompactionFilter(nil)
	cfilter.Close()
	options_with_filter.Close()

	t.Log("phase: compaction_filter_factory")
	options_with_filter_factory := NewOptions()
	options_with_filter_factory.SetCreateIfMissing(true)
	cff := &testCompactionFilterFactory{t: t}
	factory := NewCompactionFilterFactory(cff)
	db.Close()
	stat = DestroyDB(options_with_filter_factory, &dbname)
	t.Logf("compaction_filter_factory: DestroyDB: status = %s", stat)
	options_with_filter_factory.SetCompactionFilterFactory(factory)
	db = checkCompaction(t, &dbname, options_with_filter_factory, ropts, woptions)
	options_with_filter_factory.SetCompactionFilterFactory(nil)
	options_with_filter_factory.Close()

	t.Log("phase: merge_operator")
	cmop := &testMergeOperator{t: t}
	merge_operator := NewMergeOperator(cmop)
	db.Close()
	stat = DestroyDB(options, &dbname)
	t.Logf("merge_operator: DestroyDB: status = %s", stat)
	options.SetMergeOperator(merge_operator)
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2: err: open: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo"), []byte("foovalue"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("foovalue"))
	stat = db.Merge(woptions, []byte("foo"), []byte("barvalue"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("fake"))
	stat = db.Merge(woptions, []byte("bar"), []byte("barvalue"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("bar"), []byte("fake"))

	t.Log("phase: columnfamilies")
	db.Close()
	stat = DestroyDB(options, &dbname)
	t.Logf("columnfamilies: DestroyDB: status = %s", stat)
	db_options := NewOptions()
	db_options.SetCreateIfMissing(true)
	db, stat, _ = Open(db_options, &dbname)
	if !stat.Ok() {
		t.Fatalf("columnfamilies: err: open: stat = %s", stat)
	}
	var cfh *ColumnFamilyHandle
	cf1 := "cf1"
	default_s := "default"
	cfh, stat = db.CreateColumnFamily(&db_options.ColumnFamilyOptions, &cf1)
	if !stat.Ok() {
		t.Fatalf("columnfamilies: err: CreateColumnFamily: stat = %s", stat)
	}
	cfh.Close()
	db.Close()
	var cfss []string
	cfss, stat = ListColumnFamilies(&db_options.DBOptions, &dbname) 
	if !stat.Ok() {
		t.Fatalf("columnfamilies: err: ListColumnFamilies: stat = %s", stat)
	}
	checkCondition(t, default_s == cfss[0])
	checkCondition(t, cf1 == cfss[1])
	checkCondition(t, 2 == len(cfss))

	cf_options := NewColumnFamilyOptions()
	cfds := []*ColumnFamilyDescriptor{NewColumnFamilyDescriptor(&default_s, cf_options), NewColumnFamilyDescriptor(&cf1, cf_options)}
	var cfhs []*ColumnFamilyHandle
	db, stat, cfhs = Open(db_options, &dbname, cfds...)
	if !stat.Ok() {
		t.Fatalf("columnfamilies: err: open with cfds: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo"), []byte("hello"), cfhs[1])
	if !stat.Ok() {
		t.Fatalf("columnfamilies:err: put err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"), cfhs[1])
	stat = db.Delete(woptions, []byte("foo"), cfhs[1])
	if !stat.Ok() {
		t.Fatalf("columnfamilies:err: Delete err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), nil, cfhs[1])
	wb = NewWriteBatch()
	wb.Put([]byte("baz"), []byte("a"), cfhs[1])
	wb.Clear()
	wb.Put([]byte("bar"), []byte("b"), cfhs[1])
	wb.Put([]byte("box"), []byte("c"), cfhs[1])
	wb.Delete([]byte("bar"), cfhs[1])
	stat = db.Write(woptions, wb)
	if !stat.Ok() {
		t.Fatalf("columnfamilies: err: writebatch: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("baz"), nil, cfhs[1])
	db.checkGet(t, ropts, []byte("bar"), nil, cfhs[1])
	db.checkGet(t, ropts, []byte("box"), []byte("c"), cfhs[1])
	wb.Close()

	iter = db.NewIterator(ropts, cfhs[1])
	checkCondition(t, !iter.Valid())
	iter.SeekToFirst()
	checkCondition(t, iter.Valid())
	i := 0
	for ; iter.Valid(); iter.Next() {
		i++
	}
	checkCondition(t, 1 == i)
	stat = iter.Status()
	if !stat.Ok() {
		t.Fatalf("columnfamilies: err: iter: stat = %s", stat)
	}
	iter.Close()
	stat = db.DropColumnFamily(cfhs[1])
	if !stat.Ok() {
		t.Fatalf("columnfamilies: err: DropColumnFamily: stat = %s", stat)
	}
	for _, cfd := range cfds {
		cfd.Close()
	}
	db.Close()
	stat = DestroyDB(options, &dbname)
	t.Logf("columnfamilies: last: DestroyDB: status = %s", stat)
	db_options.Close()
	cf_options.Close()

	t.Log("phase: prefix")
	// Create new database
	options.SetAllowMmapReads(true)
	options.SetPrefixExtractor(NewFixedPrefixTransform(3))
	options.SetMemtableFactory(NewHashSkipListRepFactory(uint64(5000), int32(4), int32(4)))
	pto := NewPlainTableOptions()
	pto.SetUserKeyLen(4)
	pto.SetBloomBitsPerKey(10)
	pto.SetHashTableRatio(0.75)
	pto.SetIndexSparseness(16)
	options.SetTableFactory(pto.NewPlainTableFactory())
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("prefix: err: open: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo1"), []byte("foo"))
	if !stat.Ok() {
		t.Fatalf("prefix:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo2"), []byte("foo"))
	if !stat.Ok() {
		t.Fatalf("prefix:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo3"), []byte("foo"))
	if !stat.Ok() {
		t.Fatalf("prefix:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("bar1"), []byte("bar"))
	if !stat.Ok() {
		t.Fatalf("prefix:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("bar2"), []byte("bar"))
	if !stat.Ok() {
		t.Fatalf("prefix:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("bar3"), []byte("bar"))
	if !stat.Ok() {
		t.Fatalf("prefix:err: put err: stat = %s", stat)
	}

	iter = db.NewIterator(ropts)
	checkCondition(t, !iter.Valid())
	iter.Seek([]byte("bar"))
	stat = iter.Status()
	if !stat.Ok() {
		t.Fatalf("prefix: err: iter: stat = %s", stat)
	}
	checkCondition(t, iter.Valid())

	checkIter(t, iter, "bar1", "bar")
	iter.Next()
	checkIter(t, iter, "bar2", "bar")
	iter.Next()
	checkIter(t, iter, "bar3", "bar")
	stat = iter.Status()
	if !stat.Ok() {
		t.Fatalf("prefix: err: checkIter: stat = %s", stat)
	}
	iter.Close()
	db.Close()
	stat = DestroyDB(options, &dbname)
	t.Logf("prefix: DestroyDB: status = %s", stat)

	t.Log("phase: cuckoo_options")
	cuckoo_options := NewCuckooTableOptions()
	cuckoo_options.SetHashTableRatio(0.5)
	cuckoo_options.SetMaxSearchDepth(200)
	cuckoo_options.SetCuckooBlockSize(10)
	cuckoo_options.SetIdentityAsFirstHash(true)
	cuckoo_options.SetUseModuleHash(false)
	options.SetTableFactory(cuckoo_options.NewCuckooTableFactory())
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("prefix: err: open: stat = %s", stat)
	}
	cuckoo_options.Close()

	t.Log("phase: iterate_upper_bound")
	// Create new empty database
	db.Close()
	stat = DestroyDB(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("iterate_upper_bound: DestroyDB: status = %s", stat)
	}
	options.SetPrefixExtractor(nil)
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("iterate_upper_bound: err: open: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("a"), []byte("0"))
	if !stat.Ok() {
		t.Fatalf("iterate_upper_bound:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo"), []byte("bar"))
	if !stat.Ok() {
		t.Fatalf("iterate_upper_bound:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo1"), []byte("bar1"))
	if !stat.Ok() {
		t.Fatalf("iterate_upper_bound:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("g1"), []byte("0"))
	if !stat.Ok() {
		t.Fatalf("iterate_upper_bound:err: put err: stat = %s", stat)
	}

	// testing basic case with no iterate_upper_bound and no prefix_extractor
	ropts.SetIterateUpperBound(nil)
	iter = db.NewIterator(ropts)
	iter.Seek([]byte("foo"))
	checkCondition(t, iter.Valid())
	checkIter(t, iter, "foo", "bar")
	iter.Next()
	checkIter(t, iter, "foo1", "bar1")
	iter.Next()
	checkIter(t, iter, "g1", "0")
	iter.Close()

	// testing iterate_upper_bound and forward iterator
	// to make sure it stops at bound
	// iterate_upper_bound points beyond the last expected entry
	ropts.SetIterateUpperBound([]byte("foo2"))
	iter = db.NewIterator(ropts)
	iter.Seek([]byte("foo"))
	checkCondition(t, iter.Valid())
	checkIter(t, iter, "foo", "bar")
	iter.Next()
	checkIter(t, iter, "foo1", "bar1")
	iter.Next()
	// should stop here...
	checkCondition(t, !iter.Valid())
	iter.Close()

	t.Log("phase: cleanup")
	db.Close()
	options.Close()
	table_options.Close()
	ropts.Close()
	woptions.Close()
	cache.Close()
	// Keep cmp from being garbage collected eariler
	cmp.Close()
}
