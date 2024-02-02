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
func (service *LotteryImpl) Spin(scaleFactor float64) (rsa.PublicKey, error) {
	hasher := service.Services.HasherService
	stake := service.Services.StakeService
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
