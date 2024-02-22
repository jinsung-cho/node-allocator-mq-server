package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"server/internal/models"
	"server/pkg/handler"
)

func GetWorkflowInfos(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:" + string(pythonPort) + "/api/v1/info")

	if handler.CheckHttpError(w, err, "Argo Integration Server[failed]") {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var workflowResp models.WorkflowInfos
	json.Unmarshal(body, &workflowResp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflowResp)
}
