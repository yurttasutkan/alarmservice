package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
)

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
	WaterLeak         bool    `db:"w_leak"`
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
	WaterLeak         bool    `db:"w_leak"`
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
	WaterLeak         bool    `db:"w_leak"`
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
	ID               int64  `db:"id"`
	DevEui           string `db:"dev_eui"`
	AlarmId          int64  `db:"alarm_id"`
	DefrostTime      int64  `db:"defronst_time"`
	DefrostFrequency int64  `db:"defrost_frequency"`
	AlarmTime        int64  `db:"alarm_time"`
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
	WaterLeak      bool      `db:"w_leak"`
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
func DoorAlarm(a AlarmWithDates, deviceName string, zonename string, alarmType string, date string) error {
	currentTime := time.Now().Add(time.Hour * 3)
	db := DB()
	var u User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return HandlePSQLError(Select, err, "alarm log insert error")
	}

	notification := Notification{
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
	if err != nil {
		log.Println("CreateNotification error")
	}

	if a.Sms {
		numbers := []string{u.PhoneNumber}
		numbersString := NumbersArrayToString(numbers)

		sms1N := OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + " deki " + deviceName + " sensörü açıldı"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email {

		var user User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return HandlePSQLError(Select, err, "alarm log insert error")
		}

		SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+" deki "+deviceName+" sensörü açıldı")
	}

	if a.Notification {
		n := FirebaseNotificationData{
			Title: "Vaps",
			Body:  currentTime.Format("2006-01-02 15:04:05") + " tarihinde " + zonename + " deki " + deviceName + " sensörü açıldı",
			Time:  300000,
			Delay: false,
		}

		SendFirebaseNotification(u, n)
	}
	return nil
}

func AlarmButton(a AlarmWithDates, deviceName string, zonename string) error {
	db := DB()
	currentTime := time.Now().Add(time.Hour * 3)
	var u User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return HandlePSQLError(Select, err, "alarm log insert error")
	}

	if a.Sms {
		numbers := []string{u.PhoneNumber}
		numbersString := NumbersArrayToString(numbers)

		sms1N := OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + " deki " + deviceName + " sensöründen çağrı var"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email {

		var user User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return HandlePSQLError(Select, err, "alarm log insert error")
		}

		SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+" deki "+deviceName+" sensöründen çağrı var")
	}

	if a.Notification {
		n := FirebaseNotificationData{
			Title: "Vaps",
			Body:  zonename + " deki " + deviceName + " sensöründen çağrı var" + " tarih: " + currentTime.Format("2006-01-02 15:04:05"),
			Time:  300000,
			Delay: false,
		}

		SendFirebaseNotification(u, n)
	}

	return nil
}
func WaterLeakAlarm(a AlarmWithDates, deviceName string, zonename string) error {
	db := DB()
	currentTime := time.Now().Add(time.Hour * 3)
	var u User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return HandlePSQLError(Select, err, "alarm log insert error")
	}

	if a.Sms {
		numbers := []string{u.PhoneNumber}
		numbersString := NumbersArrayToString(numbers)

		sms1N := OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + " deki " + deviceName + " sensöründe kaçak var"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email {
		var user User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return HandlePSQLError(Select, err, "alarm log insert error")
		}

		SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+" deki "+deviceName+" sensöründe kaçak var")
	}

	if a.Notification {
		n := FirebaseNotificationData{
			Title: "Vaps",
			Body:  currentTime.Format("2006-01-02 15:04:05") + " tarihinde " + zonename + " deki " + deviceName + " sensöründe kaçak var",
			Time:  300000,
			Delay: false,
		}

		SendFirebaseNotification(u, n)
	}

	return nil
}

func EmergencyAlarm(a AlarmWithDates, deviceName string, zonename string) error {

	db := DB()

	currentTime := time.Now().Add(time.Hour * 3)
	var u User

	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return HandlePSQLError(Select, err, "alarm log insert error")
	}

	if a.Sms {
		numbers := []string{u.PhoneNumber}
		numbersString := NumbersArrayToString(numbers)

		sms1N := OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = zonename + "deki" + deviceName + " sensöründe acil durum var"
		sms1N.Type = "normal"
		sms1N.Send1N()
	}
	if a.Email {

		var user User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return HandlePSQLError(Select, err, "alarm log insert error")
		}

		SendEmail(user.Email, currentTime.Format("2006-01-02 15:04:05")+" tarihinde "+zonename+"deki"+deviceName+" sensöründe acil durum var")
	}

	if a.Notification {
		n := FirebaseNotificationData{
			Title: "Vaps",
			Body:  currentTime.Format("2006-01-02 15:04:05") + " tarihinde " + zonename + "deki" + deviceName + " sensöründe acil durum var",
			Time:  300000,
			Delay: false,
		}

		SendFirebaseNotification(u, n)
	}
	return nil
}

