# STsCache

## Introduction
This is the source code of STsCache (using three serialization methods: custom, Protobuf, and JSON) and Bench.

Bench, built on [TSBS](https://github.com/timescale/tsbs) Expanded, measures throughput, latency, and hit ratio, utilizing [YCSB-TS](https://github.com/TSDBBench/YCSB-TS) and distributions like Latest and Zipfian for workload generation.

Additionally, we also provide the source code of a Grafana integrated with the STsCache client. For detailed usage of Grafana, please refer to [Grafana](https://github.com/grafana/grafana).

## Dependencies

The code runs under Ubuntu 20.04.

* C++17

* clang 12.0.1 x86_64-pc-linux-gnu

* cmake version 3.10

* gcc version 9.4.0

* go 1.21.3

database:

* InfluxDB 1.8.10
* TimescaleDB 2.6.0


## Build

**STsCache:**   in `STsCache/` directory

```
$ cmake -B build -DCMAKE_C_COMPILER:FILEPATH=clang-12 -DCMAKE_CXX_COMPILER:FILEPATH=clang++-12 ./
$ cd build && make
```

**Bench:**  in `bench/` directory

```
$ go mod tidy
$ make
```



## Run STsCache

```
$ ~/STsCache/build/bin/STsCache [-p port] [-D ssd device] [-m max slab memory]
```

**Parameters:**

* -p : port used for data communication
* -D : storage path
* -m : allocated memory size


## Run Bench

Our work supports InfluxDB and TimescaleDB now. Each database has its own `IoT` and `DevOps` datasets. Please refer to [TSBS](https://github.com/timescale/tsbs) for details on these datasets.

### Generate and load data

The operations of generating data and loading it into database are the same as  [TSBS](https://github.com/timescale/tsbs).

### Generate query workloads

in `bench/bin` directory

We use query generation from two datasets on two databases respectively as an example.

**InfluxDB-IoT-small:**

```
$ ./tsbs_generate_queries --seed=123 \
        --timestamp-start="2022-01-01T08:00:00Z" \ 
        --timestamp-end="2022-12-31T00:00:01Z"  \
        --db-name="iot_small" \
        --use-case="iot" \
        --query-type="iot-queries" \
        --format="influx" \
        --queries=100000 \
        --random-tag=true \
        --scale=100 \
        --tag-num=10 \
        --file="path_to_workloads/iot_workloads.txt"
```

**TimescaleDB-DevOps-medium：**

```
$ ./tsbs_generate_queries --seed=123 \
        --timestamp-start="2022-01-01T08:00:00Z" \
        --timestamp-end="2022-12-31T00:00:01Z"  
        --db-name="devops_medium" \
        --use-case="devops" 
        --query-type="cpu-queries" \
        --format="timescaledb" \
        --timescale-use-time-bucket=true \
        --queries=100000 \
        --random-tag=true \
        --scale=500 \
        --tag-num=50 \
        --file="path_to_workloads/devops_workloads.txt"
```

**Parameters:**

* use-case: iot / devops
* query-type: iot-queries / cpu-queries
* format: influx / timescaledb
* queries: number of queries
* random-tag: whether the queried tags arranged in order
* scale: total device scale: 100 for small, 500 for medium, 1000 for large
* tag-num: tag query scale: default 10% of scale, can be modified
* file: storage path of query workload file 



### Run query workloads

in `bench/bin` directory

We assume that the address of the database server is 192.168.1.103 and the STsCache server is 192.168.1.101.

***First start STsCache service.***
```
$ ~/STsCache/build/bin/STsCache [-p port] [-D ssd device] [-m max slab memory]
```

**InfluxDB-IoT-small:**

```
$ ./tsbs_run_queries_influx \
        --urls="http://192.168.1.103:8086" \
        --cache-url="192.168.1.101:11211" \
        --workers=64 \
        --use-cache="stscache" \
        --db-name="iot_small"  \
        --file="path_to_workloads/iot_workloads.txt" \
        --print-interval=1000 \
        --max-queries=100000
```

**TimescaleDB-DevOps-medium：**

```
$ ./tsbs_run_queries_timescaledb \
        --hosts="192.168.1.103" \
        --pass="your_password"  \
        --cache-url="192.168.1.101:11211" \
        --workers=64 \
        --use-cache="db" \
        --db-name="devops_medium" \
        --file="path_to_workloads/devops_workloads.txt" \
        --print-interval=10000 \
        --max-queries=100000
```

There are some differences in the way database addresses are processed between the two databases.

**Parameters:**

* urls / hosts: database address and port
* cache-url: cache server address and port
* workers: Number of concurrent requests to make
* use-cache: db / stscache, It determines whether to query directly from the database or first from the cache.
* file: File name to read queries from
* print-interval: Print timing stats to stderr after this many queries (0 to disable)
* burn-in: Number of queries to ignore before collecting statistics
* max-queries: Limit the number of queries to send, 0 = no limit

## Custom Serialization and Deserialization Methods

### Data Points Encoding

```
┌───────────────┬─────────────────────┬───────────────┬─────┬──────────────┬─────────────────────┬─────────────────┬─────┐
│ timestamp1    │   field values1     │  timestamp2   │ ... │ timestampn   │  field values(n)    │ timestamp(n+1)  │ ... │
└───────────────┴─────────────────────┴───────────────┴─────┴──────────────┴─────────────────────┴─────────────────┴─────┘
```


### Serialization value

```
┌────────────────────────┬─────────────────────┬─────────────────┬─────┬─────────────────────┬─────────────────────┬─────────────────────────┬─────┐
│ NumbleofSemanticSeries │  Semantic Series1   │  Data Points1   │ ... │ SemanticSeries(n)   │   Data Points(n)    │  Semantic Series(n+1)   │ ... │
└────────────────────────┴─────────────────────┴─────────────────┴─────┴─────────────────────┴─────────────────────┴─────────────────────────┴─────┘
```

### Example

The detailed test code is located at /TSBS4Cache/InfluxDB-client/v2/example_test.go
```azure
func TestExample(t *testing.T) {
	var influxdbConn, _ = NewHTTPClient(HTTPConfig{
		Addr: "http://192.168.1.103:8086",
	})
	dbName := "devops_small"
	stscacheConn := stscache.New("192.168.1.102:11211")
	queryString := `SELECT mean("usage_nice") FROM "cpu" WHERE ("hostname"='host_0' OR "hostname"='host_1') AND TIME >= '2022-06-01T00:00:00Z' AND TIME < '2022-06-01T10:00:00Z' GROUP BY "hostname",time(1m)`

	InitClient(influxdbConn, dbName)
	qry := NewQuery(queryString, dbName, "s")
	queryResult, err := influxdbConn.Query(qry)
	if err != nil {
		panic(err)
	}
	// Querystring --> WriteSemanticKey and Time-range [startTime, endTime].
	writeSemanticKey := GetWriteSemanticKey(queryString)
	startTime, endTime := GetTimeRange(queryString)

	// Query result --> Serialization value
	serializationValue := CustomSerialization(queryResult, queryString)

	err = stscacheConn.Set(&stscache.Item{
		Key:        writeSemanticKey,
		Value:      serializationValue,
		Time_start: startTime,
		Time_end:   endTime,
	})
	if err != nil {
		panic(err)
	}
	// Querystring --> ReadSemanticKey
	readSemanticKey := GetReadSemanticKey(queryString)
	values, _, err := stscacheConn.Get(readSemanticKey, startTime, endTime)

	//Serialization value -->  query result object.
	queryCacheResult, flagNum, flagArr, timeRangeArr, tagArr := CustomDeSerialization(values, readSemanticKey)
	fmt.Println(queryCacheResult)
	fmt.Println(flagNum)
	fmt.Println(flagArr)
	fmt.Println(timeRangeArr)
	fmt.Println(tagArr)
}
```
