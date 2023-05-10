package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	MQTTClient "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	utilities "github.com/stefannilsson/solaredgedc/common"
	mapping "github.com/stefannilsson/solaredgedc/datamapping"
	logger "github.com/stefannilsson/solaredgedc/logger"
	modbus "github.com/stefannilsson/solaredgedc/poller"
	mqtt "github.com/stefannilsson/solaredgedc/publisher"
)

const (
	GRACEFUL_SHUTDOWN_TIMEOUT_MS = 2000  // ms to let modbus poller / mqtt publisher to gracefully disconnect.
	DELAY_UNSUCCESSFUL_POLLS_MS  = 1000  // ms to be delayed between attempts if no successful modbus reads finished
	DELAY_FAILING_MODBUS_LINK_MS = 40000 // ms to be delayed before modbus reconnect attempt if modbus connection fails
)

var errorLog *logrus.Entry
var infoLog *logrus.Entry

func main() {
	// Allow application to be profiled via app argument '-trace'
	traceMode := flag.Bool("trace", false, "Trace application and write trace info to trace-*.out")
	// Must be run before any Loggers get instansiated.
	_, modbusConfig, mqttConfig := ParseArgumentsConfig()
	//TODO: Make log level configurable via flags/env.

	// Enable tracing to file if we recieved the '-flag' argument.
	if *traceMode {
		utilities.StartTrace()
	}

	// Get an instance of the logger.
	errorLog, infoLog, _ = logger.GetLoggers("main")
	infoLog.Println("SolarEdge Data Collector started.")

	// Initialize and try connect MQTT publisher
	mqttClient := mqtt.NewTelemetryMqtt(&mqtt.MqttConfig{
		URI:      mqttConfig.uri,
		ClientId: mqttConfig.clientId,
		Username: mqttConfig.username,
		Password: mqttConfig.password,
		Qos:      mqttConfig.qos,
		Topic:    mqttConfig.topic,
	})

	// Initialize and try connect Modbus poller
	modbusClient := modbus.NewPoller(&modbus.ModbusConfiguration{
		Hostname: modbusConfig.hostname,
		Port:     modbusConfig.port,
		SlaveId:  modbusConfig.slaveId,
	})

	//
	HandleSigInt(mqttClient)

	// Let's keep on polling all Modbus registers - for ever and ever.
	// MQTT Publisher maintains its own internal buffer if MQTT connection is temporarily down.
	for {

		registerValues := modbus.PollRegisters(modbusClient)
		// TODO: Implement check and indicator from PollRegister(...) if total read time was more than X amount of ms. Could be an issue if some registers took a very long time to read.

		// if no successfully register values read, let's sleep for a second and try again.
		if len(*registerValues) == 0 {
			time.Sleep(DELAY_FAILING_MODBUS_LINK_MS * time.Millisecond)

			// Disconnect Modbus
			modbusClient.TCPClientHandler.Close()
			time.Sleep(GRACEFUL_SHUTDOWN_TIMEOUT_MS * time.Millisecond)

			// Initialize and try connect Modbus poller
			modbusClient = modbus.NewPoller(&modbus.ModbusConfiguration{
				Hostname: modbusConfig.hostname,
				Port:     modbusConfig.port,
				SlaveId:  modbusConfig.slaveId,
			})
			time.Sleep(DELAY_UNSUCCESSFUL_POLLS_MS * time.Millisecond)
			continue
		}

		// key/value map
		// key : Modbus/SunSpec Register name
		// value : Scaled (in case of a numeric data type). Possible data types: {int16, uint16, uint32, string}
		// scaledValues := modbuspoller.ModbusRegistries{}
		parsedValues := mapping.ParseValues(registerValues)

		// Populate JSON with successfully read registers mapped to standard model.
		json := mapping.SerializeToJson(parsedValues)

		// publish JSON to MQTT broker...
		// (if MQTT broker is currently down, we'll use Paho MQTT library's internal buffer to send messages once online again.)
		mqttClient.Publish(mqttConfig.topic, byte(mqttConfig.qos), false, json)

		// and wait for some time before polling registers again.
		time.Sleep(time.Duration(modbusConfig.pollInterval) * time.Millisecond)
	}
}

func HandleSigInt(mqttClient MQTTClient.Client) {
	// give modbus & mqtt client some time to gracefully disconnect in case of CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// sig is a ^C, handle it
		errorLog.Errorf("Shutting down... Exiting in %v ms.", GRACEFUL_SHUTDOWN_TIMEOUT_MS)

		// Disconnect MQTT
		mqttClient.Disconnect(500)

		// Let's give clients some time to
		time.Sleep(GRACEFUL_SHUTDOWN_TIMEOUT_MS * time.Millisecond)

		os.Exit(-1)
	}()
}
