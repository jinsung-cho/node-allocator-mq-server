package controller

import (
	"backend/util"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func ParseYamlFile(w http.ResponseWriter, r *http.Request) {
	body, requestErr := ioutil.ReadAll(r.Body)
	if util.CheckHttpError(w, requestErr, "Check request") {
		return
	}
	hashBytes := md5.Sum(body)
	hash := hex.EncodeToString(hashBytes[:])
	workflowInfo := util.ParseWorkflowInfo(body, hash)

	byteCh := make(chan []byte)

	go util.Subscribe(byteCh, hash)

	publishErr := util.Publish(workflowInfo)
	if util.CheckHttpError(w, publishErr, "Check Publish logic") {
		return
	}

	select {
	case res := <-byteCh:
		result := util.ModifyWorkflow(body, res)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	case <-time.After(3 * time.Second):
		http.Error(w, "timeout", http.StatusRequestTimeout)
		return
	}
}

func RunWorkflow(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RunWorkflow")
	fmt.Fprint(w, "RunWorkflow")
}
