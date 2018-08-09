package main

import (
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	addr      string
	username  string
	password  string
	wg        *sync.WaitGroup
	rawClient mqtt.Client
}

func NewMQTTClient(addr, username, password string) *MQTTClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(addr)
	opts.SetUsername(username)
	opts.SetPassword(password)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	var wg sync.WaitGroup

	return &MQTTClient{
		addr:      addr,
		username:  username,
		password:  password,
		wg:        &wg,
		rawClient: client,
	}
}

func (self *MQTTClient) Sub(topic string, callback func(client mqtt.Client, msg mqtt.Message)) {
	if token := self.rawClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		callback(client, msg)
		self.wg.Done()
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func (self *MQTTClient) Pub(topic string, message string) {
	if token := self.rawClient.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	self.wg.Add(1)

}

func (self *MQTTClient) Wait() {
	self.wg.Wait()
}
