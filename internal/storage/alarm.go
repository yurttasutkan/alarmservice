package storage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
)

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

func CheckThreshold(alarm AlarmWithDates, data float32, device als.Device, alarmType string, date string, db sqlx.Ext) error {
	if data < alarm.MinTreshold || data > alarm.MaxTreshold {

		switch alarm.ZoneCategoryId {
		case 1:
			var coldRoom ColdRoomRestrictions
			err := sqlx.Get(db, &coldRoom, `select * from cold_room_restrictions where alarm_id = $1`, alarm.ID)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
			if float64(coldRoom.AlarmTime) > ((float64(coldRoom.DefrostTime) * 3.5) / 5) {
				_, err := db.Exec(`update cold_room_restrictions set alarm_time = alarm_time -12 where alarm_id = $1`, alarm.ID)
				if err != nil {
					return HandlePSQLError(Select, err, "alarm log insert error")
				}

				ExecuteAlarm(alarm, data, device, alarmType, date, db)

			} else {
				_, err := db.Exec(`update cold_room_restrictions set alarm_time = alarm_time +1 where alarm_id = $1`, alarm.ID)
				if err != nil {
					return HandlePSQLError(Select, err, "alarm log insert error")
				}
			}
			break
		case 0:

			err := ExecuteAlarm(alarm, data, device, alarmType, date, db)

			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
			break
		case 2:
			err := ExecuteAlarm(alarm, data, device, alarmType, date, db)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
			break
		default:
			err := ExecuteAlarm(alarm, data, device, alarmType, date, db)
			if err != nil {
				return HandlePSQLError(Select, err, "alarm log insert error")
			}
			break
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
