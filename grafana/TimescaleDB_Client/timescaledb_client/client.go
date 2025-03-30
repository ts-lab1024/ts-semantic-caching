package timescaledb_client

import (
	"database/sql"
	"fmt"
	stscache "github.com/grafana/grafana/TimescaleDB_Client/stscache_client"
	"log"
	"strings"
	"sync"
	"time"
)

const pgxDriver = "pgx"

var DB = "devops_small"
var DbName = ""
var STsCacheURL string
var UseCache = "db"
var MaxThreadNum = 64

const STRINGBYTELENGTH = 32

var CacheHash = make(map[string]int)

func GetCacheHashValue(fields string) int {
	CacheNum := len(STsConnArr)
	if CacheNum == 0 {
		CacheNum = 1
	}
	if _, ok := CacheHash[fields]; !ok {
		value := len(CacheHash) % CacheNum
		CacheHash[fields] = value
	}
	hashValue := CacheHash[fields]
	return hashValue
}

var mtx sync.Mutex

var STsConnArr []*stscache.Client

func InitStsConnsArr(urlArr []string) []*stscache.Client {
	conns := make([]*stscache.Client, 0)
	for i := 0; i < len(urlArr); i++ {
		conns = append(conns, stscache.New(urlArr[i]))
	}
	return conns
}

//func init() {
//	STsConnArr = InitStsConnsArr([]string{"192.168.1.101"})
//}

func STsCacheClientSeg(conn *sql.DB, queryString string, semanticSegment string) ([][][]interface{}, uint64, uint8) {

	CacheNum := len(STsConnArr)

	if CacheNum == 0 {
		CacheNum = 1
	}

	byteLength := uint64(0)
	hitKind := uint8(0)

	queryTemplate, startTime, endTime, tags := GetQueryTemplate(queryString)

	partialSegment := ""
	fields := ""
	metric := ""
	partialSegment, fields, metric = SplitPartialSegment(semanticSegment)

	starSegment := GetStarSegment(metric, partialSegment)

	CacheIndex := 0
	fields = "time[int64],name[string]," + fields
	colLen := strings.Split(fields, ",")
	datatypes := DataTypeFromColumn(len(colLen))

	values, _, err := STsConnArr[CacheIndex].Get(semanticSegment, startTime, endTime)
	if err != nil {
		rows, err := conn.Query(queryString)
		if err != nil {
			log.Println(queryString)
		}

		var dataArray [][]interface{} = nil
		if !ResponseIsEmpty(rows) {
			dataArray = RowsToInterface(rows, len(datatypes))
			remainValues, _ := ResponseInterfaceToByteArrayWithParams(dataArray, datatypes, tags, metric, partialSegment)
			err = STsConnArr[CacheIndex].Set(&stscache.Item{Key: starSegment, Time_start: startTime, Time_end: endTime, Value: remainValues})
		} else {
			singleSemanticSegment := GetSingleSegment(metric, partialSegment, tags)
			emptyValues := make([]byte, 0)
			zero, _ := Int64ToByteArray(int64(0))
			nb, _ := Int64ToByteArray(int64(len(singleSemanticSegment)))
			emptyValues = append(emptyValues, nb...)
			for _, ss := range singleSemanticSegment {
				emptyValues = append(emptyValues, []byte(ss)...)
				emptyValues = append(emptyValues, []byte(" ")...)
				emptyValues = append(emptyValues, zero...)
			}
			err = STsConnArr[CacheIndex].Set(&stscache.Item{Key: starSegment, Time_start: startTime, Time_end: endTime, Value: emptyValues})
		}

		return [][][]interface{}{dataArray}, byteLength, hitKind

	} else {
		convertedResponse, flagNum, flagArr, timeRangeArr, tagArr := ByteArrayToResponseWithDatatype(values, datatypes)
		if flagNum == 0 {
			hitKind = 2
			return convertedResponse, byteLength, hitKind
		} else {
			hitKind = 1
			remainQueryString, minTime, maxTime := RemainQuery(queryTemplate, flagArr, timeRangeArr, tagArr)
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
			remainRows, err := conn.Query(remainQueryString)
			if err != nil {
				log.Println(remainQueryString)
			}
			if ResponseIsEmpty(remainRows) {
				hitKind = 2
				singleSemanticSegment := GetSingleSegment(metric, partialSegment, remainTags)
				emptyValues := make([]byte, 0)
				zero, _ := Int64ToByteArray(int64(0))
				nb, _ := Int64ToByteArray(int64(len(singleSemanticSegment)))
				emptyValues = append(emptyValues, nb...)
				for _, ss := range singleSemanticSegment {
					emptyValues = append(emptyValues, []byte(ss)...)
					emptyValues = append(emptyValues, []byte(" ")...)
					emptyValues = append(emptyValues, zero...)
				}
				err = STsConnArr[CacheIndex].Set(&stscache.Item{Key: starSegment, Time_start: startTime, Time_end: endTime, Value: emptyValues})
				return convertedResponse, byteLength, hitKind
			}

			remainDataArray := RowsToInterface(remainRows, len(datatypes))
			remainByteArr, _ := ResponseInterfaceToByteArrayWithParams(remainDataArray, datatypes, remainTags, metric, partialSegment)
			err = STsConnArr[CacheIndex].Set(&stscache.Item{Key: starSegment, Time_start: minTime, Time_end: maxTime, Value: remainByteArr})
			totalResp := MergeRemainResponse(convertedResponse, remainDataArray, timeRangeArr)
			return totalResp, byteLength, hitKind
		}

	}

}

