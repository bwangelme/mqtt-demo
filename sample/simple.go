package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	exit <- true
}

var exit = make(chan bool)

func main() {
	n := 5
	opts := mqtt.NewClientOptions()

	opts.AddBroker("tcp://vm.bwangel.me:1883")
	opts.SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	for i := 0; i < n; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
	}

	for i := 0; i < n; i++ {
		<-exit
	}

	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
