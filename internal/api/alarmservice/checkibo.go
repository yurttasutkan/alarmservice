package alarmservice

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
)

func (a *AlarmServerAPI) CheckAlarmibo(ctx context.Context, req *als.CheckAlarmIbo) (*empty.Empty, error) {

	fmt.Println("IBO WORKKK")
	return &empty.Empty{}, nil

}
