# node-allocator-mq-server

## 프로젝트 소개
![image](https://github.com/jinsung-cho/node-allocator-mq-server/assets/57334203/f6fdafc6-4ea1-411e-97b5-0d392a1f0141)



- Argo workflow 실험을 위한 최적 노드 배치와 실험을 진행할 수 있는 API를 제공하는 go 및 python으로 작성된 백엔드 서버입니다.
- 제공하는 API의 상세 내용은 다음과 같습니다.
  - 워크로드 수행을 위한 매니패스트 파일 파싱 및 최적 노드 배치 결과 반환을 위한 API 제공
  - 최적 노드 배치 결과에 따른 수정된 워크로드 매니패스트 파일 실험연계를 위한 API 제공

## Requirements
- golang version 1.17.0 or higher
- python version 3.6.0 or higher
- running argo workflow server
  
## 디렉토리 구조

프로젝트의 디렉토리 구조는 다음과 같습니다:

```
├── argoIntegrationServer
│   ├── main.py
│   └── requirements.txt
├── go.mod
├── go.sum
├── internal
│   ├── api
│   ├── manifest
│   │   └── yaml.go
│   ├── models
│   │   └── struct.go
│   └── mq
│       └── pubsub.go
├── main.go
├── nodeAllocator
│   ├── golang
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── main.go
│   └── python
│       ├── main.py
│       └── requirements.txt
├── pkg
│   └── handler
│       └── handler.go
├── rabbitmq-docker
│   └── docker-compose.yaml
├── README.md
├── .env-sample
└── start.sh
```
- **argoIntegrationServer** : argo workflow의 workflow list 수신 및 workflow 실험을 시작할 수 있도록 go 기반 API 서버와 연계하여 동작하는 python 기반 서버 코드와 Package 버전 관리 파일이 포함된 디렉토리
- **internal** : 본 프로젝트에 한정되어 동작하는 코드
    - **api** : workflow 처리 및 실행 대한 API 요청 처리 코드가 위치한 디렉토리
    - **manifest** : API 요청에 포함된 매니페스트 파일 파싱 및 수정을 위한 코드가 위치한 디렉토리
    - **models** : stuct 구조가 정의된 코드가 위치한 디렉토리
    - **mq** : mq에 pub/sub 요청을 위한 코드가 위치한 디렉토리
- **main.go** : go 기반 API 서버 실행을 위한 시작 코드
- **nodeAllocator** : nodeSelector를 추가하는 알고리즘을 시뮬레이션 할 수 있도록 동작 시키는 코드(실제로는 케이웨어에서 동작해야하는 코드)
    - **golang** : 원본 매니페스트에서 자원 요청량에 따른 nodeSelector를 추가하기 위한 golang 코드가 포함된 디렉토리
    - **python** : 원본 매니페스트에서 자원 요청량에 따른 nodeSelector를 추가하기 위한 python 코드가 포함된 디렉토리
- **rabbitmq-docker** : rabbitMQ를 컨테이너로 실행시키기 위한 dock-compose.yaml 파일이 포함된 디렉토리
- **start.sh** : rabbitmq, go 서버, python 서버, nodeAllocator의 실행을 자동화 하기 위한 스크립트 파일
- **pkg** : 본 프로젝트 및 외부 프로젝트에서 공용으로 사용되는 핸들러 함수가 포함된 디렉토리
    - **handler.go** : 에러 처리 및 에러에 대한 http response를 위한 핸들러 함수


## REST API
**POST - {SERVER_IP}:{SERVER_PORT}/yaml**
![image](https://github.com/jinsung-cho/node-allocator-mq-server/assets/57334203/c813c095-5958-47c0-9430-dbb9b9448d6f)


- body (required)
  - argo workflow json
    - example   
      ```json
            # example
            {
              "apiVersion": "argoproj.io/v1alpha1",
              "kind": "Workflow",
              "metadata": {
                "generateName": "boston-housing-pipeline-",
                "annotations": {
                  "pipelines.kubeflow.org/kfp_sdk_version": "1.8.20",
                  "pipelines.kubeflow.org/pipeline_compilation_time": "2023-05-18T12:23:37.248431",
                  "pipelines.kubeflow.org/pipeline_spec": "{\"description\": \"An example pipeline that trains and logs a regression model.\", \"name\": \"Boston Housing Pipeline\"}"
                },
                "labels": {
                  "pipelines.kubeflow.org/kfp_sdk_version": "1.8.20"
                }
              },
              "spec": {
                "entrypoint": "boston-housing-pipeline",
                "templates": [
                중략...
                  {
                    "name": "deploy-model",
                    "container": {
                      "args": ["--model", "/tmp/inputs/input-0/data"],
                      "image": "gnovack/boston_pipeline_deploy_model:latest",
                      "resources": {
                        "limits": {
                          "cpu": "2",
                          "memory": "2G"
                        },
                        "requests": {
                          "cpu": "1",
                          "memory": "1G"
                        }
                      }
                    },
                    "inputs": {
                      "artifacts": [
                        {
                          "name": "train-model-model",
                          "path": "/tmp/inputs/input-0/data"
                        }
                      ]
                    },
                    "metadata": {
                      "labels": {
                        "pipelines.kubeflow.org/kfp_sdk_version": "1.8.20",
                        "pipelines.kubeflow.org/pipeline-sdk-type": "kfp",
                        "pipelines.kubeflow.org/enable_caching": "true"
                      }
                    }
                  },
                 생략...
          ```


- response (nodeSelector 관련 내용 추가)
  - argo workflow json with nodeSelector
    - example 
      ```json
            {
                "apiVersion": "argoproj.io/v1alpha1",
                "kind": "Workflow",
                "metadata": {
                    "annotations": {
                        "pipelines.kubeflow.org/kfp_sdk_version": "1.8.20",
                        "pipelines.kubeflow.org/pipeline_compilation_time": "2023-05-18T12:23:37.248431",
                        "pipelines.kubeflow.org/pipeline_spec": "{\"description\": \"An example pipeline that trains and logs a regression model.\", \"name\": \"Boston Housing Pipeline\"}"
                    },
                    "generateName": "boston-housing-pipeline-",
                    "labels": {
                        "pipelines.kubeflow.org/kfp_sdk_version": "1.8.20"
                    }
                },
                "spec": {
                    "arguments": {
                        "parameters": []
                    },
                    "entrypoint": "boston-housing-pipeline",
                    "serviceAccountName": "pipeline-runner",
                    "templates": [
                    중략...
                        {
                            "container": {
                                "args": [
                                    "--model",
                                    "/tmp/inputs/input-0/data"
                                ],
                                "image": "gnovack/boston_pipeline_deploy_model:latest",
                                "resources": {
                                    "limits": {
                                        "cpu": "2",
                                        "memory": "2G"
                                    },
                                    "requests": {
                                        "cpu": "1",
                                        "memory": "1G"
                                    }
                                }
                            },
                            "inputs": {
                                "artifacts": [
                                    {
                                        "name": "train-model-model",
                                        "path": "/tmp/inputs/input-0/data"
                                    }
                                ]
                            },
                            "metadata": {
                                "labels": {
                                    "pipelines.kubeflow.org/enable_caching": "true",
                                    "pipelines.kubeflow.org/kfp_sdk_version": "1.8.20",
                                    "pipelines.kubeflow.org/pipeline-sdk-type": "kfp"
                                }
                            },
                            "name": "deploy-model",
                            "nodeSelector": {
                                "private": "5"
                            }
                        },
                    생략...
      ```
            
**POST - {SERVER_IP}:{SERVER_PORT}/run**
![image](https://github.com/jinsung-cho/node-allocator-mq-server/assets/57334203/f2c06c02-81ce-4e81-ad10-f0834d458f61)


- body (requiered)
  - argo workflow json with nodeSelector
- response
  - succeed - status 200
  - failed - status 500

```
## 실행방법
1. `$ git clone https://github.com/jinsung-cho/node-allocator-mq-server.git`
2. `$ mv .env-sample .env`
3. `$ vim .env`
4. .env 파일 수정
```

```
######## 초기 상태 ########
MQ_ID=rabbit
MQ_PASSWD=rabbit
MQ_IP=localhost
MQ_PORT=5672
MQ_RESOURCE_QUE=queue1
ARGO_WORKFLOW_IP=1.1.1.1
ARGO_WORKFLOW_PORT=8888
GO_SERVER_PORT=11111
PYTHON_SERVER_PORT=22222


######## 각 파라미터 설명 ########
MQ_ID=rabbitmq 계정의 ID (rabbimq/docker-compose.yaml 파일의 RABBITMQ_DEFAULT_USER 동일해야함)
MQ_PASSWD=rabbitmq 계정의 Passwd (rabbimq/docker-compose.yaml 파일의 RABBITMQ_DEFAULT_PASS 동일해야함)
MQ_IP=rabbitmq가 동작하는 PC의 IP
MQ_PORT=rabbitmq의 포트 (default: 5672, docker-compose.yaml 파일의 ports 설정중 5672 포트와 연결되어야함)
MQ_RESOURCE_QUE=go API 서버에서 nodeAllocator로 데이터를 전달할 때 사용되는 queue의 이름
ARGO_WORKFLOW_IP=실행중인 Argo workflow의 IP
ARGO_WORKFLOW_PORT=실행중인 Argo workflow의 Port
GO_SERVER_PORT=go API 서버에서 사용할 Port
PYTHON_SERVER_PORT=Argo workflow와 연결되는 python API 서버에서 사용할 Port (GO_SERVER_PORT와 중복되어선 안됨)


######## 수정 예 ########
MQ_ID=root
MQ_PASSWD=rabbit
MQ_IP=localhost
MQ_PORT=5672
MQ_RESOURCE_QUE=resource
ARGO_WORKFLOW_IP=192.0.0.1
ARGO_WORKFLOW_PORT=30000
GO_SERVER_PORT=8080
PYTHON_SERVER_PORT=8888
```
```
5. `$ ./start.sh`
```
