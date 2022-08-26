package alarm

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/yurttasutkan/alarmservice/internal/storage"
)

// GetAlarm gets the alarm via alarmID given by GetAlarmRequest.
func (a *AlarmServerAPI) GetAlarm(ctx context.Context, alReq *als.GetAlarmRequest) (*als.GetAlarmResponse, error) {
	db := storage.DB()

	var resp als.GetAlarmResponse

	err := sqlx.Get(db, &resp.Alarm, "select * from alarm_refactor where id = $1 and is_active = true", alReq.AlarmID)
	if err != nil {
		return &als.GetAlarmResponse{Alarm: resp.Alarm}, storage.HandlePSQLError(storage.Select, err, "select error")
	}
	return &resp, nil
}

func (a *AlarmServerAPI) GetAlarmLogs(ctx context.Context, req *als.GetAlarmLogsRequest) (*als.GetAlarmLogsResponse, error) {
	db := storage.DB()

	var logs []AlarmLogs
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
		return &als.GetAlarmLogsResponse{RespLog: result}, storage.HandlePSQLError(storage.Select, err, "select error")
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
			WLeak:       log.W_leak,
		}
		var err error
		item.SubmissionDate, err = ptypes.TimestampProto(log.SubmissionDate)
		if err != nil {
			return &als.GetAlarmLogsResponse{RespLog: result}, storage.HandlePSQLError(storage.Select, err, "select error")
		}
		result = append(result, &item)
	}
	return &als.GetAlarmLogsResponse{RespLog: result}, nil
}

func (a *AlarmServerAPI) GetAlarmDates(ctx context.Context, req *als.GetAlarmDatesRequest) (*als.GetAlarmDatesResponse, error) {
	db := storage.DB()

	var returnDates []*als.AlarmDateTime
	var alarmDates []AlarmDateFilter
	err := sqlx.Select(db, &alarmDates, "select * from alarm_date_time where alarm_id = $1", req.AlarmId)
	if err != nil {
		return &als.GetAlarmDatesResponse{RespDate: returnDates}, storage.HandlePSQLError(storage.Select, err, "select error")
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

func (a *AlarmServerAPI) GetAlarmList(ctx context.Context, req *als.GetAlarmListRequest) (*als.GetAlarmListResponse, error) {
	db := storage.DB()

	var returnDates []*als.AlarmDateTime
	filters := AlarmFilters{
		Limit:  int(req.Filter.Limit),
		DevEui: req.Filter.DevEui,
		UserID: req.Filter.UserID,
	}
	var returnAlarms []*als.Alarm
	var alarms []Alarm
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
	select *
	from alarm_refactor
	`+filters.SQL(), filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}
	err = sqlx.Select(db, &alarms, query, args...)
	if err != nil {
		return &als.GetAlarmListResponse{RespList: returnAlarms, RespDate: returnDates}, storage.HandlePSQLError(storage.Select, err, "select error")
	}
	for _, alarm := range alarms {
		var alarmDates []AlarmDateFilter
		err := sqlx.Select(db, &alarmDates, "select * from alarm_date_time where alarm_id = $1", alarm.ID)
		if err != nil {
			return &als.GetAlarmListResponse{RespList: returnAlarms, RespDate: returnDates}, storage.HandlePSQLError(storage.Select, err, "select error")
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
			WLeak:             alarm.W_leak,
			IsTimeLimitActive: alarm.IsTimeLimitActive,
			AlarmStartTime:    alarm.AlarmStartTime,
			AlarmStopTime:     alarm.AlarmStopTime,
			Notification:      alarm.Notification,
		}
		returnAlarms = append(returnAlarms, &al)
	}
	return &als.GetAlarmListResponse{RespList: returnAlarms, RespDate: returnDates}, nil
}

func (a *AlarmServerAPI) GetOrganizationAlarmList(ctx context.Context, req *als.GetOrganizationAlarmListRequest) (*als.GetOrganizationAlarmListResponse, error) {
	db := storage.DB()

	var returnAlarms []*als.OrganizationAlarm
	var alarms []OrganizationAlarm
	var returnDates []*als.AlarmDateTime
	err := sqlx.Select(db, &alarms, `select u.username, z.zone_name, d.name as device_name, ar.*
	from alarm_refactor as ar
		inner join public.user as u on ar.user_id = u.id
		inner join organization_user as ou on ou.user_id = u.id
		inner join device as d on d.dev_eui::text = '\x' || ar.dev_eui
		inner join zone as z on  d.dev_eui::text = any(z.devices)
		where ou.organization_id = $1 and ar.is_active = true;`, req.OrganizationID)
	if err != nil {
		return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms, RespDate: returnDates}, storage.HandlePSQLError(storage.Select, err, "select error")
	}
	for _, alarm := range alarms {
		var alarmDates []AlarmDateFilter

		err := sqlx.Select(db, &alarmDates, "select * from alarm_date_time where alarm_id = $1", alarm.ID)
		if err != nil {
			return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms, RespDate: returnDates}, storage.HandlePSQLError(storage.Select, err, "select error")
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
			WLeak:             alarm.W_leak,
			IsTimeLimitActive: alarm.IsTimeLimitActive,
			Notification:      alarm.Notification,
			DeviceName:        alarm.DeviceName,
			ZoneName:          alarm.ZoneName,
			UserName:          alarm.Username,
		}
		returnAlarms = append(returnAlarms, &al)
	}
	return &als.GetOrganizationAlarmListResponse{RespList: returnAlarms, RespDate: returnDates}, nil
}
