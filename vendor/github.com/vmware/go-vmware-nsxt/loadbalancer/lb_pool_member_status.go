package loadbalancer

type LbPoolMemberStatus struct {

	// The healthcheck failure cause when status is DOWN
	FailureCause string `json:"failure_cause,omitempty"`

	// Pool member status
	IPAddress string `json:"ip_address"`

	LastCheckTime int64 `json:"last_check_time,omitempty"`

	LastStateChangeTime int64 `json:"last_state_change_time,omitempty"`

	// Pool member port
	Port string `json:"port,omitempty"`

	// Pool member status
	Status string `json:"status"`
}
