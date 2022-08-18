package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net"

// 	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
// 	"google.golang.org/grpc"
// )

// func main() {

// 	grpcServer := grpc.NewServer()

// 	als.RegisterAlarmServerServiceServer(grpcServer, &AlarmServerAPI{})
// 	lis, err := net.Listen("tcp", "172.22.0.18:9000")
// 	if err != nil {
// 		log.Fatalf("Failed to listen on port 9000: %v", err)
// 	}
// 	err = grpcServer.Serve(lis)
// 	if err != nil {
// 		log.Fatalf("Failed to serve on port 9000: %v", err)
// 	}
// }

// type AlarmServerAPI struct {
// }

// func NewAlarmServerAPI() *AlarmServerAPI {
// 	return &AlarmServerAPI{}
// }

// func (a *AlarmServerAPI) CreateAlarm(context context.Context, alarm *als.CreateAlarmRequest) (*als.CreateAlarmResponse, error) {
// 	fmt.Printf("Alarm: %s", alarm.Alarm)
// 	return &als.CreateAlarmResponse{AlarmResp: "Alarm Response!"}, nil
// }
