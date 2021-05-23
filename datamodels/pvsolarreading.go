package models

// Inverter Device Status Values
type InverterStatus int

const (
	// Off
	ivsSTATUS_OFF = 1

	// Sleeping (auto-shutdown) – Night mode
	ivsSTATUS_SLEEPING = 2

	// Grid Monitoring/wake-up
	ivsSTATUS_STARTING = 3

	// Inverter is ON and producing power
	ivsSTATUS_MPPT = 4

	// Production (curtailed)
	ivsSTATUS_THROTTLED = 5

	// Shutting down
	ivsSTATUS_SHUTTING_DOWN = 6

	// Fault
	ivsSTATUS_FAULT = 7

	// Maintenance/setup
	ivsSTATUS_STANDBY = 8
)

type PVSolarReading struct {

	// Identifier of component being measured.
	MeterId *string

	// AC Voltage Phase A/L1 to N value (Volts)
	AC_Voltage_L1_N *float64
	// AC Voltage Phase B/L2 to N value (Volts)
	AC_Voltage_L2_N *float64
	// AC Voltage Phase C/L3 to N value (Volts)
	AC_Voltage_L3_N *float64

	// AC Power (Watts)
	AC_Power *float64

	// AC Frequency (hz)
	AC_Frequency *float64

	// AC Active Power (VA)
	AC_VA *float64 // Apparent Power

	// AC Reactive Power (VAR)
	AC_VAR *float64 // Reactive Power

	// AC Power Factor (pf, 0.0-1.0)
	AC_PF *float64 // Power Factor

	// AC Lifetime Energy production
	AC_Energy_WH *float64

	// DC Current (Amps)
	DC_Current *float64

	// DC Voltage (Volts)
	DC_Voltage *float64

	// DC Power (Watts)
	DC_Power *float64

	// Inverter Heat Sink Temperatur (°C)
	Temp_Sink *float64

	/*
	 Inverter Status:
	 {ivsSTATUS_OFF,ivsSTATUS_SLEEPING,ivsSTATUS_STARTING,ivsSTATUS_MPPT,ivsSTATUS_THROTTLED,ivsSTATUS_SHUTTING_DOWN,ivsSTATUS_FAULT,ivsSTATUS_STANDBY}
	*/
	InverterStatus *uint16

	// Unix time in milliseconds of Modbus read
	Time *int64 `json:"time"`
}
