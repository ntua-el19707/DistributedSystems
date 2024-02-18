package MessageSystem

import (
	"Logger"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
ProducerBroadcast - Produce  A Message  and  Brodacast The Message
@Param Message - interface {}
@Param RABBITMQ string
@Param EXCHANGENAME string
*/
func ProducerBroadCast(Message interface{}, RABBITMQ, EXCHANGENAME string, logger Logger.LoggerService) error {
	logger.Log("Start  BroadCastMessage")
	//Open  a  connection
	conn, channel, err := connectionMakerBroadcast(RABBITMQ, EXCHANGENAME, logger)
	if err != nil {
		logger.Error("Abbort  BroadCastMessage")
		return err
	}
	defer conn.Close()
	defer channel.Close()
	//Marshal json
	payload, err := json.Marshal(Message)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToMarshalJson, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort : %s", errMsg))
		return err
	}

	//Publish
	err = channel.Publish(
		EXCHANGENAME, //Exchange
		"",           //Routing Key
		true,         //Mandatory
		false,        //Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToPublish, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort : %s", errMsg))
		return err
	}
	logger.Log("Commit  BroadCastMessage")
	return nil
}

/*
Publish  - publish  A Message  to Queue
@Param Message - interface {}
@Param RABBITMQ string
@Param EXCHANGENAME string
*/
func Publish(Message interface{}, RABBITMQ, QUEUE string, Exclusive bool, logger Logger.LoggerService) error {
	logger.Log("Start  Publish Message ")
	//Open  a  connection
	conn, channel, err := connectionMaker(RABBITMQ, logger)
	if err != nil {
		logger.Error("Abbort  Publish Message ")
		return err
	}
	defer conn.Close()
	defer channel.Close()
	//Marshal json
	payload, err := json.Marshal(Message)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToMarshalJson, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort Publish Message  : %s", errMsg))
		return err
	}
	//Create  The  Queue
	err = CreateAQueue(channel, QUEUE, Exclusive, logger)
	if err != nil {
		return err
	}
	err = channel.Publish(
		"",
		QUEUE,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToPublish, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort  Publish Message: %s", errMsg))
		return err
	}
	logger.Log("Commit  Publish Message ")
	return nil
}
