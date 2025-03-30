package influxdb_client

import (
	"encoding/json"
	"github.com/timescale/tsbs/InfluxDB-client/models"
	jscodec "github.com/timescale/tsbs/jsoncodec"

	"log"
)

func ResponseToJson(resp *Response, starSegment string, singleSegments []string) *jscodec.SemanticMetaValue {
	if ResponseIsEmpty(resp) {
		return nil
	}

	semSeriesValues := make([]*jscodec.SemanticSeriesValue, 0)
	for i, series := range resp.Results[0].Series {
		samples := make([]*jscodec.Sample, 0)
		for _, row := range series.Values {
			sample := &jscodec.Sample{
				Timestamp: 0,
				Value:     make([]float64, 0),
			}
			timestamp, ok := row[0].(json.Number)
			if !ok {
				log.Fatal("not a json number")
			}
			sample.Timestamp, _ = timestamp.Int64()

			for k := 1; k < len(row); k++ {
				jsval, ok := row[k].(json.Number)
				if !ok {
					log.Fatal("not a json number")
				}
				flval, _ := jsval.Float64()
				sample.Value = append(sample.Value, flval)
			}

			samples = append(samples, sample)
		}
		semSeriesValues = append(semSeriesValues, &jscodec.SemanticSeriesValue{
			SeriesSegment: singleSegments[i],
			Values:        samples,
		})
	}

	return &jscodec.SemanticMetaValue{
		SemanticMeta: starSegment,
		SeriesArray:  semSeriesValues,
	}
}

func JsonToResponse(semMeta *jscodec.SemanticMetaValue) *Response {
	modelRows := make([]models.Row, 0)
	for _, seriesArray := range semMeta.SeriesArray {
		row := &models.Row{
			Name:    "",
			Tags:    nil,
			Columns: nil,
			Values:  make([][]interface{}, 0),
			Partial: false,
		}
		for _, sample := range seriesArray.Values {
			rowValues := make([]interface{}, 0)
			rowValues = append(rowValues, sample.Timestamp)
			for j := 0; j < len(sample.Value); j++ {
				rowValues = append(rowValues, sample.Value[j])
			}
			row.Values = append(row.Values, rowValues)
		}
		modelRows = append(modelRows, *row)
	}

	return &Response{
		Results: []Result{{
			StatementId: 0,
			Series:      modelRows,
			Messages:    nil,
			Err:         "",
		},
		},
		Err: "",
	}
}
