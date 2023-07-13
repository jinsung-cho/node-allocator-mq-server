from flask import Flask, request
import os
import argo_workflows
from pprint import pprint
from argo_workflows.api import workflow_service_api
from argo_workflows.model.io_argoproj_workflow_v1alpha1_workflow_create_request import (
    IoArgoprojWorkflowV1alpha1WorkflowCreateRequest,
)

app = Flask(__name__)

@app.route('/run', methods=['POST'])
def run_workflow():
    try:
        argowfIp = os.getenv("ARGO_WORKFLOW_IP")
        argowfPort = os.getenv("ARGO_WORKFLOW_PORT")

        configuration = argo_workflows.Configuration(host="https://"+argowfIp+":"+argowfPort)
        configuration.verify_ssl = False
        manifest = request.get_json()
        api_client = argo_workflows.ApiClient(configuration)
        api_instance = workflow_service_api.WorkflowServiceApi(api_client)
        api_response = api_instance.create_workflow(
            namespace="argo-test",
            body=IoArgoprojWorkflowV1alpha1WorkflowCreateRequest(workflow=manifest, _check_type=False),
            _check_return_type=False)
        return {"status": "success", "response": api_response.to_dict()}
    except Exception as e:
        return {"status": "failure", "error": str(e)}

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8008)