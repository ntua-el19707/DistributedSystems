package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestTransactions interface {
	valid(To, Total int) error
}
type SetStakeRequest struct {
	Stake float64 `json:"stake"`
}

func (r *SetStakeRequest) Parse(request *http.Request) (int16, error) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil

}

type TranferMoney struct {
	To     int     `json:"to"`
	Amount float64 `json:"amount"`
}
type TranferMsg struct {
	To  int    `json:"to"`
	Msg string `json:"msg"`
}

const errCannotentToMyself = "Cannot send  to myself"
const errReceiverDoesNotExist = "Receiver does not  exist"

func (t TranferMoney) valid(Me, Total int) error {
	return validator[float64](t.To, Me, Total, t.Amount, "amount")
}
func (t TranferMsg) valid(Me, Total int) error {
	return validator[string](t.To, Me, Total, t.Msg, "msg")
}

func validator[T comparable](To, Me, Total int, value T, field string) error {
	var zero T
	if To == Me {
		return errors.New(errCannotentToMyself)
	}
	if To >= Total {
		return errors.New(errReceiverDoesNotExist)
	}
	if value == zero {
		return errors.New(fmt.Sprintf("field  '%s' cannot  be:'%v' ", field, zero))
	}
	return nil
}
func ParseRequest[T RequestTransactions](r *http.Request, Me, Total int) (int16, error, T) {

	decoder := json.NewDecoder(r.Body)
	var payload T
	err := decoder.Decode(&payload)
	if err != nil {
		return http.StatusBadRequest, err, payload
	}
	err = payload.valid(Me, Total)
	if err != nil {
		return http.StatusBadRequest, err, payload
	}

	return http.StatusOK, nil, payload

}
