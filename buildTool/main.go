package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type LoggerService interface {
	Construct() error
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type LoggerImpl struct {
	ServiceName string
}

const (
	ColorReset   = "\033[0m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorRed     = "\033[31m"
	ColorBoldRed = "\033[31;1m"
)

func (logger *LoggerImpl) Construct() error {
	log.Printf("%s[INFO]%s[%s] logger constructed\n", ColorGreen, ColorReset, logger.ServiceName)
	return nil
}
func (logger *LoggerImpl) Valid() error {
	return nil
}
func (logger *LoggerImpl) Info(msg string) {
	log.Printf("%s[INFO]%s[%s] %s\n", ColorGreen, ColorReset, logger.ServiceName, msg)
}

func (logger *LoggerImpl) Warn(msg string) {
	log.Printf("%s[WARN]%s[%s] %s\n", ColorYellow, ColorReset, logger.ServiceName, msg)
}

func (logger *LoggerImpl) Error(msg string) {
	log.Printf("%s[ERROR]%s[%s] %s\n", ColorRed, ColorReset, logger.ServiceName, msg)
}

func (logger *LoggerImpl) Fatal(msg string) {
	log.Fatalf("%s[FATAL]%s[%s] %s\n", ColorBoldRed, ColorReset, logger.ServiceName, msg)
}

type sshServiceConnectionService interface {
	createASession(host, user string) (*ssh.Session, error)
}
type sshServiceConnectionImpl struct {
	KeyPath string
}

const errLoadPrivateKey string = "Failed to  load privateKey Due to %s "
const errParsePrivateKey string = "Failed to  parse privateKey Due to %s "
const errFailedToConnectSSH string = "Failed to connect   to %s "
const errFaileToCreateSession string = "Failed to create session Due to %s "

func (service sshServiceConnectionImpl) createASession(host, user string) (*ssh.Client, *ssh.Session, error) {
	logger := LoggerImpl{ServiceName: "ssh-service-conection"}
	info := "creating  client and session"
	logger.Info(fmt.Sprintf("start  %s", info))
	const port string = "22"
	privateKey, err := ioutil.ReadFile(service.KeyPath)
	if err != nil {
		err = errors.New(fmt.Sprintf(errLoadPrivateKey, err.Error()))
		logger.Error(fmt.Sprintf("abbort %s due to %s", info, err.Error()))
		return nil, nil, err
	}
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		err = errors.New(fmt.Sprintf(errParsePrivateKey, err.Error()))
		logger.Error(fmt.Sprintf("abbort %s due to %s", info, err.Error()))
		return nil, nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Only for testing. Use proper verification in production.
	}
	client, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		err = errors.New(fmt.Sprintf(errFailedToConnectSSH, err.Error()))
		logger.Error(fmt.Sprintf("abbort %s due to %s", info, err.Error()))
		return nil, nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		err = errors.New(fmt.Sprintf(errFaileToCreateSession, err.Error()))
		logger.Error(fmt.Sprintf("abbort %s due to %s", info, err.Error()))
		return nil, nil, err
	}
	logger.Info(fmt.Sprintf("commit %s", info))
	return client, session, nil

}
func CreateBuildSSH(host, user string) {
	privateKeyPath := "../.ssh/id_rsa"

	sshServiceConnection := sshServiceConnectionImpl{KeyPath: privateKeyPath}
	client, session, err := sshServiceConnection.createASession(host, user)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer client.Close()
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}
	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			log.Fatalf("Failed to copy session stdout to terminal stdout: %v", err)
		}
	}()
	logger := LoggerImpl{ServiceName: "create and run  build.sh"}
	info := "creating  client and session"
	logger.Info(fmt.Sprintf("start  %s", info))

	const fileName = "build.sh"
	const context = `#!/bin/bash
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
`
	err = session.Run(fmt.Sprintf("echo '%s' > %s && chmod +x ./%s &&  ./%s", context, fileName, fileName, fileName))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Abbort  %s", info))
	}
	logger.Info(fmt.Sprintf("commit  %s", info))

}

