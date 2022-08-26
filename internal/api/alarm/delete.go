package alarm

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/yurttasutkan/alarmservice/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *AlarmServerAPI) DeleteAlarm(ctx context.Context, req *als.DeleteAlarmRequest) (*empty.Empty, error) {
	db := storage.DB()

	var al Alarm
	err := sqlx.Get(db, &al, "select * from alarm_refactor where id = $1", req.AlarmID)
	if err != nil {
		return &empty.Empty{}, storage.HandlePSQLError(storage.Select, err, "select error")
	}

	res, err := db.Exec("update alarm_refactor set is_active = false where id = $1 and user_id = $2", req.AlarmID, req.UserID)
	if err != nil {
		return &empty.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return &empty.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return &empty.Empty{}, nil
	}
	return &empty.Empty{}, nil
}

func (a *AlarmServerAPI) DeleteAlarmDates(ctx context.Context, req *als.DeleteAlarmDatesRequest) (*empty.Empty, error) {
	db := storage.DB()

	_, err := db.Exec("delete from alarm_date_time where alarm_id = $1", req.AlarmId)
	if err != nil {
		return &empty.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
	}

	return &empty.Empty{}, nil
}

func (a *AlarmServerAPI) DeleteUserAlarm(ctx context.Context, req *als.DeleteUserAlarmRequest) (*empty.Empty, error) {
	db := storage.DB()

	res, err := db.Exec("update alarm_refactor set is_active = false where user_id = any($1)", pq.Array(req.UserIds))
	if err != nil {
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "get rows affected error")
	}
	_, err = db.Exec(`INSERT INTO public.alarm_change_logs(
		dev_eui, min_treshold, max_treshold, user_id, ip_address, is_deleted, sms, temperature, humadity, ec, door, w_leak)
	   select dev_eui,  min_treshold, max_treshold, user_id, '', 1, sms,temperature, humadity, ec, door, w_leak 
	   from alarm_refactor where user_id = any($1) and is_active = true`, pq.Array(req.UserIds))
	if err != nil {
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
	}
	if ra == 0 {
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}

func (a *AlarmServerAPI) DeleteSensorAlarm(ctx context.Context, req *als.DeleteSensorAlarmRequest) (*empty.Empty,error){
	db := storage.DB()

	log.Println(req.DevEuis)
	res, err := db.Exec("update alarm_refactor set is_active = false where dev_eui = any($1)", pq.Array(req.DevEuis))
	if err != nil {
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
	}
	_, err = db.Exec(`INSERT INTO public.alarm_change_logs(
		dev_eui, min_treshold, max_treshold, user_id, ip_address, is_deleted, sms, temperature, humadity, ec, door, w_leak)
	   select dev_eui,  min_treshold, max_treshold, user_id, '', 1, sms,temperature, humadity, ec, door, w_leak 
	   from alarm_refactor where dev_eui = any($1) and is_active = true`, pq.Array(req.DevEuis))
	if err != nil {
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
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

func (a *AlarmServerAPI) DeleteZoneAlarm(ctx context.Context,req *als.DeleteZoneAlarmRequest) (*empty.Empty,error){
	db := storage.DB()

	log.Println(req.Zones)
	var devEuis []string

	err := sqlx.Select(db, &devEuis, `select devices from zone where zone_id = any($1)`, pq.Array(req.Zones))
	if err!=nil{
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Select, err, "select error")
	}

	log.Println(devEuis)
	res, err := db.Exec(`update alarm_refactor set is_active = false where  '\\x' || dev_eui = any($1)`, pq.Array(devEuis))
	if err != nil {
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
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

func (a *AlarmServerAPI) DeleteAlarmDevEui(ctx context.Context, req *als.DeleteAlarmDevEuiRequest) (*empty.Empty,error){
	db := storage.DB()

	res, err := db.Exec("delete from alarm_refactor where dev_eui = $1 and user_id = $2", req.Deveui, req.UserId)
	if err != nil {
		return &emptypb.Empty{}, storage.HandlePSQLError(storage.Delete, err, "delete error")
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