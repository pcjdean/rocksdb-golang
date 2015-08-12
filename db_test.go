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
		t.Fatalf("err: get err: stat = %s", stat)
	} else if !bytes.Equal(val, valExp) {
		t.Fatalf("err: get err: expected = %v, got = %v", valExp, val)
	}
}

func checkCondition(t *testing.T, cond bool) {
	if !cond {
		t.Fatal("err: checkCondition:")
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
	t.Log("phase: create_objects")
	fmt.Printf("dbname = %s\n", dbname)
	fmt.Printf("dbbackupname = %s\n", dbbackupname)
	// cmp = rocksdb_comparator_create(NULL, CmpDestroy, CmpCompare, CmpName);
	// env = rocksdb_create_default_env();
	// cache = rocksdb_cache_create_lru(100000);

	options := NewOptions()

	// rocksdb_options_set_comparator(options, cmp);
	options.SetErrorIfExists(true);
	// rocksdb_options_set_env(options, env);
	// rocksdb_options_set_info_log(options, NULL);
	// rocksdb_options_set_write_buffer_size(options, 100000);
	// rocksdb_options_set_paranoid_checks(options, 1);
	// rocksdb_options_set_max_open_files(options, 10);
	// table_options = rocksdb_block_based_options_create();
	// rocksdb_block_based_options_set_block_cache(table_options, cache);
	// rocksdb_options_set_block_based_table_factory(options, table_options);

	options.SetCompression(NoCompression);
	options.SetCompressionOptions(-14, -1, 0)
	compressionLevels := []int{
		NoCompression, NoCompression, NoCompression, NoCompression}
	options.SetCompressionPerLevel(compressionLevels)

	t.Log("phase: Destroy")
	stat = DestroyDB(&dbname, options)
	t.Logf("DestroyDB: status = ", stat)

	t.Log("phase: open_error")
	_, stat, _ = Open(&dbname, options)
	if stat.Ok() {
		t.Error("err: open_error")
	} else {
		t.Logf("open_error: status = ", stat)
	}

	t.Log("phase: open")
	options.SetCreateIfMissing(true)
	db, stat, _ = Open(&dbname, options)
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

	// StartPhase("writebatch");
	// {
	// 	rocksdb_writebatch_t* wb = rocksdb_writebatch_create();
	// 	rocksdb_writebatch_put(wb, "foo", 3, "a", 1);
	// 	rocksdb_writebatch_clear(wb);
	// 	rocksdb_writebatch_put(wb, "bar", 3, "b", 1);
	// 	rocksdb_writebatch_put(wb, "box", 3, "c", 1);
	// 	rocksdb_writebatch_delete(wb, "bar", 3);
	// 	rocksdb_write(db, woptions, wb, &err);
	// 	CheckNoError(err);
	// 	CheckGet(db, roptions, "foo", "hello");
	// 	CheckGet(db, roptions, "bar", NULL);
	// 	CheckGet(db, roptions, "box", "c");
	// 	int pos = 0;
	// 	rocksdb_writebatch_iterate(wb, &pos, CheckPut, CheckDel);
	// 	CheckCondition(pos == 3);
	// 	rocksdb_writebatch_destroy(wb);
	// }

	// StartPhase("writebatch_rep");
	// {
	// 	rocksdb_writebatch_t* wb1 = rocksdb_writebatch_create();
	// 	rocksdb_writebatch_put(wb1, "baz", 3, "d", 1);
	// 	rocksdb_writebatch_put(wb1, "quux", 4, "e", 1);
	// 	rocksdb_writebatch_delete(wb1, "quux", 4);
	// 	size_t repsize1 = 0;
	// 	const char* rep = rocksdb_writebatch_data(wb1, &repsize1);
	// 	rocksdb_writebatch_t* wb2 = rocksdb_writebatch_create_from(rep, repsize1);
	// 	CheckCondition(rocksdb_writebatch_count(wb1) ==
	// 		rocksdb_writebatch_count(wb2));
	// 	size_t repsize2 = 0;
	// 	CheckCondition(
	// 		memcmp(rep, rocksdb_writebatch_data(wb2, &repsize2), repsize1) == 0);
	// 	rocksdb_writebatch_destroy(wb1);
	// 	rocksdb_writebatch_destroy(wb2);
	// }

	// StartPhase("iter");
	// {
	// 	rocksdb_iterator_t* iter = rocksdb_create_iterator(db, roptions);
	// 	CheckCondition(!rocksdb_iter_valid(iter));
	// 	rocksdb_iter_seek_to_first(iter);
	// 	CheckCondition(rocksdb_iter_valid(iter));
	// 	CheckIter(iter, "box", "c");
	// 	rocksdb_iter_next(iter);
	// 	CheckIter(iter, "foo", "hello");
	// 	rocksdb_iter_prev(iter);
	// 	CheckIter(iter, "box", "c");
	// 	rocksdb_iter_prev(iter);
	// 	CheckCondition(!rocksdb_iter_valid(iter));
	// 	rocksdb_iter_seek_to_last(iter);
	// 	CheckIter(iter, "foo", "hello");
	// 	rocksdb_iter_seek(iter, "b", 1);
	// 	CheckIter(iter, "box", "c");
	// 	rocksdb_iter_get_error(iter, &err);
	// 	CheckNoError(err);
	// 	rocksdb_iter_destroy(iter);
	// }

	t.Log("phase: approximate_sizes")
	rngs := []*Range{NewRange([]byte("a"), []byte("k00000000000000010000")), NewRange([]byte("k00000000000000010000"), []byte("z"))}
	n := 20000
	woptions.SetSync(false)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("k%020d", i)
		val := fmt.Sprintf("v%020d", i)
		db.Put(woptions, []byte(key), []byte(val))
		if !stat.Ok() {
			t.Fatalf("err: approximate_sizes put: stat = %s", stat)
		}
	}
	szs := db.GetApproximateSizes(rngs)
	fmt.Printf("szs = %v\n", szs)
	checkCondition(t, szs[0] > 0)
	checkCondition(t, szs[1] > 0)

	// StartPhase("property");
	// {
	// 	char* prop = rocksdb_property_value(db, "nosuchprop");
	// 	CheckCondition(prop == NULL);
	// 	prop = rocksdb_property_value(db, "rocksdb.stats");
	// 	CheckCondition(prop != NULL);
	// 	Free(&prop);
	// }

	// StartPhase("snapshot");
	// {
	// 	const rocksdb_snapshot_t* snap;
	// 	snap = rocksdb_create_snapshot(db);
	// 	rocksdb_delete(db, woptions, "foo", 3, &err);
	// 	CheckNoError(err);
	// 	rocksdb_readoptions_set_snapshot(roptions, snap);
	// 	CheckGet(db, roptions, "foo", "hello");
	// 	rocksdb_readoptions_set_snapshot(roptions, NULL);
	// 	CheckGet(db, roptions, "foo", NULL);
	// 	rocksdb_release_snapshot(db, snap);
	// }

	// StartPhase("repair");
	// {
	// 	// If we do not compact here, then the lazy deletion of
	// 	// files (https://reviews.facebook.net/D6123) would leave
	// 	// around deleted files and the repair process will find
	// 	// those files and put them back into the database.
	// 	rocksdb_compact_range(db, NULL, 0, NULL, 0);
	// 	rocksdb_close(db);
	// 	rocksdb_options_set_create_if_missing(options, 0);
	// 	rocksdb_options_set_error_if_exists(options, 0);
	// 	rocksdb_repair_db(options, dbname, &err);
	// 	CheckNoError(err);
	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);
	// 	CheckGet(db, roptions, "foo", NULL);
	// 	CheckGet(db, roptions, "bar", NULL);
	// 	CheckGet(db, roptions, "box", "c");
	// 	rocksdb_options_set_create_if_missing(options, 1);
	// 	rocksdb_options_set_error_if_exists(options, 1);
	// }

	// StartPhase("filter");
	// for (run = 0; run < 2; run++) {
	// 	// First run uses custom filter, second run uses bloom filter
	// 	CheckNoError(err);
	// 	rocksdb_filterpolicy_t* policy;
	// 	if (run == 0) {
	// 		policy = rocksdb_filterpolicy_create(
	// 			NULL, FilterDestroy, FilterCreate, FilterKeyMatch, NULL, FilterName);
	// 	} else {
	// 		policy = rocksdb_filterpolicy_create_bloom(10);
	// 	}

	// 	rocksdb_block_based_options_set_filter_policy(table_options, policy);

	// 	// Create new database
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options, dbname, &err);
	// 	rocksdb_options_set_block_based_table_factory(options, table_options);
	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo", 3, "foovalue", 8, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "bar", 3, "barvalue", 8, &err);
	// 	CheckNoError(err);
	// 	rocksdb_compact_range(db, NULL, 0, NULL, 0);

	// 	fake_filter_result = 1;
	// 	CheckGet(db, roptions, "foo", "foovalue");
	// 	CheckGet(db, roptions, "bar", "barvalue");
	// 	if (phase == 0) {
	// 		// Must not find value when custom filter returns false
	// 		fake_filter_result = 0;
	// 		CheckGet(db, roptions, "foo", NULL);
	// 		CheckGet(db, roptions, "bar", NULL);
	// 		fake_filter_result = 1;

	// 		CheckGet(db, roptions, "foo", "foovalue");
	// 		CheckGet(db, roptions, "bar", "barvalue");
	// 	}
	// 	// Reset the policy
	// 	rocksdb_block_based_options_set_filter_policy(table_options, NULL);
	// 	rocksdb_options_set_block_based_table_factory(options, table_options);
	// }

	// StartPhase("compaction_filter");
	// {
	// 	rocksdb_options_t* options_with_filter = rocksdb_options_create();
	// 	rocksdb_options_set_create_if_missing(options_with_filter, 1);
	// 	rocksdb_compactionfilter_t* cfilter;
	// 	cfilter = rocksdb_compactionfilter_create(NULL, CFilterDestroy,
	// 		CFilterFilter, CFilterName);
	// 	// Create new database
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options_with_filter, dbname, &err);
	// 	rocksdb_options_set_compaction_filter(options_with_filter, cfilter);
	// 	db = CheckCompaction(db, options_with_filter, roptions, woptions);

	// 	rocksdb_options_set_compaction_filter(options_with_filter, NULL);
	// 	rocksdb_compactionfilter_destroy(cfilter);
	// 	rocksdb_options_destroy(options_with_filter);
	// }

	// StartPhase("compaction_filter_factory");
	// {
	// 	rocksdb_options_t* options_with_filter_factory = rocksdb_options_create();
	// 	rocksdb_options_set_create_if_missing(options_with_filter_factory, 1);
	// 	rocksdb_compactionfilterfactory_t* factory;
	// 	factory = rocksdb_compactionfilterfactory_create(
	// 		NULL, CFilterFactoryDestroy, CFilterCreate, CFilterFactoryName);
	// 	// Create new database
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options_with_filter_factory, dbname, &err);
	// 	rocksdb_options_set_compaction_filter_factory(options_with_filter_factory,
	// 		factory);
	// 	db = CheckCompaction(db, options_with_filter_factory, roptions, woptions);

	// 	rocksdb_options_set_compaction_filter_factory(
	// 		options_with_filter_factory, NULL);
	// 	rocksdb_options_destroy(options_with_filter_factory);
	// }

	// StartPhase("compaction_filter_v2");
	// {
	// 	rocksdb_compactionfilterfactoryv2_t* factory;
	// 	rocksdb_slicetransform_t* prefix_extractor;
	// 	prefix_extractor = rocksdb_slicetransform_create(
	// 		NULL, CFV2PrefixExtractorDestroy, CFV2PrefixExtractorTransform,
	// 		CFV2PrefixExtractorInDomain, CFV2PrefixExtractorInRange,
	// 		CFV2PrefixExtractorName);
	// 	factory = rocksdb_compactionfilterfactoryv2_create(
	// 		prefix_extractor, prefix_extractor, CompactionFilterFactoryV2Destroy,
	// 		CompactionFilterFactoryV2Create, CompactionFilterFactoryV2Name);
	// 	// Create new database
	// 	rocksdb_close(db);
	// 	rocksdb_destroy_db(options, dbname, &err);
	// 	rocksdb_options_set_compaction_filter_factory_v2(options, factory);
	// 	db = rocksdb_open(options, dbname, &err);
	// 	CheckNoError(err);
	// 	// Only foo2 is GC'd, foo3 is changed.
	// 	rocksdb_put(db, woptions, "foo1", 4, "no gc", 5, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo2", 4, "gc", 2, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "foo3", 4, "change", 6, &err);
	// 	CheckNoError(err);
	// 	// All bars are GC'd.
	// 	rocksdb_put(db, woptions, "bar1", 4, "no gc", 5, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "bar2", 4, "gc all", 6, &err);
	// 	CheckNoError(err);
	// 	rocksdb_put(db, woptions, "bar3", 4, "no gc", 5, &err);
	// 	CheckNoError(err);
	// 	// Compact the DB to garbage collect.
	// 	rocksdb_compact_range(db, NULL, 0, NULL, 0);

	// 	// Verify foo entries.
	// 	CheckGet(db, roptions, "foo1", "no gc");
	// 	CheckGet(db, roptions, "foo2", NULL);
	// 	CheckGet(db, roptions, "foo3", "changed");
	// 	// Verify bar entries were all deleted.
	// 	CheckGet(db, roptions, "bar1", NULL);
	// 	CheckGet(db, roptions, "bar2", NULL);
	// 	CheckGet(db, roptions, "bar3", NULL);
	// }

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
