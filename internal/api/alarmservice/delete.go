package alarmservice

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *AlarmServerAPI) DeleteAlarm(ctx context.Context, req *als.DeleteAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	var currentAlarm s.Alarm
	// Get the previous values of the alarm
	err := sqlx.Get(db, &currentAlarm, "select * from alarm_refactor2 where id = $1", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	// Log the delete action
	err = s.LogAudit(db, currentAlarm.ID, req.UserID, "DELETE", currentAlarm, nil)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Insert, err, "insert error")
	}

	// Delete from `alarm_refactor2`
	res, err := db.Exec("DELETE FROM alarm_refactor2 WHERE id = $1", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	// Check if the alarm was actually deleted
	ra, err := res.RowsAffected()
	if err != nil {
		return &empty.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return &empty.Empty{}, errors.New("no rows deleted: alarm may not exist")
	}

	// Log the alarm deletion event
	reqAlarm := &als.Alarm{
		Id:                currentAlarm.ID,
		DevEui:            currentAlarm.DevEui,
		MinTreshold:       currentAlarm.MinTreshold,
		MaxTreshold:       currentAlarm.MaxTreshold,
		Sms:               currentAlarm.Sms,
		Email:             currentAlarm.Email,
		Notification:      currentAlarm.Notification,
		Temperature:       currentAlarm.Temperature,
		Humadity:          currentAlarm.Humadity,
		Ec:                currentAlarm.Ec,
		Door:              currentAlarm.Door,
		WLeak:             currentAlarm.WaterLeak,
		UserID:            currentAlarm.UserId,
		IpAddress:         currentAlarm.IpAddress,
		IsTimeLimitActive: currentAlarm.IsTimeLimitActive,
		AlarmStartTime:    currentAlarm.AlarmStartTime,
		AlarmStopTime:     currentAlarm.AlarmStopTime,
		ZoneCategoryID:    currentAlarm.ZoneCategoryId,
		IsActive:          currentAlarm.IsActive,
	}
	s.CreateAlarmLog(ctx, db, reqAlarm, currentAlarm.UserId, currentAlarm.IpAddress, 1)

	// Deactivate automation rules related to this alarm
	_, err = db.Exec("UPDATE alarm_automation_rules SET is_active = false WHERE alarm_id = $1", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Update, err, "update error")
	}

	// Delete related records in `alarm_date_time`
	_, err = db.Exec("DELETE FROM alarm_date_time WHERE alarm_id = $1", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	return &empty.Empty{}, nil
}

// Implements the RPC method DeleteAlarmDates.
// Deletes the AlarmDateTime according to the AlarmID given by the request.
func (a *AlarmServerAPI) DeleteAlarmDates(ctx context.Context, req *als.DeleteAlarmDatesRequest) (*empty.Empty, error) {
	db := s.DB()

	_, err := db.Exec("delete from alarm_date_time where alarm_id = $1", req.AlarmId)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	return &empty.Empty{}, nil
}

