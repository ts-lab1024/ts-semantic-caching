package influxdb_client

import (
	"fmt"
	stscache "github.com/timescale/tsbs/TimescaleDB_Client/stscache_client"
	"log"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	minute = time.Minute
	hour   = time.Hour
	day    = 24 * time.Hour
	week   = 7 * day
	month  = 4*week + 2*day
)

var TimeDuration = []time.Duration{
	6 * hour, 12 * hour, 1 * day, 2 * day, 3 * day, 5 * day, 1 * week, 2 * week, 3 * week, 1 * month,
}

var iniStartTime int64 = TimeStringToInt64("2022-09-01T00:00:00Z")
var iniEndTime int64 = TimeStringToInt64("2022-10-01T00:00:00Z")

var qryTemplate string = `SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM "cpu" WHERE (%s) AND TIME >= '%s' AND TIME < '%s' GROUP BY "hostname",time(5m)`
var sgmtTemplate string = `{%s}#{usage_system[float64],usage_idle[float64],usage_nice[float64]}#{empty}#{mean,5m}`

var qryWorkloads [][]string = make([][]string, 0)

func generateTimeInterval() (startTimeStr, endTimeStr string) {
	duration := TimeDuration[rand.Intn(len(TimeDuration))]

	startTime := iniStartTime
	endTime := iniStartTime + duration.Microseconds()/1e6

	return TimeInt64ToString(startTime), TimeInt64ToString(endTime)
}

func generateDevopsQueryTags(num int) (string, string) {
	queryTags := make([]string, 0)
	segmentTags := make([]string, 0)

	for i := 0; i < num; i++ {
		queryTags = append(queryTags, fmt.Sprintf("hostname = 'host_%d'", i))
		segmentTags = append(segmentTags, fmt.Sprintf("(cpu.hostname=host_%d)", i))
	}

	return strings.Join(queryTags, " or "), strings.Join(segmentTags, "")
}

func generateDevopsQueryWorkloads(seriesNum, workloadNum int) {
	for i := 0; i < workloadNum; i++ {
		st, et := generateTimeInterval()
		qryTags, sgmtTags := generateDevopsQueryTags(seriesNum)
		qryStr := fmt.Sprintf(qryTemplate, qryTags, st, et)
		sgmtStr := fmt.Sprintf(sgmtTemplate, sgmtTags)

		qryWorkloads = append(qryWorkloads, []string{qryStr, sgmtStr})
	}
}

func TestPartialHit(t *testing.T) {
	var conn, _ = NewHTTPClient(HTTPConfig{
		Addr: "http://192.168.1.101:8086",
	})
	tableName := "devops_small"
	cacheConn := stscache.New("192.168.1.102:11211")

	//queryString := "SELECT mean(usage_system),mean(usage_idle),mean(usage_nice) FROM \"cpu\" WHERE (hostname = 'host_1' or hostname = 'host_2' or hostname = 'host_3' or hostname = 'host_4' or hostname = 'host_5' or hostname = 'host_6' or hostname = 'host_7' or hostname = 'host_8' or hostname = 'host_9' or hostname = 'host_10') " +
	//	"AND TIME >= '2022-08-30T10:00:00Z' AND TIME < '2022-08-30T10:30:00Z' GROUP BY \"hostname\",time(1m)"
	//
	//segment := "{(cpu.hostname=host_1)(cpu.hostname=host_2)(cpu.hostname=host_3)(cpu.hostname=host_4)(cpu.hostname=host_5)(cpu.hostname=host_6)(cpu.hostname=host_7)(cpu.hostname=host_8)(cpu.hostname=host_9)(cpu.hostname=host_10)}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,1m}"

	seriesNum := 50
	generateDevopsQueryWorkloads(seriesNum, 20000)

	setQryTags, setSgmtTags := generateDevopsQueryTags(40)

	setQryStr := fmt.Sprintf(qryTemplate, setQryTags, TimeInt64ToString(iniStartTime), TimeInt64ToString(iniEndTime))
	setSgmtStr := fmt.Sprintf(sgmtTemplate, setSgmtTags)
	//setQryStr = strings.Replace(setQryStr, "$", TimeInt64ToString(iniStartTime), 1)
	//setQryStr = strings.Replace(setQryStr, "$", TimeInt64ToString(iniEndTime), 1)

	qry := NewQuery(setQryStr, tableName, "s")
	resp, err := conn.Query(qry)
	if err != nil {
		panic(err)
	}

	partialSegment, fields, metric := SplitPartialSegment(setSgmtStr)
	starSegment := GetStarSegment(metric, partialSegment)
	_, startTime, endTime, tags := GetQueryTemplate(setQryStr)
	//singleSegments := GetSingleSegment(metric, partialSegment, tags)
	fields = "time[int64]," + fields
	datatypes := GetDataTypeArrayFromSF(fields)
	data := ResponseToByteArrayWithParams(resp, datatypes, tags, metric, partialSegment)

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
	threadNum := 64
	semaphore := make(chan struct{}, threadNum)
	totalLatency := uint64(0)
	totalCount := uint64(0)
	lastLatency := uint64(0)

	for _, qs := range qryWorkloads {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(qs []string) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()

			// todo cache
			now := time.Now()

			//q := NewQuery(qs[0], tableName, "s")
			//_, err := conn.Query(q)
			//if err != nil {
			//	panic(err)
			//}

			_, _, _ = sTsCacheClientOnlyGet(conn, tableName, cacheConn, qs[0], qs[1])

			atomic.AddUint64(&totalLatency, uint64(time.Since(now).Nanoseconds()))
			atomic.AddUint64(&totalCount, 1)

			interval := uint64(100)
			if atomic.LoadUint64(&totalCount)%interval == 0 {
				fmt.Printf("after %d queries, avg latency is : %d ns, \tlast %d queries latency is %d ns\n", atomic.LoadUint64(&totalCount), atomic.LoadUint64(&totalLatency)/atomic.LoadUint64(&totalCount), interval, (atomic.LoadUint64(&totalLatency)-lastLatency)/interval)
				lastLatency = atomic.LoadUint64(&totalLatency)
			}

		}(qs)
	}

	go func() {
		wg.Wait()
	}()

	avgLatency := totalLatency / totalCount
	fmt.Println("total count : ", totalCount)
	fmt.Println("avg latency (ns) : ", avgLatency)
}

