package Register

import (
	"Logger"
	"Service"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"entitys"
	"fmt"
	"net/http"
	"time"
)

type RegisterService interface {
	Service.Service
	Register()
}

type RegisterProviders struct {
	LoggerService Logger.LoggerService
}

const serviceName = "register-service"

func (p *RegisterProviders) Construct() error {
	if p.LoggerService == nil {
		p.LoggerService = &Logger.Logger{ServiceName: serviceName}
		err := p.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	return nil
}

type RegisterImpl struct {
	Me        string
	Who       string
	MyPk      rsa.PublicKey
	MyId      string
	UriPublic string
	Providers RegisterProviders
}

func (r *RegisterImpl) Construct() error {
	return r.Providers.Construct()
}

func (r *RegisterImpl) Register() {
	providers := &r.Providers
	logger := providers.LoggerService
	to := r.Who

	var body entitys.ClientRequestBody
	body.PublicKey = r.MyPk
	body.Client.Id = r.MyId
	body.Client.Uri = r.Me
	body.Client.UriPublic = r.UriPublic

	payload, err := json.Marshal(body)
	if err != nil {
		logger.Fatal(err.Error())
	}

	to = fmt.Sprintf("%s/api/v1/register", to)
	request, err := http.NewRequest("POST", to, bytes.NewBuffer(payload))
	if err != nil {
		logger.Fatal(err.Error())
	}
	client := &http.Client{}
	for {
		res, err := client.Do(request)
		if err != nil {
			//	logger.Fatal(err.Error())
		} else {
			if res.StatusCode == http.StatusOK {
				break
			}
		}
		//try again in  3 seconds
		time.Sleep(time.Second * 3)

	}

}
