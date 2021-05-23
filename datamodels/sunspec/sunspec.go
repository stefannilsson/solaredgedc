package sunspec

// SunSpec register addresses:
// https://www.solaredge.com/sites/default/files/sunspec-implementation-technical-note.pdf

const (
	Dt_uint16 = iota
	Dt_uint32
	Dt_int16
	Dt_string
	Dt_acc32
)

type ModbusAddress struct {
	// E.g. 40000
	Address uint16

	// Only needed for 'Type: string'
	Size uint16

	/**
	* <b>Dt_uint16</b>
	* <b>Dt_uint32</b>
	* <b>Dt_int16</b>
	* <b>Dt_string</b>
	* <b>Dt_acc32</b>
	 */
	Type int

	Value interface{}
}

// Inverter Statuses
const (
	Ivs_I_STATUS_OFF           = 1
	Ivs_I_STATUS_SLEEPING      = 2
	Ivs_I_STATUS_STARTING      = 3
	Ivs_I_STATUS_MPPT          = 4
	Ivs_I_STATUS_THROTTLED     = 5
	Ivs_I_STATUS_SHUTTING_DOWN = 6
	Ivs_I_STATUS_FAULT         = 7
	Ivs_I_STATUS_STANDBY       = 8
)

/*func GetModbusRegisters() map[string]ModbusAddress {
	return Registers
}*/

// const registers map[string]ModbusAddress =
var Registers = map[string]ModbusAddress{
	/** [SUNSPEC : COMMON BLOCK] **/
	/*
		// Value = "SunS" (0x53756e53). Uniquely identifies this as a SunSpec MODBUS Map
		"C_SunSpec_ID": {Address: 40000, Type: Dt_uint32},

		// Value = 0x0001. Uniquely identifies this as a SunSpec Common Model Block
		"C_SunSpec_DID": {Address: 40002, Type: Dt_uint16},

		// 65 = Length of block in 16-bit registers
		"C_SunSpec_Length": {Address: 40003, Type: Dt_uint16},

		// Value Registered with SunSpec = "SolarEdge"
		"C_Manufacturer": {Address: 40004, Size: 32, Type: Dt_string},

		// SolarEdge Specific Value
		"C_Model": {Address: 40020, Size: 32, Type: Dt_string},

		// SolarEdge Specific Value
		"C_Version": {Address: 40044, Size: 16, Type: Dt_string},

		// SolarEdge Unique Value
		"C_SerialNumber": {Address: 40052, Size: 32, Type: Dt_string},

		// MODBUS Unit ID
		"C_DeviceAddress": {Address: 40068, Type: Dt_uint16},
		/*
		/** END of [SUNSPEC : COMMON BLOCK] **/

	"C_SerialNumber": {Address: 40052, Size: 4, Type: Dt_string}, // Default: Size: 32

	/** [SolarEdge Specific Registers] **/
	// 101 = single phase, 102 = split phase, 103 = three phase
	"C_SunSpec_DID": {Address: 40069, Type: Dt_uint16},

	// 50 = Length of model block
	"C_SunSpec_Length": {Address: 40070, Type: Dt_uint16},

	// AC Total Current value
	"I_AC_Current": {Address: 40071, Type: Dt_uint16},

	// AC Phase A Current value
	"I_AC_CurrentA": {Address: 40072, Type: Dt_uint16},

	// AC Phase B Current value
	"I_AC_CurrentB": {Address: 40073, Type: Dt_uint16},

	// AC Phase C Current value
	"I_AC_CurrentC": {Address: 40074, Type: Dt_uint16},

	// AC Current scale factor
	"I_AC_Current_SF": {Address: 40075, Type: Dt_int16},

	// AC Voltage Phase AB value
	"I_AC_VoltageAB": {Address: 40076, Type: Dt_uint16},

	// AC Voltage Phase BC value
	"I_AC_VoltageBC": {Address: 40077, Type: Dt_uint16},

	// AC Voltage Phase CA value
	"I_AC_VoltageCA": {Address: 40078, Type: Dt_uint16},

	// AC Voltage Phase A to N value
	"I_AC_VoltageAN": {Address: 40079, Type: Dt_uint16},

	// AC Voltage Phase B to N value
	"I_AC_VoltageBN": {Address: 40080, Type: Dt_uint16},

	// AC Voltage Phase C to N value
	"I_AC_VoltageCN": {Address: 40081, Type: Dt_uint16},

	// AC Voltage scale factor
	"I_AC_Voltage_SF": {Address: 40082, Type: Dt_int16},

	// AC Power value
	"I_AC_Power": {Address: 40083, Type: Dt_int16},

	// AC Power scale factor
	"I_AC_Power_SF": {Address: 40084, Type: Dt_int16},

	// AC Frequency value
	"I_AC_Frequency": {Address: 40085, Type: Dt_uint16},

	// Scale factor
	"I_AC_Frequency_SF": {Address: 40086, Type: Dt_int16},

	// Apparent Power
	"I_AC_VA": {Address: 40087, Type: Dt_int16},

	// Scale factor
	"I_AC_VA_SF": {Address: 40088, Type: Dt_int16},

	// Reactive Power
	"I_AC_VAR": {Address: 40089, Type: Dt_int16},

	// Scale factor
	"I_AC_VAR_SF": {Address: 40090, Type: Dt_int16},

	// Power Factor (%)
	"I_AC_PF": {Address: 40091, Type: Dt_int16},

	// Scale factor
	"I_AC_PF_SF": {Address: 40092, Type: Dt_int16},

	// AC Lifetime Energy production (WattHours)
	"I_AC_Energy_WH": {Address: 40093, Type: Dt_acc32},

	// Scale factor
	"I_AC_Energy_WH_SF": {Address: 40095, Type: Dt_int16}, // Data type typo 'uint16' where it _should_ be 'int16'

	// DC Current value (Amps)
	"I_DC_Current": {Address: 40096, Type: Dt_uint16},

	// Scale factor
	"I_DC_Current_SF": {Address: 40097, Type: Dt_int16},

	// DC Voltage value (Volts)
	"I_DC_Voltage": {Address: 40098, Type: Dt_uint16},

	// Scale factor
	"I_DC_Voltage_SF": {Address: 40099, Type: Dt_int16},

	// DC Power value (Watts)
	"I_DC_Power": {Address: 40100, Type: Dt_int16},

	// Scale factor
	"I_DC_Power_SF": {Address: 40101, Type: Dt_int16},

	// Heat Sink Temperature (Degrees C)
	"I_Temp_Sink": {Address: 40103, Type: Dt_int16},

	// Scale factor
	"I_Temp_SF": {Address: 40106, Type: Dt_int16},

	// Operating State
	"I_Status": {Address: 40107, Type: Dt_uint16},

	// Vendor-defined operating state and error codes. For error description, meaning and troubleshooting, refer to the SolarEdge Installation Guide.
	"I_Status_Vendor": {Address: 40108, Type: Dt_uint16},

	/** End of [SolarEdge Specific Registers] **/
}
