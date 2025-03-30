package influxdb_client

import (
	"encoding/json"
	"fmt"
	stscache "github.com/timescale/tsbs/TimescaleDB_Client/stscache_client"
	jscodec "github.com/timescale/tsbs/jsoncodec"
	"google.golang.org/protobuf/proto"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTransmitExperimentProtobuf(t *testing.T) {
	var conn, _ = NewHTTPClient(HTTPConfig{
		Addr: "http://192.168.1.101:8086",
	})
	tableName := "devops_small"
	//queryString := "SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM \"cpu\" WHERE (hostname = 'host_1' or hostname = 'host_2' or hostname = 'host_3' or hostname = 'host_4' or hostname = 'host_5' or hostname = 'host_6' or hostname = 'host_7' or hostname = 'host_8' or hostname = 'host_9' or hostname = 'host_10') " +
	//	"AND TIME >= '2022-08-30T10:00:00Z' AND TIME < '2022-08-30T10:30:00Z' GROUP BY \"hostname\",time(1m)"
	//
	//segment := "{(cpu.hostname=host_1)(cpu.hostname=host_2)(cpu.hostname=host_3)(cpu.hostname=host_4)(cpu.hostname=host_5)(cpu.hostname=host_6)(cpu.hostname=host_7)(cpu.hostname=host_8)(cpu.hostname=host_9)(cpu.hostname=host_10)}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,1m}"

	queryTags, segmentTags := generateDevopsQueryTags(10)
	queryString := fmt.Sprintf("SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM \"cpu\" WHERE (%s) "+
		"AND TIME >= '2022-08-30T10:00:00Z' AND TIME < '2022-09-02T10:00:00Z' GROUP BY \"hostname\",time(1m)", queryTags)

	segment := fmt.Sprintf("{%s}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,1m}", segmentTags)

	partialSegment, fields, metric := SplitPartialSegment(segment)
	starSegment := GetStarSegment(metric, partialSegment)
	_, startTime, endTime, tags := GetQueryTemplate(queryString)
	singleSegments := GetSingleSegment(metric, partialSegment, tags)
	fields = "time[int64]," + fields
	//datatypes := GetDataTypeArrayFromSF(fields)
	//startTime := TimeStringToInt64("2022-09-30T09:00:00Z")
	//endTime := TimeStringToInt64("2022-10-31T10:00:00Z")

	qry := NewQuery(queryString, tableName, "s")
	resp, err := conn.Query(qry)
	if err != nil {
		panic(err)
	}

	semMeta := ResponseToProto(resp, starSegment, singleSegments)
	data, err := proto.Marshal(semMeta)
	fmt.Println(len(data))
	dlb, _ := Int64ToByteArray(int64(len(data) + 8))
	data = append(dlb, data...)
	//fmt.Println(len(data))

	//seriesNum := len(semMeta.SeriesArray)
	//sampleNum := len(semMeta.SeriesArray[0].Values)

	//fmt.Println("old: series num : ", len(semMeta.SeriesArray))
	//fmt.Println("old: sample num : ", len(semMeta.SeriesArray[0].Values))

	cacheConn := stscache.New("192.168.1.102:11211")

	err = cacheConn.Set(&stscache.Item{
		Key:        starSegment,
		Value:      data,
		Time_start: startTime,
		Time_end:   endTime,
	})
	if err != nil {
		return
	}

	getNum := 64
	totalByteLength := uint64(0)
	totalLatency := uint64(0)
	var wg sync.WaitGroup
	for i := 0; i < getNum; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			//fmt.Println(i)
			st := time.Now()
			values, _, err := cacheConn.Get(segment, startTime, endTime)
			//resp, _, _, _, _ = ByteArrayToResponseWithDatatype(values, datatypes)
			//fmt.Println(len(values))
			if len(values) < 2 {
				return
			}
			err = proto.Unmarshal(values[:len(values)-2], semMeta)
			//println("************************")
			if err != nil {
				panic(err)
			}

			//assert.Equal(t, seriesNum, len(semMeta.SeriesArray))
			//assert.Equal(t, sampleNum, len(semMeta.SeriesArray[0].Values))
			//fmt.Println("new: series num : ", len(semMeta.SeriesArray))
			//fmt.Println("new: sample num : ", len(semMeta.SeriesArray[0].Values))

			_ = ProtoToResponse(semMeta)

			atomic.AddUint64(&totalLatency, uint64(time.Since(st).Nanoseconds()))
			atomic.AddUint64(&totalByteLength, uint64(len(values)))
		}()
	}

	wg.Wait()

	avgByteLength := totalByteLength / uint64(getNum)
	avgLatency := totalLatency / uint64(getNum)
	fmt.Println("avg byte length : ", avgByteLength)
	fmt.Println("avg latency (ns) : ", avgLatency)
}

