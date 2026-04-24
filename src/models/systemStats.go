package models

type SystemStats struct {
	CPUPercent  float64
	RAMUsedGB   float64
	RAMTotalGB  float64
	RAMPercent  float64
	DiskUsedGB  float64
	DiskTotalGB float64
	DiskPercent float64
}
