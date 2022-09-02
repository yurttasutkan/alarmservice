package alarm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

//AlarmServerAPI implements the Alarm server API.
type AlarmServerAPI struct {
}

//Creates a new AlarmServerAPI
func NewAlarmServerAPI() *AlarmServerAPI {
	return &AlarmServerAPI{}
}

type Alarm struct {
	ID                int64   `db:"id"`
	DevEui            string  `db:"dev_eui"`
	MinTreshold       float32 `db:"min_treshold"`
	MaxTreshold       float32 `db:"max_treshold"`
	Sms               bool    `db:"sms"`
	Email             bool    `db:"email"`
	Notification      bool    `db:"notification"`
	Temperature       bool    `db:"temperature"`
	Humadity          bool    `db:"humadity"`
	Ec                bool    `db:"ec"`
	Door              bool    `db:"door"`
	W_leak            bool    `db:"w_leak"`
	UserId            int64   `db:"user_id"`
	IpAddress         string  `db:"ip_address"`
	IsTimeLimitActive bool    `db:"is_time_limit_active"`
	AlarmStartTime    float32 `db:"alarm_start_time"`
	AlarmStopTime     float32 `db:"alarm_stop_time"`
	ZoneCategoryId    int64   `db:"zone_category"`
	IsActive          bool    `db:"is_active"`
}

type OrganizationAlarm struct {
	ID                int64   `db:"id"`
	DevEui            string  `db:"dev_eui"`
	MinTreshold       float32 `db:"min_treshold"`
	MaxTreshold       float32 `db:"max_treshold"`
	Sms               bool    `db:"sms"`
	Email             bool    `db:"email"`
	Notification      bool    `db:"notification"`
	Temperature       bool    `db:"temperature"`
	Humadity          bool    `db:"humadity"`
	Ec                bool    `db:"ec"`
	Door              bool    `db:"door"`
	W_leak            bool    `db:"w_leak"`
	UserId            int64   `db:"user_id"`
	IpAddress         string  `db:"ip_address"`
	IsTimeLimitActive bool    `db:"is_time_limit_active"`
	AlarmStartTime    float32 `db:"alarm_start_time"`
	AlarmStopTime     float32 `db:"alarm_stop_time"`
	ZoneCategoryId    int64   `db:"zone_category"`
	IsActive          bool    `db:"is_active"`
	ZoneName          string  `db:"zone_name"`
	DeviceName        string  `db:"device_name"`
	Username          string  `db:"username"`
}
type AlarmWithDates struct {
	ID                int64   `db:"id"`
	DevEui            string  `db:"dev_eui"`
	MinTreshold       float32 `db:"min_treshold"`
	MaxTreshold       float32 `db:"max_treshold"`
	Sms               bool    `db:"sms"`
	Notification      bool    `db:"notification"`
	Email             bool    `db:"email"`
	Temperature       bool    `db:"temperature"`
	Humadity          bool    `db:"humadity"`
	Ec                bool    `db:"ec"`
	Door              bool    `db:"door"`
	W_leak            bool    `db:"w_leak"`
	UserId            int64   `db:"user_id"`
	IpAddress         string  `db:"ip_address"`
	IsTimeLimitActive bool    `db:"is_time_limit_active"`
	ZoneCategoryId    int64   `db:"zone_category"`
	AlarmDay          int64   `db:"alarm_day"`
	AlarmStartTime2   float32 `db:"alarm_start_time"`
	AlarmStopTime2    float32 `db:"alarm_stop_time"`
	AlarmStartTime    float32 `db:"start_time"`
	AlarmEndTime      float32 `db:"end_time"`
	IsActive          bool    `db:"is_active"`
}

