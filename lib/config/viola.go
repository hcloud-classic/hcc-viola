package config

type viola struct {
	StartingDelayTimeSec int64 `goconf:"viola:viola_start_delay_time_sec"`    // StartingDelayTimeSec : Delay time for start viola operations
	NodeAddRetryCount    int64 `goconf:"viola:viola_node_add_retry_count"`    // NodeAddRetryCount : Retry count for add compute nodes
	NodeAddRetryWaitSec  int64 `goconf:"viola:viola_node_add_retry_wait_sec"` // NodeAddWaitSec : Wait seconds when retrying to add compute nodes
}

// Viola : violin config structure
var Viola viola
