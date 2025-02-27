package alarmservice

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
)

// Implements the RPC method GetAlarm.
// Request takes alarmID as field and returns Alarm as response.
func (a *AlarmServerAPI) GetAlarm(ctx context.Context, alReq *als.GetAlarmRequest) (*als.GetAlarmResponse, error) {
	db := s.DB()
	var resp als.GetAlarmResponse
	var respAlarm s.Alarm
	var alarmDates []*als.AlarmDateTime

	err := sqlx.Get(db, &respAlarm, "select * from alarm_refactor2 where id = $1", alReq.AlarmID)
	if err != nil {
		fmt.Println(err)
		return &resp, s.HandlePSQLError(s.Select, err, "select error")
	}
	var dates []s.AlarmDateFilter
	err = sqlx.Select(db, &dates, "select * from alarm_date_time where alarm_id = $1", alReq.AlarmID)
	if err != nil {
		return &resp, s.HandlePSQLError(s.Select, err, "select error")
	}

	for _, date := range dates {
		dt := &als.AlarmDateTime{
			Id:             date.ID,
			AlarmId:        date.AlarmId,
			AlarmDay:       date.AlarmDay,
			AlarmStartTime: date.AlarmStartTime,
			AlarmEndTime:   date.AlarmEndTime,
		}
		alarmDates = append(alarmDates, dt)
	}
	al := als.Alarm{
		Id:                respAlarm.ID,
		DevEui:            respAlarm.DevEui,
		MinTreshold:       respAlarm.MinTreshold,
		MaxTreshold:       respAlarm.MaxTreshold,
		Sms:               respAlarm.Sms,
		Email:             respAlarm.Email,
		Temperature:       respAlarm.Temperature,
		Humadity:          respAlarm.Humadity,
		Ec:                respAlarm.Ec,
		Door:              respAlarm.Door,
		WLeak:             respAlarm.WaterLeak,
		IsTimeLimitActive: respAlarm.IsTimeLimitActive,
		AlarmStartTime:    respAlarm.AlarmStartTime,
		AlarmStopTime:     respAlarm.AlarmStopTime,
		Notification:      respAlarm.Notification,
		UserID:            respAlarm.UserId,
		IpAddress:         respAlarm.IpAddress,
		ZoneCategoryID:    respAlarm.ZoneCategoryId,
		IsActive:          respAlarm.IsActive,
		AlarmDateTime:     alarmDates,
		NotificationSound: respAlarm.NotificationSound,
		Distance:          respAlarm.Distance,
		Pressure:          respAlarm.Pressure,
		DefrostTime:       respAlarm.DefrostTime,
	}
	fmt.Println("GEL ALARM SONU")

	resp.Alarm = &al
	return &resp, nil
}

// Implements the RPC method GetAlarmLogs.
// Request takes DevEUI as field and returns []AlarmLogs as response.
func (a *AlarmServerAPI) GetAlarmLogs(ctx context.Context, req *als.GetAlarmLogsRequest) (*als.GetAlarmLogsResponse, error) {
	db := s.DB()
	var logs []s.AlarmLogs
	var result []*als.AlarmLogs

	err := sqlx.Select(db, &logs, `select dev_eui,min_treshold,
	max_treshold,
	user_id,
	ip_address,
	submission_date,
	is_deleted,
	sms,
	temperature,
	humadity,
	ec,
	door,
	w_leak from alarm_change_logs where dev_eui = $1`, req.DevEui)
	if err != nil {
		return &als.GetAlarmLogsResponse{RespLog: result}, s.HandlePSQLError(s.Select, err, "select error")
	}

	for _, log := range logs {
		item := als.AlarmLogs{
			DevEui:      log.DevEui,
			MinTreshold: log.MinTreshold,
			MaxTreshold: log.MaxTreshold,
			UserId:      log.UserId,
			IpAddress:   log.IpAddress,
			IsDeleted:   log.IsDeleted,
			Sms:         log.Sms,
			Temperature: log.Temperature,
			Humadity:    log.Humadity,
			Ec:          log.Ec,
			Door:        log.Door,
			WLeak:       log.WaterLeak,
		}
		var err error
		item.SubmissionDate, err = ptypes.TimestampProto(log.SubmissionDate)
		if err != nil {
			return &als.GetAlarmLogsResponse{RespLog: result}, s.HandlePSQLError(s.Select, err, "select error")
		}
		result = append(result, &item)
	}
	return &als.GetAlarmLogsResponse{RespLog: result}, nil
}

