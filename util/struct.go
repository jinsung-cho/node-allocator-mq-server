package util

type ContainerInfo struct {
	Name     string                 `json:"name"`
	Image    string                 `json:"image"`
	Limits   map[string]interface{} `json:"limits"`
	Requests map[string]interface{} `json:"requests"`
}

type Workflow struct {
	Hash       string          `json:"hash"`
	Containers []ContainerInfo `json:"containers"`
}
