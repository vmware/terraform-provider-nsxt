package loadbalancer

type LbServiceStatistics struct {
	// Timestamp when the data was last updated
	LastUpdateTimestamp int64 `json:"last_update_timestamp,omitempty"`

	// Statistics of load balancer pools
	Pools []LbPoolStatistics `json:"pools,omitempty"`

	// Load balancer service identifier
	ServiceId string `json:"service_id"`

	// Load balancer service statistics counter
	Statistics LbServiceStatisticsCounter `json:"statistics,omitempty"`

	// Statistics of load balancer virtual servers
	VirtualServes []LbVirtualServerStatistics `json:"virtual_servers,omitempty"`
}
