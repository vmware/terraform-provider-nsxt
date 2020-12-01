package loadbalancer

type LbVirtualServerStatus struct {

	// Timestamp when the data was last updated.
	LastUpdateTimestamp int64 `json:"last_update_timestamp,omitempty"`

	// Virtual server status
	Status string `json:"status"`

	// load balancer virtual server identifier
	VirtualServerId string `json:"virtual_server_id"`
}
