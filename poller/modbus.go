package modbus

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	MODBUS "github.com/goburrow/modbus"
	utilities "github.com/stefannilsson/solaredgedc/common"
	sunspec "github.com/stefannilsson/solaredgedc/datamodels/sunspec"
	"github.com/stefannilsson/solaredgedc/logger"
)

type ModbusConfiguration struct {
	Hostname          string
	Port              int
	SlaveId           int
	ConnectionTimeout int
}

type ModbusClient struct {
	Handler          MODBUS.Client
	TCPClientHandler *MODBUS.TCPClientHandler
}

type ModbusRegisters map[string]interface{}

const (
	SECONDS_BETWEEN_RECONNECTS = 10
)

var errorLog *logrus.Entry
var infoLog *logrus.Entry
var debugLog *logrus.Entry

func NewPoller(config *ModbusConfiguration) *ModbusClient {
	// Only log the warning severity or above.
	errorLog, infoLog, debugLog = logger.GetLoggers("modbus")

	handler := MODBUS.NewTCPClientHandler(fmt.Sprintf("%s:%d", config.Hostname, config.Port))
	handler.Timeout = SECONDS_BETWEEN_RECONNECTS * time.Second
	handler.SlaveId = byte(config.SlaveId)

	for {
		err := handler.Connect()
		if err != nil {
			errorLog.Println("TCP connection could not be established. Please check Modbus configuration.")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	client := MODBUS.NewClient(handler)

	// Create an return our _own_ Modbus client and its state.
	modbusClient := &ModbusClient{Handler: client, TCPClientHandler: handler}

	infoLog.Println("Modbus TCP connection successfully established.")
	return modbusClient
}

func PollRegisters(client *ModbusClient) *ModbusRegisters {

	// Successfully read Modbus registers (to be scaled later with each registers *_SF field)
	var readValues ModbusRegisters = ModbusRegisters{}

	for key, element := range sunspec.Registers {
		switch element.Type {
		case sunspec.Dt_uint16:
			result, err := client.Handler.ReadHoldingRegisters(element.Address, 1)
			if err != nil {
				errorLog.Printf("Failed to retrieve Modbus register '%s'", key)
				continue
			}
			readValues[key] = utilities.BytesToUInt16(result)
		case sunspec.Dt_int16:
			result, err := client.Handler.ReadHoldingRegisters(element.Address, 1)
			if err != nil {
				errorLog.Printf("Failed to retrieve Modbus register '%s'", key)
				continue
			}
			readValues[key] = utilities.BytesToInt16(result)
		case sunspec.Dt_uint32:
			result, err := client.Handler.ReadHoldingRegisters(element.Address, 2)
			if err != nil {
				errorLog.Printf("Failed to retrieve Modbus register '%s'", key)
				continue
			}
			readValues[key] = utilities.BytesToUint32(result)
		case sunspec.Dt_acc32:
			result, err := client.Handler.ReadHoldingRegisters(element.Address, 2)
			if err != nil {
				errorLog.Printf("Failed to retrieve Modbus register '%s'", key)
				continue
			}
			readValues[key] = utilities.BytesToUint32(result)
		case sunspec.Dt_string:
			result, err := client.Handler.ReadHoldingRegisters(element.Address, element.Size)
			if err != nil {
				errorLog.Printf("Failed to retrieve Modbus register '%s'", key)
				continue
			}
			readValues[key] = string(result)
		default:
			errorLog.Println(fmt.Sprintf("UNKNOWN datatype: int(%d)", element.Type))
		}
	}

	if len(readValues) == 0 {
		errorLog.Println("No values successfully read from registers.")
	}

	return &readValues
}
