package modelsDTOs

type DockerContainerResponseDTO struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
	Uptime    string `json:"uptime"`
	HostPort  string `json:"hostPort"`
}

// type DockerStats struct {
// 	Name      string `json:"name"`
// 	IsRunning bool   `json:"isRunning"`
// 	Uptime    string `json:"uptime"`
// 	CPUStats  struct {
// 		CPUUsage struct {
// 			TotalUsage uint64 `json:"total_usage"`
// 		} `json:"cpu_usage"`
// 		SystemUsage uint64 `json:"system_cpu_usage"`
// 	} `json:"cpu_stats"`

// 	PreCPUStats struct {
// 		CPUUsage struct {
// 			TotalUsage uint64 `json:"total_usage"`
// 		} `json:"cpu_usage"`
// 		SystemUsage uint64 `json:"system_cpu_usage"`
// 	} `json:"precpu_stats"`

// 	MemoryStats struct {
// 		Usage uint64 `json:"usage"`
// 		Limit uint64 `json:"limit"`
// 	} `json:"memory_stats"`

// 	Networks map[string]struct {
// 		RxBytes uint64 `json:"rx_bytes"`
// 		TxBytes uint64 `json:"tx_bytes"`
// 	} `json:"networks"`
// }

type DockerStats struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"is_running"`
	Uptime    string `json:"uptime"`

	CPU       float64 `json:"cpu"`
	MemoryMB  float64 `json:"memory_mb"`
	NetworkRX uint64  `json:"network_rx"`
	NetworkTX uint64  `json:"network_tx"`
	HostPort  string  `json:"host_port"`
}

type DockerStatsRaw struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_usage"`
	} `json:"cpu_stats"`

	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`

	MemoryStats struct {
		Usage uint64 `json:"usage"`
		Limit uint64 `json:"limit"`
	} `json:"memory_stats"`

	Networks map[string]struct {
		RxBytes uint64 `json:"rx_bytes"`
		TxBytes uint64 `json:"tx_bytes"`
	} `json:"networks"`
}