type AlarmDateFilter struct {
	ID             int64   `db:"id"`
	AlarmId        int64   `db:"alarm_id"`
	AlarmDay       int64   `db:"alarm_day"`
	AlarmStartTime float32 `db:"start_time"`
	AlarmEndTime   float32 `db:"end_time"`
}
type ColdRoomRestrictions struct {
	ID          int64  `db:"id"`
	DevEui      string `db:"dev_eui"`
	AlarmId     int64  `db:"alarm_id"`
	DefrostTime int64  `db:"defronst_time"`
	DefrostFrq  int64  `db:"defrost_frequency"`
	AlarmTime   int64  `db:"alarm_time"`
}
type AlarmLogs struct {
	DevEui         string    `db:"dev_eui"`
	MinTreshold    float32   `db:"min_treshold"`
	MaxTreshold    float32   `db:"max_treshold"`
	UserId         int64     `db:"user_id"`
	IpAddress      string    `db:"ip_address"`
	IsDeleted      int64     `db:"is_deleted"`
	Sms            bool      `db:"sms"`
	Temperature    bool      `db:"temperature"`
	Humadity       bool      `db:"humadity"`
	Ec             bool      `db:"ec"`
	Door           bool      `db:"door"`
	W_leak         bool      `db:"w_leak"`
	SubmissionDate time.Time `db:"submission_date"`
}

// AlarmFilters filters
type AlarmFilters struct {
	Limit  int    `db:"limit"`
	DevEui string `db:"dev_eui"`
	UserID int64  `db:"user_id"`
}

