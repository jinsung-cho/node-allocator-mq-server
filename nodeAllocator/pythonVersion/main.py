import json
import os
import random
import threading
import pika
from dotenv import load_dotenv

load_dotenv()
load_dotenv(dotenv_path="../../.env")
id = os.getenv("MQ_ID")
passwd = os.getenv("MQ_PASSWD")
ip = os.getenv("MQ_IP")
port = os.getenv("MQ_PORT")
resourceQueue = os.getenv("MQ_RESOURCE_QUE")

updateWorkflowAvailable = threading.Event()
publishAvailable = threading.Event()

class ContainerInfo:
    def __init__(self):
        self.name = ""
        self.image = ""
        self.limits = {}
        self.requests = {}
        self.nodeSelector = {}


class Workflow:
    def __init__(self):
        self.hash = ""
        self.containers = []


def subscribe(byteCh):
    global ip 
    global port 
    global id
    global passwd 
    global resourceQueue 
    global updateWorkflowAvailable

    credentials = pika.PlainCredentials(id, passwd)
    connection = pika.BlockingConnection(pika.ConnectionParameters(ip, port, "/", credentials))
    channel = connection.channel()

    channel.queue_declare(queue=resourceQueue, durable=False)

    def callback(ch, method, properties, body):
        byteCh.append(body)
        ch.basic_ack(delivery_tag=method.delivery_tag)
        updateWorkflowAvailable.set()

    channel.basic_consume(queue=resourceQueue, on_message_callback=callback)

    channel.start_consuming()

def allocateNode(container):
    container_v2 = ContainerInfo()
    container_v2.name = container["name"]
    container_v2.image = container["image"]
    container_v2.limits = container["limits"]
    container_v2.requests = container["requests"]
    
    clouds = ["private", "azure", "aws"]
    random_cloud_name = random.choice(clouds)
    random_node_num = str(random.randint(1, 10))
    container_v2.nodeSelector[random_cloud_name] = random_node_num
    
    return container_v2

def updateWorkflow(byteCh, byteChV2):
    global updateWorkflowAvailable
    global publishAvailable
    while True:
        updateWorkflowAvailable.wait()
        body = byteCh.pop(0)
        workflow = Workflow()
        json_data = json.loads(body.decode())
        workflow.hash = json_data["hash"]
        containers = json_data["containers"]
        for container in containers:
            container_v2 = allocateNode(container)
            workflow.containers.append(container_v2.__dict__)
        final_result = json.dumps(workflow.__dict__).encode()
        byteChV2.append(final_result)
        updateWorkflowAvailable.clear()
        publishAvailable.set()

def publish(byteChV2):
    global ip 
    global port 
    global id
    global passwd 
    global publishAvailable
    while True:
        publishAvailable.wait()
        body = byteChV2.pop(0)
        json_body = json.loads(body)
        hash = json_body["hash"]
        credentials = pika.PlainCredentials(id, passwd)
        connection = pika.BlockingConnection(pika.ConnectionParameters(ip, port, "/", credentials))
        channel = connection.channel()
        channel.queue_declare(queue=hash, durable=False)
        channel.basic_publish(
            exchange="",
            routing_key=hash,
            body=body,
            properties=pika.BasicProperties(content_type="application/json"),
        )
        connection.close()
        publishAvailable.clear()

if __name__ == "__main__":
    byteCh = []
    byteChV2 = []

    subscribeThread = threading.Thread(target=subscribe, args=(byteCh,))
    subscribeThread.start()

    updateThread = threading.Thread(target=updateWorkflow, args=(byteCh, byteChV2))
    updateThread.start()

    publishThread = threading.Thread(target=publish, args=(byteChV2,))
    publishThread.start()

    forever = threading.Event()
    forever.wait()