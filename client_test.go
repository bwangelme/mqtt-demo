package main

import (
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMqttPubSub(t *testing.T) {
	const TOPIC = "mytopic/test"
	post_message := "hello, xff1"

	client := NewMQTTClient("tcp://localhost:1883", "admin", "public")

	client.Sub(TOPIC, func(client mqtt.Client, msg mqtt.Message) {
		message := string(msg.Payload())
		if post_message != message {
			t.Fatalf("%s != %s\n", post_message, message)
		}
	})

	client.Pub(TOPIC, post_message)
	client.Wait()
}
