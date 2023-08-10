package main

import (
	// "fmt"
	"log"
	"productConsumer/pkg/helpers"
	"productConsumer/pkg/productModule"
)

func main(){

	if err := helpers.InitEnvSetting();err !=nil{
		panic(err)
	}

	channelName := helpers.GetEnvStr("amqp.productchannel")

	if err,_ := helpers.InitMySql();err !=nil{
		panic("SQL is down")
	}

	msgs , err := helpers.InitQueueConnect(channelName)
	if err !=nil{
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			go productModule.InitialProductController().Reduce(d.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}