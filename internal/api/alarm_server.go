package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/yurttasutkan/alarmservice/internal/storage"
)

//AlarmServerAPI implements the Alarm server API.
type AlarmServerAPI struct {
}

//Creates a new AlarmServerAPI
func NewAlarmServerAPI() *AlarmServerAPI {
	return &AlarmServerAPI{}
}

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

// Implements the RPC method CreateAlarm.
func (a *AlarmServerAPI) CreateAlarm(context context.Context, alarm *als.CreateAlarmRequest) (*als.CreateAlarmResponse, error) {
	db := storage.DB()
	var al Alarm
	var returnID int64
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
		alarm.Alarm.DevEui,
		alarm.Alarm.MinTreshold,
		alarm.Alarm.MaxTreshold,
		alarm.Alarm.Sms,
		alarm.Alarm.Email,
		alarm.Alarm.Temperature,
		alarm.Alarm.Humadity,
		alarm.Alarm.Ec,
		alarm.Alarm.Door,
		alarm.Alarm.WLeak,
		alarm.Alarm.UserID,
		alarm.Alarm.IsTimeLimitActive,
		alarm.Alarm.AlarmStartTime,
		alarm.Alarm.AlarmStopTime,
		alarm.Alarm.ZoneCategoryID,
		alarm.Alarm.Notification,
	).Scan(&returnID)
	if err != nil {
		fmt.Println(err)
	}

	resp := als.CreateAlarmResponse{
		Alarm: &als.Alarm{
			Id:                al.ID,
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
			WLeak:             al.W_leak,
			UserID:            al.UserId,
			IpAddress:         al.IpAddress,
			IsTimeLimitActive: al.IsTimeLimitActive,
			AlarmStartTime:    al.AlarmStartTime,
			AlarmStopTime:     al.AlarmStopTime,
			ZoneCategoryID:    al.ZoneCategoryId,
			IsActive:          al.IsActive,
		},
	}
	return &resp, nil
}

func (a *AlarmServerAPI) GetAlarm(context.Context, *als.GetAlarmRequest) (*als.GetAlarmResponse, error) {
	var resp *als.GetAlarmResponse
	return resp, nil
}

// func (a *AlarmServerAPI) GetAlarm(context context.Context, req *als.GetAlarmRequest) (*als.GetAlarmResponse,error){
// 	var al Alarm
// 	db := storage.DB()
// 	err := sqlx.Get(db, &al, "select * from alarm_refactor where id = $1 and is_active = true", req.AlarmID)
// 	if err != nil {
// 		log.Fatalf("Sqlx get error: %v", err)
// 	}

// 	resp := als.GetAlarmResponse{
// 		Alarm: &als.Alarm{
// 			Id:              al.ID,
// 			DevEui:          al.DevEui,
// 			MinTreshold:     al.MinTreshold,
// 			MaxTreshold:     al.MaxTreshold,
// 			Sms:             al.Sms,
// 			Email:           al.Email,
// 			Temperature:     al.Temperature,
// 			Humadity:        al.Humadity,
// 			Ec:              al.Ec,
// 			Door:            al.Door,
// 			WLeak:           al.W_leak,
// 			IsTimeScheduled: al.IsTimeLimitActive,
// 			StartTime:       al.AlarmStartTime,
// 			EndTime:         al.AlarmStopTime,
// 		},
// 	}
// 	return &resp,nil
// }

type Alarm struct {
	ID                int64   `db:"id"`
	DevEui            string  `db:"dev_eui"`
	MinTreshold       float32 `db:"min_treshold"`
	MaxTreshold       float32 `db:"max_treshold"`
	Sms               bool    `db:"sms"`
	Email             bool    `db:"email"`
	Notification      bool    `db:"notification"`
	Temperature       bool    `db:"temperature"`
	Humadity          bool    `db:"humadity"`
	Ec                bool    `db:"ec"`
	Door              bool    `db:"door"`
	W_leak            bool    `db:"w_leak"`
	UserId            int64   `db:"user_id"`
	IpAddress         string  `db:"ip_address"`
	IsTimeLimitActive bool    `db:"is_time_limit_active"`
	AlarmStartTime    float32 `db:"alarm_start_time"`
	AlarmStopTime     float32 `db:"alarm_stop_time"`
	ZoneCategoryId    int64   `db:"zone_category"`
	IsActive          bool    `db:"is_active"`
}

