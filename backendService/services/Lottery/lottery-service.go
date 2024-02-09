package Lottery

import (
	"Hasher"
	"Logger"
	"Service"
	"crypto/rsa"
	"errors"
	"math/rand"

	"Stake"
)

type LotteryService interface {
	Service.Service
	Spin(scaleFactor float64) (rsa.PublicKey, error)
	LoadStakeService(stakeService Stake.StakeService)
}

type LotteryProviders struct {
	LoggerService Logger.LoggerService
	HasherService Hasher.HashService
	StakeService  Stake.StakeService
}
type LotteryImpl struct {
	Services LotteryProviders
}

func (service LotteryImpl) Construct() error {
	return nil
}
func (service *LotteryImpl) LoadStakeService(stakeService Stake.StakeService) {
	service.Services.StakeService = stakeService
}

func (service *LotteryImpl) Spin(scaleFactor float64) (rsa.PublicKey, error) {

	hasher := service.Services.HasherService
	stake := service.Services.StakeService
	if stake == nil {
		return rsa.PublicKey{}, errors.New("Stake Service Not  Loaded ")
	}

	seed, err := hasher.Seed(stake.GetCurrentHash())

	if err != nil {
		return rsa.PublicKey{}, err
	}
	epiTis100, apo := stake.MapOfDistibutesRoundUp(scaleFactor)

	rand.Seed(seed)
	//make  spin
	spin := rand.Intn(apo + 1) //use  rand and seed
	sum := 0
	for _, key := range stake.GetWorkers() {
		tis100 := epiTis100[key]
		sum += tis100
		if sum >= spin {
			return key, nil

		}
	}
	return rsa.PublicKey{}, errors.New("Spin Max  is  bigger  than  that the distributionMap ") // if  stake service corect this  will never  been send

}

type MockLottery struct {
	SpinError     error
	SpinRsp       rsa.PublicKey
	CallSpin      int
	CallLoadStake int
}

func (service *MockLottery) Construct() error {
	return nil
}
func (service *MockLottery) LoadStakeService(stakeService Stake.StakeService) {
	service.CallLoadStake++
}
func (service *MockLottery) Spin(scaleFactor float64) (rsa.PublicKey, error) {
	service.CallSpin++
	return service.SpinRsp, service.SpinError
}