// Implements the RPC method GetAlarmDates.
// Request takes alarmID as field and returns []AlarmDateTime as response.
func (a *AlarmServerAPI) GetAlarmDates(ctx context.Context, req *als.GetAlarmDatesRequest) (*als.GetAlarmDatesResponse, error) {
	db := s.DB()
	fmt.Println("ALS = ALARM GET ALARM DATES")

	var returnDates []*als.AlarmDateTime
	var alarmDates []s.AlarmDateFilter

	err := sqlx.Select(db, &alarmDates, "select * from alarm_date_time where alarm_id = $1", req.AlarmId)
	if err != nil {
		return &als.GetAlarmDatesResponse{RespDate: returnDates}, s.HandlePSQLError(s.Select, err, "select error")
	}
	for _, date := range alarmDates {
		d := als.AlarmDateTime{
			Id:             date.ID,
			AlarmId:        date.AlarmId,
			AlarmDay:       date.AlarmDay,
			AlarmStartTime: date.AlarmStartTime,
			AlarmEndTime:   date.AlarmEndTime,
		}
		returnDates = append(returnDates, &d)
	}
	return &als.GetAlarmDatesResponse{RespDate: returnDates}, nil
}

// Implements the RPC method GetAlarmList.
// Request takes AlarmFilter as field and returns []Alarm as response.
func (a *AlarmServerAPI) GetAlarmList(ctx context.Context, req *als.GetAlarmListRequest) (*als.GetAlarmListResponse, error) {
	db := s.DB()

	filters := s.AlarmFilters{
		Limit:  int(req.Filter.Limit),
		DevEui: req.Filter.DevEui,
		UserID: req.Filter.UserID,
	}
	var returnAlarms []*als.Alarm
	var alarms []s.Alarm
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
	select *
	from alarm_refactor2
	`+filters.SQL(), filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}
	err = sqlx.Select(db, &alarms, query, args...)
	if err != nil {
		return &als.GetAlarmListResponse{RespList: returnAlarms}, s.HandlePSQLError(s.Select, err, "select error")
	}

	for _, alarm := range alarms {
		var alarmDates []*als.AlarmDateTime
		var dates []s.AlarmDateFilter
		err := sqlx.Select(db, &dates, "select * from alarm_date_time where alarm_id = $1", alarm.ID)
		if err != nil {
			return &als.GetAlarmListResponse{RespList: returnAlarms}, s.HandlePSQLError(s.Select, err, "select error")
		}
		for _, date := range dates {
			dt := &als.AlarmDateTime{
				Id:             date.ID,
				AlarmId:        date.AlarmId,
				AlarmDay:       date.AlarmDay,
				AlarmStartTime: date.AlarmStartTime,
				AlarmEndTime:   date.AlarmEndTime,
			}
			alarmDates = append(alarmDates, dt)
		}

		al := als.Alarm{
			Id:                alarm.ID,
			DevEui:            alarm.DevEui,
			MinTreshold:       alarm.MinTreshold,
			MaxTreshold:       alarm.MaxTreshold,
			Sms:               alarm.Sms,
			Email:             alarm.Email,
			Temperature:       alarm.Temperature,
			Humadity:          alarm.Humadity,
			Ec:                alarm.Ec,
			Door:              alarm.Door,
			WLeak:             alarm.WaterLeak,
			IsTimeLimitActive: alarm.IsTimeLimitActive,
			AlarmStartTime:    alarm.AlarmStartTime,
			AlarmStopTime:     alarm.AlarmStopTime,
			Notification:      alarm.Notification,
			UserID:            alarm.UserId,
			IpAddress:         alarm.IpAddress,
			ZoneCategoryID:    alarm.ZoneCategoryId,
			IsActive:          alarm.IsActive,
			Power:             alarm.Power,
			Current:           alarm.Current,
			Factor:            alarm.Factor,
			Voltage:           alarm.Voltage,
			Status:            alarm.Status,
			PowerSum:          alarm.PowerSum,
			AlarmDateTime:     alarmDates,
			NotificationSound: alarm.NotificationSound,
			Distance:          alarm.Distance,
			DefrostTime:       alarm.DefrostTime,
			Pressure:          alarm.Pressure,
		}
		returnAlarms = append(returnAlarms, &al)
	}
	return &als.GetAlarmListResponse{RespList: returnAlarms}, nil
}

// Implements the RPC method GetOrganizationAlarmList.
// Request takes organizationID as field and returns []OrganizationAlarm as response.
func (a *AlarmServerAPI) GetOrganizationAlarmList(ctx context.Context, req *als.GetOrganizationAlarmListRequest) (*als.GetOrganizationAlarmListResponse, error) {
	db := s.DB()

	var returnAlarms []*als.OrganizationAlarm
	var alarms []s.OrganizationAlarm
	var doorAlarms []s.DoorAlarm
	err := sqlx.Select(db, &alarms, `select z.zone_name, d.name as device_name, 0 AS time, ar.*
	from alarm_refactor2 as ar
		inner join device as d on d.dev_eui::text = '\x' || ar.dev_eui
		inner join zone as z on  d.dev_eui::text = any(z.devices)
		where d.organization_id = $1`, req.OrganizationID)
	if err != nil {
		return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms}, s.HandlePSQLError(s.Select, err, "select error")
	}

	err = sqlx.Select(db, &doorAlarms, `select z.zone_name, d.name as device_name, dta.* from door_time_alarm as dta
	inner join device as d on d.dev_eui::text = '\x' || dta.dev_eui
		inner join zone as z on  d.dev_eui::text = any(z.devices) where  dta.organization_id = $1`, req.OrganizationID)
	if err != nil {
		return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms}, s.HandlePSQLError(s.Select, err, "select error")
	}
	fmt.Println("1")
	for _, alarm := range alarms {
		var alarmDates []*als.AlarmDateTime
		var dates []s.AlarmDateFilter
		err := sqlx.Select(db, &dates, "select * from alarm_date_time where alarm_id = $1", alarm.ID)
		if err != nil {
			return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms}, s.HandlePSQLError(s.Select, err, "select error")
		}
		for _, date := range dates {
			dt := &als.AlarmDateTime{
				Id:             date.ID,
				AlarmId:        date.AlarmId,
				AlarmDay:       date.AlarmDay,
				AlarmStartTime: date.AlarmStartTime,
				AlarmEndTime:   date.AlarmEndTime,
			}
			alarmDates = append(alarmDates, dt)
		}
		al := als.OrganizationAlarm{
			Id:                alarm.ID,
			DevEui:            alarm.DevEui,
			MinTreshold:       alarm.MinTreshold,
			MaxTreshold:       alarm.MaxTreshold,
			Sms:               alarm.Sms,
			Email:             alarm.Email,
			Temperature:       alarm.Temperature,
			Humadity:          alarm.Humadity,
			Ec:                alarm.Ec,
			Door:              alarm.Door,
			WLeak:             alarm.WaterLeak,
			IsTimeLimitActive: alarm.IsTimeLimitActive,
			Notification:      alarm.Notification,
			DeviceName:        alarm.DeviceName,
			ZoneName:          alarm.ZoneName,
			UserName:          alarm.Username,
			UserID:            alarm.UserId,
			IpAddress:         alarm.IpAddress,
			AlarmDateTime:     alarmDates,
			Distance:          alarm.Distance,
			Time:              alarm.Time,
			IsActive:          alarm.IsActive,
			DefrostTime:       alarm.DefrostTime,
			Pressure:          alarm.Pressure,
			ZoneCategoryID:    alarm.ZoneCategoryId,
		}
		returnAlarms = append(returnAlarms, &al)
	}
	fmt.Println("3")
	for _, door := range doorAlarms {
		var alarmDates []*als.AlarmDateTime
		var dates []s.AlarmDateFilter
		err := sqlx.Select(db, &dates, "select * from door_alarm_date_time where alarm_id = $1", door.ID)
		if err != nil {
			return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms}, s.HandlePSQLError(s.Select, err, "select error")
		}
		for _, date := range dates {
			dt := &als.AlarmDateTime{
				Id:             date.ID,
				AlarmId:        date.AlarmId,
				AlarmDay:       date.AlarmDay,
				AlarmStartTime: date.AlarmStartTime,
				AlarmEndTime:   date.AlarmEndTime,
			}
			alarmDates = append(alarmDates, dt)
		}
		al := als.OrganizationAlarm{
			Id:                door.ID,
			DevEui:            door.DevEui,
			MinTreshold:       0,
			MaxTreshold:       0,
			Sms:               door.Sms,
			Email:             door.Email,
			Temperature:       false,
			Humadity:          false,
			Ec:                false,
			Door:              false,
			WLeak:             false,
			IsTimeLimitActive: door.IsTimeLimitActive,
			Notification:      door.Notification,
			DeviceName:        door.DeviceName,
			ZoneName:          door.ZoneName,
			UserName:          "",
			UserID:            door.UserId,
			IpAddress:         "",
			AlarmDateTime:     alarmDates,
			Distance:          false,
			Time:              door.Time,
		}
		returnAlarms = append(returnAlarms, &al)
	}
	fmt.Println("4")
	return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms}, nil
}
