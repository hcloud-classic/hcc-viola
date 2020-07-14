package model

type NormalAction struct {
}

type HccAction struct {
	ActionArea  string `json:"action_area"`
	ActionClass string `json:"action_class"`
	ActionScope string `json:"action_scope"`
	HccIPRange  string `json:"iprange"`
	ServerUUID  string `json:"server_uuid"`
}

type Action struct {
	ActionType   string       `json:"action_type"`
	NormalType   NormalAction `json:"normal_action"`
	HccType      HccAction    `json:"hcc_action"`
	ActionResult string       `json:"action_result"`
}

type Control struct {
	Control   Action `json:"action"`
	Publisher string `json:"publisher"`
	Receiver  string `json:"receiver"`
}

type Controls struct {
	Controls Control `json:"control"`
}
