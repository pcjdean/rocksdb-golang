Step one:

Install rocksdb 

Step two:

 1) Go to source folder
 2) Run  'go test -v'

Sample output:

=== RUN TestCMain
rocksdbgo version = 1.0
rocksdb version = 4.5
dbname = /tmp/rocksdb_go_test-5346208
dbbackupname = /tmp/rocksdb_go_test-5346208-backup
--- PASS: TestCMain (2.42s)
	db_test.go:246: phase: create_objects
	db_test.go:274: phase: Destroy
	db_test.go:276: DestroyDB: status = OK
	db_test.go:278: phase: open_error
	db_test.go:283: open_error: status = Invalid argument: /tmp/rocksdb_go_test-5346208: does not exist (create_if_missing is false)
	db_test.go:286: phase: open
	db_test.go:292: phase: get
	db_test.go:300: phase: put
	db_test.go:309: phase: backup_and_restore
	db_test.go:348: phase: compactall
	db_test.go:356: phase: compactrange
	db_test.go:363: phase: writebatch
	db_test.go:385: phase: writebatch_rep
	db_test.go:397: phase: iter
	db_test.go:415: phase: approximate_sizes
	db_test.go:431: phase: property
	db_test.go:439: phase: snapshot
	db_test.go:451: phase: repair
	db_test.go:474: phase: filter
	db_test.go:486: filter: DestroyDB: status = OK
	db_test.go:486: filter: DestroyDB: status = OK
	db_test.go:521: phase: compaction_filter
	db_test.go:527: compaction_filter: DestroyDB: status = OK
	db_test.go:535: phase: compaction_filter_factory
	db_test.go:542: compaction_filter_factory: DestroyDB: status = OK
	db_test.go:548: phase: merge_operator
	db_test.go:553: merge_operator: DestroyDB: status = OK
	db_test.go:575: phase: columnfamilies
	db_test.go:578: columnfamilies: DestroyDB: status = OK
	db_test.go:658: columnfamilies: last: DestroyDB: status = OK
	db_test.go:662: phase: prefix
	db_test.go:723: prefix: DestroyDB: status = OK
	db_test.go:725: phase: cuckoo_options
	db_test.go:739: phase: iterate_upper_bound
	db_test.go:795: phase: cleanup
PASS
ok  	_/media/sf_VBoxShare/Projects/Go/shareme/db/src/rocksdb	2.430s

