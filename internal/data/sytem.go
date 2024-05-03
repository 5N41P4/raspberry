package data

// Data structs for the API / UI

type DiskUsage struct {
	Total   float64 `json:"total"`
	Free    float64 `json:"free"`
	Used    float64 `json:"used"`
	Percent float64 `json:"percent"`
}

type CpuUsage struct {
	AvgLoad1  float64 `json:"load1"`
	AvgLoad5  float64 `json:"load5"`
	AvgLoad15 float64 `json:"load15"`
}

type MemUsage struct {
	Memory float64 `json:"memory"`
}
