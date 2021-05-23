package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/stefannilsson/solaredgedc/logger"
)

type MqttConfig struct {
	URI          string
	Port         int
	ClientId     string
	Username     string
	Password     string
	CleanSession bool
	Qos          int
	Topic        string
	PollInterval int
}

func NewTelemetryMqtt(mqttConfig *MqttConfig) MQTT.Client {
	errorLog, _, _ := logger.GetLoggers("mqtt")

	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttConfig.URI)
	opts.SetClientID(mqttConfig.ClientId)
	opts.SetUsername(mqttConfig.Username)
	opts.SetPassword(mqttConfig.Password)
	opts.SetCleanSession(mqttConfig.CleanSession)
	//TODO: Implement (optional) file based buffer
	/*if *store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(*store))
	}*/

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		errorLog.Println(token.Error())
		panic(token.Error())
	}

	return client
}
