package SystemInfo

import (
	"Logger"
	"RabbitMqService"
	"Service"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"entitys"
	"errors"
	"fmt"
	"log"
	"sync"
)

// System Info  Service
// >> 1  - Coordinator
//
//	AddWorker  ,  BroadCastClients
//
// >> 2 -  workers
//
//	ConsumeClientInfo
type SystemInfoService interface {
	Service.Service
	AddWorker(body entitys.ClientRequestBody) error
	IsFull() bool
	BroadCast(sFm, sFc float64) error
	Consume() (error, float64, float64)
	IsOk() bool
	GetWorkers() []rsa.PublicKey
	NodeDetails(key rsa.PublicKey) (entitys.ClientInfo, int)
	Who(index int) (rsa.PublicKey, error)
	ClientList(key rsa.PublicKey) ([]entitys.ClientInfo, int)
}
type SystemInfoProviders struct {
	LoggerService   Logger.LoggerService
	RabbitMqService RabbitMqService.RabbitMqService
	vld             bool
}

const serviceName string = "System-Info-Service"

// -- Error Templates --
const ErrNoProvider string = "There  is no '%s' provider for  the service"
const ErrNoConstructed string = "The   provider is not constructed"
const ErrWorkerListFull string = "The   Workers  slice  next index is  %d  but size %d "

func (s *SystemInfoImpl) IsFull() bool {
	var rsp bool
	s.mu.Lock()
	defer s.mu.Unlock()
	rsp = s.index == s.ExpectedWorkers
	return rsp
}

