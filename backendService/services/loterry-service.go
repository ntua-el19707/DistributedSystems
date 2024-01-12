package services 

import  (
	"crypto/rsa"
	"math/rand"
	"errors"
	"fmt"
)
type  lotteryService interface   {
	Service 
	spin()(rsa.PublicKey ,error)
} 

type  lotteryProviders struct {
	logger LogerService 
	hasher hashService
	stake  stakeService 
} 
type  lotteryImpl struct {
	services lotteryProviders 
}
func  (service  lotteryImpl ) construct() error    {
	return  nil
}
func  (service * lotteryImpl ) spin() (rsa.PublicKey ,error )  {

	seed , err  := service.services.hasher.Seed(service.services.stake.getCurrentHash() )
	if  err != nil {
		return rsa.PublicKey{} ,  err
	} 
    epiTis100 , apo := service.services.stake.MapOfDistibutesRoundUp(0.0)

	rand.Seed(seed)
	//make  spin 
	spin := rand.Intn(apo + 1) //use  rand and seed 
	sum := 0 
	fmt.Println(spin)

	for  _ , key  := range service.services.stake.getWorkers() {
		tis100 :=  epiTis100[key]
		sum +=  tis100 
		if  sum >= spin {
			return key  ,nil 

		} 
	} 
	return  rsa.PublicKey{} ,errors.New("Spin Max  is  bigger  than  that the distributionMap ") // if  stake service corect this  will never  been send




}
