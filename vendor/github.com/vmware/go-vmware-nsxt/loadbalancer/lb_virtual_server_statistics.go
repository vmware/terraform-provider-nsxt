package loadbalancer

type LbVirtualServerStatistics struct {
	// Timestamp when the data was last updated
	LastUpdateTimestamp int64 `json:"last_update_timestamp,omitempty"`

	// Load balancer virtual server identifier
	VirtualServerId string `json:"virtual_server_id"`

	// Virtual server statistics counter
	Statistics LbStatisticsCounter `json:"statistics"`
}
