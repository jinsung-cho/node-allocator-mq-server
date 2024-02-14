package controller

import (
	"backend/util"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
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
	body, _ := ioutil.ReadAll(r.Body)
	env_err := godotenv.Load(".env")
	util.FailOnError(env_err, ".env Load fail")
	pythonPort := os.Getenv("PYTHON_SERVER_PORT")
	resp, err := http.Post("http://localhost:" + string(pythonPort) + "/api/v1/run", "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	var responseJson map[string]interface{}
	json.Unmarshal(responseBody, &responseJson)

	if responseJson["status"] == "success" {
		fmt.Fprintln(w, "Workflow executed successfully.")
	} else {
		fmt.Fprintln(w, "Failed to execute workflow.")
	}
}
