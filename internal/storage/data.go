package storage

// LSN50V2JSON parse codec
type LSN50V2JSON struct {
	Battery     float32 `json:"BatV"`
	ADC         float32 `json:"ADC_CH0V"`
	Humidity    string  `json:"Hum_SHT"`
	Temperature string  `json:"TempC_SHT"`
	Exp         string  `json:"Ext_sensor"`
	TempDS      string  `json:"TempC_DS"`
	Status      string  `json:"Digital_IStatus"`
	DoorStatus  string  `json:"Door_status"`
	Trigger     string  `json:"EXTI_Trigger"`
	TempC1      string  `json:"TempC1"`
	WorkMode    string  `json:"Work_mode"`
}

// LSE01JSON parse soil codec
type LSE01JSON struct {
	Battery         float32 `json:"BatV"`
	TempC           string  `json:"TempC_DS18B20"`
	ConductSoil     float32 `json:"conduct_SOIL"`
	TemperatureSoil string  `json:"temp_SOIL"`
	WaterSoil       string  `json:"water_SOIL"`
}

type LDS01JSON struct {
	Battery              float32 `json:"BatV"`
	DoorStatus           int64   `json:"door_open_status"`
	DoorOpenTimes        int64   `json:"door_open_times"`
	LastDoorOpenDuration int64   `json:"last_door_open_duration"`
}
type LWL01JSON struct {
	Battery               float32 `json:"BatV"`
	WaterStatus           int64   `json:"WATER_LEAK_STATUS"`
	WaterLeekTimes        int64   `json:"WATER_LEAK_TIMES"`
	LastWaterLeekDuration int64   `json:"LAST_WATER_LEAK_DURATION"`
}

// Milesight EM-300-TH
type EM300THJSON struct {
	Battery     float32 `json:"battery"`
	Humidity    float32 `json:"humidity"`
	Temperature float32 `json:"temperature"`
}

// Milesight Ws191
type WS101JSON struct {
	Alarm int `json:"press"`
}

// Milesight EM300-ZLD
type EM300ZLDJSON struct {
	WaterLeek int64 `json:"water_leak"`
}