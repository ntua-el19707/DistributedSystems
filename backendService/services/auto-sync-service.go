package services

import "crypto/rsa"

type AutoSyncService interface {
	Service
	addWorker(key rsa.PublicKey, who string)
	GetWorkers() []rsa.PublicKey
	GetWorker(who string) rsa.PublicKey
	IExist()
	PublishTransaction()
	PublishBlock()
}

type AutoSyncServiceImpl struct {
	workers []rsa.PublicKey
	clients map[string]rsa.PublicKey
}

/*
  - addWorker - add  a worker in The List @Service AutoSyncService @Implemetation AutoSyncServiceImpl
    @Param key  rsa.PublicKey
    @Param who  string
*/
func (service *AutoSyncServiceImpl) addWorker(key rsa.PublicKey, who string) {
	//TO DO  check if  key  exist in map  if  yes are also  in  workers list
	service.workers = append(service.workers, key)
	service.clients[who] = key
}

/*
  - GetWorkers - get the  worker  List @Service AutoSyncService @Implemetation AutoSyncServiceImpl
    @Returns  []  rsa.Publickey
*/
func (service AutoSyncServiceImpl) GetWorkers() []rsa.PublicKey {
	return service.workers
}

/*
  - GetWorker - get  a worker from the map  @Service AutoSyncService @Implemetation AutoSyncServiceImpl
    @Returns  rsa.PublicKey
*/
func (service AutoSyncServiceImpl) GetWorker(who string) rsa.PublicKey {
	return service.clients[who]
}

// BootStrapAutoSync -  go function IMPORTANT
// Say  that Node  exist  synchronize  node with other Nodes
func BootStrapAutoSync(service AutoSyncService) {

	service.IExist()
	for {
		//until  the  proccess  DEATH
	}
}
