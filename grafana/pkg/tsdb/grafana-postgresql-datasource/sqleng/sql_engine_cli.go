package sqleng

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/grafana/TimescaleDB_Client/timescaledb_client"
	"runtime/debug"
	"sync"
	"time"
)

func (e *DataSourceHandler) executeQueryCLI(query backend.DataQuery, wg *sync.WaitGroup, queryContext context.Context,
	ch chan DBDataResponse, queryJson QueryJson) {
	defer wg.Done()
	queryResult := DBDataResponse{
		dataResponse: backend.DataResponse{},
		refID:        query.RefID,
	}

	logger := e.log.FromContext(queryContext)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("ExecuteQuery panic", "error", r, "stack", string(debug.Stack()))
			if theErr, ok := r.(error); ok {
				queryResult.dataResponse.Error = theErr
			} else if theErrString, ok := r.(string); ok {
				queryResult.dataResponse.Error = errors.New(theErrString)
			} else {
				queryResult.dataResponse.Error = fmt.Errorf("unexpected error - %s", e.userError)
			}
			ch <- queryResult
		}
	}()

	if queryJson.RawSql == "" {
		panic("Query model property rawSql should not be empty at this point")
	}

	timeRange := query.TimeRange

	errAppendDebug := func(frameErr string, err error, query string) {
		var emptyFrame data.Frame
		emptyFrame.SetMeta(&data.FrameMeta{
			ExecutedQueryString: query,
		})
		queryResult.dataResponse.Error = fmt.Errorf("%s: %w", frameErr, err)
		queryResult.dataResponse.Frames = data.Frames{&emptyFrame}
		ch <- queryResult
	}

	// global substitutions
	interpolatedQuery := Interpolate(query, timeRange, e.dsInfo.JsonData.TimeInterval, queryJson.RawSql)

	// data source specific substitutions
	interpolatedQuery, err := e.macroEngine.Interpolate(&query, timeRange, interpolatedQuery)
	if err != nil {
		errAppendDebug("interpolation failed", e.TransformQueryError(logger, err), interpolatedQuery)
		return
	}

	//rows, err := e.db.QueryContext(queryContext, interpolatedQuery)

	//fmt.Println("SQL ", interpolatedQuery)

	/* todo 生成语义段 */
	// SELECT time_bucket('5 minute', time) as bucket,name,avg(latitude) FROM readings WHERE name IN ('truck_0') AND time >= '2022-12-29 06:00:00 +0000' AND time < '2022-12-30 06:00:00 +0000' GROUP BY name,bucket ORDER BY name,bucket
	segment := SQLsegment(interpolatedQuery)
	//fmt.Println(segment)

	values, columnNames, _, err := timescaledb_client.GrafanaSTsCacheClientSeg(e.db, interpolatedQuery, segment)
	//fmt.Println(timescaledb_client.ResultInterfaceToString(values[0], datatypes))
	if err != nil {
		errAppendDebug("query error", e.TransformQueryError(logger, err), interpolatedQuery)
		return
	}

	qm, err := e.newProcessCfgCLI(query, queryContext, interpolatedQuery, columnNames)
	if err != nil {
		errAppendDebug("failed to get configurations", err, interpolatedQuery)
		return
	}

	frames := make([]*data.Frame, 0, len(values))
	for i, value := range values {
		if len(value) == 0 {
			continue
		}
		frame, err := FrameFromRowsInterface(values[i], e.rowLimit, columnNames)
		if err != nil {
			errAppendDebug("convert frame from rows error", err, interpolatedQuery)
			return
		}

		if frame.Meta == nil {
			frame.Meta = &data.FrameMeta{}
		}

		frame.Meta.ExecutedQueryString = interpolatedQuery

		// If no rows were returned, clear any previously set `Fields` with a single empty `data.Field` slice.
		// Then assign `queryResult.dataResponse.Frames` the current single frame with that single empty Field.
		// This assures 1) our visualization doesn't display unwanted empty fields, and also that 2)
		// additionally-needed frame data stays intact and is correctly passed to our visulization.
		if frame.Rows() == 0 {
			frame.Fields = []*data.Field{}
			queryResult.dataResponse.Frames = data.Frames{frame}
			ch <- queryResult
			return
		}

		if err := convertSQLTimeColumnsToEpochMSCLI(frame, qm); err != nil {
			errAppendDebug("converting time columns failed", err, interpolatedQuery)
			return
		}

		if qm.Format == dataQueryFormatSeries {
			// time series has to have time column
			if qm.timeIndex == -1 {
				errAppendDebug("db has no time column", errors.New("no time column found"), interpolatedQuery)
				return
			}

			// Make sure to name the time field 'Time' to be backward compatible with Grafana pre-v8.
			frame.Fields[qm.timeIndex].Name = data.TimeSeriesTimeFieldName

			for i := range qm.columnNames {
				if i == qm.timeIndex || i == qm.metricIndex {
					continue
				}

				if t := frame.Fields[i].Type(); t == data.FieldTypeString || t == data.FieldTypeNullableString {
					continue
				}

				var err error
				if frame, err = convertSQLValueColumnToFloat(frame, i); err != nil {
					errAppendDebug("convert value to float failed", err, interpolatedQuery)
					return
				}
			}

			tsSchema := frame.TimeSeriesSchema()
			if tsSchema.Type == data.TimeSeriesTypeLong {
				var err error
				originalData := frame
				frame, err = data.LongToWide(frame, qm.FillMissing)
				if err != nil {
					errAppendDebug("failed to convert long to wide series when converting from dataframe", err, interpolatedQuery)
					return
				}

				// Before 8x, a special metric column was used to name time series. The LongToWide transforms that into a metric label on the value field.
				// But that makes series name have both the value column name AND the metric name. So here we are removing the metric label here and moving it to the
				// field name to get the same naming for the series as pre v8
				if len(originalData.Fields) == 3 {
					for _, field := range frame.Fields {
						if len(field.Labels) == 1 { // 7x only supported one label
							name, ok := field.Labels["metric"]
							if ok {
								field.Name = name
								field.Labels = nil
							}
						}
					}
				}
			}
			if qm.FillMissing != nil {
				// we align the start-time
				startUnixTime := qm.TimeRange.From.Unix() / int64(qm.Interval.Seconds()) * int64(qm.Interval.Seconds())
				alignedTimeRange := backend.TimeRange{
					From: time.Unix(startUnixTime, 0),
					To:   qm.TimeRange.To,
				}

				var err error
				frame, err = sqlutil.ResampleWideFrame(frame, qm.FillMissing, alignedTimeRange, qm.Interval)
				if err != nil {
					logger.Error("Failed to resample dataframe", "err", err)
					frame.AppendNotices(data.Notice{Text: "Failed to resample dataframe", Severity: data.NoticeSeverityWarning})
				}
			}
		}
		frames = append(frames, frame)
	}

	//queryResult.dataResponse.Frames = data.Frames{frame}
	queryResult.dataResponse.Frames = frames
	ch <- queryResult
}
