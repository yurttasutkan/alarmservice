package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
)

func ExecuteAlarm(alarm AlarmWithDates, data float32, device als.Device, alarmType string, date string, db sqlx.Ext) error {

	var user User
	var zoneName string
	var message string
	err := sqlx.Get(db, &user, "select * from public.user where id = $1", alarm.UserId)
	if err != nil {
		return HandlePSQLError(Select, err, "alarm log insert error")
	}

	err = sqlx.Get(db, &zoneName, "select zone_name from zone where '\\x' || $1 = any(zone.devices)", device.DevEui)
	if err != nil {
		fmt.Println("get zone anme error")
	}
	switch alarmType {
	case "isi":
		message = date + " tarihinde " + zoneName + " ortamındaki " + device.Name + " isimli sensör " + alarmType + " kritik alarm seviyerini gecti. şu an ki değeri: " + fmt.Sprintf("%.2f", data)
		break
	}

	notification := Notification{
		SenderId:   0,
		ReceiverId: alarm.UserId,
		Message:    message,
		CategoryId: 1,
		IsRead:     false,
		SendTime:   time.Now(),
		SenderIp:   "system",
		ReaderIp:   "",
		IsDeleted:  false,
		DeviceName: device.Name,
		DevEui:     alarm.DevEui,
	}
	err = CreateNotification(notification)
	if err != nil {
		log.Println("CreateNotification Error")
	}

	// Check what kind of alarm it is going to send
	if alarm.Sms {
		numbers := []string{user.PhoneNumber}
		numbersString := NumbersArrayToString(numbers)

		sms1N := OneToN{}
		sms1N.UserID = 40584
		sms1N.Username = "905322424400"
		sms1N.Password = "001Sye44"
		sms1N.Sender = "VERITEL"
		sms1N.Numbers = numbersString
		sms1N.Message = message
		sms1N.Type = "normal"
		sms1N.Send1N()
	}

	if alarm.Email {
		SendEmail(user.Email, message)
	}
	if alarm.Notification {
		n := FirebaseNotificationData{
			Title: "Vaps",
			Body:  message,
			Time:  300000,
			Delay: false,
		}

		SendFirebaseNotification(user, n)
	}
	return nil
}
