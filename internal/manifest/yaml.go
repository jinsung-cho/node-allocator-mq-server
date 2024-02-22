package manifest

import (
	"encoding/json"

	"github.com/tidwall/gjson"

	"server/internal/models"
	"server/pkg/handler"
)

func yaml2jsonSpec(yamlFile []byte) map[string]interface{} {
	var data map[string]interface{}
	yamlErr := json.Unmarshal(yamlFile, &data)
	handler.CheckErrorAndPanic(yamlErr, "Failed unmarshal yaml")

	jsonData := make(map[string]interface{})
	jsonData["wolkflow"] = data["spec"].(map[string]interface{})["templates"]

	jsonBytes, jsonMarshalErr := json.Marshal(jsonData)
	handler.CheckErrorAndPanic(jsonMarshalErr, "Failed Marshal json")

	var templates map[string]interface{}
	jsonUnmarshalErr := json.Unmarshal(jsonBytes, &templates)
	handler.CheckErrorAndPanic(jsonUnmarshalErr, "Failed UnMarshal json")
	return templates
}

func parseResource(resource map[string]interface{}) (map[string]interface{}, map[string]interface{}) {
	cSpec, containerMarshalErr := json.Marshal(resource)
	handler.CheckErrorAndPanic(containerMarshalErr, "Failed Marshal json")

	requestSpec := gjson.Get(string(cSpec), "resources.requests").Value().(map[string]interface{})
	limitSpec := gjson.Get(string(cSpec), "resources.limits").Value().(map[string]interface{})

	return requestSpec, limitSpec
}

func parseContainerInfo(workInfo map[string]interface{}) models.ContainerInfo {
	containerInfo := models.ContainerInfo{}
	containerInfo.Name = workInfo["name"].(string)
	containerSpec := workInfo["container"].(map[string]interface{})
	containerInfo.Image = containerSpec["image"].(string)

	_, existResource := containerSpec["resources"]
	if existResource {
		request, limit := parseResource(containerSpec)
		containerInfo.Requests = request
		containerInfo.Limits = limit
	}

	return containerInfo
}

func ParseWorkflowInfo(b []byte, hash string) []byte {
	templates := yaml2jsonSpec(b)

	result := models.Workflow{}
	result.Containers = []models.ContainerInfo{}

	workflowList := templates["wolkflow"].([]interface{})
	for _, work := range workflowList {
		workInfo := work.(map[string]interface{})
		_, existContainer := workInfo["container"]

		if existContainer {
			containerInfo := parseContainerInfo(workInfo)
			result.Containers = append(result.Containers, containerInfo)
		}
	}
	result.Hash = hash

	finalResult, _ := json.Marshal(result)
	return finalResult
}

func ModifyWorkflow(origin []byte, new []byte) map[string]interface{} {
	var modifiedWorkflow models.Workflow
	_ = json.Unmarshal(new, &modifiedWorkflow)

	templates := gjson.Get(string(origin), "spec.templates").Value().([]interface{})

	var modifiedContainerInfo []models.ContainerInfo
	modifiedContainerInfo = modifiedWorkflow.Containers
	for _, container := range templates {
		containerMap, _ := container.(map[string]interface{})
		for _, mdContainer := range modifiedContainerInfo {
			if containerMap["name"] == mdContainer.Name {
				containerMap["nodeSelector"] = mdContainer.NodeSelector
			}
		}
	}
	var tmp map[string]interface{}
	_ = json.Unmarshal(origin, &tmp)
	tmp["spec"].(map[string]interface{})["templates"] = templates

	return tmp
}