var TotalLatency int64
var TotalCount int64

func init() {
	STsConnArr = InitStsConnsArr([]string{"192.168.1.102:11211"})
}

func GrafanaSTsCacheClientSeg(conn *sql.DB, queryString string, semanticSegment string) ([][][]interface{}, []string, []string, error) {

	queryTemplate, startTime, endTime, tags := GetQueryTemplate(queryString)

	partialSegment := ""
	fields := ""
	metric := ""
	partialSegment, fields, metric = SplitPartialSegment(semanticSegment)

	starSegment := GetStarSegment(metric, partialSegment)

	fields = "time[int64]," + fields
	columns := strings.Split(fields, ",")
	for i, col := range columns {
		columns[i] = col[:strings.Index(col, "[")]
	}
	datatypes := DataTypeFromColumn(len(columns))

	values, _, err := STsConnArr[0].Get(semanticSegment, startTime, endTime)
	if err != nil {
		rows, err := conn.Query(queryString)
		if err != nil {
			log.Println(queryString)
		}

		var dataArray [][]interface{} = nil
		if !ResponseIsEmpty(rows) {
			dataArray = RowsToInterface(rows, len(datatypes))
			remainValues, _ := ResponseInterfaceToByteArrayWithParams(dataArray, datatypes, tags, metric, partialSegment)
			err = STsConnArr[0].Set(&stscache.Item{Key: starSegment, Time_start: startTime, Time_end: endTime, Value: remainValues})
			if err != nil {
				fmt.Println(err)
			}
		} else {
			singleSemanticSegment := GetSingleSegment(metric, partialSegment, tags)
			emptyValues := make([]byte, 0)
			zero, _ := Int64ToByteArray(int64(0))
			nb, _ := Int64ToByteArray(int64(len(singleSemanticSegment)))
			emptyValues = append(emptyValues, nb...)
			for _, ss := range singleSemanticSegment {
				emptyValues = append(emptyValues, []byte(ss)...)
				emptyValues = append(emptyValues, []byte(" ")...)
				emptyValues = append(emptyValues, zero...)
			}
			err = STsConnArr[0].Set(&stscache.Item{Key: starSegment, Time_start: startTime, Time_end: endTime, Value: emptyValues})
		}

		return [][][]interface{}{dataArray}, columns, datatypes, err

	} else {
		convertedResponse, flagNum, flagArr, timeRangeArr, tagArr := ByteArrayToResponseWithDatatype(values, datatypes)
		if flagNum == 0 {
			return convertedResponse, columns, datatypes, err
		} else {
			remainQueryString, minTime, maxTime := RemainQuery(queryTemplate, flagArr, timeRangeArr, tagArr)
			remainTags := make([]string, 0)
			for i, tag := range tagArr {
				if flagArr[i] == 1 {
					remainTags = append(remainTags, fmt.Sprintf("%s=%s", tag[0], tag[1]))
				}
			}
			//if maxTime-minTime <= int64(time.Minute.Seconds()) {
			//	return convertedResponse, columns, datatypes, err
			//}
			remainRows, err := conn.Query(remainQueryString)
			if err != nil {
				log.Println(remainQueryString)
			}
			if ResponseIsEmpty(remainRows) {
				singleSemanticSegment := GetSingleSegment(metric, partialSegment, remainTags)
				emptyValues := make([]byte, 0)
				zero, _ := Int64ToByteArray(int64(0))
				nb, _ := Int64ToByteArray(int64(len(singleSemanticSegment)))
				emptyValues = append(emptyValues, nb...)
				for _, ss := range singleSemanticSegment {
					emptyValues = append(emptyValues, []byte(ss)...)
					emptyValues = append(emptyValues, []byte(" ")...)
					emptyValues = append(emptyValues, zero...)
				}
				err = STsConnArr[0].Set(&stscache.Item{Key: starSegment, Time_start: startTime, Time_end: endTime, Value: emptyValues})
				return convertedResponse, columns, datatypes, err
			}

			remainDataArray := RowsToInterface(remainRows, len(datatypes))
			remainByteArr, _ := ResponseInterfaceToByteArrayWithParams(remainDataArray, datatypes, remainTags, metric, partialSegment)

			err = STsConnArr[0].Set(&stscache.Item{Key: starSegment, Time_start: minTime, Time_end: maxTime, Value: remainByteArr})

			totalResp := MergeRemainResponse2(convertedResponse, remainDataArray, timeRangeArr)
			return totalResp, columns, datatypes, err
		}

	}

}
