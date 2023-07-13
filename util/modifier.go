package util

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v3"
)

func yaml2json(yamlFile []byte) []byte {
	var data map[string]interface{}
	yamlErr := yaml.Unmarshal(yamlFile, &data)
	failOnError(yamlErr, "Failed unmarshal yaml")

	jsonBytes, jsonMarshalErr := json.Marshal(data)
	failOnError(jsonMarshalErr, "Failed Marshal json")

	return jsonBytes
}

func json2yaml(jsonFile []byte) []byte {
	var data map[string]interface{}
	_ = json.Unmarshal(jsonFile, &data)

	yamlData, _ := yaml.Marshal(data)
	return yamlData
}

// func ModifyWorkflow(byteCh <-chan []byte) {
// 	for body := range byteCh {
// 		var modifiedWorkflow Workflow
// 		unmarshalErr := json.Unmarshal(body, &modifiedWorkflow)
// 		failOnError(unmarshalErr, "Unmarshal")

// 		// Origin JSON file
// 		jsonData := yaml2json(yamlFile)
// 		templates := gjson.Get(string(jsonData), "spec.templates").Value().([]interface{})

// 		var modifiedContainerInfo []ContainerInfo
// 		modifiedContainerInfo = modifiedWorkflow.Containers
// 		for _, container := range templates {
// 			containerMap, _ := container.(map[string]interface{})
// 			for _, mdContainer := range modifiedContainerInfo {
// 				if containerMap["name"] == mdContainer.Name {
// 					containerMap["nodeSelector"] = mdContainer.NodeSelector
// 				}
// 			}
// 		}
// 		var tmp map[string]interface{}
// 		_ = yaml.Unmarshal(yamlFile, &tmp)
// 		tmp["spec"].(map[string]interface{})["templates"] = templates

// 		yamlData, _ := yaml.Marshal(tmp)
// 		_ = ioutil.WriteFile("./modify_result/"+modifiedWorkflow.Filename+".yaml", yamlData, 0644)
// 	}
// }