/*
*

	Construct - @SystemInfoProviders
*/
func (p *SystemInfoProviders) Construct() error {
	if p.LoggerService == nil {
		p.LoggerService = &Logger.Logger{ServiceName: serviceName}
		err := p.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
*

	valid  - @SystemInfoProviders
*/
func (p *SystemInfoProviders) valid() error {
	if !p.vld {
		return errors.New(ErrNoConstructed)

	}
	if p.LoggerService == nil {
		errMsg := fmt.Sprintf(ErrNoProvider, "Logger Service")
		p.vld = false
		return errors.New(errMsg)

	}
	if p.RabbitMqService == nil {
		errMsg := fmt.Sprintf(ErrNoProvider, "RabbitMqService")
		p.vld = false
		return errors.New(errMsg)

	}
	return nil
}

type SystemInfoImpl struct {
	Workers         []rsa.PublicKey
	Clients         map[string]entitys.ClientInfo
	Providers       SystemInfoProviders
	index           int
	mu              sync.Mutex
	ExpectedWorkers int
	Coordinator     bool
	ok              bool
}

/*
*

	Construct - @SystemInfoImpl  Contract  @SystemInfoService
*/
func (s *SystemInfoImpl) Construct() error {
	if s.Coordinator {
		if s.ExpectedWorkers <= 0 {
			errMsg := "Must atlest have one worker"
			return errors.New(errMsg)
		}
	}
	s.Workers = make([]rsa.PublicKey, s.ExpectedWorkers)
	s.index = 0
	s.Clients = make(map[string]entitys.ClientInfo)

	err := s.Providers.Construct()
	if err != nil {
		return err
	}
	s.Providers.vld = true
	return s.Providers.valid()

}

/*
*

	AddWorker - @SystemInfoImpl  AddWorker @SystemInfoService
	@Param  body ClientRequestBody
*/
func (s *SystemInfoImpl) AddWorker(body entitys.ClientRequestBody) error {
	//Get Services
	providers := &s.Providers
	err := providers.valid()
	if err != nil {
		return err
	}
	if s.IsFull() {
		errMsg := "Workers list  is  already full"
		return errors.New(errMsg)
	}
	//Get LoggerService
	logger := providers.LoggerService
	logger.Log(fmt.Sprintf("Start  adding node  with id:%s", body.Client.Id))
	s.mu.Lock()
	logger.Log(fmt.Sprintf("Lock  adding node  with id:%s", body.Client.Id))
	Unlock := func() {
		logger.Log(fmt.Sprintf("Unlock  adding node  with id:%s", body.Client.Id))
		s.mu.Unlock()
	}
	defer Unlock()
	if s.index >= len(s.Workers) {
		errMsg := fmt.Sprintf(ErrWorkerListFull, s.index, len(s.Workers))
		return errors.New(errMsg)
	}
	s.Workers[s.index] = body.PublicKey
	nodeDetails := entitys.ClientInfo{Id: body.Client.Id, IndexId: s.index, Uri: body.Client.Uri}
	publicKey := s.Workers[s.index]
	key := getHash(publicKey)

	s.Clients[key] = nodeDetails
	s.index++
	logger.Log(fmt.Sprintf("Commit  adding node  with id:%s", body.Client.Id))
	return nil
}
func (s *SystemInfoImpl) BroadCast(sFm, sFc float64) error {
	if !s.Coordinator {
		errMsg := "Only  Coordinator can  Broadcast"
		return errors.New(errMsg)

	}
	providers := &s.Providers

	err := providers.valid()
	if err != nil {
		return err
	}
	logger := providers.LoggerService
	lagoudaki := providers.RabbitMqService
	if !s.IsFull() {
		errMsg := "Workers list  is not full still  waiting "
		return errors.New(errMsg)
	}
	logger.Log("Start  BroadCasting System Info")
	s.mu.Lock()
	logger.Log("Lock")
	Unlock := func() {
		logger.Log("Unlock")
		s.mu.Unlock()
	}
	defer Unlock()
	clients := make([]entitys.ClientRequestBody, s.ExpectedWorkers)
	for i := 0; i < len(clients); i++ {
		publicKey := s.Workers[i]
		key := getHash(publicKey)

		row := entitys.ClientRequestBody{PublicKey: s.Workers[i], Client: s.Clients[key]}
		clients[i] = row
	}
	var payload entitys.RabbitMqSystemInfoPack
	payload.Clients = clients
	payload.ExpectedWorkers = s.ExpectedWorkers
	payload.ScaleFactorMsg = sFm
	payload.ScaleFactorCoin = sFc

	err = lagoudaki.BroadCastSystemInfo(payload)
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort due: %s", err.Error()))
		return err
	}
	logger.Log("Commit  BroadCasting System Info")
	return nil

}
func (s *SystemInfoImpl) Consume() (error, float64, float64) {
	providers := &s.Providers

	err := providers.valid()
	if err != nil {
		return err, 0, 0
	}
	logger := providers.LoggerService
	lagoudaki := providers.RabbitMqService
	logger.Log("Start  Consuming System Info")
	logger.Log("Waiting  for  systemInfo ... ")

	info := lagoudaki.ConsumeNextSystemInfo()
	logger.Log("consume  systemInfo ... ")
	if s.Coordinator {
		logger.Log("Commit  Consuming System Info")
		s.ok = true
		return nil, info.ScaleFactorMsg, info.ScaleFactorCoin
	}
	s.ExpectedWorkers = info.ExpectedWorkers
	err = s.Construct()

	if err != nil {
		logger.Fatal(err.Error())
		return err, 0, 0
	}
	for _, body := range info.Clients {
		log.Println(body)
		err = s.AddWorker(body)
		if err != nil {
			logger.Fatal(err.Error())
			return err, 0, 0
		}
	}

	s.ok = true
	logger.Log("Commit  Consuming System Info")
	return nil, info.ScaleFactorMsg, info.ScaleFactorCoin

}
func (s *SystemInfoImpl) GetWorkers() []rsa.PublicKey {
	return s.Workers
}

func (s *SystemInfoImpl) IsOk() bool {
	return s.ok
}
func (s *SystemInfoImpl) NodeDetails(key rsa.PublicKey) (entitys.ClientInfo, int) {
	hashedPublicKey := getHash(key)
	return s.Clients[hashedPublicKey], len(s.Workers)
}
func (s *SystemInfoImpl) Who(index int) (rsa.PublicKey, error) {
	var key rsa.PublicKey
	if index >= len(s.Workers) {
		return key, errors.New("Worker does not exist")
	}
	key = s.Workers[index]
	return key, nil
}

func (s *SystemInfoImpl) ClientList(key rsa.PublicKey) ([]entitys.ClientInfo, int) {
	var list []entitys.ClientInfo
	hashedPublicKey := getHash(key)
	for key, info := range s.Clients {
		if key != hashedPublicKey {
			list = append(list, info)

		}
	}
	return list, len(list)

}

type BroadCastFuncParams struct {
	ScaleFactorMessage float64
	ScaleFactorCoins   float64
}
type ConsumeResponse struct {
	Error            error
	ScaleFactorMsg   float64
	ScaleFactorCoins float64
}
type NodeDetailsResponse struct {
	Info  entitys.ClientInfo
	Total int
}
type WhoResponse struct {
	Key   rsa.PublicKey
	Error error
}
type ClientListResponse struct {
	Clients []entitys.ClientInfo
	Total   int
}
type MockSystemInfoService struct {
	ErrConstruct        error
	ErrAddWorker        error
	ErrBroadCast        error
	IsFullResponse      bool
	ConsumeRsp          ConsumeResponse
	IsOkResposne        bool
	Workers             []rsa.PublicKey
	NodeDetailsRsp      NodeDetailsResponse
	WhoRsp              WhoResponse
	ClientListRsp       ClientListResponse
	CallAddWorker       int
	CallConstruct       int
	CallIsFull          int
	CallBroadCast       int
	CallConsume         int
	CallIsOk            int
	CallGetWorkers      int
	CallNodeDetails     int
	CallWho             int
	CallClientList      int
	CallAddWorkerWith   []entitys.ClientRequestBody
	CallBroadCastWith   []BroadCastFuncParams
	CallNodeDetailsWith []rsa.PublicKey
	CallWhoWith         []int
	CallClientListWith  []rsa.PublicKey
}

// Construct
func (mock *MockSystemInfoService) Construct() error {
	mock.CallConstruct++
	return mock.ErrConstruct
}
func (mock *MockSystemInfoService) AddWorker(body entitys.ClientRequestBody) error {
	mock.CallAddWorker++
	mock.CallAddWorkerWith = append(mock.CallAddWorkerWith, body)
	return mock.ErrAddWorker
}
func (mock *MockSystemInfoService) IsFull() bool {
	mock.CallIsFull++
	return mock.IsFullResponse
}
func (mock *MockSystemInfoService) BroadCast(sFm, sFc float64) error {
	mock.CallBroadCast++
	params := BroadCastFuncParams{
		ScaleFactorMessage: sFm,
		ScaleFactorCoins:   sFc,
	}
	mock.CallBroadCastWith = append(mock.CallBroadCastWith, params)
	return mock.ErrBroadCast
}
func (mock *MockSystemInfoService) Consume() (error, float64, float64) {
	mock.CallConsume++
	return mock.ConsumeRsp.Error, mock.ConsumeRsp.ScaleFactorMsg, mock.ConsumeRsp.ScaleFactorCoins
}
func (mock *MockSystemInfoService) IsOk() bool {
	mock.CallIsOk++
	return mock.IsOkResposne
}
func (mock *MockSystemInfoService) GetWorkers() []rsa.PublicKey {
	mock.CallGetWorkers++
	return mock.Workers
}
func (mock *MockSystemInfoService) NodeDetails(key rsa.PublicKey) (entitys.ClientInfo, int) {
	mock.CallNodeDetails++
	mock.CallNodeDetailsWith = append(mock.CallNodeDetailsWith, key)
	return mock.NodeDetailsRsp.Info, mock.NodeDetailsRsp.Total
}
func (mock *MockSystemInfoService) Who(index int) (rsa.PublicKey, error) {
	mock.CallWho++
	mock.CallWhoWith = append(mock.CallWhoWith, index)

	return mock.WhoRsp.Key, mock.WhoRsp.Error
}
func (mock *MockSystemInfoService) ClientList(key rsa.PublicKey) ([]entitys.ClientInfo, int) {
	mock.CallClientList++
	mock.CallClientListWith = append(mock.CallClientListWith, key)
	return mock.ClientListRsp.Clients, mock.ClientListRsp.Total
}
func getHash(key rsa.PublicKey) string {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		log.Fatal(err)
	}

	// Hash the public key bytes using SHA-256
	hash := sha256.Sum256(publicKeyBytes)
	return hex.EncodeToString(hash[:])
}
