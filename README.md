####Step one:

 * Install rocksdb 
 * The tested version is 4.5 (latest)
 * Install gflags
 * Install snappy.
 * Install zlib.
 * Install bzip2:

####Step two:

 1. Git clone https://github.com/pcjdean/rocksdbgo.git
 2. Go to source folder
 3. Run  'go test -v'

*Note: It's tested on Ubuntu so far.*

###Features of this implementation:

 1. It's based totally on the C++ head files of rocksdb. So, it potentially can support all the features of rocksdb.
 2. There is a C wrapper between the go and C++ rocksdb. Some logics can be freely added on this layer.
 3. Everything is pure golang on the go layer including callback classes and garbage collection of DB related objects. 
 4. DB related objects can be closed manually and garbage collected.
 5. It supports most of features already and can be easily extended with the current code struct.

###Usage:
Check TestCMain() in db_test.go
