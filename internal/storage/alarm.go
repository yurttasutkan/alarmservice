package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
)

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


func CreateAlarmLog(ctx context.Context, db sqlx.Ext, a *als.Alarm, userID int64, ipAddress string, isDeleted int64) error {
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
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) `, a.DevEui, a.MinTreshold,
		a.MaxTreshold, userID, ipAddress, isDeleted, a.Sms, a.Temperature, a.Humadity, a.Ec, a.Door, a.WLeak)
	if err != nil {
		return HandlePSQLError(Insert, nil, "insert error")
	}

	return nil
}
