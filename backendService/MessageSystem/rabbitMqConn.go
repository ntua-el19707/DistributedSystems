package MessageSystem

import (
	"Logger"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
  - connectionMaker - Create  a conection ,channel and exachane
    @Param RABBITMQ string
    @Param logger  Logger.LoggerService
*/
func connectionMaker(RABBITMQ string, logger Logger.LoggerService) (*amqp.Connection, *amqp.Channel, error) {
	logger.Log("Start  Creating  a connection ")
	//Open  a  connection
	conn, err := amqp.Dial(RABBITMQ)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToDial, RABBITMQ, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort : %s", errMsg))
		return nil, nil, err
	}
	//Create a channel
	channel, err := conn.Channel()
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToCreateChannel, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort : %s", errMsg))
		conn.Close()
		return nil, nil, err
	}
	logger.Log("Commit  Creating  a connection ")
	return conn, channel, nil
}

/*
  - connectionMakerBroadcast - Create  a conection ,channel and exachane
    @Paran RABBITMQ string
    @Param EXCHANGENAME string
*/
func connectionMakerBroadcast(RABBITMQ, EXCHANGENAME string, logger Logger.LoggerService) (*amqp.Connection, *amqp.Channel, error) {
	logger.Log("Start  Creating  a  broadcast connection ")
	//Create  connection and  channel
	conn, channel, err := connectionMaker(RABBITMQ, logger)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToDial, RABBITMQ, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort : %s", errMsg))
		return nil, nil, err
	}
	//Declare  The Exchange
	err = channel.ExchangeDeclare(
		EXCHANGENAME,
		"topic",
		true,  //Durable
		false, //Auto-Deleted
		false, //Internal
		false, //No-wait
		nil,
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToDeclareExchange, err.Error())
		err = errors.New(errMsg)
		logger.Error(fmt.Sprintf("Abbort : %s", errMsg))
		conn.Close()
		channel.Close()
		return nil, nil, err
	}
	logger.Log("Commit  Creating  a  broadcast connection ")
	return conn, channel, nil
}

/*
*
CreateAndBindQueue - create  and  bind queue  in to an exachange  Topic
@Param  channel *  amqp.Channel
@Param  QUEUE string
@Param	EXCHANGENAME string
@Param	logger Logger.LoggerService
*/
func CreateAndBindQueue(channel *amqp.Channel, QUEUE, EXCHANGENAME string, logger Logger.LoggerService) error {
	logger.Log("Start Binding Queue  to the  Topic ")

	_, err := channel.QueueDeclare(
		QUEUE, // Queue name
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive (set to false)
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToCreateQueue, err.Error())
		logger.Error(fmt.Sprintf("Abbort Binding Queue  to the  Topic due to %s", errMsg))
		err = errors.New(errMsg)
		return err
	}
	err = channel.QueueBind(
		QUEUE,        // Queue name
		"#",          // Routing pattern (matches all topics)
		EXCHANGENAME, // Exchange name
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToBindQueue, err.Error())
		logger.Error(fmt.Sprintf("Abbort Binding Queue  to the  Topic due to %s", errMsg))
		err = errors.New(errMsg)
		return err
	}
	logger.Log("Commit Binding Queue  to the  Topic ")
	return nil
}

/*
*
CreateAQueue  - create  a  Queue
@Param  channel *  amqp.Channel
@Param  QUEUE string
@Param	Exclusive bool
@Param	logger Logger.LoggerService
*/
func CreateAQueue(channel *amqp.Channel, QUEUE string, Exclusive bool, logger Logger.LoggerService) error {
	logger.Log("Start create a Queue ")

	_, err := channel.QueueDeclare(
		QUEUE,     // Queue name
		true,      // Durable
		false,     // Delete when unused
		Exclusive, // Exclusive (set to false)
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedToCreateQueue, err.Error())
		logger.Error(fmt.Sprintf("Abbort create  a Queue :%s", errMsg))
		err = errors.New(errMsg)
		return err
	}
	logger.Log("Commit create a  Queue")
	return nil
}
