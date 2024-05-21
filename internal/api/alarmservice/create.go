package alarmservice

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Implements the RPC method CreateAlarm.
// Inserts into alarm_refactor with parameters given by request and returns the created Alarm as response.
func (a *AlarmServerAPI) CreateAlarm(context context.Context, alarm *als.CreateAlarmRequest) (*als.CreateAlarmResponse, error) {
	db := s.DB()
	var returnID int64
	var alarmDates []s.AlarmDateFilter
	al := alarm.Alarm
	fmt.Println(al.DevEui, len(al.AlarmDateTime))
	pqInt64Array := pq.Int64Array(al.UserID)
	err := db.QueryRowx(`
	insert into alarm_refactor2 (
		dev_eui,
		min_treshold,
		max_treshold,
		sms,
		email,
		temperature,
		humadity,
		ec,
		door,
		w_leak,
		user_id,
		is_time_limit_active,
		alarm_start_time,
		alarm_stop_time,
		zone_category,
		notification,
		notification_sound,
		distance
	) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) returning id`,
		al.DevEui,
		al.MinTreshold,
		al.MaxTreshold,
		al.Sms,
		al.Email,
		al.Temperature,
		al.Humadity,
		al.Ec,
		al.Door,
		al.WLeak,
		pqInt64Array,
		al.IsTimeLimitActive,
		al.AlarmStartTime,
		al.AlarmStopTime,
		al.ZoneCategoryID,
		al.Notification,
		al.NotificationSound,
		al.Distance,
	).Scan(&returnID)
	if err != nil {
		return nil, s.HandlePSQLError(s.Insert, err, "insert error")
	}

	// If Zone Category is 1, initialize als struct in order to use for CreateColdRoomRestrictionsRequest
	if alarm.Alarm.ZoneCategoryID == 1 {
		err := s.CreateColdRoomRestrictions(al, returnID, db)
		if err != nil {
			log.Println(err)
		}
		err = s.CreateUtku(al, returnID, db)
		if err != nil {
			log.Println(err)
		}
	}
	for _, alarmDateTime := range al.AlarmDateTime {
		fmt.Println("Create alarm date time ", al.DevEui)
		dt := s.AlarmDateFilter{
			AlarmId:        returnID,
			AlarmDay:       alarmDateTime.AlarmDay,
			AlarmStartTime: alarmDateTime.AlarmStartTime,
			AlarmEndTime:   alarmDateTime.AlarmEndTime,
		}
		alarmDates = append(alarmDates, dt)
	}

	dates, err := s.CreateAlarmDates(db, alarmDates)
	if err != nil {
		log.Println(err)
	}

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
		},
	}

	s.CreateAlarmLog(context, db, resp.Alarm, resp.Alarm.UserID, resp.Alarm.IpAddress, 1)

	return &resp, nil
}

// Implements the RPC method UpdateAlarm.
// Updates alarm_refactor table with the parameters given by request.
func (a *AlarmServerAPI) UpdateAlarm(ctx context.Context, req *als.UpdateAlarmRequest) (*empty.Empty, error) {
	db := s.DB()
	var alarmDates []s.AlarmDateFilter

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
	user_id = $9
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
		log.Println(err)
	}

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

	return &emptypb.Empty{}, nil
}
