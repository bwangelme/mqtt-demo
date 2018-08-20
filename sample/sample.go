package main

import (
	"flag"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
Options:
	[-help]                      Display help
	[-a pub|sub]                 Action pub (publish) or sub (subscribe)
	[-m <message>]               Payload to send
	[-n <number>]                Number of messages to send or receive
	[-q 0|1|2]                   Quality of Service
	[-clean]                     CleanSession (true if -clean is present)
	[-id <clientid>]             CliendID
	[-user <user>]               User
	[-password <password>]       Password
	[-broker <uri>]              Broker URI
	[-topic <topic>]             Topic
	[-store <path>]              Store Directory

*/

func main() {
	action := flag.String("action", "", "Action to Publish/Subscribe (required)")
	message := flag.String("message", "Hello, EMQTT!", "Message to send")
	number := flag.Int("num", 1, "Message send number")
	qos := flag.Int("qos", 1, "Quality of Service (0|1|2)")
	broker := flag.String("broker", "tcp://vm.bwangel.me:1883", "Broker Addr(ex tcp://vm.bwangel.me:1883)")
	topic := flag.String("topic", "", "Message Topic")

	flag.Parse()

	if *action != "pub" && *action != "sub" {
		fmt.Println("action must be pub or sub")
		return
	}

	if *topic == "" {
		fmt.Println("topic cannot be empty")
		return
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(*broker)

	if *action == "sub" {
		message := make(chan [2]string)

		opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
			content := msg.Payload()
			topic := msg.Topic()

			message <- [2]string{topic, string(content)}
		})
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalln(token.Error())
		}

		if token := client.Subscribe(*topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
			log.Fatalln(token.Error())
		}

		for msgCount := 0; msgCount < *number; msgCount++ {
			msg := <-message
			topic, content := msg[0], msg[1]
			fmt.Printf("[%s]Recv: %s\n", topic, content)
		}

	} else {
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalln(token.Error())
		}

		fmt.Println("Start to send message")

		for msgCount := 0; msgCount <= *number; msgCount++ {
			if token := client.Publish(*topic, byte(*qos), false, *message); token.Wait() && token.Error() != nil {
				log.Fatalln(token.Error())
			}
		}
	}
}
