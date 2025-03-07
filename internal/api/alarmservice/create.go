package alarmservice

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Implements the RPC method CreateAlarm.
// Inserts into alarm_refactor2 and logs the change in the audit logs.
func (a *AlarmServerAPI) CreateAlarm(context context.Context, req *als.CreateAlarmRequest) (*als.CreateAlarmResponse, error) {
	db := s.DB()
	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback()

	var returnID int64
	var alarmDates []s.AlarmDateFilter
	al := req.Alarm

	// Insert alarm into alarm_refactor2
	pqInt64Array := pq.Int64Array(al.UserID)
	err = tx.QueryRowx(`
		insert into alarm_refactor2 (
			dev_eui, min_treshold, max_treshold, sms, email, temperature, humadity, ec, door, w_leak,
			user_id, is_time_limit_active, alarm_start_time, alarm_stop_time, zone_category, notification,
			notification_sound, distance, pressure
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		returning id`,
		al.DevEui, al.MinTreshold, al.MaxTreshold, al.Sms, al.Email, al.Temperature, al.Humadity, al.Ec,
		al.Door, al.WLeak, pqInt64Array, al.IsTimeLimitActive, al.AlarmStartTime, al.AlarmStopTime,
		al.ZoneCategoryID, al.Notification, al.NotificationSound, al.Distance, al.Pressure,
	).Scan(&returnID)
	if err != nil {
		return nil, s.HandlePSQLError(s.Insert, err, "insert error")
	}

	// Log the creation in the audit log
	// Creating a struct for previous values (which is nil here as it's a new record)
	var previousValue interface{} = nil

	// New value contains the values of the alarm being created
	newAlarm := als.Alarm{
		Id:                returnID,
		DevEui:            al.DevEui,
		MinTreshold:       al.MinTreshold,
		MaxTreshold:       al.MaxTreshold,
		Sms:               al.Sms,
		Email:             al.Email,
		Notification:      al.Notification,
		Temperature:       al.Temperature,
		Humadity:          al.Humadity,
		Ec:                al.Ec,
		Door:              al.Door,
		WLeak:             al.WLeak,
		UserID:            al.UserID,
		IpAddress:         al.IpAddress,
		IsTimeLimitActive: al.IsTimeLimitActive,
		AlarmStartTime:    al.AlarmStartTime,
		AlarmStopTime:     al.AlarmStopTime,
		ZoneCategoryID:    al.ZoneCategoryID,
		IsActive:          al.IsActive,
		AlarmDateTime:     nil, // You can append dates later
		NotificationSound: al.NotificationSound,
		Distance:          al.Distance,
		DefrostTime:       al.DefrostTime,
		Pressure:          al.Pressure,
	}

	// Log the creation in the audit log
	if err := s.LogAudit(db, newAlarm.Id, req.UserId, "INSERT", previousValue, newAlarm); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("could not log audit: %v", err)
	}

	// Handle specific logic for ZoneCategoryID = 1
	if al.ZoneCategoryID == 1 {
		if err := s.CreateColdRoomRestrictions(al, returnID, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := s.CreateUtku(al, returnID, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Handle alarm date times
	for _, alarmDateTime := range al.AlarmDateTime {
		dt := s.AlarmDateFilter{
			AlarmId:        returnID,
			AlarmDay:       alarmDateTime.AlarmDay,
			AlarmStartTime: alarmDateTime.AlarmStartTime,
			AlarmEndTime:   alarmDateTime.AlarmEndTime,
		}
		alarmDates = append(alarmDates, dt)
	}

	// Create alarm dates
	dates, err := s.CreateAlarmDates(tx, alarmDates)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Construct response
	resp := als.CreateAlarmResponse{
		Alarm: &als.Alarm{
			Id:                returnID,
			DevEui:            al.DevEui,
			MinTreshold:       al.MinTreshold,
			MaxTreshold:       al.MaxTreshold,
			Sms:               al.Sms,
			Email:             al.Email,
			Notification:      al.Notification,
			Temperature:       al.Temperature,
			Humadity:          al.Humadity,
			Ec:                al.Ec,
			Door:              al.Door,
			WLeak:             al.WLeak,
			UserID:            al.UserID,
			IpAddress:         al.IpAddress,
			IsTimeLimitActive: al.IsTimeLimitActive,
			AlarmStartTime:    al.AlarmStartTime,
			AlarmStopTime:     al.AlarmStopTime,
			ZoneCategoryID:    al.ZoneCategoryID,
			IsActive:          al.IsActive,
			AlarmDateTime:     dates,
			NotificationSound: al.NotificationSound,
			Distance:          al.Distance,
			DefrostTime:       al.DefrostTime,
			Pressure:          al.Pressure,
		},
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return &resp, nil
}

// Implements the RPC method UpdateAlarm.
// Updates alarm_refactor table with the parameters given by request.
func (a *AlarmServerAPI) UpdateAlarm(ctx context.Context, req *als.UpdateAlarmRequest) (*empty.Empty, error) {
	db := s.DB()
	var alarmDates []s.AlarmDateFilter

	// Get the previous values of the alarm
	currentAlarm, err := db.Exec("select * from alarm_refactor2 where id = $1", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	alarm := req.Alarm
	pqInt64Array := pq.Int64Array(alarm.UserID)
	res, err := db.Exec(`update alarm_refactor2 
	set   min_treshold = $1,
	max_treshold = $2,
	sms    = $3,
	email = $4,
	notification = $5,
	is_time_limit_active = $6,
	notification_sound = $8,
	user_id = $9,
	is_active = $10,
	defrost_time = $11
	where id = $7`,
		alarm.MinTreshold,
		alarm.MaxTreshold,
		alarm.Sms,
		alarm.Email,
		alarm.Notification,
		alarm.IsTimeLimitActive,
		req.AlarmID,
		alarm.NotificationSound,
		pqInt64Array,
		alarm.IsActive,
		alarm.DefrostTime,
	)
	if err != nil {
		log.Println(err)
	}
	_, err = db.Exec("delete from alarm_date_time where alarm_id = $1", req.Alarm.Id)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}
	for _, alarmDateTime := range req.Alarm.AlarmDateTime {
		dt := s.AlarmDateFilter{
			AlarmId:        alarmDateTime.AlarmId,
			AlarmDay:       alarmDateTime.AlarmDay,
			AlarmStartTime: alarmDateTime.AlarmStartTime,
			AlarmEndTime:   alarmDateTime.AlarmEndTime,
		}
		alarmDates = append(alarmDates, dt)
	}
	_, err = s.CreateAlarmDates(db, alarmDates)
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Update, err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return &emptypb.Empty{}, nil
	}

	// Fetch the updated alarm from the database
	var updatedAlarm s.Alarm
	err = sqlx.Get(db, &updatedAlarm, "select * from alarm_refactor2 where id = $1", req.AlarmID)
	if err != nil {
		return nil, s.HandlePSQLError(s.Select, err, "fetch updated alarm error")
	}
	// Log the creation in the audit log
	if err := s.LogAudit(db, alarm.Id, req.UserId, "UPDATE", currentAlarm, updatedAlarm); err != nil {
		return nil, fmt.Errorf("could not log audit: %v", err)
	}

	return &emptypb.Empty{}, nil
}
