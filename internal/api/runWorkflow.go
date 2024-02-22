package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"server/pkg/handler"
)

func RunWorkflow(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	resp, err := http.Post("http://localhost:"+string(pythonPort)+"/api/v1/run", "application/json", bytes.NewBuffer(body))

	if handler.CheckHttpError(w, err, "Argo Integration Server[failed]") {
		return
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	var responseJson map[string]interface{}
	json.Unmarshal(responseBody, &responseJson)

	if responseJson["status"] == "success" {
		fmt.Fprintln(w, "Workflow executed successfully.")
	} else {
		fmt.Fprintln(w, "Failed to execute workflow.")
	}
}
