package MessageSystem

import (
	"Logger"
	"encoding/json"
	"errors"
	"fmt"
)

type ConsumerMsgResp[T interface{}] struct {
	Payload T
	Err     error
}

func CreateAndBind(RABBITMQ, QUEUE, EXCHANGE string, logger Logger.LoggerService) error {
	whatToDo := fmt.Sprintf("create and  bind queue %s to topic %s", QUEUE, EXCHANGE)

	logger.Log(fmt.Sprintf("start %s", whatToDo))
	conn, channel, err := connectionMakerBroadcast(RABBITMQ, EXCHANGE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("abbort %s due  to %s ", whatToDo, err.Error()))
		return err
	}

	defer conn.Close()
	defer channel.Close()
	err = CreateAndBindQueue(channel, QUEUE, EXCHANGE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("abbort %s due  to %s ", whatToDo, err.Error()))
		return err
	}

	logger.Log(fmt.Sprintf("commit %s", whatToDo))
	return nil

}

/*
  - Cunsumer - Create  a consumer for T type messages
    @generic T - interface{}{}
    @Param messageBus chan  ConsumerMsgResp[T]
    @Param RABBITMQ string
    @Param QUEUE  string
    @Param EXCHANGENAME string
    @Param  logger Logger.LoggerService
*/
func Consumer[T interface{}](messageBus chan ConsumerMsgResp[T], RABBITMQ, QUEUE, EXCHANGENAME string, logger Logger.LoggerService) {
	logger.Log("Start  Consuming")
	var Payload T // Zero val
	conn, channel, err := connectionMakerBroadcast(RABBITMQ, EXCHANGENAME, logger)
	if err != nil {
		logger.Error("Abbort  Consuming")
		message := ConsumerMsgResp[T]{Payload: Payload, Err: err}
		messageBus <- message
		return
	}
	defer conn.Close()
	defer channel.Close()
	//Declare  Queue
	err = CreateAndBindQueue(channel, QUEUE, EXCHANGENAME, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort  Consuming due to :%s ", err.Error()))
		message := ConsumerMsgResp[T]{Payload: Payload, Err: err}
		messageBus <- message
		return
	}
	msgs, err := channel.Consume(
		QUEUE, // Queue name
		"",    // Consumer
		false, // Auto-Acknowledge set to false
		true,  // Exclusive (set to true to ensure only one consumer at a time)
		false, // No-local
		false, // No-Wait
		nil,   // Arguments
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToRegisterConsumer, err.Error())
		logger.Error(fmt.Sprintf("Abbort  Consuming due to :%s ", errMsg))
		err = errors.New(errMsg)
		message := ConsumerMsgResp[T]{Payload: Payload, Err: err}
		messageBus <- message
		return
	}
	for msg := range msgs {
		var m T
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			errMsg := fmt.Sprintf("Failed  to unmarshal json due %s", err.Error())
			logger.Error(errMsg)
		} else {
			logger.Log("Received  Message")
			messageBus <- ConsumerMsgResp[T]{Payload: m, Err: nil}
		}
		msg.Ack(false)
	}
}

/*
  - CunsumeOne - Create  a consumer for T type messages
    @generic T - interface{}{}
    @Param RABBITMQ string
    @Param QUEUE  string
    @Param  logger Logger.LoggerService
    @Returns T ,error
*/
func ConsumeOne[T interface{}](RABBITMQ, QUEUE string, logger Logger.LoggerService) (T, error) {
	logger.Log("Start  Consuming")
	var Payload T // Zero val
	conn, channel, err := connectionMaker(RABBITMQ, logger)
	if err != nil {
		logger.Error("Abbort  Consuming")
		return Payload, err
	}
	defer conn.Close()
	defer channel.Close()
	msgs, err := channel.Consume(
		QUEUE, // Queue name
		"",    // Consumer
		false, // Auto-Acknowledge set to false
		true,  // Exclusive (set to true to ensure only one consumer at a time)
		false, // No-local
		false, // No-Wait
		nil,   // Arguments
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToRegisterConsumer, err.Error())
		logger.Error(fmt.Sprintf("Abbort  Consuming due to :%s ", errMsg))
		err = errors.New(errMsg)
		return Payload, err
	}
	var errReturn error
	for msg := range msgs {
		err := json.Unmarshal(msg.Body, &Payload)
		if err != nil {
			errMsg := fmt.Sprintf("Failed  to unmarshal json due %s", err.Error())
			logger.Error(errMsg)
			errReturn = err
		} else {
			logger.Log("Received  Message")
		}
		msg.Ack(false)
		break
	}
	return Payload, errReturn

}
