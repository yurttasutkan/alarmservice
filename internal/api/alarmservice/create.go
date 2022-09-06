package alarmservice

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/pkg/errors"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Implements the RPC method CreateAlarm.
//Inserts into alarm_refactor with parameters given by request and returns the created Alarm as response.
func (a *AlarmServerAPI) CreateAlarm(context context.Context, alarm *als.CreateAlarmRequest) (*als.CreateAlarmResponse, error) {
	db := s.DB()
	var returnID int64

	al := alarm.Alarm
	err := db.QueryRowx(`
	insert into alarm_refactor (
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
		notification
	) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) returning id`,
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
		al.UserID,
		al.IsTimeLimitActive,
		al.AlarmStartTime,
		al.AlarmStopTime,
		al.ZoneCategoryID,
		al.Notification,
	).Scan(&returnID)
	if err != nil {
		return nil, s.HandlePSQLError(s.Insert, err, "insert error")
	}

	// If Zone Category is 1, initialize als struct in order to use for CreateColdRoomRestrictionsRequest
	if alarm.Alarm.ZoneCategoryID == 1 {
		err := s.CreateColdRoomRestrictions(al, returnID)
		if err != nil {
			log.Println(err)
		}
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
		},
	}

	s.CreateAlarmLog(context, db, resp.Alarm, resp.Alarm.UserID, resp.Alarm.IpAddress, 1)

	return &resp, nil
}

//Implements the RPC method CreateAlarmDates.
//Inserts into alarm_date_time with parameters given by request and returns the created AlarmDateTime as response.
func (a *AlarmServerAPI) CreateAlarmDates(ctx context.Context, req *als.CreateAlarmDatesRequest) (*als.CreateAlarmDatesResponse, error) {
	db := s.DB()

	var returnDates []*als.AlarmDateTime

	if len(req.ReqFilter) > 0 {
		for _, date := range req.ReqFilter {
			var returnID int64

			err := db.QueryRowx(`insert into 
			alarm_date_time(alarm_id, alarm_day, start_time, end_time) values ($1, $2, $3, $4) returning id`,
				date.AlarmId, date.AlarmDay, date.AlarmStartTime, date.AlarmEndTime).Scan(&returnID)

			if err != nil {
				return &als.CreateAlarmDatesResponse{RespDateTime: returnDates}, s.HandlePSQLError(s.Insert, err, "insert error")
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
	return &als.CreateAlarmDatesResponse{RespDateTime: returnDates}, nil
}

//Implements the RPC method CreateAlarmLog.
//Inserts into alarm_change_logs with parameters given by request.
func (a *AlarmServerAPI) CreateAlarmLog(ctx context.Context, req *als.CreateAlarmLogRequest) (*empty.Empty, error) {
	db := s.DB()
	al := req.Alarm

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
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) `, al.DevEui, al.MinTreshold,
		al.MaxTreshold, req.UserID, req.IpAddress, req.IsDeleted, al.Sms, al.Temperature, al.Humadity, al.Ec, al.Door, al.WLeak)
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Insert, err, "insert error")
	}

	return &empty.Empty{}, nil
}

//Implements the RPC method UpdateAlarm.
//Updates alarm_refactor table with the parameters given by request.
func (a *AlarmServerAPI) UpdateAlarm(ctx context.Context, req *als.UpdateAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	alarm := req.Alarm
	res, err := db.Exec(`update alarm_refactor 
	set   min_treshold = $1,
	max_treshold = $2,
	sms    = $3,
	email = $4,
	notification = $5,
	is_time_limit_active = $6
	where id = $7`,
		alarm.MinTreshold,
		alarm.MaxTreshold,
		alarm.Sms,
		alarm.Email,
		alarm.Notification,
		alarm.IsTimeLimitActive,
		req.AlarmID,
	)

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
