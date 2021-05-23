# SolarEdge Modbus TCP Poller (and MQTT publisher)
> A Modbus poller for SolarEdge PV Solar Inverters.

# Application flow
1) Poll data via Modbus TCP at regular interval (`MODBUS_POLLINTERVAL`)
2) Mapp the [SunSpec](./datamodels/sunspec/sunspec.go) data into a more user-friendly [PVSolarReading](./datamodels/pvsolarreading.go) structure.
3) Publish mapped data to the provided MQTT broker (`MQTT_URI`) & topic (`MQTT_TOPIC`).


# Usage
## Environment variables available:
```shell
MODBUS_HOSTNAME=192.168.0.100
MODBUS_PORT=1502
MODBUS_SLAVEID=1
MODBUS_POLLINTERVAL=5000
LOG_LEVEL=INFO
MQTT_URI=tcp://iot.eclipse.org:1883
MQTT_USERNAME=modbuspublisher
MQTT_PASSWORD=h4ck3rPassw0rd
MQTT_TOPIC=pvsolar/7E16A12F
MQTT_QOS=1
```

## Sample MQTT data
```json
{
    "MeterId": "7E16A12F",
    "AC_Voltage_L1_N": 229,
    "AC_Voltage_L2_N": 238.2,
    "AC_Voltage_L3_N": 238,
    "AC_Power": 8482,
    "AC_Frequency": 50.07,
    "AC_VA": 8493,
    "AC_VAR": -516.27,
    "AC_PF": -99.809,
    "AC_Energy_WH": 15445719,
    "DC_Current": 11.08,
    "DC_Voltage": 776.53,
    "DC_Power": 8614,
    "Temp_Sink": 52.29,
    "InverterStatus": 5,
    "time": 1621811112386
}
```

## License
[MIT](https://github.com/stefannilsson/solaredgedc/blob/master/LICENSE)