type exportedPortAndId struct {
	Id   string
	Port int
}

func CreateCopyWorkersAndRunSSH(host, user, rabbitMqUri, publicIp, hostNestwork, hostCoordinator string, coordinatorPort int, exportedPortAndId []exportedPortAndId) {
	privateKeyPath := "../.ssh/id_rsa"

	sshServiceConnection := sshServiceConnectionImpl{KeyPath: privateKeyPath}
	client, session, err := sshServiceConnection.createASession(host, user)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer client.Close()
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}
	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			log.Fatalf("Failed to copy session stdout to terminal stdout: %v", err)
		}
	}()
	logger := LoggerImpl{ServiceName: "create and run  copyWorker.sh and makeWorkers.sh"}
	info := "creating  scripts  for copyWorker.sh  && makeWorker.sh and thne run"
	logger.Info(fmt.Sprintf("start  %s", info))

	const fileName = "copyWorker.sh"
	const context = `#!/bin/bash
# Check if all parameters are provided
if [ "$#" -ne 7 ]; then
    echo "Usage: $0 <port-server> <port-cordinator> <port-forwarded> <id_string> <ip-private-node> <ip-private-coordinator> <folder>"
    exit 1
fi

# Assign parameters to variables
port_server=$1
port_coordinator=$2
port_forwarded=$3
id_string=$4
ip_private_node=$5
ip_private_coordinator=$6
ip_public=%s
folder=$7
rabbitMq=%s
# Check if ports are numbers
if ! [[ "$port_server" =~ ^[0-9]+$ ]]; then
    echo "Error: Server Port must be a number."
    exit 1
fi

if ! [[ "$port_coordinator" =~ ^[0-9]+$ ]]; then
    echo "Error: Coordinator Port must be a number."
    exit 1
fi

if ! [[ "$port_forwarded" =~ ^[0-9]+$ ]]; then
    echo "Error: Forwarded Port must be a number."
    exit 1
fi

# Remove existing folder if it exists
if [ -d "$folder" ]; then
    echo "Removing existing folder: $folder"
    rm -r "$folder"
fi

# Copy folder
echo "Copying folder to $folder"
cp -r DistributedSystems/backendService "$folder"
cat <<EOF > "$folder/.env"
serverPort=$port_server
coordinator=false
nodeId=$id_string
hostCoordinator=$ip_private_coordinator
myNetwork=$ip_private_node
coordinatorPort=$port_coordinator
publicUri=http://$ip_public:$port_forwarded
rabbitMQ=$rabbitMq
EOF
`
	fileNameMakeWorkers := "makeWorkers.sh"
	contextMakeWorkers := "#!/bin/bash"
	for i, portAndId := range exportedPortAndId {
		contextMakeWorkers = fmt.Sprintf("%s\n  ./%s %d %d %d %s %s %s %s", contextMakeWorkers, fileName, coordinatorPort+i, coordinatorPort, portAndId.Port, portAndId.Id, hostNestwork, hostCoordinator, fmt.Sprintf("worker-%d", i+1))
	}
	contentWorker := fmt.Sprintf(context, publicIp, rabbitMqUri)
	err = session.Run(fmt.Sprintf("echo '%s' > %s && chmod +x ./%s && echo '%s' > %s && chmod +x ./%s && ./%s", contentWorker, fileName, fileName, contextMakeWorkers, fileNameMakeWorkers, fileNameMakeWorkers, fileNameMakeWorkers))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Abbort  %s", info))
	}
	logger.Info(fmt.Sprintf("commit  %s", info))

}
func generateID(size int) string {
	logger := LoggerImpl{ServiceName: "id-generator"}
	info := fmt.Sprintf("generating  id of size %d", size)
	logger.Info(fmt.Sprintf("start  %s", info))
	id := make([]byte, size)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	const allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range id {
		id[i] = allChars[seededRand.Intn(len(allChars))]
	}
	logger.Info(fmt.Sprintf("commit  %s", info))
	strId := string(id)
	logger.Info(fmt.Sprintf("unique  string id: %s  created", strId))
	return strId
}
func generateExportPortAndIdPair(coordinatorPort, total int) ([]exportedPortAndId, []string) {
	logger := LoggerImpl{ServiceName: "id&port-generator"}
	info := fmt.Sprintf("generating  id and  port  pair  total %d", total)
	logger.Info(fmt.Sprintf("start  %s", info))
	idSize := 8
	coordinatorId := generateID(idSize)
	var ids []string
	list := []exportedPortAndId{exportedPortAndId{Port: coordinatorPort, Id: coordinatorId}}
	for i := 1; i < total; i++ {
		row := exportedPortAndId{Id: generateID(idSize), Port: coordinatorPort + i}
		list = append(list, row)
		ids = append(ids, row.Id)
	}
	logger.Info(fmt.Sprintf("commit  %s", info))
	return list, ids
}

