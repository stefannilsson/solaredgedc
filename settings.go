package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	utilities "github.com/stefannilsson/solaredgedc/common"
)

const (
	DEFAULT_MQTT_QOS      = 1
	DEFAULT_POLL_INTERVAL = 15000
	DEFAULT_MODBUS_PORT   = 502
)

const (
	LOG_LEVEL_ERROR = iota
	LOG_LEVEL_WARNING
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
)

type LogFlags struct {
	logLevel int
}

type ModbusFlags struct {
	hostname     string // hostname and/or IP address to Modbus server.
	port         int    // default to : '502'
	slaveId      int    // normally '1' in case of
	pollInterval int64  // number of 'ms' between polls
}

type MqttFlags struct {
	uri      string // e.g. 'tcp://127.0.0.1:1883'
	clientId string
	username string
	password string
	qos      int
	topic    string
}

/*
	Parse application arguments/flags, and read optional ENVironment variables.
	App flags override ENVironment variables, which in turns overrides the default settings.
*/
func ParseArgumentsConfig() (*LogFlags, *ModbusFlags, *MqttFlags) {
	// init long-living app config variables w/ default settings.
	logging := LogFlags{logLevel: LOG_LEVEL_INFO}
	modbus := ModbusFlags{port: DEFAULT_MODBUS_PORT}
	mqtt := MqttFlags{}

	// Modbus config parsing
	envModbusHostname := os.Getenv("MODBUS_HOSTNAME")
	flagModbusHostname := flag.String("modbus_hostname", "", "Modbus TCP hostname/IP address")

	envModbusPort := os.Getenv("MODBUS_PORT")
	flagModbusPort := flag.Int("modbus_port", 0, "Modbus TCP port")

	envModbusSlaveId := os.Getenv("MODBUS_SLAVEID")
	flagModbusSlaveId := flag.Int("modbus_slaveid", 0, "Modbus TCP SlaveID")

	envModbusPollInterval := os.Getenv("MODBUS_POLLINTERVAL")
	flagModbusPollInterval := flag.Int64("modbus_pollinterval", -1, "Modbus Poll interval (number of 'ms' between registers polls)")

	// MQTT config parsing
	envMqttUri := os.Getenv("MQTT_URI")
	flagMqttUri := flag.String("mqtt_uri", "", "The broker URI. ex: tcp://10.10.1.1:1883")

	envMqttClientId := os.Getenv("MQTT_CLIENTID")
	flagMqttClientId := flag.String("mqtt_clientid", "", "The ClientID (optional)")

	envMqttUsername := os.Getenv("MQTT_USERNAME")
	flagMqttUsername := flag.String("mqtt_username", "", "The User (optional)")

	envMqttPassword := os.Getenv("MQTT_PASSWORD")
	flagMqttPassword := flag.String("mqtt_password", "", "The password")

	envMqttQos := os.Getenv("MQTT_QOS")
	flagMqttQos := flag.Int("mqtt_qos", -1, "The Quality of Service {0,1,2} (default 1)")

	envMqttTopic := os.Getenv("MQTT_TOPIC")
	flagMqttTopic := flag.String("mqtt_topic", "", "The topic name to/from which to publish/subscribe")

	// Log config parsing
	envLogLevel := os.Getenv("LOG_LEVEL") // {DEBUG, INFO, WARNING, ERROR}
	flagLog := flag.String("log_level", "", "Log level - {DEBUG, INFO, WARNING, ERROR}")

	flag.Parse()

	// Modbus :: Hostname selection
	if *flagModbusHostname != "" {
		modbus.hostname = *flagModbusHostname
	} else if envModbusHostname != "" {
		modbus.hostname = envModbusHostname
	} else {
		panic("No Modbus TCP hostname provided.")
	}

	// Modbus :: Port selection
	if *flagModbusPort != 0 {
		modbus.port = *flagModbusPort
	} else if envModbusPort != "" {
		modbus.port, _ = strconv.Atoi(envModbusPort)
	}

	// Modbus :: Slave id selection
	if *flagModbusSlaveId != 0 {
		modbus.slaveId = *flagModbusSlaveId
	} else if envModbusSlaveId != "" {
		modbus.slaveId, _ = strconv.Atoi(envModbusSlaveId)
	}

	// Modbus :: Slave id selection
	if *flagModbusPollInterval != -1 {
		modbus.pollInterval = *flagModbusPollInterval
	} else if envModbusPollInterval != "" {
		envModbusPollIntervalInt64, _ := strconv.Atoi(envModbusPollInterval)
		modbus.pollInterval = int64(envModbusPollIntervalInt64)
	} else {
		modbus.pollInterval = DEFAULT_POLL_INTERVAL // default is to poll every 15 second.
	}

	// Log level selection
	switch strings.ToUpper(*flagLog) {
	case "DEBUG":
		logging.logLevel = LOG_LEVEL_DEBUG
	case "INFO":
		logging.logLevel = LOG_LEVEL_INFO
	case "WARNING":
		logging.logLevel = LOG_LEVEL_WARNING
	case "ERROR":
		logging.logLevel = LOG_LEVEL_ERROR
	default:
		switch envLogLevel {
		case "DEBUG":
			logging.logLevel = LOG_LEVEL_DEBUG
		case "INFO":
			logging.logLevel = LOG_LEVEL_INFO
		case "WARNING":
			logging.logLevel = LOG_LEVEL_WARNING
		case "ERROR":
			logging.logLevel = LOG_LEVEL_ERROR
		default:
			logging.logLevel = LOG_LEVEL_WARNING
		}
	}

	// MQTT :: URI select
	if *flagMqttUri != "" {
		mqtt.uri = *flagMqttUri
	} else if envMqttUri != "" {
		mqtt.uri = envMqttUri
	} else {
		panic("No MQTT hostname provided.")
	}

	// MQTT :: ClientId select
	if *flagMqttClientId != "" {
		mqtt.clientId = *flagMqttClientId
	} else if envMqttClientId != "" {
		mqtt.clientId = envMqttClientId
	} else {
		mqtt.clientId = fmt.Sprintf("auto-%s", utilities.RandomString(12)) // Default to random /[a-z0-9]{12}/
	}

	// MQTT :: Username select.
	if *flagMqttUsername != "" {
		mqtt.username = *flagMqttUsername
	} else if envMqttUsername != "" {
		mqtt.username = envMqttUsername
	}

	// MQTT :: Password select
	if *flagMqttPassword != "" {
		mqtt.password = *flagMqttPassword
	} else if envMqttPassword != "" {
		mqtt.password = envMqttPassword
	}

	// MQTT :: QoS select
	if *flagMqttQos != -1 {
		switch *flagMqttQos {
		case 0, 1, 2:
			mqtt.qos = *flagMqttQos
		default:
			panic("Unknown QoS specified.")
		}
	} else if envMqttQos != "" {
		qos, _ := strconv.Atoi(envMqttQos)
		switch qos {
		case 0, 1, 2:
			mqtt.qos = qos
		default:
			panic("Unknown QoS specified.")
		}
	} else {
		mqtt.qos = DEFAULT_MQTT_QOS
	}

	// MQTT :: Topic select
	if *flagMqttTopic != "" {
		mqtt.topic = *flagMqttTopic
	} else if envMqttTopic != "" {
		mqtt.topic = envMqttTopic
	} else {
		panic("No MQTT topic provided.")
	}

	return &logging, &modbus, &mqtt
}
