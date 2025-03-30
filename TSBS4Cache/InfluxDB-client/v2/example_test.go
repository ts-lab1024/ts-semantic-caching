package influxdb_client

import (
	"fmt"
	stscache "github.com/timescale/tsbs/InfluxDB-client/memcache"
	"testing"
)

func InitClient(influxdbConn Client, dbName string) {
	TagKV = GetTagKV(influxdbConn, dbName)
	Fields = GetFieldKeys(influxdbConn, dbName)
}

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
