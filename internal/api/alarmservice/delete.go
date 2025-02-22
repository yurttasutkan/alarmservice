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

// Implements the RPC method DeleteAlarm.
// Deletes the alarm according to userID and alarmID given by the request.
func (a *AlarmServerAPI) DeleteAlarm(ctx context.Context, req *als.DeleteAlarmRequest) (*empty.Empty, error) {
	db := s.DB()
	var al s.Alarm
	err := sqlx.Get(db, &al, "select * from alarm_refactor2 where id = $1", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	res, err := db.Exec("delete from  alarm_refactor2 where id = $1 ", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return &empty.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return &empty.Empty{}, nil
	}

	reqAlarm := &als.Alarm{
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
		WLeak:             al.WaterLeak,
		UserID:            al.UserId,
		IpAddress:         al.IpAddress,
		IsTimeLimitActive: al.IsTimeLimitActive,
		AlarmStartTime:    al.AlarmStartTime,
		AlarmStopTime:     al.AlarmStopTime,
		ZoneCategoryID:    al.ZoneCategoryId,
		IsActive:          al.IsActive,
	}
	s.CreateAlarmLog(ctx, db, reqAlarm, al.UserId, al.IpAddress, 1)

	_, err = db.Exec("update alarm_automation_rules set is_active = false where alarm_id = $1 ", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}
	_, err = db.Exec("delete from alarm_date_time where alarm_id = $1", req.AlarmID)
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
            RETURNING id
        )
        SELECT id FROM updated_rows;
    `
		// Execute the query and retrieve the result
		var idSlice []int64
		err := db.Select(&idSlice, query, i)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// Convert aggregatedIds to a slice of int64
		_, err = db.Exec(`delete from  public.alarm_refactor2
	WHERE id = ANY($1)
	AND cardinality(user_id) = 0`, pq.Int64Array(idSlice))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

// Implements the RPC method DeleteSensorAlarm.
// Deletes the Alarm according to the DevEui given by the request.
func (a *AlarmServerAPI) DeleteSensorAlarm(ctx context.Context, req *als.DeleteSensorAlarmRequest) (*empty.Empty, error) {
	db := s.DB()

	var alarmIds []int64
	err := sqlx.Select(db, &alarmIds, "select id from alarm_refactor2 where dev_eui = any($1) ", pq.Array(req.DevEuis))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	log.Println(req.DevEuis)
	_, err = db.Exec("delete from alarm_refactor2 where dev_eui = any($1)", pq.Array(req.DevEuis))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	_, err = db.Exec("delete from alarm_date_time where alarm_id = any($1)", pq.Array(alarmIds))
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

	log.Println(req.Zones)
	var devEuis []string

	err := sqlx.Select(db, &devEuis, `select devices from zone where zone_id = any($1)`, pq.Array(req.Zones))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Select, err, "select error")
	}

	log.Println(devEuis)
	res, err := db.Exec(`update alarm_refactor2 set is_active = false where  '\\x' || dev_eui = any($1)`, pq.Array(devEuis))
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	log.Println(ra)
	if ra == 0 {
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}

// Implements the RPC method DeleteAlarmDevEui.
// Deletes the alarm corresponding to the DevEui and UserID given in the request.
func (a *AlarmServerAPI) DeleteAlarmDevEui(ctx context.Context, req *als.DeleteAlarmDevEuiRequest) (*empty.Empty, error) {
	db := s.DB()

	res, err := db.Exec("delete from alarm_refactor2 where dev_eui = $1 and user_id = $2", req.Deveui, req.UserId)
	if err != nil {
		return &emptypb.Empty{}, s.HandlePSQLError(s.Delete, err, "delete error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}
