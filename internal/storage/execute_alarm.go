package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

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

	// Check what kind of alarm it is going to send
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
