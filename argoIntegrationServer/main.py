import os
import time
import datetime
import argo_workflows
from flask_cors import CORS
from flask import Flask, request
from argo_workflows.api import workflow_service_api
from argo_workflows.model.io_argoproj_workflow_v1alpha1_workflow_create_request import (
    IoArgoprojWorkflowV1alpha1WorkflowCreateRequest,
)

app = Flask(__name__)
CORS(app, resources={r"/api*": {"origins": "*"}})
def load_argo_info():

    # get workflow informations
    argowfIp = os.getenv("ARGO_WORKFLOW_IP")
    argowfPort = os.getenv("ARGO_WORKFLOW_PORT")
    configuration = argo_workflows.Configuration(host="https://"+argowfIp+":"+argowfPort)
    configuration.verify_ssl = False
    
    api_client = argo_workflows.ApiClient(configuration)
    api_instance = workflow_service_api.WorkflowServiceApi(api_client)
    return api_instance

@app.route('/api/v1/info', methods=['GET'])
def get_workflow_info():
    ##### preprocessing workflow information ######
    # argo_table after preprocessing              #
    # {                                           #
    #     '0':{                                   #
    #         'name' : 'workflow_1',              #
    #         'status' : 'falied',                #
    #         'duration' : '6m 21s'               #
    #     },                                      #
    #     '1':{                                   # 
    #         ...                                 #
    #     }                                       # 
    # }                                           # 
    ###############################################
    
    api_instance= load_argo_info()
    namespace = os.getenv("ARGO_NAMESPACE")
    workflow_list = api_instance.list_workflows(namespace, _check_return_type=False).to_dict()
    
    if workflow_list['items']!=None:
        try:
            argo_table={"status":"succeeded",
                        "items":{}}
            
            for i in range(len(workflow_list['items'])):
                argo_table['items'][i]={}
                workflow_name = workflow_list['items'][i]['metadata']['name']
                workflow_status = workflow_list['items'][i]['status']['phase']
                target = workflow_list['items'][i]['status']['finishedAt']

                start_time=time.mktime(datetime.datetime.strptime(workflow_list['items'][i]['status']['startedAt'], '%Y-%m-%dT%H:%M:%SZ').timetuple()) 
                if target != None:
                    end_time = time.mktime(datetime.datetime.strptime(workflow_list['items'][i]['status']['finishedAt'], '%Y-%m-%dT%H:%M:%SZ').timetuple()) 
                else:
                    now = str(datetime.datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ'))
                    end_time = time.mktime(datetime.datetime.strptime(now,'%Y-%m-%dT%H:%M:%SZ').timetuple())
                
                split_duration = str(datetime.timedelta(seconds=end_time - start_time)).split(':')
                split_duration = [str(int(x)) for x in split_duration]
                if split_duration[0]=='0':
                    if split_duration[1]=='0':
                        workflow_duration = split_duration[2] + 's'
                    else:
                        workflow_duration = split_duration[1] + 'm ' + split_duration[2] + 's'
                else:
                    workflow_duration = split_duration[0] + 'h '+ split_duration[1] + 'm ' + split_duration[2] + 's'
                argo_table['items'][i]['name'] = workflow_name
                argo_table['items'][i]['status'] = workflow_status
                argo_table['items'][i]['duration'] = workflow_duration        
            return argo_table
        except Exception as e:
            return {"status": "failure", "error": str(e) + ' : check your argo workflow service domain.'}
    else:
        print("There is no pipeline on Argo Workflow. Check Argo workflow status.")
        return {"status": "None pipeline", "error" : 'Warning : There is no any pipeline on argo workflow.'}
        
@app.route('/api/v1/run', methods=['POST'])
def run_workflow():
    try:
        manifest = request.get_json()
        api_instance = load_argo_info()
        namespace = os.getenv("ARGO_NAMESPACE")

        api_response = api_instance.create_workflow(
            namespace=namespace,
            body=IoArgoprojWorkflowV1alpha1WorkflowCreateRequest(workflow=manifest, _check_type=False),
            _check_return_type=False)
        return {"status": "success", "response": api_response.to_dict()}
    except Exception as e:
        return {"status": "failure", "error": str(e)}
    
if __name__ == '__main__':
    port = os.getenv("PYTHON_SERVER_PORT")
    app.run(host='0.0.0.0', port=port)