type nodeGenEnvParam struct {
	HostNetwork        string
	CoordinatorNetwork string
	ExportedPortsAndId []exportedPortAndId
}
type nodeGenEnvParamList []nodeGenEnvParam

func slicer[T any](perSlice int, list []T) [][]T {
	var listRsp [][]T
	start := 0
	end := len(list)
	for start < end {
		endItem := start + perSlice
		if endItem > end {
			endItem = end
		}
		listRsp = append(listRsp, list[start:endItem])
		start = start + perSlice
	}
	return listRsp
}
func shellRunner(shellCommands, name string, host HostAndUser) {
	logger := LoggerImpl{ServiceName: "name"}
	info := fmt.Sprintf("start sending  commands %s", shellCommands)
	logger.Info(fmt.Sprintf("start  %s", info))
	privateKeyPath := "../.ssh/id_rsa"

	sshServiceConnection := sshServiceConnectionImpl{KeyPath: privateKeyPath}
	client, session, err := sshServiceConnection.createASession(host.Host, host.User)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer client.Close()
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}
	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			log.Fatalf("Failed to copy session stdout to terminal stdout: %v", err)
		}
	}()
	err = session.Run(shellCommands)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Abbort  %s", info))
	}
	logger.Info(fmt.Sprintf("commit  %s", info))
}
func (list *nodeGenEnvParamList) Gen(host, user, pass string, networks []string, coordinatorNet, publicIp string, port, perNetwork int) string {
	logger := LoggerImpl{ServiceName: "id&port-generator"}
	info := fmt.Sprintf("generating  node gen env  param list %d ", len(networks)+1)
	logger.Info(fmt.Sprintf("start  %s", info))
	allNetworks := append([]string{coordinatorNet}, networks...)
	total := len(allNetworks) * perNetwork
	listPortAndId, ids := generateExportPortAndIdPair(port, total)
	*list = make([]nodeGenEnvParam, 0)
	slices := slicer(perNetwork, listPortAndId)
	if len(slices) != len(allNetworks) {
		logger.Fatal("slices is not eual to networks")
	}
	for i, network := range allNetworks {
		row := nodeGenEnvParam{HostNetwork: network, CoordinatorNetwork: coordinatorNet, ExportedPortsAndId: slices[i]}
		*list = append(*list, row)
	}
	logger.Info(fmt.Sprintf("commit  %s", info))
	natContent := `!/bin/bash/
echo "Enabling ipv4 forwarding (cleaning old rules)"

# flushing old rules -- USE WITH CARE
iptables --flush
iptables --table nat --flush
iptables --delete-chain
iptables --table nat --delete-chain

#grap networks  from /etc/hosts
node0Public=$(cat /etc/hosts | grep publicIpNode-0 | awk '{print $1}')`
	for _, network := range allNetworks {

		natContent = fmt.Sprintf("%s\n%s=$(cat /etc/hosts | grep %s | awk '{print $1}')", natContent, network, network)
	}
	for _, params := range *list {
		if params.HostNetwork != params.CoordinatorNetwork {
			natContent = fmt.Sprintf("%s\n\n#Forward Ports  Of Network %s\nnode=$%s", natContent, params.HostNetwork, params.HostNetwork)
			for i, portAndId := range params.ExportedPortsAndId {
				subPart := `# Forward %s:%d -> $node0Public:%d and  $node-0:%d
port=%d
portServer=%d
iptables -t nat -A PREROUTING -d $node0Public -p tcp --dport $port -j DNAT --to-destination $node:$portServer
iptables -t nat -A PREROUTING -d $node0 -p tcp --dport $port -j DNAT --to-destination $node:$portServer
iptables -t nat -A POSTROUTING -d $node -p tcp --dport $portServer -j SNAT --to-source $node0
iptables -t nat -A OUTPUT -p tcp --dport $port -d $node0 -j DNAT --to-destination $node:$portServer
echo "net&Port $node:$portServer  IS  Now forwarted to  $node0:$port  &&  $node0Public:$port"`
				subPartFill := fmt.Sprintf(subPart, params.HostNetwork, i+port, portAndId.Port, portAndId.Port, portAndId.Port, port+i)
				natContent = fmt.Sprintf("%s\n\n%s", natContent, subPartFill)
			}
		}
	}
	natContent = fmt.Sprintf("%s\n\necho 1 > /proc/sys/net/ipv4/ip_forward", natContent)
	createRabbitMqSripts(host, user, pass, ids)
	return natContent
}

