package loadbalancer

type LbPoolStatistics struct {
	// Timestamp when the data was last updated
	LastUpdateTimestamp int64 `json:"last_update_timestamp,omitempty"`

	// Statistics of load balancer pool members
	Members []LbPoolMemberStatistics `json:"members,omitempty"`

	// Load balancer pool identifier
	PoolId string `json:"pool_id"`

	// Load balancer statistics counter
	Statistics LbStatisticsCounter `json:"statistics"`
}
