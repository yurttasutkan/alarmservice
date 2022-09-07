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

	err := sqlx.Get(db, &resp.Alarm, "select * from alarm_refactor where id = $1 and is_active = true", alReq.AlarmID)
	if err != nil {
		return &als.GetAlarmResponse{Alarm: resp.Alarm}, s.HandlePSQLError(s.Select, err, "select error")
	}
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
	from alarm_refactor
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
			AlarmDateTime:     alarmDates,
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
	err := sqlx.Select(db, &alarms, `select u.username, z.zone_name, d.name as device_name, ar.*
	from alarm_refactor as ar
		inner join public.user as u on ar.user_id = u.id
		inner join organization_user as ou on ou.user_id = u.id
		inner join device as d on d.dev_eui::text = '\x' || ar.dev_eui
		inner join zone as z on  d.dev_eui::text = any(z.devices)
		where ou.organization_id = $1 and ar.is_active = true;`, req.OrganizationID)
	if err != nil {
		return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms}, s.HandlePSQLError(s.Select, err, "select error")
	}
	

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
		}
		returnAlarms = append(returnAlarms, &al)
	}
	return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms}, nil
}
