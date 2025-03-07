package storage

import (
	"context"
	"encoding/json"
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
	filters = append(filters, fmt.Sprint(" and ", f.UserID, " = any(user_id)"))
	if f.Limit != 0 {
		filters = append(filters, fmt.Sprint(" LIMIT ", f.Limit))
	}

	return " where is_active = true and  " + strings.Join(filters, " ")
}

func CreateAlarmLog(ctx context.Context, db sqlx.Ext, a *als.Alarm, userID []int64, ipAddress string, isDeleted int64) error {
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

// LogAudit logs changes into the alarm_audit_log table
func LogAudit(db *sqlx.DB, alarmID, userID int64, changeType string, previousValue, newValue interface{}) error {
	// Convert old and new values to JSON
	previousJSON, _ := json.Marshal(previousValue)
	newJSON, _ := json.Marshal(newValue)

	query := `
		INSERT INTO alarm_audit_log (alarm_id, change_type, changed_by, old_values, new_values)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.Exec(query, alarmID, changeType, userID, previousJSON, newJSON)
	return err
}
