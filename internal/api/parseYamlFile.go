package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"server/internal/manifest"
	"server/internal/mq"
	"server/pkg/handler"
)

func ParseYamlFile(w http.ResponseWriter, r *http.Request) {
	body, requestErr := ioutil.ReadAll(r.Body)
	if handler.CheckHttpError(w, requestErr, "Check request") {
		return
	}
	hashBytes := md5.Sum(body)
	hash := hex.EncodeToString(hashBytes[:])
	workflowInfo := manifest.ParseWorkflowInfo(body, hash)

	byteCh := make(chan []byte)

	go mq.Subscribe(byteCh, hash)

	publishErr := mq.Publish(workflowInfo)
	if handler.CheckHttpError(w, publishErr, "Check Publish logic") {
		return
	}

	select {
	case res := <-byteCh:
		result := manifest.ModifyWorkflow(body, res)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	case <-time.After(3 * time.Second):
		http.Error(w, "timeout", http.StatusRequestTimeout)
		return
	}
}
