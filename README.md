# DistributedSystems
This  my  Solution  For  the project  of  the Distributed Systems  2023-2024
My solution will be  implemented  with  golang as  the  backend  language 

#BUILD THE SYSTEM
1)  BUILD FRONTEND
2)  COPY frontend to backend
3)  Build the backend
commands

build.sh
#!/bin/bash
echo  "start build  frontend"
cd DistributedSystems/frontend
npm i
npm run build
echo  "commit build  frontend"
cd  ~
cd  DistributedSystems/backendService
echo  "starrt  copy  frontend"
rm -r staticServer
mkdir staticServer
cp  -r  ../frontend/dist/frontend/* staticServer/
echo  "commit  copy  frontend"
echo  "start  building "
export PATH="$PATH:/usr/local/go/bin"
go  build
echo  "commit  building "

RUN build.sh  to build the  server 
Next Step 
#Add  a file .env on each node(1 coordinator and  other workers)

#coorinators 
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

#worker 
serverPort=<server -port>
coordinator=false
nodeId=<node-id>
hostCoordinator= <private  - ipv4  of  coordinator   ex. 10.0.1.4  >
myNetwork=       <private  - ipv4  of  node  ex. 10.0.1.4 >
coordinatorPort=<coordinator - port >
publicUri=http://<public-ip>:<exposed-port>
rabbitMQ=amqp://<rabbit-mq-user>:<rabbit-mq-pass>@<rabbit-mq-network>:5672/

Next Step
After copping  the folders  and add .env  fo each node then 
run c

   
