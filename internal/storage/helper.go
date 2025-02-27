package storage

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

// LogAudit logs changes into the audit logs table.
func (s *AlarmServerAPI) LogAudit(db *sqlx.DB, userID int64, changeType, tableName string, recordID int64, previousValue, newValue interface{}, ipAddress, reason string) error {
	previousJSON, _ := json.Marshal(previousValue)
	newJSON, _ := json.Marshal(newValue)

	query := `
        INSERT INTO audit_logs (user_id, change_type, table_name, record_id, previous_value, new_value, ip_address, reason)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := db.Exec(query, userID, changeType, tableName, recordID, previousJSON, newJSON, ipAddress, reason)
	return err
}

// AlarmServerAPI implements the Alarm server API.
type AlarmServerAPI struct {
}
