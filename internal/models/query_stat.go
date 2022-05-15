package models

type QueryStat struct {
	Query      string  `json:"query"`
	Calls      int64   `json:"calls"`
	TotalTime  float64 `json:"total_time"`
	MeanTime   float64 `json:"mean_time"`
	Percentage float64 `json:"percentage"`
}