func sTsCacheClientOnlyGet(dbConn Client, tableName string, cacheConn *stscache.Client, queryString string, semanticSegment string) (*Response, uint64, uint8) {

	byteLength := uint64(0)
	hitKind := uint8(0)

	queryTemplate, startTime, endTime, tags := GetQueryTemplate(queryString)

	partialSegment := ""
	fields := ""
	metric := ""
	partialSegment, fields, metric = SplitPartialSegment(semanticSegment)

	//starSegment := GetStarSegment(metric, partialSegment)

	fields = "time[int64]," + fields
	datatypes := GetDataTypeArrayFromSF(fields)

	values, _, err := cacheConn.Get(semanticSegment, startTime, endTime)
	if err != nil {
		q := NewQuery(queryString, tableName, "s")
		resp, err := dbConn.Query(q)
		if err != nil {
			log.Println(queryString)
		}

		if !ResponseIsEmpty(resp) {
			_ = ResponseToByteArrayWithParams(resp, datatypes, tags, metric, partialSegment)
			//err = cacheConn.Set(&stscache.Item{Key: starSegment, Value: remainValues, Time_start: startTime, Time_end: endTime, NumOfTables: numOfTab})
		} else {
			//singleSemanticSegment := GetSingleSegment(metric, partialSegment, tags)
			//emptyValues := make([]byte, 0)
			//for _, ss := range singleSemanticSegment {
			//	zero, _ := Int64ToByteArray(int64(0))
			//	emptyValues = append(emptyValues, []byte(ss)...)
			//	emptyValues = append(emptyValues, []byte(" ")...)
			//	emptyValues = append(emptyValues, zero...)
			//}
			//numOfTab := int64(len(singleSemanticSegment))
			//err = cacheConn.Set(&stscache.Item{Key: starSegment, Value: emptyValues, Time_start: startTime, Time_end: endTime, NumOfTables: numOfTab})
		}

		return resp, byteLength, hitKind

	} else {
		convertedResponse, flagNum, flagArr, timeRangeArr, tagArr := ByteArrayToResponseWithDatatype(values, datatypes)
		if flagNum == 0 {
			hitKind = 2
			return convertedResponse, byteLength, hitKind
		} else {
			hitKind = 1

			remainQueryString, minTime, maxTime := RemainQueryString(queryString, queryTemplate, flagArr, timeRangeArr, tagArr)

			remainTags := make([]string, 0)
			for i, tag := range tagArr {
				if flagArr[i] == 1 {
					remainTags = append(remainTags, fmt.Sprintf("%s=%s", tag[0], tag[1]))
				}

			}

			if maxTime-minTime <= int64(time.Minute.Seconds()) {
				hitKind = 2
				return convertedResponse, byteLength, hitKind
			}

			remainQuery := NewQuery(remainQueryString, tableName, "s")
			remainResp, err := dbConn.Query(remainQuery)
			if err != nil {
				log.Println(remainQueryString)
			}
			if ResponseIsEmpty(remainResp) {
				hitKind = 1

				//singleSemanticSegment := GetSingleSegment(metric, partialSegment, remainTags)
				//emptyValues := make([]byte, 0)
				//for _, ss := range singleSemanticSegment {
				//	zero, _ := Int64ToByteArray(int64(0))
				//	emptyValues = append(emptyValues, []byte(ss)...)
				//	emptyValues = append(emptyValues, []byte(" ")...)
				//	emptyValues = append(emptyValues, zero...)
				//}
				//
				//numOfTab := int64(len(singleSemanticSegment))
				//err = cacheConn.Set(&stscache.Item{Key: starSegment, Value: emptyValues, Time_start: startTime, Time_end: endTime, NumOfTables: numOfTab})
				return convertedResponse, byteLength, hitKind
			}

			_ = RemainResponseToByteArrayWithParams(remainResp, datatypes, remainTags, metric, partialSegment)
			_ = len(remainResp.Results)

			//err = cacheConn.Set(&stscache.Item{
			//	Key:         starSegment,
			//	Value:       remainByteArr,
			//	Time_start:  minTime,
			//	Time_end:    maxTime,
			//	NumOfTables: int64(numOfTableR),
			//})
			totalResp := MergeRemainResponse(remainResp, convertedResponse)

			return totalResp, byteLength, hitKind
		}

	}

}
