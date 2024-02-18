package MessageSystem

import (
	"Logger"
	"fmt"
	"testing"
	"time"
)

func callExpector[T comparable](obj1, obj2 T, t *testing.T, prefix, what string) {
	if obj1 != obj2 {
		t.Errorf("%s  Expected  '%s' to get %v but %v ", prefix, what, obj1, obj2)
	}
}
func expectorNoErr(t *testing.T, err error, prefixOld string) {
	if err != nil {
		t.Errorf("%s Expect no Err  but  got %v", prefixOld, err)
	}
}

func TestConsumer(t *testing.T) {
	const prefix = "----"
	fmt.Printf("*  test  for  consumers\n")

	const RABBITMQ = "amqp://v:123456@127.0.0.1:5672/"
	const EXCHANGE = "collors"
	const QUEUE = "colorQueue"
	logger := &Logger.MockLogger{}
	err := CreateAndBind(RABBITMQ, QUEUE, EXCHANGE, &Logger.MockLogger{})
	expectorNoErr(t, err, "")
	type msgType struct {
		Msg string `json:"msg"`
	}
	msg1 := msgType{Msg: "green"}
	msg2 := msgType{Msg: "yellow"}
	msg3 := msgType{Msg: "red"}

	err = ProducerBroadCast(msg1, RABBITMQ, EXCHANGE, logger)
	expectorNoErr(t, err, "")
	err = ProducerBroadCast(msg2, RABBITMQ, EXCHANGE, logger)
	expectorNoErr(t, err, "")
	err = ProducerBroadCast(msg3, RABBITMQ, EXCHANGE, logger)
	expectorNoErr(t, err, "")
	TestConsumeOne := func(prefixOld string) {
		itShouldConsumeOne := func(prefixOld string) {
			first, err := ConsumeOne[msgType](RABBITMQ, QUEUE, logger)
			expectorNoErr(t, err, "")
			callExpector(msg1, first, t, prefixOld, "first message ")
			second, err := ConsumeOne[msgType](RABBITMQ, QUEUE, logger)
			expectorNoErr(t, err, "")
			callExpector(msg2, second, t, prefixOld, "second  message ")
			third, err := ConsumeOne[msgType](RABBITMQ, QUEUE, logger)
			expectorNoErr(t, err, "")
			callExpector(msg3, third, t, prefixOld, "third message ")
		}
		itShouldConsumeOne(prefixOld)
	}
	TestConsumeOne(prefix)
}

func TestSendOnePackageTo3Conusmers(t *testing.T) {
	logger := &Logger.MockLogger{}
	type msgType struct {
		Msg string `json:"msg"`
	}
	msg := msgType{Msg: "hello"}

	const RABBITMQ = "amqp://v:123456@127.0.0.1:5672/"
	const EXCHANGENAME = "helloExchange"
	const QUEUE = "helloQueue"
	consumer := 0
	consumerChannel := make(chan int)
	con := func(who string) {
		resp := make(chan ConsumerMsgResp[msgType])
		go Consumer[msgType](resp, RABBITMQ, QUEUE+who, EXCHANGENAME, logger)
		msg := <-resp
		if msg.Payload.Msg == "hello" {
			consumerChannel <- 1
		} else {
			consumerChannel <- 0
		}

	}
	go con("consumer 1")
	go con("consumer 2")
	go con("consumer 3")
	time.Sleep(2 * time.Second)
	err := ProducerBroadCast(msg, RABBITMQ, EXCHANGENAME, logger)
	if err != nil {
		t.Errorf("Failed to Produce msg ")
	}

	for i := 0; i < 3; i++ {
		consumer += <-consumerChannel
	}
	if consumer == 3 {
		fmt.Println("it should  send  1 json package  to  3 consumers ")
	} else {
		t.Errorf("The 3 consumer  do not receive The same  message 'hello'")
	}
}
func TestSendThreePackageTo3Conusmers(t *testing.T) {
	logger := &Logger.MockLogger{}
	type msgType struct {
		Msg string `json:"msg"`
	}
	msg1 := msgType{Msg: "green"}
	msg2 := msgType{Msg: "red"}
	msg3 := msgType{Msg: "yellow"}

	const RABBITMQ = "amqp://v:123456@127.0.0.1:5672/"
	const EXCHANGENAME = "helloExchange2"
	const QUEUE = "helloQueue"
	consumer := 0
	consumerChannel := make(chan int)
	con := func(who string) {
		resp := make(chan ConsumerMsgResp[msgType])
		go Consumer[msgType](resp, RABBITMQ, QUEUE+who, EXCHANGENAME, logger)
		msg := <-resp
		if msg.Payload.Msg == "green" {
			msg := <-resp
			if msg.Payload.Msg == "red" {
				msg := <-resp
				if msg.Payload.Msg == "yellow" {
					consumerChannel <- 1

				} else {
					consumerChannel <- 0
				}
			} else {
				consumerChannel <- 0
			}
		} else {
			consumerChannel <- 0
		}

	}
	go con("consumer 1 1")
	go con("consumer 2 2")
	go con("consumer 3 3")
	time.Sleep(2 * time.Second)
	err := ProducerBroadCast(msg1, RABBITMQ, EXCHANGENAME, logger)
	if err != nil {
		t.Errorf("Failed to Produce msg ")
	}
	time.Sleep(2 * time.Second)
	err = ProducerBroadCast(msg2, RABBITMQ, EXCHANGENAME, logger)
	if err != nil {
		t.Errorf("Failed to Produce msg ")
	}
	time.Sleep(2 * time.Second)
	err = ProducerBroadCast(msg3, RABBITMQ, EXCHANGENAME, logger)
	if err != nil {
		t.Errorf("Failed to Produce msg ")
	}

	for i := 0; i < 3; i++ {
		consumer += <-consumerChannel
	}
	if consumer == 3 {
		fmt.Println("it should  send  3 json package  to  3 consumers and  all should be Get by FIFO ")
	} else {
		t.Errorf("The 3 consumer  do not receive The same  message 'hello'")
	}
}
