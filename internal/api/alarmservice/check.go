package alarmservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *AlarmServerAPI) CheckAlarm(ctx context.Context, req *als.CheckAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	currentTime := time.Now().Add(time.Hour * 3)

	var alarms []s.AlarmWithDates
	weekday := time.Now().Weekday() + 1

	err := sqlx.Select(db, &alarms, `select alrm.*, alrmDate.alarm_day, alrmDate.start_time, alrmDate.end_time  from alarm_refactor as alrm 
	inner join alarm_date_time alrmDate on alrm.id = alrmDate.alarm_id where dev_eui = $1
	and ( alrmDate.alarm_day = 0 or alrmDate.alarm_day = $2 ) and is_active = true`, req.Device.DevEui, int(weekday))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	switch req.Device.DeviceType {
	case 1:
		data := s.LSN50V2JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)
		for _, element := range alarms {
			if s.CheckAlarmTime(element) {
				if element.Temperature {
					temp, err := strconv.ParseFloat(data.Temperature, 32)
					if err != nil {
						fmt.Println("parse error")
					}
					err = s.CheckThreshold(element, float32(temp), *req.Device, "ısı", currentTime.Format("2006-01-02 15:04:05"), db)
					if err != nil {
						log.Println("CheckThreshold error")
					}
				} else if element.Humadity {
					temp, err := strconv.ParseFloat(data.Humidity, 32)
					if err != nil {
						fmt.Println("parse error")
					}
					err = s.CheckThreshold(element, float32(temp), *req.Device, "nem", currentTime.Format("2006-01-02 15:04:05"), db)
					if err != nil {
						log.Println("CheckThreshold error")
					}
				}
			}
		}

	case 2:
		data := s.LSE01JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)
		if data.TemperatureSoil != "0.00" && data.WaterSoil != "0.00" {
			for _, element := range alarms {
				if s.CheckAlarmTime(element) {
					if element.Temperature {
						temp, err := strconv.ParseFloat(data.TemperatureSoil, 32)
						if err != nil {
							fmt.Println("parse error")
						}
						err = s.CheckThreshold(element, float32(temp), *req.Device, "ısı", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("CheckThreshold error")
						}

					} else if element.Humadity {
						temp, err := strconv.ParseFloat(data.WaterSoil, 32)
						if err != nil {
							fmt.Println("parse error")
						}
						err = s.CheckThreshold(element, float32(temp), *req.Device, "nem", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("CheckThreshold error")
						}
					} else if element.Ec {
						err = s.CheckThreshold(element, float32(data.ConductSoil), *req.Device, "ec", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("CheckThreshold error")
						}
					}
				}

			}
		}

	case 3:
		data := s.LDS01JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)
		for _, element := range alarms {
			if element.Door {
				if data.DoorStatus == 1 {
					if s.CheckAlarmTime(element) {
						err = s.ExecuteAlarm(element, 0, *req.Device, "door", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("doorAlarm error")
						}
					}
				}
			}
		}

	case 4:
		data := s.LWL01JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)

		for _, element := range alarms {
			if element.WaterLeak {
				if data.WaterStatus == 1 {
					if s.CheckAlarmTime(element) {
						err = s.ExecuteAlarm(element, 0, *req.Device, "kacak", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("waterLeakAlarm error")
						}
					}
				}
			}
		}

	case 10:
		data := s.LWL01JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)

		for _, element := range alarms {
			if element.WaterLeak {
				if data.WaterStatus == 1 {
					if s.CheckAlarmTime(element) {
						err = s.ExecuteAlarm(element, 0, *req.Device, "acil durum", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("emergencyAlarm error")
						}
					}
				}
			}
		}

	case 12:
		data := s.EM300THJSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)
		if data.Temperature != 0 && data.Humidity != 0 {
			for _, element := range alarms {
				if s.CheckAlarmTime(element) {
					if element.Temperature {
						err = s.CheckThreshold(element, float32(data.Temperature), *req.Device, "ısı", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("CheckThreshold error")
						}
					} else if element.Humadity {
						err = s.CheckThreshold(element, float32(data.Humidity), *req.Device, "nem", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("CheckThreshold error")
						}
					}
				}

			}
		}
	case 14:
		data := s.WS101JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)

		for _, element := range alarms {

			if element.WaterLeak {
				if data.Alarm == 1 {
					if s.CheckAlarmTime(element) {
						err = s.ExecuteAlarm(element, 0, *req.Device, "button", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("alarmButton error")
						}
					}
				}

			}
		}
	case 18:
		data := s.EM300ZLDJSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)

		for _, element := range alarms {
			if element.WaterLeak {
				if data.WaterLeek == 1 {
					if s.CheckAlarmTime(element) {
						err = s.ExecuteAlarm(element, 0, *req.Device, "kacak", currentTime.Format("2006-01-02 15:04:05"), db)
						if err != nil {
							log.Println("waterLeakAlarm error")
						}
					}
				}
			}
		}
	}

	return &emptypb.Empty{}, nil
}
