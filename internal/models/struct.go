package models

type ContainerInfo struct {
	Name         string                 `json:"name"`
	Image        string                 `json:"image"`
	Limits       map[string]interface{} `json:"limits"`
	Requests     map[string]interface{} `json:"requests"`
	NodeSelector map[string]interface{} `json:"nodeSelector"`
}

type Workflow struct {
	Hash       string          `json:"hash"`
	Containers []ContainerInfo `json:"containers"`
}

type WorkflowInfoItem struct {
	Duration string `json:"duration"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type WorkflowInfos struct {
	Items  map[string]WorkflowInfoItem `json:"items"`
	Status string                      `json:"status"`
}
