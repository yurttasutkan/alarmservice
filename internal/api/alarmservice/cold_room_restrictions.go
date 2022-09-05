package alarmservice

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	s "github.com/yurttasutkan/alarmservice/internal/storage"
)

//Implements the RPC method CreateColdRoomRestrictions.
//Inserts into cold_room_restrictions table with given parameters in the request.
func (a *AlarmServerAPI) CreateColdRoomRestrictions(ctx context.Context, req *als.CreateColdRoomRestrictionsRequest) (*empty.Empty, error) {
	db := s.DB()
	coldRes := req.ColdRes
	_, err := db.Exec(`insert into cold_room_restrictions(
		dev_eui,
		alarm_id,
		defronst_time,
		defrost_frequency,
		alarm_time
	) values ($1, $2, $3, $4, $5)`, coldRes.DevEui, coldRes.AlarmID, coldRes.DefrostTime, coldRes.DefrostFrq, coldRes.AlarmTime)
	if err != nil {
		log.Fatalf("Insert error %v", err)
	}

	return &empty.Empty{}, nil
}