type OrganizationAlarm struct {
	ID                int64   `db:"id"`
	DevEui            string  `db:"dev_eui"`
	MinTreshold       float32 `db:"min_treshold"`
	MaxTreshold       float32 `db:"max_treshold"`
	Sms               bool    `db:"sms"`
	Email             bool    `db:"email"`
	Notification      bool    `db:"notification"`
	Temperature       bool    `db:"temperature"`
	Humadity          bool    `db:"humadity"`
	Ec                bool    `db:"ec"`
	Door              bool    `db:"door"`
	W_leak            bool    `db:"w_leak"`
	UserId            int64   `db:"user_id"`
	IpAddress         string  `db:"ip_address"`
	IsTimeLimitActive bool    `db:"is_time_limit_active"`
	AlarmStartTime    float32 `db:"alarm_start_time"`
	AlarmStopTime     float32 `db:"alarm_stop_time"`
	ZoneCategoryId    int64   `db:"zone_category"`
	IsActive          bool    `db:"is_active"`
	ZoneName          string  `db:"zone_name"`
	DeviceName        string  `db:"device_name"`
	Username          string  `db:"username"`
}
type AlarmWithDates struct {
	ID                int64   `db:"id"`
	DevEui            string  `db:"dev_eui"`
	MinTreshold       float32 `db:"min_treshold"`
	MaxTreshold       float32 `db:"max_treshold"`
	Sms               bool    `db:"sms"`
	Notification      bool    `db:"notification"`
	Email             bool    `db:"email"`
	Temperature       bool    `db:"temperature"`
	Humadity          bool    `db:"humadity"`
	Ec                bool    `db:"ec"`
	Door              bool    `db:"door"`
	W_leak            bool    `db:"w_leak"`
	UserId            int64   `db:"user_id"`
	IpAddress         string  `db:"ip_address"`
	IsTimeLimitActive bool    `db:"is_time_limit_active"`
	ZoneCategoryId    int64   `db:"zone_category"`
	AlarmDay          int64   `db:"alarm_day"`
	AlarmStartTime2   float32 `db:"alarm_start_time"`
	AlarmStopTime2    float32 `db:"alarm_stop_time"`
	AlarmStartTime    float32 `db:"start_time"`
	AlarmEndTime      float32 `db:"end_time"`
	IsActive          bool    `db:"is_active"`
}
type AlarmDateFilter struct {
	ID             int64   `db:"id"`
	AlarmId        int64   `db:"alarm_id"`
	AlarmDay       int64   `db:"alarm_day"`
	AlarmStartTime float32 `db:"start_time"`
	AlarmEndTime   float32 `db:"end_time"`
}
type ColdRoomRestrictions struct {
	ID          int64  `db:"id"`
	DevEui      string `db:"dev_eui"`
	AlarmId     int64  `db:"alarm_id"`
	DefrostTime int64  `db:"defronst_time"`
	DefrostFrq  int64  `db:"defrost_frequency"`
	AlarmTime   int64  `db:"alarm_time"`
}
type AlarmLogsStruct struct {
	DevEui         string    `db:"dev_eui"`
	MinTreshold    float32   `db:"min_treshold"`
	MaxTreshold    float32   `db:"max_treshold"`
	UserId         int64     `db:"user_id"`
	IpAddress      string    `db:"ip_address"`
	IsDeleted      int64     `db:"is_deleted"`
	Sms            bool      `db:"sms"`
	Temperature    bool      `db:"temperature"`
	Humadity       bool      `db:"humadity"`
	Ec             bool      `db:"ec"`
	Door           bool      `db:"door"`
	W_leak         bool      `db:"w_leak"`
	SubmissionDate time.Time `db:"submission_date"`
}

// AlarmFilters filters
type AlarmFilters struct {
	Limit  int    `db:"limit"`
	DevEui string `db:"dev_eui"`
	UserID int64  `db:"user_id"`
}

// SMSRequestBody ...
type SMSRequestBody struct {
	From      string `json:"from"`
	Text      string `json:"text"`
	To        string `json:"to"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}
