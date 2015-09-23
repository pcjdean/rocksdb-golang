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
	db.CompactRange(nil, nil)

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

// Custom CompactionFilterV2 filter
type testCompactionFilterV2 struct {
	t *testing.T
}

func (tcf *testCompactionFilterV2) Name() string {
	return "TestCompactionFilterV2"
}

func (tcf *testCompactionFilterV2) Filter(level int, keys, exvals [][]byte) (newvals [][]byte, valchgs []bool, removed []bool) {
	l := len(keys)
	if 0 < l {
		removed = make([]bool, l)
		valchgs = make([]bool, l)
	}

	for i, _ := range keys {
		lv := len(exvals[i])
		// If any value is "gc", it's removed.
		if 2 == lv && bytes.Equal([]byte("gc"), exvals[i]) {
			removed[i] = true;
		} else if 6 == lv && bytes.Equal([]byte("gc all"), exvals[i]) {
			// If any value is "gc all", all keys are removed.
			for j, _ := range keys {
				removed[j] = true;
			}
			// tcf.t.Logf("testCompactionFilter - gc all - level = %v, keys = %s, exvals = %s, newvals = %s, valchgs = %v, removed = %v", level, keys, exvals, newvals, valchgs, removed)
			return;
		} else if 6 == lv && bytes.Equal([]byte("change"), exvals[i]) {
			// If value is "change", set changed value to "changed".
			newvals = append(newvals, []byte("changed"))
			valchgs[i] = true
		} else {
			// Otherwise, no keys are removed.
		}
	}

	// tcf.t.Logf("testCompactionFilter level = %v, keys = %s, exvals = %s, newvals = %s, valchgs = %v, removed = %v", level, keys, exvals, newvals, valchgs, removed)
	return
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

// Custom CompactionFilterFactoryV2 filter
type testCompactionFilterFactoryV2 struct {
	t *testing.T
	istf ISliceTransform
}

func (tcff *testCompactionFilterFactoryV2) Name() string {
	return "TestCompactionFilterV2"
}

func (tcff *testCompactionFilterFactoryV2) CreateCompactionFilterV2(context *CompactionFilterContext) ICompactionFilterV2 {
	// tcff.t.Logf("testCompactionFilterFactoryV2 context = %v", context)
	cff := &testCompactionFilterV2{t: tcff.t}
	return cff
}

func (tcff *testCompactionFilterFactoryV2) GetPrefixExtractor() ISliceTransform {
	return tcff.istf
}

func (tcff *testCompactionFilterFactoryV2) SetPrefixExtractor(prextrc ISliceTransform) {
	tcff.istf = prextrc
}

// Test from rocksdb's c_test.c.
func TestCMain(t *testing.T) {
	var (
		stat *Status
		db *DB
	)

	dbname := fmt.Sprintf("%s/rocksdb_go_test-%d", os.TempDir(), os.Geteuid);
	dbbackupname := fmt.Sprintf("%s/rocksdb_go_test-%d-backup", os.TempDir(), os.Geteuid);
	t.Log("phase: create_objects")
	fmt.Printf("dbname = %s\n", dbname)
	fmt.Printf("dbbackupname = %s\n", dbbackupname)
	// cmp = rocksdb_comparator_create(NULL, CmpDestroy, CmpCompare, CmpName);
	// env = rocksdb_create_default_env();
	cache := NewLRUCache(100000)

	options := NewOptions()

	// rocksdb_options_set_comparator(options, cmp);
	options.SetErrorIfExists(true);
	// rocksdb_options_set_env(options, env);
	// rocksdb_options_set_info_log(options, NULL);
	options.SetWriteBufferSize(100000)
	// rocksdb_options_set_paranoid_checks(options, 1);
	// rocksdb_options_set_max_open_files(options, 10);
	table_options := NewBlockBasedTableOptions()
	table_options.SetBlockCache(cache)
	options.SetTableFactory(table_options.NewBlockBasedTableFactory())

	options.SetCompression(NoCompression);
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

	// rocksdb_readoptions_set_verify_checksums(roptions, 1);
	// rocksdb_readoptions_set_fill_cache(roptions, 0);

	db.checkGet(t, ropts, []byte("foo"), nil)

	t.Log("phase: put")
	woptions := NewWriteOptions()
	woptions.SetSync(true)
	stat = db.Put(woptions, []byte("foo"), []byte("hello"))
	if !stat.Ok() {
		t.Fatalf("err: put err: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))

	// StartPhase("backup_and_restore");
	// {
	// 	rocksdb_destroy_db(options, dbbackupname, &err);
	// 	CheckNoError(err);

	// 	rocksdb_backup_engine_t *be = rocksdb_backup_engine_open(options, dbbackupname, &err);
	// 	CheckNoError(err);

	// 	rocksdb_backup_engine_create_new_backup(be, db, &err);
	// 	CheckNoError(err);

	// 	rocksdb_delete(db, woptions, "foo", 3, &err);
	// 	CheckNoError(err);

	// 	rocksdb_close(db);

	// 	rocksdb_destroy_db(options, dbname, &err);
	// 	CheckNoError(err);

	// 	rocksdb_restore_options_t *restore_options = rocksdb_restore_options_create();
	// 	rocksdb_restore_options_set_keep_log_files(restore_options, 0);
	// 	rocksdb_backup_engine_restore_db_from_latest_backup(be, dbname, dbname, restore_options, &err);
	// 	CheckNoError(err);
	// 	rocksdb_restore_options_destroy(restore_options);

	// 	rocksdb_options_set_error_if_exists(options, 0);
	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);
	// 	rocksdb_options_set_error_if_exists(options, 1);

	// 	CheckGet(db, roptions, "foo", "hello");

	// 	rocksdb_backup_engine_close(be);
	// }

	t.Log("phase: compactall")
	stat = db.CompactRange(nil, nil)
	if !stat.Ok() {
		t.Fatalf("err: compactall: stat = %s", stat)
	}
	db.checkGet(t, ropts, []byte("foo"), []byte("hello"))

	t.Log("phase: compactrange")
	stat = db.CompactRange([]byte("a"), []byte("z"))
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
	stat = db.CompactRange(nil, nil)
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
		db.CompactRange(nil, nil)

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

	t.Log("phase: compaction_filter_v2")
	tstf := &testSliceTransform{t: t}
	tcfv2 := &testCompactionFilterFactoryV2{t: t}
	factoryv2 := NewCompactionFilterFactoryV2(tcfv2, tstf)
	db.Close()
	stat = DestroyDB(options, &dbname)
	t.Logf("compaction_filter_v2: DestroyDB: status = %s", stat)
	options.SetCompactionFilterFactoryV2(factoryv2)
	db, stat, _ = Open(options, &dbname)
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2: err: open: stat = %s", stat)
	}
	// Only foo2 is GC'd, foo3 is changed.
	stat = db.Put(woptions, []byte("foo1"), []byte("no gc"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo2"), []byte("gc"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("foo3"), []byte("change"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	// All bars are GC'd.
	stat = db.Put(woptions, []byte("bar1"), []byte("no gc"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("bar2"), []byte("gc all"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	stat = db.Put(woptions, []byte("bar3"), []byte("no gc"))
	if !stat.Ok() {
		t.Fatalf("compaction_filter_v2:err: put err: stat = %s", stat)
	}
	// Compact the DB to garbage collect.
	db.CompactRange(nil, nil)

	// Verify foo entries.
	db.checkGet(t, ropts, []byte("foo1"), []byte("no gc"))
	db.checkGet(t, ropts, []byte("foo2"), nil)
	db.checkGet(t, ropts, []byte("foo3"), []byte("changed"))
	// Verify bar entries were all deleted.
	db.checkGet(t, ropts, []byte("bar1"), nil)
	db.checkGet(t, ropts, []byte("bar2"), nil)
	db.checkGet(t, ropts, []byte("bar3"), nil)

	// StartPhase("merge_operator");
	// {
	// 	rocksdb_mergeoperator_t* merge_operator;
	// 	merge_operator = rocksdb_mergeoperator_create(
	// 		NULL, MergeOperatorDestroy, MergeOperatorFullMerge,
	// 		MergeOperatorPartialMerge, NULL, MergeOperatorName);
	// 	// Create new database
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options, dbname, &err);
	// 	rocksdb_options_set_merge_operator(options, merge_operator);
	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo", 3, "foovalue", 8, &err);
	// 	CheckNoError(err);
	// 	CheckGet(db, roptions, "foo", "foovalue");
	// 	rocksdb_merge(db, woptions, "foo", 3, "barvalue", 8, &err);
	// 	CheckNoError(err);
	// 	CheckGet(db, roptions, "foo", "fake");

	// 	// Merge of a non-existing value
	// 	rocksdb_merge(db, woptions, "bar", 3, "barvalue", 8, &err);
	// 	CheckNoError(err);
	// 	CheckGet(db, roptions, "bar", "fake");

	// }

	// StartPhase("columnfamilies");
	// {
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options, dbname, &err);
	// 	CheckNoError(err)

	// 	rocksdb_options_t* db_options = rocksdb_options_create();
	// 	rocksdb_options_set_create_if_missing(db_options, 1);
	// 	db = rocksdb_open(db_options, dbname, &err);
	// 	CheckNoError(err)
	// 	rocksdb_column_family_handle_t* cfh;
	// 	cfh = rocksdb_create_column_family(db, db_options, "cf1", &err);
	// 	rocksdb_column_family_handle_destroy(cfh);
	// 	CheckNoError(err);
	// 	rocksdb_close(db);

	// 	size_t cflen;
	// 	char** column_fams = rocksdb_list_column_families(db_options, dbname, &cflen, &err);
	// 	CheckNoError(err);
	// 	CheckEqual("default", column_fams[0], 7);
	// 	CheckEqual("cf1", column_fams[1], 3);
	// 	CheckCondition(cflen == 2);
	// 	rocksdb_list_column_families_destroy(column_fams, cflen);

	// 	rocksdb_options_t* cf_options = rocksdb_options_create();

	// 	const char* cf_names[2] = {"default", "cf1"};
	// 	const rocksdb_options_t* cf_opts[2] = {cf_options, cf_options};
	// 	rocksdb_column_family_handle_t* handles[2];
	// 	db = rocksdb_open_column_families(db_options, dbname, 2, cf_names, cf_opts, handles, &err);
	// 	CheckNoError(err);

	// 	rocksdb_put_cf(db, woptions, handles[1], "foo", 3, "hello", 5, &err);
	// 	CheckNoError(err);

	// 	CheckGetCF(db, roptions, handles[1], "foo", "hello");

	// 	rocksdb_delete_cf(db, woptions, handles[1], "foo", 3, &err);
	// 	CheckNoError(err);

	// 	CheckGetCF(db, roptions, handles[1], "foo", NULL);

	// 	rocksdb_writebatch_t* wb = rocksdb_writebatch_create();
	// 	rocksdb_writebatch_put_cf(wb, handles[1], "baz", 3, "a", 1);
	// 	rocksdb_writebatch_clear(wb);
	// 	rocksdb_writebatch_put_cf(wb, handles[1], "bar", 3, "b", 1);
	// 	rocksdb_writebatch_put_cf(wb, handles[1], "box", 3, "c", 1);
	// 	rocksdb_writebatch_delete_cf(wb, handles[1], "bar", 3);
	// 	rocksdb_write(db, woptions, wb, &err);
	// 	CheckNoError(err);
	// 	CheckGetCF(db, roptions, handles[1], "baz", NULL);
	// 	CheckGetCF(db, roptions, handles[1], "bar", NULL);
	// 	CheckGetCF(db, roptions, handles[1], "box", "c");
	// 	rocksdb_writebatch_destroy(wb);

	// 	rocksdb_iterator_t* iter = rocksdb_create_iterator_cf(db, roptions, handles[1]);
	// 	CheckCondition(!rocksdb_iter_valid(iter));
	// 	rocksdb_iter_seek_to_first(iter);
	// 	CheckCondition(rocksdb_iter_valid(iter));

	// 	int i;
	// 	for (i = 0; rocksdb_iter_valid(iter) != 0; rocksdb_iter_next(iter)) {
	// 		i++;
	// 	}
	// 	CheckCondition(i == 1);
	// 	rocksdb_iter_get_error(iter, &err);
	// 	CheckNoError(err);
	// 	rocksdb_iter_destroy(iter);

	// 	rocksdb_drop_column_family(db, handles[1], &err);
	// 	CheckNoError(err);
	// 	for (i = 0; i < 2; i++) {
	// 		rocksdb_column_family_handle_destroy(handles[i]);
	// 	}
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options, dbname, &err);
	// 	rocksdb_options_destroy(db_options);
	// 	rocksdb_options_destroy(cf_options);
	// }

	// StartPhase("prefix");
	// {
	// 	// Create new database
	// 	rocksdb_options_set_allow_mmap_reads(options, 1);
	// 	rocksdb_options_set_prefix_extractor(options, rocksdb_slicetransform_create_fixed_prefix(3));
	// 	rocksdb_options_set_hash_skip_list_rep(options, 5000, 4, 4);
	// 	rocksdb_options_set_plain_table_factory(options, 4, 10, 0.75, 16);

	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);

	// 	rocksdb_put(db, woptions, "foo1", 4, "foo", 3, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo2", 4, "foo", 3, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo3", 4, "foo", 3, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "bar1", 4, "bar", 3, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "bar2", 4, "bar", 3, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "bar3", 4, "bar", 3, &err);
	// 	CheckNoError(err);

	// 	rocksdb_iterator_t* iter = rocksdb_create_iterator(db, roptions);
	// 	CheckCondition(!rocksdb_iter_valid(iter));

	// 	rocksdb_iter_seek(iter, "bar", 3);
	// 	rocksdb_iter_get_error(iter, &err);
	// 	CheckNoError(err);
	// 	CheckCondition(rocksdb_iter_valid(iter));

	// 	CheckIter(iter, "bar1", "bar");
	// 	rocksdb_iter_next(iter);
	// 	CheckIter(iter, "bar2", "bar");
	// 	rocksdb_iter_next(iter);
	// 	CheckIter(iter, "bar3", "bar");
	// 	rocksdb_iter_get_error(iter, &err);
	// 	CheckNoError(err);
	// 	rocksdb_iter_destroy(iter);

	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options, dbname, &err);
	// }

	// StartPhase("cuckoo_options");
	// {
	// 	rocksdb_cuckoo_table_options_t* cuckoo_options;
	// 	cuckoo_options = rocksdb_cuckoo_options_create();
	// 	rocksdb_cuckoo_options_set_hash_ratio(cuckoo_options, 0.5);
	// 	rocksdb_cuckoo_options_set_max_search_depth(cuckoo_options, 200);
	// 	rocksdb_cuckoo_options_set_cuckoo_block_size(cuckoo_options, 10);
	// 	rocksdb_cuckoo_options_set_identity_as_first_hash(cuckoo_options, 1);
	// 	rocksdb_cuckoo_options_set_use_module_hash(cuckoo_options, 0);
	// 	rocksdb_options_set_cuckoo_table_factory(options, cuckoo_options);

	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);

	// 	rocksdb_cuckoo_options_destroy(cuckoo_options);
	// }

	// StartPhase("iterate_upper_bound");
	// {
	// 	// Create new empty database
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options, dbname, &err);
	// 	CheckNoError(err);

	// 	rocksdb_options_set_prefix_extractor(options, NULL);
	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);

	// 	rocksdb_put(db, woptions, "a",    1, "0",    1, &err); CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo",  3, "bar",  3, &err); CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo1", 4, "bar1", 4, &err); CheckNoError(err);
	// 	rocksdb_put(db, woptions, "g1",   2, "0",    1, &err); CheckNoError(err);

	// 	// testing basic case with no iterate_upper_bound and no prefix_extractor
	// 	{
	// 		rocksdb_readoptions_set_iterate_upper_bound(roptions, NULL, 0);
	// 		rocksdb_iterator_t* iter = rocksdb_create_iterator(db, roptions);

	// 		rocksdb_iter_seek(iter, "foo", 3);
	// 		CheckCondition(rocksdb_iter_valid(iter));
	// 		CheckIter(iter, "foo", "bar");

	// 		rocksdb_iter_next(iter);
	// 		CheckCondition(rocksdb_iter_valid(iter));
	// 		CheckIter(iter, "foo1", "bar1");

	// 		rocksdb_iter_next(iter);
	// 		CheckCondition(rocksdb_iter_valid(iter));
	// 		CheckIter(iter, "g1", "0");

	// 		rocksdb_iter_destroy(iter);
	// 	}

	// 	// testing iterate_upper_bound and forward iterator
	// 	// to make sure it stops at bound
	// 	{
	// 		// iterate_upper_bound points beyond the last expected entry
	// 		rocksdb_readoptions_set_iterate_upper_bound(roptions, "foo2", 4);

	// 		rocksdb_iterator_t* iter = rocksdb_create_iterator(db, roptions);

	// 		rocksdb_iter_seek(iter, "foo", 3);
	// 		CheckCondition(rocksdb_iter_valid(iter));
	// 		CheckIter(iter, "foo", "bar");

	// 		rocksdb_iter_next(iter);
	// 		CheckCondition(rocksdb_iter_valid(iter));
	// 		CheckIter(iter, "foo1", "bar1");

	// 		rocksdb_iter_next(iter);
	// 		// should stop here...
	// 		CheckCondition(!rocksdb_iter_valid(iter));

	// 		rocksdb_iter_destroy(iter);
	// 	}
	// }

	// StartPhase("cleanup");
	// rocksdb_close(db);
	// rocksdb_options_destroy(options);
	// rocksdb_block_based_options_destroy(table_options);
	// rocksdb_readoptions_destroy(roptions);
	// rocksdb_writeoptions_destroy(woptions);
	// rocksdb_cache_destroy(cache);
	// rocksdb_comparator_destroy(cmp);
	// rocksdb_env_destroy(env);

	// fprintf(stderr, "PASS\n");
}