// SMSRequestBody ...
type SMSRequestBody struct {
	From      string `json:"from"`
	Text      string `json:"text"`
	To        string `json:"to"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// SQL function to convert filters to SQL line
func (f AlarmFilters) SQL() string {
	var filters []string

	if f.DevEui != "" {
		filters = append(filters, fmt.Sprint(" dev_eui =  '", f.DevEui+"'"))
	}
	filters = append(filters, fmt.Sprint(" and user_id = ", f.UserID))
	if f.Limit != 0 {
		filters = append(filters, fmt.Sprint(" LIMIT ", f.Limit))
	}

	return " where is_active = true and  " + strings.Join(filters, " ")
}

func (a *AlarmServerAPI) CheckAlarm(ctx context.Context, req *als.CheckAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	var device_type_id int64
	var zoneName string
	currentTime := time.Now().Add(time.Hour * 3)

	_ap := req.Application

	spID, err := uuid.FromString(_ap.ServiceProfileId)
	if err != nil {
		log.Printf("UUID parse error")
	}

	ap := s.Application{
		ID:                   _ap.Id,
		Name:                 _ap.Name,
		Description:          _ap.Description,
		OrganizationID:       _ap.OrganizationId,
		ServiceProfileID:     spID,
		PayloadCodec:         _ap.PayloadCodec,
		PayloadEncoderScript: _ap.PayloadEncoderScript,
		PayloadDecoderScript: _ap.PayloadDecoderScript,
		MQTTTLSCert:          _ap.MQTTTLSCert,
	}

	err = sqlx.Get(db, &device_type_id, `select id from public.device_type_tb  as dtt
	inner join public.device_profile as dp on  dp.name = dtt.device_profile_name 
	inner join public.device as d on d.device_profile_id = dp.device_profile_id
	where d.device_profile_id =  $1 limit 1`, req.Device.DeviceProfileId)
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}
	var alarms []AlarmWithDates
	weekday := time.Now().Weekday() + 1

	err = sqlx.Select(db, &alarms, `select alrm.*, alrmDate.alarm_day, alrmDate.start_time, alrmDate.end_time  from alarm_refactor as alrm 
	inner join alarm_date_time alrmDate on alrm.id = alrmDate.alarm_id where dev_eui = $1
	and ( alrmDate.alarm_day = 0 or alrmDate.alarm_day = $2 ) and is_active = true`, req.Device.DevEui, int(weekday))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	err = sqlx.Get(db, &zoneName, "select zone_name from zone where '\\x' || $1 = any(zone.devices)", req.Device.DevEui)
	if err != nil {
		fmt.Println("get zone anme error")
	}

	switch device_type_id {
	case 1:
		data := s.LSN50V2JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)

		for _, element := range alarms {
			if element.Temperature {
				temp, err := strconv.ParseFloat(data.Temperature, 32)
				if err != nil {
					fmt.Println("parse error")
				}
				err = checkAlarmTimeSchedule(element, float32(temp), req.DeviceName, ap, zoneName, "ısı", currentTime.Format("2006-01-02 15:04:05"), db)
				if err!=nil{
					log.Println("checkAlarmTimeSchedule error")
				}
			} else if element.Humadity {
				temp, err := strconv.ParseFloat(data.Humidity, 32)
				if err != nil {
					fmt.Println("parse error")
				}
				err = checkAlarmTimeSchedule(element, float32(temp), req.DeviceName, ap, zoneName, "nem", currentTime.Format("2006-01-02 15:04:05"), db)
				if err!=nil{
					log.Println("checkAlarmTimeSchedule error")
				}
			}
		}
	case 2:
		data := s.LSE01JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)
		if data.TemperatureSoil != "0.00" && data.WaterSoil != "0.00" {
			for _, element := range alarms {
				if element.Temperature {
					temp, err := strconv.ParseFloat(data.TemperatureSoil, 32)
					if err != nil {
						fmt.Println("parse error")
					}
					err = checkAlarmTimeSchedule(element, float32(temp), req.DeviceName, ap, zoneName, "ısı", currentTime.Format("2006-01-02 15:04:05"), db)
					if err!=nil{
						log.Println("checkAlarmTimeSchedule error")
					}
				} else if element.Humadity {
					temp, err := strconv.ParseFloat(data.WaterSoil, 32)
					if err != nil {
						fmt.Println("parse error")
					}
					err = checkAlarmTimeSchedule(element, float32(temp), req.DeviceName, ap, zoneName, "nem", currentTime.Format("2006-01-02 15:04:05"), db)
					if err!=nil{
						log.Println("checkAlarmTimeSchedule error")
					}
					} else if element.Ec {
					err = checkAlarmTimeSchedule(element, float32(data.ConductSoil), req.DeviceName, ap, zoneName, "ec", currentTime.Format("2006-01-02 15:04:05"), db)
					if err!=nil{
						log.Println("checkAlarmTimeSchedule error")
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
					if checkAlarmTime(element) {
						err = doorAlarm(element, req.DeviceName, zoneName, "kapı", currentTime.Format("2006-01-02 15:04:05"))
						if err!=nil{
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
			if element.W_leak {
				if data.WaterStatus == 1 {
					if checkAlarmTime(element) {
						err = waterLeakAlarm(element, req.DeviceName, zoneName)
						if err!=nil{
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
			if element.W_leak  {
				if data.WaterStatus == 1 {
					if checkAlarmTime(element) {
						err = emergencyAlarm(element, req.DeviceName, zoneName)
						if err!=nil{
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

				if element.Temperature  {
					err = checkAlarmTimeSchedule(element, float32(data.Temperature), req.DeviceName, ap, zoneName, "ısı", currentTime.Format("2006-01-02 15:04:05"), db)
					if err!=nil{
						log.Println("checkAlarmTimeSchedule error")
					}
				} else if element.Humadity  {
					err = checkAlarmTimeSchedule(element, float32(data.Humidity), req.DeviceName, ap, zoneName, "nem", currentTime.Format("2006-01-02 15:04:05"), db)
					if err!=nil{
						log.Println("checkAlarmTimeSchedule error")
					}
				}
			}
		}
	case 14:
		data := s.WS101JSON{}
		json.Unmarshal([]byte(req.ObjectJSON), &data)

		for _, element := range alarms {

			if element.W_leak  {
				if data.Alarm == 1 {
					if checkAlarmTime(element) {
						err = alarmButton(element, req.DeviceName, zoneName)
						if err!=nil{
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
			if element.W_leak  {
				if data.WaterLeek == 1 {
					if checkAlarmTime(element) {
						err = waterLeakAlarm(element, req.DeviceName, zoneName)
						if err!=nil{
							log.Println("waterLeakAlarm error")
						}
					}
				}
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func doorAlarm(a AlarmWithDates, deviceName string, zonename string, alarmType string, date string) error {
	currentTime := time.Now().Add(time.Hour * 3)
	db := s.DB()
	var u s.User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return s.HandlePSQLError(s.Select, err, "alarm log insert error")
	}
	// n := FirebaseNotificationData{
	// 	Title: "Vaps",
	// 	Body:  zonename + " deki " + deviceName + " kapı sensörü açıldı",
	// }

	// for _, element := range u {
	// SendFirebaseNotification(u, n)

	notification := s.Notification{
		SenderId:   0,
		ReceiverId: a.UserId,
		Message:    date + " tarihinde " + zonename + " ortamındaki " + deviceName + " isimli sensör " + alarmType + " sensörü açıldı",
		CategoryId: 1,
		IsRead:     false,
		SendTime:   time.Now(),
		SenderIp:   "system",
		ReaderIp:   "",
		IsDeleted:  false,
		DeviceName: deviceName,
		DevEui:     a.DevEui,
	}
	err = CreateNotification(notification)
	if err!=nil{
		log.Println("CreateNotification error")
	}

	if a.Sms  {
		numbers := []string{u.PhoneNumber}
		numbersString := s.NumbersArrayToString(numbers)

		sms1N := s.OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + " deki " + deviceName + " sensörü açıldı"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email  {

		var user s.User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return s.HandlePSQLError(s.Select, err, "alarm log insert error")
		}

		s.SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+" deki "+deviceName+" sensörü açıldı")
	}

	if a.Notification  {
		n := s.FirebaseNotificationData{
			Title: "Vaps",
			Body:  currentTime.Format("2006-01-02 15:04:05") + " tarihinde " + zonename + " deki " + deviceName + " sensörü açıldı",
			Time:  300000,
			Delay: false,
		}

		s.SendFirebaseNotification(u, n)
	}
	return nil
}
func alarmButton(a AlarmWithDates, deviceName string, zonename string) error {
	db := s.DB()
	currentTime := time.Now().Add(time.Hour * 3)
	var u s.User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return s.HandlePSQLError(s.Select, err, "alarm log insert error")
	}
	// n := FirebaseNotificationData{
	// 	Title: "Vaps",
	// 	Body:  zonename + "deki" + deviceName + " sensöründe kaçak var",
	// }

	// for _, element := range u {
	// SendFirebaseNotification(u, n)
	if a.Sms  {
		numbers := []string{u.PhoneNumber}
		numbersString := s.NumbersArrayToString(numbers)

		sms1N := s.OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + " deki " + deviceName + " sensöründen çağrı var"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email  {

		var user s.User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return s.HandlePSQLError(s.Select, err, "alarm log insert error")
		}

		s.SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+" deki "+deviceName+" sensöründen çağrı var")
	}

	if a.Notification  {
		n := s.FirebaseNotificationData{
			Title: "Vaps",
			Body:  zonename + " deki " + deviceName + " sensöründen çağrı var" + " tarih: " + currentTime.Format("2006-01-02 15:04:05"),
			Time:  300000,
			Delay: false,
		}

		s.SendFirebaseNotification(u, n)
	}

	return nil
}
func waterLeakAlarm(a AlarmWithDates, deviceName string, zonename string) error {
	db := s.DB()
	currentTime := time.Now().Add(time.Hour * 3)
	var u s.User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return s.HandlePSQLError(s.Select, err, "alarm log insert error")
	}
	// n := FirebaseNotificationData{
	// 	Title: "Vaps",
	// 	Body:  zonename + "deki" + deviceName + " sensöründe kaçak var",
	// }

	// for _, element := range u {
	// SendFirebaseNotification(u, n)

	if a.Sms  {
		numbers := []string{u.PhoneNumber}
		numbersString := s.NumbersArrayToString(numbers)

		sms1N := s.OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + " deki " + deviceName + " sensöründe kaçak var"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email  {

		var user s.User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return s.HandlePSQLError(s.Select, err, "alarm log insert error")
		}

		s.SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+" deki "+deviceName+" sensöründe kaçak var")
	}

	if a.Notification  {
		n := s.FirebaseNotificationData{
			Title: "Vaps",
			Body:  currentTime.Format("2006-01-02 15:04:05") + " tarihinde " + zonename + " deki " + deviceName + " sensöründe kaçak var",
			Time:  300000,
			Delay: false,
		}

		s.SendFirebaseNotification(u, n)
	}

	return nil
}

func emergencyAlarm(a AlarmWithDates, deviceName string, zonename string) error {
	db := s.DB()
	currentTime := time.Now().Add(time.Hour * 3)
	var u s.User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return s.HandlePSQLError(s.Select, err, "alarm log insert error")
	}
	// n := FirebaseNotificationData{
	// 	Title: "Vaps",
	// 	Body:  zonename + "deki" + deviceName + " sensöründe acil durum var",
	// }

	// for _, element := range u {
	// SendFirebaseNotification(u, n)

	if a.Sms  {
		numbers := []string{u.PhoneNumber}
		numbersString := s.NumbersArrayToString(numbers)

		sms1N := s.OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + "deki" + deviceName + " sensöründe acil durum var"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email  {

		var user s.User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return s.HandlePSQLError(s.Select, err, "alarm log insert error")
		}

		s.SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+"deki"+deviceName+" sensöründe acil durum var")
	}

	if a.Notification  {
		n := s.FirebaseNotificationData{
			Title: "Vaps",
			Body:  currentTime.Format("2006-01-02 15:04:05") + " tarihinde " + zonename + "deki" + deviceName + " sensöründe acil durum var",
			Time:  300000,
			Delay: false,
		}

		s.SendFirebaseNotification(u, n)
	}
	return nil
}

func checkAlarmTimeSchedule(a AlarmWithDates, v float32, deviceName string, ap s.Application, zoneName string, alarmType string, date string, db sqlx.Ext) error {
	if a.IsTimeLimitActive {
		hours, minutes, _ := time.Now().Clock()
		result := strconv.Itoa(hours+3) + "." + strconv.Itoa(minutes)
		t, err := strconv.ParseFloat(result, 32)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		if a.AlarmEndTime > a.AlarmStartTime {
			if a.AlarmStartTime < float32(t) && float32(t) < a.AlarmEndTime {
				checkThreshold(a, v, deviceName, ap, zoneName, alarmType, date, db)
			}
		} else {
			if a.AlarmStartTime < float32(t) && float32(t) < 24 {
				checkThreshold(a, v, deviceName, ap, zoneName, alarmType, date, db)
			}
			if 0 < float32(t) && float32(t) < a.AlarmEndTime {
				checkThreshold(a, v, deviceName, ap, zoneName, alarmType, date, db)
			}
		}
	} else {
		checkThreshold(a, v, deviceName, ap, zoneName, alarmType, date, db)
	}
	return nil
}

func checkThreshold(a AlarmWithDates, v float32, deviceName string, ap s.Application, zoneName string, alarmType string, date string, db sqlx.Ext) error {
	if v < a.MinTreshold || v > a.MaxTreshold {
		switch a.ZoneCategoryId {
		case 1:
			var coldRoom ColdRoomRestrictions
			err := sqlx.Get(db, &coldRoom, `select * from cold_room_restrictions where alarm_id = $1`, a.ID)
			if err != nil {
				return s.HandlePSQLError(s.Select, err, "alarm log insert error")
			}
			if float64(coldRoom.AlarmTime) > ((float64(coldRoom.DefrostTime) * 3.5) / 5) {
				_, err := db.Exec(`update cold_room_restrictions set alarm_time = alarm_time -12 where alarm_id = $1`, a.ID)
				if err != nil {
					return s.HandlePSQLError(s.Select, err, "alarm log insert error")
				}
				ExecuteAlarm(a, v, deviceName, ap, zoneName, alarmType, date, db)
			} else {
				_, err := db.Exec(`update cold_room_restrictions set alarm_time = alarm_time +1 where alarm_id = $1`, a.ID)
				if err != nil {
					return s.HandlePSQLError(s.Select, err, "alarm log insert error")
				}
			}

		case 0:
			err := ExecuteAlarm(a, v, deviceName, ap, zoneName, alarmType, date, db)
			if err != nil {
				return s.HandlePSQLError(s.Select, err, "alarm log insert error")
			}

		case 2:
			if a.DevEui == "24e124136b325291" {
				fmt.Println("GELDİİİİİİİİ")
			}
			err := ExecuteAlarm(a, v, deviceName, ap, zoneName, alarmType, date, db)
			if err != nil {
				return s.HandlePSQLError(s.Select, err, "alarm log insert error")
			}
		default:
			if a.DevEui == "24e124136b325291" {
				fmt.Println("GELDİİİİİİİİ")
			}
			err := ExecuteAlarm(a, v, deviceName, ap, zoneName, alarmType, date, db)
			if err != nil {
				return s.HandlePSQLError(s.Select, err, "alarm log insert error")
			}
		}

	} else {
		switch a.ZoneCategoryId {
		case 1:
			_, err := db.Exec(`update cold_room_restrictions set alarm_time = 0 where alarm_id = $1`, a.ID)
			if err != nil {
				return s.HandlePSQLError(s.Select, err, "alarm log insert error")
			}
		default:
		}

	}
	return nil
}

func ExecuteAlarm(a AlarmWithDates, v float32, deviceName string, ap s.Application, zoneName string, alarmType string, date string, db sqlx.Ext) error {
	var u s.User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return s.HandlePSQLError(s.Select, err, "alarm log insert error")
	}

	notification := s.Notification{
		SenderId:   0,
		ReceiverId: a.UserId,
		Message:    date + " tarihinde " + zoneName + " ortamındaki " + deviceName + " isimli sensör " + alarmType + " kritik alarm seviyerini gecti. şu an ki değeri: " + fmt.Sprintf("%.2f", v),
		CategoryId: 1,
		IsRead:     false,
		SendTime:   time.Now(),
		SenderIp:   "system",
		ReaderIp:   "",
		IsDeleted:  false,
		DeviceName: deviceName,
		DevEui:     a.DevEui,
	}
	err = CreateNotification(notification)
	if err!=nil{
		log.Println("CreateNotification Error")
	}
	if a.Sms  {
		numbers := []string{u.PhoneNumber}
		numbersString := s.NumbersArrayToString(numbers)

		sms1N := s.OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = date + " tarihinde " + zoneName + " ortamındaki " + deviceName + " isimli sensör " + alarmType + " kritik alarm seviyerini gecti. şu an ki değeri: " + fmt.Sprintf("%.2f", v)
		sms1N.Type = "normal"
		sms1N.Send1N()

		// _, err := db.Exec(`insert into
		// values($1, $2, $3, $4, $5)`, u.ID, a.ID, a.DevEui, "sms", fmt.Sprintf("%.2f", v))
		// if err != nil {
		// 	return handlePSQLError(Select, err, "alarm log insert error")
		// }
	}

	if a.Email  {

		var user s.User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return s.HandlePSQLError(s.Select, err, "alarm log insert error")
		}

		s.SendEmail(user.Email, date+" tarihinde "+zoneName+" ortamındaki "+deviceName+" isimli sensör "+alarmType+" kritik alarm seviyerini gecti. şu an ki değeri: "+fmt.Sprintf("%.2f", v))
	}
	if a.Notification  {
		n := s.FirebaseNotificationData{
			Title: "Vaps",
			Body:  date + " tarihinde " + zoneName + " ortamındaki " + deviceName + " isimli sensör " + alarmType + " kritik alarm seviyerini gecti. şu an ki değeri: " + fmt.Sprintf("%.2f", v),
			Time:  300000,
			Delay: false,
		}

		s.SendFirebaseNotification(u, n)
	}
	return nil
}

func checkAlarmTime(a AlarmWithDates) bool {
	if a.IsTimeLimitActive {
		hours, minutes, _ := time.Now().Clock()
		result := strconv.Itoa(hours+3) + "." + strconv.Itoa(minutes)
		t, err := strconv.ParseFloat(result, 32)
		if err != nil {
			println("ParseFloat error")
		}
		if a.AlarmEndTime > a.AlarmStartTime {
			if a.AlarmStartTime < float32(t) && float32(t) < a.AlarmEndTime {
				return true
			}
		} else {
			if a.AlarmStartTime < float32(t) && float32(t) < 24 {
				return true
			}
			if 0 < float32(t) && float32(t) < a.AlarmEndTime {
				return true
			}
		}
	} else {
		return true
	}
	return false
}

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}