func TestTransmitExperimentCustom(t *testing.T) {
	var conn, _ = NewHTTPClient(HTTPConfig{
		Addr: "http://192.168.1.101:8086",
	})
	tableName := "devops_small"
	//queryString := "SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM \"cpu\" WHERE (hostname = 'host_1' or hostname = 'host_2' or hostname = 'host_3' or hostname = 'host_4' or hostname = 'host_5' or hostname = 'host_6' or hostname = 'host_7' or hostname = 'host_8' or hostname = 'host_9' or hostname = 'host_10') " +
	//	"AND TIME >= '2022-08-30T10:00:00Z' AND TIME < '2022-08-30T10:30:00Z' GROUP BY \"hostname\",time(1m)"
	//
	//segment := "{(cpu.hostname=host_1)(cpu.hostname=host_2)(cpu.hostname=host_3)(cpu.hostname=host_4)(cpu.hostname=host_5)(cpu.hostname=host_6)(cpu.hostname=host_7)(cpu.hostname=host_8)(cpu.hostname=host_9)(cpu.hostname=host_10)}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,1m}"

	queryTags, segmentTags := generateDevopsQueryTags(100)
	queryString := fmt.Sprintf("SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM \"cpu\" WHERE (%s) "+
		"AND TIME >= '2022-08-30T10:00:00Z' AND TIME < '2022-08-31T10:00:00Z' GROUP BY \"hostname\",time(1m)", queryTags)

	segment := fmt.Sprintf("{%s}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,1m}", segmentTags)

	partialSegment, fields, metric := SplitPartialSegment(segment)
	starSegment := GetStarSegment(metric, partialSegment)
	_, startTime, endTime, tags := GetQueryTemplate(queryString)
	//singleSegments := GetSingleSegment(metric, partialSegment, tags)
	fields = "time[int64]," + fields
	datatypes := GetDataTypeArrayFromSF(fields)

	//startTime := TimeStringToInt64("2022-09-30T09:00:00Z")
	//endTime := TimeStringToInt64("2022-10-31T10:00:00Z")

	qry := NewQuery(queryString, tableName, "s")
	resp, err := conn.Query(qry)
	if err != nil {
		panic(err)
	}

	data := ResponseToByteArrayWithParams(resp, datatypes, tags, metric, partialSegment)
	//fmt.Println(len(data))

	cacheConn := stscache.New("192.168.1.102:11211")

	err = cacheConn.Set(&stscache.Item{
		Key:        starSegment,
		Value:      data,
		Flags:      0,
		Expiration: 0,
		CasID:      0,
		Time_start: startTime,
		Time_end:   endTime,
	})
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	getNum := 64
	totalByteLength := uint64(0)
	totalLatency := uint64(0)
	for i := 0; i < getNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//fmt.Println(i)
			st := time.Now()
			values, _, _ := cacheConn.Get(segment, startTime, endTime)
			resp, _, _, _, _ = ByteArrayToResponseWithDatatype(values, datatypes)
			//fmt.Println(resp)
			//fmt.Println(len(values))
			//println("************************")
			atomic.AddUint64(&totalLatency, uint64(time.Since(st).Nanoseconds()))
			atomic.AddUint64(&totalByteLength, uint64(len(values)))
		}()

	}

	wg.Wait()
	avgByteLength := totalByteLength / uint64(getNum)
	avgLatency := totalLatency / uint64(getNum)
	fmt.Println("avg byte length : ", avgByteLength)
	fmt.Println("avg latency (ns) : ", avgLatency)
}

