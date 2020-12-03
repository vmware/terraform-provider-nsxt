package loadbalancer

type LbServiceStatus struct {

	// Ids of load balancer service related active transport nodes
	ActiveTransportNodes []string `json:"active_transport_nodes,omitempty"`

	// Cpu usage in percentage
	CPUUsage int64 `json:"cpu_usage"`

	// Cpu usage in percentage
	ErrorMessage string `json:"error_message,omitempty"`

	// Timestamp when the data was last updated
	LastUpdateTimestamp int64 `json:"last_update_timestamp,omitempty"`

	// Memory usage in percentage
	MemoryUsage int64 `json:"memory_usage"`

	// status of load balancer pools
	Pools []LbPoolStatus `json:"pools,omitempty"`

	// Load balancer service identifier
	ServiceId string `json:"service_id"`

	// Status of load balancer service
	ServiceStatus string `json:"service_status"`

	// Ids of load balancer service related standby transport nodes
	StandbyTransportNodes []string `json:"standby_transport_nodes,omitempty"`

	// status of load balancer virtual servers
	VirtualServers []LbVirtualServerStatus `json:"virtual_servers,omitempty"`
}