// Implements the RPC method DeleteUserAlarm.
// Deletes the Alarm according to the UserID given by the request.
func (a *AlarmServerAPI) DeleteUserAlarm(ctx context.Context, req *als.DeleteUserAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	for _, i := range req.UserIds {
		query := `
		WITH updated_rows AS (
			UPDATE public.alarm_refactor2
			SET user_id = array_remove(user_id, $1::bigint)
			WHERE $1 = ANY(user_id)
			RETURNING *
		)
		SELECT * FROM updated_rows;
		`

		// Fetch updated alarm records before deleting
		var alarms []s.Alarm
		err := sqlx.Select(db, &alarms, query, i)
		if err != nil {
			fmt.Println("Update error:", err)
			return nil, err
		}

		// Log the delete action before actually deleting
		for _, al := range alarms {
			if len(al.UserId) == 1 { // If it's the last user, it will be deleted
				err = s.LogAudit(db, al.ID, req.UserSentId, "DELETE", al, nil)
				if err != nil {
					return &empty.Empty{}, s.HandlePSQLError(s.Insert, err, "insert error")
				}
			}
		}

		// Delete alarms where `user_id` is now empty
		_, err = db.Exec(`DELETE FROM public.alarm_refactor2
			WHERE id = ANY($1) AND cardinality(user_id) = 0`, pq.Array(getAlarmIDs(alarms)))
		if err != nil {
			fmt.Println("Delete error:", err)
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

// Helper function to extract IDs from alarms
func getAlarmIDs(alarms []s.Alarm) []int64 {
	var ids []int64
	for _, alarm := range alarms {
		ids = append(ids, alarm.ID)
	}
	return ids
}

func (a *AlarmServerAPI) DeleteSensorAlarm(ctx context.Context, req *als.DeleteSensorAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	// Fetch all alarms that match the given DevEUIs
	var alarms []s.Alarm
	err := db.Select(&alarms, "SELECT * FROM alarm_refactor2 WHERE dev_eui = ANY($1)", pq.Array(req.DevEuis))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	// If no alarms found, return early
	if len(alarms) == 0 {
		return &emptypb.Empty{}, errors.New("no alarms found for given DevEUIs")
	}

	// Extract alarm IDs for later deletions
	var alarmIds []int64
	for _, al := range alarms {
		alarmIds = append(alarmIds, al.ID)
	}

	// Log the delete action before actual deletion
	for _, al := range alarms {
		err = s.LogAudit(db, al.ID, req.UserId, "DELETE", al, nil)
		if err != nil {
			return &empty.Empty{}, s.HandlePSQLError(s.Insert, err, "insert error")
		}
	}

	// Delete alarms from `alarm_refactor2`
	res, err := db.Exec("DELETE FROM alarm_refactor2 WHERE dev_eui = ANY($1)", pq.Array(req.DevEuis))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	// Check if alarms were deleted
	ra, err := res.RowsAffected()
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return &emptypb.Empty{}, errors.New("no rows deleted: alarms may not exist")
	}

	// Delete related records from `alarm_date_time`
	_, err = db.Exec("DELETE FROM alarm_date_time WHERE alarm_id = ANY($1)", pq.Array(alarmIds))
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}
	// _, err = db.Exec(`INSERT INTO public.alarm_change_logs(
	// 	dev_eui, min_treshold, max_treshold, user_id, ip_address, is_deleted, sms, temperature, humadity, ec, door, w_leak)
	//    select dev_eui,  min_treshold, max_treshold, user_id, '', 1, sms,temperature, humadity, ec, door, w_leak
	//    from alarm_refactor2 where dev_eui = any($1) and is_active = true`, pq.Array(req.DevEuis))
	// if err != nil {
	// 	return &emptypb.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	// }
	// ra, err := res.RowsAffected()
	// if err != nil {
	// 	return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	// }
	// if ra == 0 {
	// 	return &emptypb.Empty{}, nil
	// }
	return &emptypb.Empty{}, nil
}

// Implements the RPC method DeleteZoneAlarm.
// Deletes alarms that are in the given zone by the request.
func (a *AlarmServerAPI) DeleteZoneAlarm(ctx context.Context, req *als.DeleteZoneAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	log.Println("Zones received:", req.Zones)

	// Get device EUIs from the given zones
	var devEuis []string
	err := db.Select(&devEuis, `SELECT devices FROM zone WHERE zone_id = ANY($1)`, pq.Array(req.Zones))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	if len(devEuis) == 0 {
		return &emptypb.Empty{}, errors.New("no devices found in given zones")
	}

	log.Println("Device EUIs from zones:", devEuis)

	// Fetch alarms that will be updated
	var alarms []s.Alarm
	err = db.Select(&alarms, `SELECT * FROM alarm_refactor2 WHERE ('\\x' || dev_eui) = ANY($1)`, pq.Array(devEuis))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	if len(alarms) == 0 {
		return &emptypb.Empty{}, errors.New("no alarms found for the devices in given zones")
	}

	// Log the update action before modifying alarms
	for _, al := range alarms {
		updatedAlarm := al
		updatedAlarm.IsActive = false // Simulating the update

		err = s.LogAudit(db, al.ID, req.UserId, "UPDATE", al, updatedAlarm)
		if err != nil {
			return &emptypb.Empty{}, s.HandlePSQLError(s.Insert, err, "insert error")
		}
	}

	// Update alarms to set `is_active = false`
	res, err := db.Exec(`UPDATE alarm_refactor2 SET is_active = false WHERE ('\\x' || dev_eui) = ANY($1)`, pq.Array(devEuis))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Update, err, "update error")
	}

	// Check affected rows
	ra, err := res.RowsAffected()
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	log.Println("Rows updated:", ra)

	if ra == 0 {
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}

// Implements the RPC method DeleteAlarmDevEui.
// Deletes the alarm corresponding to the DevEui and UserID given in the request.
func (a *AlarmServerAPI) DeleteAlarmDevEui(ctx context.Context, req *als.DeleteAlarmDevEuiRequest) (*empty.Empty, error) {
	db := s.DB()

	// Fetch the alarm before deletion
	var alarm s.Alarm
	err := db.Get(&alarm, "SELECT * FROM alarm_refactor2 WHERE dev_eui = $1 AND user_id = $2", req.Deveui, req.UserId)
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	// Log the delete action before deleting the record
	err = s.LogAudit(db, alarm.ID, req.UserId, "DELETE", alarm, nil)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Insert, err, "insert error")
	}

	// Delete the alarm
	res, err := db.Exec("DELETE FROM alarm_refactor2 WHERE dev_eui = $1 AND user_id = $2", req.Deveui, req.UserId)
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	// Check if the alarm was actually deleted
	ra, err := res.RowsAffected()
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return &emptypb.Empty{}, errors.New("no alarm deleted: alarm may not exist")
	}

	return &emptypb.Empty{}, nil
}
