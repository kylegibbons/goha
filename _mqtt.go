package goha

import (
	"context"
	"fmt"
	"log"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient is
type MQTTClient struct {
	ctx           context.Context
	Broker        string
	ClientID      string
	Username      string
	Password      string
	CleanSession  bool
	chSubUpdate   chan mqttSubscriptionChanged
	subscriptions map[string][]chan MQTTMessage
	client        *MQTT.Client
}

// MQTTMessage is
type MQTTMessage struct {
	topic   string
	message string
}

// Init is
func (mc *MQTTClient) Init(ctx context.Context) {
	mc.subscriptions = make(map[string][]chan MQTTMessage)
	mc.chSubUpdate = make(chan mqttSubscriptionChanged, 1000)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(mc.Broker)
	opts.SetClientID(mc.ClientID)
	opts.SetUsername(mc.Username)
	opts.SetPassword(mc.Password)
	opts.SetCleanSession(mc.CleanSession)

	opts.SetDefaultPublishHandler(mc.messageHandler)

	mc.client = MQTT.NewClient(opts)
	if token := *mc.client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	go mc.manageSubscriptions(ctx)
}

func (mc *MQTTClient) messageHandler(client MQTT.Client, msg MQTT.Message) {
	for topic, subs := range mc.subscriptions {
		if topic == msg.Topic() {
			for _, ch := range subs {
				ch <- MQTTMessage{
					topic:   msg.Topic(),
					message: string(msg.Payload()),
				}
			}
			return
		}
	}
}

type mqttSubscriptionChanged struct {
	subscribe bool
	topic     string
	ch        chan MQTTMessage
}

func (mc *MQTTClient) manageSubscriptions(ctx context.Context) {

	log.Printf("Starting MQTT Subscription Manager")

	for {
		select {
		case <-ctx.Done():
			return
		case subChange := <-mc.chSubUpdate:
			log.Printf("Processing MQTT subscription change")

			if subscribers, ok := mc.subscriptions[subChange.topic]; ok {
				subFound := false

				for i, subscriber := range subscribers {
					if subChange.ch == subscriber && !subChange.subscribe {
						subFound = true

						subscribers[i] = subscribers[len(subscribers)-1] // Copy last element to index i.
						subscribers[len(subscribers)-1] = nil            // Erase last element (write zero value).
						subscribers = subscribers[:len(subscribers)-1]

					}
				}

				if !subFound {
					mc.subscriptions[subChange.topic] = append(mc.subscriptions[subChange.topic], subChange.ch)
				}

				if len(subscribers) == 0 {
					delete(mc.subscriptions, subChange.topic)
				}

				if len(subscribers) > 0 {
					mc.subscriptions[subChange.topic] = subscribers
				}

			} else {
				mc.subscriptions[subChange.topic] = []chan MQTTMessage{subChange.ch}
			}

			//log.Printf("%+v", em.subscriptions)
		}
	}
}

// Subscribe is
func (mc *MQTTClient) Subscribe(topic string, ch chan MQTTMessage) {

	if token := client.Subscribe(*topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	mc.chSubUpdate <- mqttSubscriptionChanged{
		subscribe: true,
		topic:     topic,
		ch:        ch,
	}
}

// Unsubscribe is
func (mc *MQTTClient) Unsubscribe(topic string, ch chan MQTTMessage) {
	mc.chSubUpdate <- mqttSubscriptionChanged{
		subscribe: false,
		topic:     topic,
		ch:        ch,
	}
}