func TestTransmitExperimentJson(t *testing.T) {
	var conn, _ = NewHTTPClient(HTTPConfig{
		Addr: "http://192.168.1.101:8086",
	})
	tableName := "devops_small"
	//queryString := "SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM \"cpu\" WHERE (hostname = 'host_1' or hostname = 'host_2' or hostname = 'host_3' or hostname = 'host_4' or hostname = 'host_5' or hostname = 'host_6' or hostname = 'host_7' or hostname = 'host_8' or hostname = 'host_9' or hostname = 'host_10') " +
	//	"AND TIME >= '2022-08-30T10:00:00Z' AND TIME < '2022-08-30T10:30:00Z' GROUP BY \"hostname\",time(1m)"
	//
	//segment := "{(cpu.hostname=host_1)(cpu.hostname=host_2)(cpu.hostname=host_3)(cpu.hostname=host_4)(cpu.hostname=host_5)(cpu.hostname=host_6)(cpu.hostname=host_7)(cpu.hostname=host_8)(cpu.hostname=host_9)(cpu.hostname=host_10)}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,1m}"

	queryTags, segmentTags := generateDevopsQueryTags(10)
	queryString := fmt.Sprintf("SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM \"cpu\" WHERE (%s) "+
		"AND TIME >= '2022-08-30T10:00:00Z' AND TIME < '2022-08-30T22:00:00Z' GROUP BY \"hostname\",time(1m)", queryTags)

	segment := fmt.Sprintf("{%s}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,1m}", segmentTags)

	partialSegment, fields, metric := SplitPartialSegment(segment)
	starSegment := GetStarSegment(metric, partialSegment)
	_, startTime, endTime, tags := GetQueryTemplate(queryString)
	singleSegments := GetSingleSegment(metric, partialSegment, tags)
	fields = "time[int64]," + fields
	//datatypes := GetDataTypeArrayFromSF(fields)
	//startTime := TimeStringToInt64("2022-09-30T09:00:00Z")
	//endTime := TimeStringToInt64("2022-10-31T10:00:00Z")

	qry := NewQuery(queryString, tableName, "s")
	resp, err := conn.Query(qry)
	if err != nil {
		panic(err)
	}

	semMeta := ResponseToJson(resp, starSegment, singleSegments)
	data, err := json.MarshalIndent(semMeta, "", " ")
	dlb, _ := Int64ToByteArray(int64(len(data) + 8))
	data = append(dlb, data...)
	//data, _ := json.Marshal(semMeta)
	//fmt.Println(len(data))

	cacheConn := stscache.New("192.168.1.102:11211")

	err = cacheConn.Set(&stscache.Item{
		Key:        starSegment,
		Value:      data,
		Flags:      0,
		Expiration: 0,
		CasID:      0,
		Time_start: startTime,
		Time_end:   endTime,
	})
	if err != nil {
		return
	}

	getNum := 64

	totalByteLength := uint64(0)
	totalLatency := uint64(0)
	var wg sync.WaitGroup
	for i := 0; i < getNum; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			//fmt.Println(i)
			st := time.Now()
			values, _, err := cacheConn.Get(segment, startTime, endTime)
			if len(values) < 2 {
				return
			}
			var newSemMeta jscodec.SemanticMetaValue
			err = json.Unmarshal(values, &newSemMeta)
			if err != nil {
				panic(err)
			}

			//newResp := JsonToResponse(&newSemMeta)
			//fmt.Println(newResp)

			_ = JsonToResponse(&newSemMeta)

			//println("************************")
			if err != nil {
				panic(err)
			}

			atomic.AddUint64(&totalLatency, uint64(time.Since(st).Nanoseconds()))
			atomic.AddUint64(&totalByteLength, uint64(len(values)))
		}()
	}

	wg.Wait()

	avgByteLength := totalByteLength / uint64(getNum)
	avgLatency := totalLatency / uint64(getNum)
	fmt.Println("avg byte length : ", avgByteLength)
	fmt.Println("avg latency (ns) : ", avgLatency)
}