type HostAndUser struct {
	Host string
	User string
}

func buildNodes(publicIp, host, user, pass string, perNode int) {

	var node_0 HostAndUser = HostAndUser{Host: "node0", User: "user"}
	var node_1 HostAndUser = HostAndUser{Host: "node1", User: "user"}
	var node_2 HostAndUser = HostAndUser{Host: "node2", User: "user"}
	var node_3 HostAndUser = HostAndUser{Host: "node3", User: "ubuntu"}
	var node_4 HostAndUser = HostAndUser{Host: "node4", User: "ubuntu"}

	networksNodes := []HostAndUser{node_1, node_2, node_3, node_4}
	allNodes := []HostAndUser{node_0, node_1, node_2, node_3, node_4}
	networks := make([]string, len(networksNodes))
	for i := 0; i < len(networksNodes); i++ {
		networks[i] = networksNodes[i].Host
	}
	coordinatorPort := 8000
	coordinatorNet := "node0"
	var paramList nodeGenEnvParamList
	natcont := paramList.Gen(host, user, pass, networks, coordinatorNet, publicIp, coordinatorPort, perNode)
	fileName := "nat.sh"

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(natcont)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	wait := make(chan int, len(networksNodes))
	exec := func(hostPair HostAndUser) {

		var nodeEnv nodeGenEnvParam
		for _, p := range paramList {
			if p.HostNetwork == hostPair.Host {
				nodeEnv = p
			}
		}
		var rabbitMq = fmt.Sprintf("amqp://%s:%s@%s:5672/", user, pass, host)
		if nodeEnv.HostNetwork == hostPair.Host {
			CreateCopyWorkersAndRunSSH(hostPair.Host, hostPair.User, rabbitMq, publicIp, nodeEnv.HostNetwork, nodeEnv.CoordinatorNetwork, coordinatorPort, nodeEnv.ExportedPortsAndId)
		}
		wait <- 1
	}
	for _, pair := range allNodes {
		go exec(pair)
	}
	for i := 0; i < len(networksNodes); i++ {
		<-wait
	}
	//	close(wait)
	total := len(allNodes) * perNode
	setCoordinator(total)

}
func setCoordinator(totalWorkers int) {
	shell := `cd worker-1 && 
sed -i 's/coordinator=false/coordinator=true/' .env && 
echo "workers=%d" >> .env &&
echo "scaleFactorCoin=10" >> .env &&
echo "scaleFactorMsg=10" >> .env &&
echo "CapicityBlockMsg=5" >> .env &&
echo "CapicityBlockCoin=6" >> .env &&
echo "perNode=1000" >> .env && cd .. &&  \
[ -d "coordinator" ] && rm -r coordinator; \&& mv worker-1/ coordinator`
	shellFill := fmt.Sprintf(shell, totalWorkers)
	var node_0 HostAndUser = HostAndUser{Host: "node0", User: "user"}
	shellRunner(shellFill, "set coordinator env", node_0)
}
func PullFromRepo(machines []HostAndUser) {
	logger := LoggerImpl{ServiceName: "pull changes from repo"}
	info := fmt.Sprintf("pull changes  for repo  %d machines ", len(machines))
	logger.Info(fmt.Sprintf("start  %s", info))
	wait := make(chan int, len(machines))
	executioner := func(pair HostAndUser) {
		logger := LoggerImpl{ServiceName: "pull changes from repo  for " + pair.Host}
		info := fmt.Sprintf("pull changes  for repo ")
		logger.Info(fmt.Sprintf("start  %s", info))
		privateKeyPath := "../.ssh/id_rsa"

		sshServiceConnection := sshServiceConnectionImpl{KeyPath: privateKeyPath}
		client, session, err := sshServiceConnection.createASession(pair.Host, pair.User)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer client.Close()
		defer session.Close()

		stdout, err := session.StdoutPipe()
		if err != nil {
			log.Fatalf("Failed to create stdout pipe: %v", err)
		}
		go func() {
			if _, err := io.Copy(os.Stdout, stdout); err != nil {
				log.Fatalf("Failed to copy session stdout to terminal stdout: %v", err)
			}

		}()
		err = session.Run(fmt.Sprintf(`cd DistributedSystems/ && GIT_SSH_COMMAND="ssh -v" git pull`))
		if err != nil {
			//logger.Fatal(fmt.Sprintf("Abbort  %s", info))
		}
		logger.Info(fmt.Sprintf("commit %s", info))

	}
	exec := func(pair HostAndUser) {
		executioner(pair)
		wait <- 1
	}
	logger.Info("Deploy  sessions in parrallel")
	for i := 0; i < len(machines); i++ {
		go exec(machines[i])
	}
	for i := 0; i < len(machines); i++ {
		logger.Info("wait  to collect   session in parrallel")
		<-wait
		logger.Info("collect   session in parrallel")
	}
	close(wait)
	logger.Info(fmt.Sprintf("commit  %s", info))

}
func CheckOut(machines []HostAndUser) {
	logger := LoggerImpl{ServiceName: "repo checkout  branch"}
	info := fmt.Sprintf("checkBranch   for repo  %d machines ", len(machines))
	logger.Info(fmt.Sprintf("start  %s", info))
	wait := make(chan int, len(machines))
	executioner := func(pair HostAndUser) {
		logger := LoggerImpl{ServiceName: "pull changes from repo  for " + pair.Host}
		info := fmt.Sprintf("checkout  for repo ")
		logger.Info(fmt.Sprintf("start  %s", info))
		privateKeyPath := "../.ssh/id_rsa"

		sshServiceConnection := sshServiceConnectionImpl{KeyPath: privateKeyPath}
		client, session, err := sshServiceConnection.createASession(pair.Host, pair.User)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer client.Close()
		defer session.Close()

		stdout, err := session.StdoutPipe()
		if err != nil {
			log.Fatalf("Failed to create stdout pipe: %v", err)
		}
		go func() {
			if _, err := io.Copy(os.Stdout, stdout); err != nil {
				log.Fatalf("Failed to copy session stdout to terminal stdout: %v", err)
			}

		}()
		err = session.Run(fmt.Sprintf(`cd DistributedSystems/ && git checkout changeStake`))
		if err != nil {
			//logger.Fatal(fmt.Sprintf("Abbort  %s", info))
		}
		logger.Info(fmt.Sprintf("commit %s", info))
	}
	exec := func(pair HostAndUser) {
		executioner(pair)
		wait <- 1
	}
	logger.Info("Deploy  sessions in parrallel")
	for i := 0; i < len(machines); i++ {
		go exec(machines[i])
	}
	for i := 0; i < len(machines); i++ {
		logger.Info("wait  to collect   session in parrallel")
		<-wait
		logger.Info("collect   session in parrallel")
	}
	close(wait)
	logger.Info(fmt.Sprintf("commit  %s", info))

}
func createRabbitMqSripts(host, user, pass string, id []string) {
	logger := LoggerImpl{ServiceName: "create python rabbit scripts"}
	info := fmt.Sprintf("cretae sripts for create  & purge  ")
	logger.Info(fmt.Sprintf("start  %s", info))
	dump := func(list []string) string {
		var strRsp string
		for _, s := range list {
			if strRsp != "" {
				strRsp = fmt.Sprintf("%s , ", strRsp)
			}
			strRsp = fmt.Sprintf("%s'%s'", strRsp, s)
		}
		return strRsp
	}
	pythonCreateBindTemplate := `#!/usr/bin/python
import pika

def create_queue_and_bind(exchange_name, queue_name, routing_key):
    # Establish connection
    credentials = pika.PlainCredentials('%s', '%s')
    parameters = pika.ConnectionParameters('%s', 5672, '/', credentials)
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
    nodes = [%s]
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
`
	pythonPurgeTempate := `#!/usr/bin/python
import pika
#purge  queue  to delete all unreaded  messages 
def purge_messages(queue_name, queue_id):
    credentials = pika.PlainCredentials('%s', '%s')
    parameters = pika.ConnectionParameters('%s', 5672, '/', credentials)
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()

    # Purge messages from the specified queue
    channel.queue_purge(queue=queue_name)

    print(f"Messages purged from queue '{queue_name}' for ID '{queue_id}'")
    connection.close()

def main():
    queue_names = ["transactionCoins", "SystemInfo" , "transactionMsg","StakeCoins" , "StakeMsg", "BlockCoins", "BlockMsg"]
    queue_ids = [%s] # nodeIds
    for queue_name in  queue_names:
        for queue_id in queue_ids:
            queue_name1 = queue_name+'-'+queue_id
            purge_messages(queue_name1, queue_id)

if __name__ == '__main__':
    main()
`
	createScript := fmt.Sprintf(pythonCreateBindTemplate, user, pass, host, dump(id))
	purgeScript := fmt.Sprintf(pythonPurgeTempate, user, pass, host, dump(id))
	logger.Info(fmt.Sprintf("commit  %s", info))
	const fileNameCreate = "createBound.py"
	const fileNamePurge = "rabbitMqClean.py"
	info = fmt.Sprintf("creating  files for  %s , %s ", fileNameCreate, fileNamePurge)
	logger.Info(fmt.Sprintf("start %s", info))
	FileWriter(createScript, fileNameCreate)
	FileWriter(purgeScript, fileNamePurge)
	logger.Info(fmt.Sprintf("commit %s", info))
	logger.Warn(fmt.Sprintf("you  should create  ./%s  when  you  first build  and evry time you start the servers  you should ./%s ", fileNameCreate, fileNamePurge))
}
func FileWriter(context, filename string) {
	logger := LoggerImpl{ServiceName: "create file" + filename}
	info := fmt.Sprintf("creatint file")
	logger.Info(fmt.Sprintf("start  %s", info))
	file, err := os.Create(filename)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating file: %s", err.Error()))
	}
	defer file.Close()
	_, err = file.WriteString(context)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error writing to file:%s", err.Error()))
	}
	logger.Info(fmt.Sprintf("commit  %s", info))

}

