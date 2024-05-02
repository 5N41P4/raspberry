package data

// Data structs for the API / UI

type DiskUsage struct {
	Total   float64 `json:"total"`
	Free    float64 `json:"free"`
	Used    float64 `json:"used"`
	Percent float64 `json:"percent"`
}

type CpuUsage struct {
	Usage float64 `json:"usage"`
}

type MemUsage struct {
	Memory float64 `json:"memory"`
}
