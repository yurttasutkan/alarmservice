package storage

import (
	"fmt"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
)

func CreateAlarmDates(db sqlx.Queryer, alarmDates []AlarmDateFilter) ([]*als.AlarmDateTime, error) {
	fmt.Println("create alarm date")

	var returnDates []*als.AlarmDateTime

	if len(alarmDates) > 0 {
		for _, date := range alarmDates {
			var returnID int64

			err := db.QueryRowx(`insert into
			alarm_date_time(alarm_id, alarm_day, start_time, end_time) values ($1, $2, $3, $4) returning id`,
				date.AlarmId, date.AlarmDay, date.AlarmStartTime, date.AlarmEndTime).Scan(&returnID)
			if err != nil {
				return returnDates, HandlePSQLError(Insert, nil, "insert error")
			}
			createdDate := als.AlarmDateTime{
				Id:             returnID,
				AlarmId:        date.AlarmId,
				AlarmDay:       date.AlarmDay,
				AlarmStartTime: date.AlarmStartTime,
				AlarmEndTime:   date.AlarmEndTime,
			}
			returnDates = append(returnDates, &createdDate)
		}
	}

	return returnDates, nil

}