type HostAndCommands struct {
	Host  HostAndUser
	Shell string
}

func Run(commands []HostAndCommands) {
	node := HostAndUser{Host: "node0", User: "user"}
	shellRunner("cd  buildTool/ && chmod +x ./rabbitMqClean.py  && ./rabbitMqClean.py", "clean-up", node)
	type ChannelRsp struct {
		Host  string
		Shell string
	}
	wait := make(chan ChannelRsp, len(commands))
	exec := func(command HostAndCommands) {
		shellRunner(command.Shell, "runner  app", command.Host)
		wait <- ChannelRsp{Host: command.Host.Host, Shell: command.Shell}
	}
	for _, command := range commands {
		go exec(command)
	}
	for i := 0; i < len(commands); i++ {
		<-wait
	}

}
func main() {
	// Define SSH server configurations
	//	host := "10.0.1.6"
	//	user := "ubuntu"
	//	ids := []string{"aa", "ab", "ac"}

	logger := LoggerImpl{ServiceName: "main"}
	args := os.Args
	if len(args) != 2 {
		logger.Fatal("-usage  ./name <mode> | help")
	}
	var node_1 HostAndUser = HostAndUser{Host: "node1", User: "user"}
	var node_2 HostAndUser = HostAndUser{Host: "node2", User: "user"}
	var node_3 HostAndUser = HostAndUser{Host: "node3", User: "ubuntu"}
	var node_4 HostAndUser = HostAndUser{Host: "node4", User: "ubuntu"}
	var node_0 HostAndUser = HostAndUser{Host: "node0", User: "user"}
	networksNodes := []HostAndUser{node_0, node_1, node_2, node_3, node_4}
	switch strings.ToLower(args[1]) {
	case "build":
		wait := make(chan int, len(networksNodes))
		exec := func(hostPair HostAndUser) {
			CreateBuildSSH(hostPair.Host, hostPair.User)
			wait <- 1
		}
		for _, n := range networksNodes {
			go exec(n)
		}
		for i := 0; i < len(networksNodes); i++ {
			<-wait
		}
		close(wait)

	case "copy":
		buildNodes("public-ip", "rabbitmqhost", "rabbitUser", "rabbitPass", 2)
	case "pull":
		PullFromRepo(networksNodes)
	case "checkout":
		CheckOut(networksNodes)
	case "run-5":
		comandWorker1 := "cd worker-1/ && ./backendService"
		//		comandCoordinator := "cd coordinator/ && ./backendService"
		commands := []HostAndCommands{
			//		HostAndCommands{Host: node_0, Shell: comandCoordinator},
			HostAndCommands{Host: node_1, Shell: comandWorker1},
			HostAndCommands{Host: node_2, Shell: comandWorker1},
			HostAndCommands{Host: node_3, Shell: comandWorker1},
			HostAndCommands{Host: node_4, Shell: comandWorker1},
		}
		Run(commands)
	case "run-10":
		comandWorker1 := "cd worker-1/ && ./backendService"
		comandWorker2 := "cd worker-2/ && ./backendService"
		comandCoordinator := "cd coordinator/ && ./backendService"
		commands := []HostAndCommands{
			HostAndCommands{Host: node_0, Shell: comandCoordinator},
			HostAndCommands{Host: node_1, Shell: comandWorker1},
			HostAndCommands{Host: node_2, Shell: comandWorker1},
			HostAndCommands{Host: node_3, Shell: comandWorker1},
			HostAndCommands{Host: node_4, Shell: comandWorker1},
			HostAndCommands{Host: node_0, Shell: comandWorker2},
			HostAndCommands{Host: node_1, Shell: comandWorker2},
			HostAndCommands{Host: node_2, Shell: comandWorker2},
			HostAndCommands{Host: node_3, Shell: comandWorker2},
			HostAndCommands{Host: node_4, Shell: comandWorker2},
		}
		Run(commands)

	case "help":
		modes := []string{"copy", "build", "pull", "checkout", "help"}
		for _, m := range modes {
			logger.Info(fmt.Sprintf("mode %s", m))
		}
	default:
		logger.Fatal("-usage  ./name <mode> | help")
	}
}
