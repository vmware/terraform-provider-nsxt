package loadbalancer

type LbStatisticsCounter struct {
	// Number of bytes in
	BytesIn int64 `json:"bytes_in"`

	// The average number of inbound bytes per second, the number is averaged over the last 5 one-second intervals.
	BytesInRate float64 `json:"bytes_in_rate,omitempty"`

	// Number of bytes out
	BytesOut int64 `json:"bytes_out"`

	// The average number of outbound bytes per second, the number is averaged over the last 5 one-second intervals.
	BytesOutRate float64 `json:"bytes_out_rate,omitempty"`

	// Number of current sessions
	CurrentSessions int64 `json:"current_sessions"`

	// The average number of current sessions per second, the number is averaged over the last 5 one-second intervals.
	CurrentSessionRate float64 `json:"current_session_rate,omitempty"`

	// Number of http requests
	HttpRequests int64 `json:"http_requests"`

	// The average number of http requests per second, the number is averaged over the last 5 one-second intervals.
	HttpRequestRate float64 `json:"http_request_rate,omitempty"`

	// Number of maximum sessions
	MaxSessions int64 `json:"max_sessions"`

	// Number of packets in
	PacketsIn int64 `json:"packets_in,omitempty"`

	// The average number of inbound packets per second, the number is averaged over the last 5 one-second intervals.
	PacketsInRate float64 `json:"packets_in_rate,omitempty"`

	// Number of packets out
	PacketsOut int64 `json:"packets_out,omitempty"`

	// The average number of outbound packets per second, the number is averaged over the last 5 one-second intervals.
	PacketsOutRate float64 `json:"packets_out_rate,omitempty"`

	// Number of source IP persistence entries
	SourceIPPersistenceEntrySize int64 `json:"source_ip_persistence_entry_size,omitempty"`

	// Number of total sessions
	TotalSessions int64 `json:"total_sessions"`
}
