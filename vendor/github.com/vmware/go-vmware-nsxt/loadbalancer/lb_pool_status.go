package loadbalancer

type LbPoolStatus struct {

	// Timestamp when the data was last updated
	LastUpdateTimestamp int64 `json:"last_update_timestamp,omitempty"`

	// Status of load balancer pool members
	Members []LbPoolMemberStatus `json:"members,omitempty"`

	// Load balancer pool identifier
	PoolId string `json:"pool_id"`

	// Virtual server status
	Status string `json:"status"`
}
