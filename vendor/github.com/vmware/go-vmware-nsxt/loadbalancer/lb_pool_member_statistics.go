package loadbalancer

type LbPoolMemberStatistics struct {
	// Pool member IP Address
	IPAddress string `json:"ip_address"`

	// Pool member port
	Port string `json:"port,omitempty"`

	// Pool member statistics counter
	Statistics LbStatisticsCounter `json:"statistics"`
}
