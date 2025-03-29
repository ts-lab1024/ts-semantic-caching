package jscodec

type SemanticMetaValue struct {
	SemanticMeta string                 `json:"semantic_meta"`
	SeriesArray  []*SemanticSeriesValue `json:"series_array"`
}

type SemanticSeriesValue struct {
	SeriesSegment string    `json:"series_segment"`
	Values        []*Sample `json:"values"`
}

type Sample struct {
	Timestamp int64     `json:"timestamp"`
	Value     []float64 `json:"value"`
}
