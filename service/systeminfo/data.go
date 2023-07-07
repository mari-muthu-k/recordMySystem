package systeminfo

type SystemInfo struct {
	Id                string    `json:"id"`
	HostName          string    `json:"hostName"`
	BatteryPercentage float64   `json:"batteryPercentage"`
	MemoryUsage       float64   `json:"memoryUsage"`
	Temperature       float64   `json:"temperature"`
	CpuPercentage     float64   `json:"cpuPercentage"`
}

type GetSystemInfoData struct {
	HostName          []string    `json:"hostName"`
	BatteryPercentage []float64   `json:"batteryPercentage"`
	MemoryUsage       []float64   `json:"memoryUsage"`
	Temperature       []float64   `json:"temperature"`
	CpuPercentage     []float64   `json:"cpuPercentage"`
}