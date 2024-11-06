package storage

import (
	"time"

	"github.com/lib/pq"
)

type Alarm struct {
	ID                int64         `db:"id"`
	DevEui            string        `db:"dev_eui"`
	MinTreshold       float32       `db:"min_treshold"`
	MaxTreshold       float32       `db:"max_treshold"`
	Sms               bool          `db:"sms"`
	Email             bool          `db:"email"`
	Notification      bool          `db:"notification"`
	Temperature       bool          `db:"temperature"`
	Humadity          bool          `db:"humadity"`
	Ec                bool          `db:"ec"`
	Door              bool          `db:"door"`
	WaterLeak         bool          `db:"w_leak"`
	UserId            pq.Int64Array `db:"user_id"`
	IpAddress         string        `db:"ip_address"`
	IsTimeLimitActive bool          `db:"is_time_limit_active"`
	AlarmStartTime    float32       `db:"alarm_start_time"`
	AlarmStopTime     float32       `db:"alarm_stop_time"`
	ZoneCategoryId    int64         `db:"zone_category"`
	IsActive          bool          `db:"is_active"`
	Pressure          bool          `db:"pressure"`
	Current           float32       `db:"current"`
	Factor            float32       `db:"factor"`
	Power             float32       `db:"power"`
	Voltage           float32       `db:"voltage"`
	Status            int64         `db:"status"`
	PowerSum          float32       `db:"power_sum"`
	NotificationSound string        `db:"notification_sound"`
	Distance          bool          `db:"distance"`
}
type DoorAlarm struct {
	ID                int64         `db:"id"`
	DevEui            string        `db:"dev_eui"`
	Sms               bool          `db:"sms"`
	Email             bool          `db:"email"`
	Notification      bool          `db:"notification"`
	Time              int64         `db:"time"`
	UserId            pq.Int64Array `db:"user_id"`
	IsActive          bool          `db:"is_active"`
	ZoneName          string        `db:"zone_name"`
	DeviceName        string        `db:"device_name"`
	SubmissionDate    time.Time     `db:"submission_time"`
	OrganizationId    int64         `db:"organization_id"`
	IsTimeLimitActive bool          `db:"is_time_limit_active"`
}
type OrganizationAlarm struct {
	ID                int64         `db:"id"`
	DevEui            string        `db:"dev_eui"`
	MinTreshold       float32       `db:"min_treshold"`
	MaxTreshold       float32       `db:"max_treshold"`
	Sms               bool          `db:"sms"`
	Email             bool          `db:"email"`
	Notification      bool          `db:"notification"`
	Temperature       bool          `db:"temperature"`
	Humadity          bool          `db:"humadity"`
	Ec                bool          `db:"ec"`
	Door              bool          `db:"door"`
	WaterLeak         bool          `db:"w_leak"`
	UserId            pq.Int64Array `db:"user_id"`
	IpAddress         string        `db:"ip_address"`
	IsTimeLimitActive bool          `db:"is_time_limit_active"`
	AlarmStartTime    float32       `db:"alarm_start_time"`
	AlarmStopTime     float32       `db:"alarm_stop_time"`
	ZoneCategoryId    int64         `db:"zone_category"`
	IsActive          bool          `db:"is_active"`
	ZoneName          string        `db:"zone_name"`
	DeviceName        string        `db:"device_name"`
	Username          string        `db:"username"`
	Pressure          bool          `db:"pressure"`
	Current           float32       `db:"current"`
	Factor            float32       `db:"factor"`
	Power             float32       `db:"power"`
	Voltage           float32       `db:"voltage"`
	Status            int64         `db:"status"`
	PowerSum          float32       `db:"power_sum"`
	NotificationSound string        `db:"notification_sound"`
	Distance          bool          `db:"distance"`
	Time              int64         `db:"time"`
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
	WaterLeak         bool    `db:"w_leak"`
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
	Distance          bool    `db:"distance"`
}

type AlarmDateFilter struct {
	ID             int64   `db:"id"`
	AlarmId        int64   `db:"alarm_id"`
	AlarmDay       int64   `db:"alarm_day"`
	AlarmStartTime float32 `db:"start_time"`
	AlarmEndTime   float32 `db:"end_time"`
}
type ColdRoomRestrictions struct {
	ID               int64  `db:"id"`
	DevEui           string `db:"dev_eui"`
	AlarmId          int64  `db:"alarm_id"`
	DefrostTime      int64  `db:"defronst_time"`
	DefrostFrequency int64  `db:"defrost_frequency"`
	AlarmTime        int64  `db:"alarm_time"`
}
type UtkuStruct struct {
	ID          int64   `db:"id"`
	DevEui      string  `db:"dev_eui"`
	AlarmId     int64   `db:"alarm_id"`
	LocalMaxVal float32 `db:"local_max_value"`
	Counter     int64   `db:"counter"`
	CntLimit    float32 `db:"cnt_limit"`
}

type AlarmLogs struct {
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
	WaterLeak      bool      `db:"w_leak"`
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
