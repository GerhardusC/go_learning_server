package datacollection

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"testing-server/cliArgs"
	"testing-server/dbInteractions"

	"github.com/eclipse/paho.mqtt.golang"
)

func messageHandler (client mqtt.Client, message mqtt.Message) {
	msg := string(message.Payload())

	measurement := dbInteractions.DBRowMeasurement[string]{
		// Timestamp gets automatically assigned at insertion, saving operation making it 0 here.
		Timestamp: 0,
		Topic: message.Topic(),
		Value: msg,
	}

	_, err := strconv.ParseFloat(msg, 64)

	var tableName string

	if err != nil {
		log.Println("Failed to parse message to float, saving as log")
		tableName = "LOGS"
	} else {
		tableName = "MEASUREMENTS"
	}

	measurement.WriteToTable(tableName)
}

func CollectData () {
	clientOpts := mqtt.NewClientOptions()
	clientOpts.SetProtocolVersion(4)
	clientOpts.SetOrderMatters(false)
	clientOpts.SetAutoReconnect(true)
	clientOpts.AddBroker(fmt.Sprintf("mqtt://%s:1883", cliargs.HostIP))
	client := mqtt.NewClient(clientOpts)

	connectionToken := client.Connect()
	connectionToken.Wait()

	if connectionToken.Error() != nil {
		log.Println(connectionToken.Error())
	}

	topic, err := validateTopic(cliargs.BaseTopic)

	if err != nil {
		log.Println(err.Error(), "\nUsing fallback topic:", topic)
	} else {
		log.Println("Using topic:", topic)
	}

	subscriptionToken := client.Subscribe(topic, 0, messageHandler)

	for {
		<- subscriptionToken.Done()
		if subscriptionToken.Error() != nil {
			log.Println(subscriptionToken.Error())
		}
	}
}

func validateTopic (baseTopic string) (string, error) {
	if baseTopic == "/#" {
		return baseTopic, nil
	}

	var startsWithSlash bool
	var endsWithSlash bool

	length := len(baseTopic)
	lastIndex := length - 1

	if length <= 1 {
		return "/#", errors.New("Failed to form base topic entirely, defaulting to '/#'")
	}

	if baseTopic[0] == '/' {
		startsWithSlash = true
	}		

	if baseTopic[lastIndex] == '/' {
		endsWithSlash = true
	}

	// Have to check the last character separately, because we won't iterate over it.
	err := checkChar(rune(baseTopic[lastIndex]))

	if err != nil {
		return "/#", err
	}

	for index, char := range baseTopic {
		if index == lastIndex {
			break
		}

		err := checkChar(char)

		if err != nil {
			return "/#", err
		}

		nextChar := baseTopic[index + 1]

		if char == '/' && nextChar == '/' {
			return "/#", errors.New("You cannot have two slashes following each other.")
		}

	}

	switch {
	case startsWithSlash && endsWithSlash:
		return fmt.Sprintf("%s#", baseTopic), nil
	case startsWithSlash && !endsWithSlash:
		return fmt.Sprintf("%s/#", baseTopic), nil
	case !startsWithSlash && endsWithSlash:
		return fmt.Sprintf("/%s#", baseTopic), nil
	case !startsWithSlash && !endsWithSlash:
		return fmt.Sprintf("/%s/#", baseTopic), nil
	}

	return "/#", errors.New("Infallible, but defaulting.")
}

func checkChar (char rune) error {
	switch char {
	case '#':
		return errors.New("Your topic may not contain the '#' character")
	case '*':
		return errors.New("Your topic may not contain the '*' character")
	}
	return nil
}
