# DistributedSystems
This  my  Solution  For  the project  of  the Distributed Systems  2023-2024
My solution will be  implemented  with  golang as  the  backend  language 

#BUILD THE SYSTEM
1)  BUILD FRONTEND
2)  COPY frontend to backend
3)  Build backend
commands

build.sh
```
#!/bin/bash
cd DistributedSystems/frontend
npm i
npm run build
cd  ../..
cd  DistributedSystems/backendService
rm -r staticServer
mkdir staticServer
cp  -r  ../frontend/dist/frontend/* staticServer/
go  build
```

#RUN build.sh  to build the  server outsite of  DistributedSystems
```
chmod +x ./build.sh
./build.sh
```

Next Step 
#Add  a file .env on each node(1 coordinator and  other workers)


#coordinator Enviroments
```
serverPort=<server -port>
coordinator=true
nodeId=<node-id>
hostCoordinator= <private  - ipv4  of  coordinator   ex. 10.0.1.4  >
myNetwork=       <private  - ipv4  of  node  ex. 10.0.1.4 >
coordinatorPort=<coordinator - port >
publicUri=http://<public-ip>:<exposed-port>
rabbitMQ=amqp://<rabbit-mq-user>:<rabbit-mq-pass>@<rabbit-mq-network>:5672/
workers=<total - workers >
scaleFactorCoin= <stake-coins-of-node>
scaleFactorMsg=<stake-msg-of-node>
CapicityBlockMsg=<capacity of block in BlockChainMsg >
CapicityBlockCoin=<capacity of block in BlockChainCoin >
perNode=<how - much  - per - node >
```
#worker Enviroments 
```
serverPort=<server -port>
coordinator=false
nodeId=<node-id>
hostCoordinator= <private  - ipv4  of  coordinator   ex. 10.0.1.4  >
myNetwork=       <private  - ipv4  of  node  ex. 10.0.1.4 >
coordinatorPort=<coordinator - port >
publicUri=http://<public-ip>:<exposed-port>
rabbitMQ=amqp://<rabbit-mq-user>:<rabbit-mq-pass>@<rabbit-mq-network>:5672/
```

Next Step
After copping  the  folder  backend  and  make  a folder for eac node and set enviroments  

use  docker-compose to lunch a rabbit mq 
```
version: '3'
services:
  rabbitmq:
    image: rabbitmq:3-management
    container_name: some-rabbit
    hostname: my-rabbit
    ports:
      - 5672:5672
      - 15672:15672
    environment:
      - RABBITMQ_DEFAULT_USER=<username>
      - RABBITMQ_DEFAULT_PASS=<password>
    command: sh -c "rabbitmq-plugins enable rabbitmq_management && rabbitmq-server"
    restart: on-failure
```

Next Step
You can use  this  Python  Script  For  Creating and Purging Queues 


#createBound.py -  create  queues  and  binds them  to topic 
```
#!/usr/bin/python
import pika

def create_queue_and_bind(exchange_name, queue_name, routing_key):
    # Establish connection
    credentials = pika.PlainCredentials('<username>', '<password>')
    parameters = pika.ConnectionParameters('<host>', 5672, '/', credentials)
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()
    try:
        # Declare queue
        channel.queue_declare(queue=queue_name, durable=True)

        # Bind queue to exchange
        channel.queue_bind(exchange=exchange_name, queue=queue_name, routing_key=routing_key)

        print(f"Queue '{queue_name}' created and bound to exchange '{exchange_name}' with routing key '{routing_key}'")

    except Exception as e:
        print(f"Failed to create and bind queue: {e}")
        raise

    finally:
        # Close connection
        connection.close()

def main():
    #nodes id
    nodes = ['node-ids' , ...'] #your  nodes IDS 
    #queues  and  topics
    queues = ['transactionCoins', 'transactionMsg', 'BlockCoins', 'BlockMsg', 'SystemInfo', 'StakeCoins', 'StakeMsg']
    topics = ['TCOINS', 'TMSG', "BCOIN", 'BMSG', 'SINFO', 'STCOIN', 'STMSG']
    for node in nodes:
        for i in range(len(queues)):
            queue_name_with_node = f"{queues[i]}-{node}"
            exchange_name = f"{topics[i]}"
            try:
                create_queue_and_bind(exchange_name, queue_name_with_node, "#")
                print("Queue creation and binding successful")
            except Exception as e:
                print(f"Error: {e}")

if __name__ == "__main__":
    main()
```
#cleanRabbitMq.py -  Purge existing  messages  from queues  goo  to run  prior system  lunch 
```
#!/usr/bin/python
import pika
#purge  queue  to delete all unreaded  messages
def purge_messages(queue_name, queue_id):
    credentials = pika.PlainCredentials('<username>', '<password>')
    parameters = pika.ConnectionParameters('<host>', 5672, '/', credentials)
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()

    # Purge messages from the specified queue
    channel.queue_purge(queue=queue_name)

    print(f"Messages purged from queue '{queue_name}' for ID '{queue_id}'")
    connection.close()

def main():
    queue_names = ["transactionCoins", "SystemInfo" , "transactionMsg","StakeCoins" , "StakeMsg", "BlockCoins", "BlockMsg"]
    queue_ids = ['node-id' ,...] # nodeIds
    for queue_name in  queue_names:
        for queue_id in queue_ids:
            queue_name1 = queue_name+'-'+queue_id
            purge_messages(queue_name1, queue_id)

if __name__ == '__main__':
    main()
```

# Start  each node  
```
./backendService
```
