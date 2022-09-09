package storage

import (
	"log"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
)

//Implements the RPC method CreateColdRoomRestrictions.
//Inserts into cold_room_restrictions table with given parameters in the request.
func CreateColdRoomRestrictions(alarm *als.Alarm, alarmID int64, db sqlx.Ext) error {
	coldRes := ColdRoomRestrictions{
		DevEui:           alarm.DevEui,
		AlarmId:          alarmID,
		DefrostTime:      alarm.ColdRoomTime,
		DefrostFrequency: alarm.ColdRoomFreq,
		AlarmTime:        0,
	}
	_, err := db.Exec(`insert into cold_room_restrictions(
		dev_eui,
		alarm_id,
		defronst_time,
		defrost_frequency,
		alarm_time
	) values ($1, $2, $3, $4, $5)`, coldRes.DevEui, coldRes.AlarmId, coldRes.DefrostTime, coldRes.DefrostFrequency, coldRes.AlarmTime)
	if err != nil {
		log.Fatalf("Insert error %v", err)
	}
	return nil
}

func CreateUtku(alarm *als.Alarm, alarmID int64, db sqlx.Ext) error {
	utkuObject := UtkuStruct{
		DevEui:  alarm.DevEui,
		AlarmId: alarmID,
	}
	_, err := db.Exec(`insert into utku_table(
		dev_eui,
		alarm_id
	) values ($1, $2)`, utkuObject.DevEui, utkuObject.AlarmId)
	if err != nil {
		log.Fatalf("Insert error %v", err)
	}
	return nil
}
