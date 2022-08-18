package api

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	"github.com/yurttasutkan/alarmservice/internal/storage"
)

//AlarmServerAPI implements the Alarm server API.
type AlarmServerAPI struct {
}

//Creates a new AlarmServerAPI
func NewAlarmServerAPI() *AlarmServerAPI {
	return &AlarmServerAPI{}
}

// Creates the given Alarm.
func (a *AlarmServerAPI) CreateAlarm(context context.Context, id *als.CreateAlarmRequest) (*empty.Empty, error) {
	var testID int64
	db := storage.DB()

	err := sqlx.Get(db, &testID, "select user_id from alarm_refactor where user_id = $1", id.UserID)
	if err != nil {
		log.Printf("SQL get error: %v", err)
	}
	log.Printf("User id is: %v", testID)

	return &empty.Empty{}, nil
}

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
