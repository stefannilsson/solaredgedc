package mapping

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"

	utilities "github.com/stefannilsson/solaredgedc/common"
	models "github.com/stefannilsson/solaredgedc/datamodels"
	logger "github.com/stefannilsson/solaredgedc/logger"
	modbus "github.com/stefannilsson/solaredgedc/poller"
)

var errorLog, infoLog, _ = logger.GetLoggers("mapping")

/* Read Modbus registers, then cast to proper types and scale values accordingly */
func ParseValues(registerValues *modbus.ModbusRegisters) map[string]interface{} {
	scaledValues := map[string]interface{}{}

	// let's timestamp the readings asap
	scaledValues["Time"] = utilities.TimeNowInUnixMs()

	values := *registerValues
	regexpEndWithSF := regexp.MustCompile("_SF$")

	for key, val := range values {
		// No need to "scale" the ScaleFactor itself.
		if regexpEndWithSF.MatchString(key) {
			continue
		}

		// No need to process 'strings' furthers.
		switch val.(type) {
		case string:
			scaledValues[key] = val
			continue
		}

		// Scale Factor key ref
		var sfKey string

		// Custom Scale Factor fields (not applying the standard form "{REGISTERNAME}_SF" )
		switch key {
		case "I_AC_Current", "I_AC_CurrentA", "I_AC_CurrentB", "I_AC_CurrentC":
			sfKey = "I_AC_Current_SF"
		case "I_AC_VoltageAB", "I_AC_VoltageBC", "I_AC_VoltageCA", "I_AC_VoltageAN", "I_AC_VoltageBN", "I_AC_VoltageCN":
			sfKey = "I_AC_Voltage_SF"
		case "I_Temp_Sink":
			sfKey = "I_Temp_SF"
		default:
			sfKey = fmt.Sprintf("%s_SF", key)
		}

		if sf, found := values[sfKey]; found {
			// if we've recieved a field with a corresponding scale factor register, we'll assume a numeric data type.
			var float64Value float64
			switch val.(type) {
			case int16:
				float64Value = float64(val.(int16))
			case uint16:
				float64Value = float64(val.(uint16))
			case uint32:
				float64Value = float64(val.(uint32))
			}

			var scaleFactor int16
			switch sf.(type) {
			case int16:
				scaleFactor = sf.(int16)
			default:
				panic("Unhandled ScaleFactory data type")
			}

			// let's scale read value with it's corresponding scale factor according to SunSpec doc.
			scaledValue := float64Value * math.Pow10(int(scaleFactor))
			scaledValues[key] = scaledValue
		} else {
			// just store the read value as scaledValue if no corresponding _SF found.
			scaledValues[key] = val
		}
	}

	return scaledValues
}

/* Move successfully parsed values into the common PVSolar */
func SerializeToJson(parsedValues map[string]interface{}) []byte {
	// Map to our standard PV Solar data model.
	C_SerialNumber := new(string)
	I_AC_VoltageAN := new(float64)
	I_AC_VoltageBN := new(float64)
	I_AC_VoltageCN := new(float64)
	I_AC_Power := new(float64)
	I_AC_Frequency := new(float64)
	I_AC_VA := new(float64)
	I_AC_VAR := new(float64)
	I_AC_PF := new(float64)
	I_AC_Energy_WH := new(float64)
	I_DC_Current := new(float64)
	I_DC_Voltage := new(float64)
	I_DC_Power := new(float64)
	I_Temp_Sink := new(float64)
	I_Status := new(uint16)
	Time := new(int64)

	// cast successfully scaled values into their right data type and prepare for JSON marshalling.
	if value, ok := parsedValues["C_SerialNumber"]; ok {
		*C_SerialNumber = value.(string)
	}

	if value, ok := parsedValues["I_AC_VoltageAN"]; ok {
		*I_AC_VoltageAN = value.(float64)
	}

	if value, ok := parsedValues["I_AC_VoltageBN"]; ok {
		*I_AC_VoltageBN = value.(float64)
	}

	if value, ok := parsedValues["I_AC_VoltageCN"]; ok {
		*I_AC_VoltageCN = value.(float64)
	}

	if value, ok := parsedValues["I_AC_Power"]; ok {
		*I_AC_Power = value.(float64)
	}

	if value, ok := parsedValues["I_AC_Frequency"]; ok {
		*I_AC_Frequency = value.(float64)
	}

	if value, ok := parsedValues["I_AC_VA"]; ok {
		*I_AC_VA = value.(float64)
	}

	if value, ok := parsedValues["I_AC_VAR"]; ok {
		*I_AC_VAR = value.(float64)
	}

	if value, ok := parsedValues["I_AC_PF"]; ok {
		*I_AC_PF = value.(float64)
	}

	if value, ok := parsedValues["I_AC_Energy_WH"]; ok {
		*I_AC_Energy_WH = value.(float64)
	}

	if value, ok := parsedValues["I_DC_Current"]; ok {
		*I_DC_Current = value.(float64)
	}

	if value, ok := parsedValues["I_DC_Voltage"]; ok {
		*I_DC_Voltage = value.(float64)
	}

	if value, ok := parsedValues["I_DC_Power"]; ok {
		*I_DC_Power = value.(float64)
	}

	if value, ok := parsedValues["I_Temp_Sink"]; ok {
		*I_Temp_Sink = value.(float64)
	}

	if value, ok := parsedValues["I_Status"]; ok {
		*I_Status = value.(uint16)
	}

	if value, ok := parsedValues["Time"]; ok {
		*Time = value.(int64)
	}

	// map to common data model
	pvRead := &models.PVSolarReading{
		MeterId:         C_SerialNumber,
		AC_Voltage_L1_N: I_AC_VoltageAN,
		AC_Voltage_L2_N: I_AC_VoltageBN,
		AC_Voltage_L3_N: I_AC_VoltageCN,
		AC_Power:        I_AC_Power,
		AC_Frequency:    I_AC_Frequency,
		AC_VA:           I_AC_VA,
		AC_VAR:          I_AC_VAR,
		AC_PF:           I_AC_PF,
		AC_Energy_WH:    I_AC_Energy_WH,
		DC_Current:      I_DC_Current,
		DC_Voltage:      I_DC_Voltage,
		DC_Power:        I_DC_Power,
		Temp_Sink:       I_Temp_Sink,
		InverterStatus:  I_Status,
		Time:            Time,
	}

	// serialize to JSON payload
	json, err := json.Marshal(pvRead)
	if err != nil {
		errorLog.Errorln(err.Error())
	}

	return json
}
