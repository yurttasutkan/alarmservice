package main

import (
	"context"
	"log"
	"net"

	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"google.golang.org/grpc"
)

func main() {

	grpcServer := grpc.NewServer()

	als.RegisterAlarmServerServiceServer(grpcServer, &AlarmServerAPI{})
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}
	go grpcServer.Serve(lis)
}

type AlarmServerAPI struct {
}

func NewAlarmServerAPI() *AlarmServerAPI {
	return &AlarmServerAPI{}
}

func (a *AlarmServerAPI) CreateAlarm(context context.Context, alarm *als.CreateAlarmRequest) (*als.CreateAlarmResponse, error) {
	log.Printf("Alarm: %s", alarm.Alarm)
	return &als.CreateAlarmResponse{AlarmResp: "Alarm Response!"}, nil
}