func ExecuteAlarm(a AlarmWithDates, v float32, deviceName string, zoneName string, alarmType string, date string, db sqlx.Ext) error {
	var u User
	err := sqlx.Get(db, &u, "select * from public.user where id = $1", a.UserId)
	if err != nil {
		return HandlePSQLError(Select, err, "alarm log insert error")
	}

	notification := Notification{
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
	if err != nil {
		log.Println("CreateNotification Error")
	}
	if a.Sms {
		numbers := []string{u.PhoneNumber}
		numbersString := NumbersArrayToString(numbers)

		sms1N := OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = date + " tarihinde " + zoneName + " ortamındaki " + deviceName + " isimli sensör " + alarmType + " kritik alarm seviyerini gecti. şu an ki değeri: " + fmt.Sprintf("%.2f", v)
		sms1N.Type = "normal"
		sms1N.Send1N()
	}

	if a.Email {

		var user User
		err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, a.UserId)
		if err != nil {
			return HandlePSQLError(Select, err, "alarm log insert error")
		}

		SendEmail(user.Email, date+" tarihinde "+zoneName+" ortamındaki "+deviceName+" isimli sensör "+alarmType+" kritik alarm seviyerini gecti. şu an ki değeri: "+fmt.Sprintf("%.2f", v))
	}
	if a.Notification {
		n := FirebaseNotificationData{
			Title: "Vaps",
			Body:  date + " tarihinde " + zoneName + " ortamındaki " + deviceName + " isimli sensör " + alarmType + " kritik alarm seviyerini gecti. şu an ki değeri: " + fmt.Sprintf("%.2f", v),
			Time:  300000,
			Delay: false,
		}

		SendFirebaseNotification(u, n)
	}
	return nil
}

func CreateNotification(notification Notification) error {
	db := DB()
	_, err := db.Exec(`insert into notifications(sender_id, 
		receiver_id,
		message,
		category_id,
		read_time,
		deleted_time,
		sender_ip,
		reader_ip,
		is_deleted,
		device_name,
		dev_eui) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		notification.SenderId,
		notification.ReceiverId,
		notification.Message,
		notification.CategoryId,
		notification.ReadTime,
		notification.DeletedTime,
		notification.SenderIp,
		notification.ReaderIp,
		notification.IsDeleted,
		notification.DeviceName,
		notification.DevEui)
	if err != nil {
		return HandlePSQLError(Insert, err, "insert error")
	}

	return nil
}

func CheckAlarmTimeSchedule(a AlarmWithDates, v float32, deviceName string, zoneName string, alarmType string, date string, db sqlx.Ext) error {
	if a.IsTimeLimitActive {
		hours, minutes, _ := time.Now().Clock()
		result := strconv.Itoa(hours+3) + "." + strconv.Itoa(minutes)
		t, err := strconv.ParseFloat(result, 32)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		if a.AlarmEndTime > a.AlarmStartTime {
			if a.AlarmStartTime < float32(t) && float32(t) < a.AlarmEndTime {
				CheckThreshold(a, v, deviceName, zoneName, alarmType, date, db)
			}
		} else {
			if a.AlarmStartTime < float32(t) && float32(t) < 24 {
				CheckThreshold(a, v, deviceName,  zoneName, alarmType, date, db)
			}
			if 0 < float32(t) && float32(t) < a.AlarmEndTime {
				CheckThreshold(a, v, deviceName, zoneName, alarmType, date, db)
			}
		}
	} else {
		CheckThreshold(a, v, deviceName, zoneName, alarmType, date, db)
	}
	return nil
}

func CheckThreshold(a AlarmWithDates, v float32, deviceName string, zoneName string, alarmType string, date string, db sqlx.Ext) error {
	if v < a.MinTreshold || v > a.MaxTreshold {
		switch a.ZoneCategoryId {
		case 1:
			var coldRoom ColdRoomRestrictions
			err := sqlx.Get(db, &coldRoom, `select * from cold_room_restrictions where alarm_id = $1`, a.ID)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
			if float64(coldRoom.AlarmTime) > ((float64(coldRoom.DefrostTime) * 3.5) / 5) {
				_, err := db.Exec(`update cold_room_restrictions set alarm_time = alarm_time -12 where alarm_id = $1`, a.ID)
				if err != nil {
					return HandlePSQLError(Select, err, "alarm log insert error")
				}
				ExecuteAlarm(a, v, deviceName, zoneName, alarmType, date, db)
			} else {
				_, err := db.Exec(`update cold_room_restrictions set alarm_time = alarm_time +1 where alarm_id = $1`, a.ID)
				if err != nil {
					return HandlePSQLError(Select, err, "alarm log insert error")
				}
			}

		case 0:
			err := ExecuteAlarm(a, v, deviceName, zoneName, alarmType, date, db)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}

		case 2:
			err := ExecuteAlarm(a, v, deviceName, zoneName, alarmType, date, db)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
		default:

			err := ExecuteAlarm(a, v, deviceName, zoneName, alarmType, date, db)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
		}

	} else {
		switch a.ZoneCategoryId {
		case 1:
			_, err := db.Exec(`update cold_room_restrictions set alarm_time = 0 where alarm_id = $1`, a.ID)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
		default:
		}

	}
	return nil
}

func CheckAlarmTime(a AlarmWithDates) bool {
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

func CreateAlarmLog(ctx context.Context, db sqlx.Ext, a *als.Alarm, userID int64, ipAddress string, isDeleted int64) error {
	fmt.Println("create alarm log")
	_, err := db.Exec(`insert into alarm_change_logs(
		dev_eui,
		min_treshold,
		max_treshold,
		user_id,
		ip_address,
		is_deleted,
		sms,
		temperature,
		humadity,
		ec,
		door,
		w_leak
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) `, a.DevEui, a.MinTreshold,
		a.MaxTreshold, userID, ipAddress, isDeleted, a.Sms, a.Temperature, a.Humadity, a.Ec, a.Door, a.WLeak)
	if err != nil {
		return HandlePSQLError(Insert, nil, "insert error")
	}

	return nil
}
