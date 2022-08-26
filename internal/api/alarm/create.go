package alarm

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/pkg/errors"
	"github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Implements the RPC method CreateAlarm.
func (a *AlarmServerAPI) CreateAlarm(context context.Context, alarm *als.CreateAlarmRequest) (*als.CreateAlarmResponse, error) {
	db := storage.DB()
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
		return nil, storage.HandlePSQLError(storage.Insert, err, "insert error")
	}

	resp := als.CreateAlarmResponse{
		Alarm: &als.Alarm{
			Id:                al.Id,
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
	return &resp, nil
}
func (a *AlarmServerAPI) CreateAlarmDates(ctx context.Context, req *als.CreateAlarmDatesRequest) (*als.CreateAlarmDatesResponse, error) {
	db := storage.DB()

	fmt.Println("create alarm date")
	var returnDates []*als.AlarmDateTime

	if len(req.ReqFilter) > 0 {
		for _, date := range req.ReqFilter {
			var returnID int64

			err := db.QueryRowx(`insert into 
			alarm_date_time(alarm_id, alarm_day, start_time, end_time) values ($1, $2, $3, $4) returning id`,
				date.AlarmId, date.AlarmDay, date.AlarmStartTime, date.AlarmEndTime).Scan(&returnID)
			if err != nil {
				return &als.CreateAlarmDatesResponse{RespDateTime: returnDates}, storage.HandlePSQLError(storage.Insert, err, "insert error")
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

func (a *AlarmServerAPI) CreateAlarmLog(ctx context.Context, req *als.CreateAlarmLogRequest) (*empty.Empty, error) {
	db := storage.DB()
	al := req.Alarm

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
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) `, al.DevEui, al.MinTreshold,
		al.MaxTreshold, req.UserID, req.IpAddress, req.IsDeleted, al.Sms, al.Temperature, al.Humadity, al.Ec, al.Door, al.WLeak)
	if err != nil {
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Insert, err, "insert error")
	}

	return &empty.Empty{}, nil
}

func (a *AlarmServerAPI) UpdateAlarm(ctx context.Context, req *als.UpdateAlarmRequest) (*empty.Empty, error) {
	db := storage.DB()

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
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Update, err, "update error")
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
