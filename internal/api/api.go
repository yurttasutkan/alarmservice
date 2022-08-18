package api

import (
	"net"

	"github.com/caarlos0/log"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/yurttasutkan/alarmservice/internal/config"
	"google.golang.org/grpc"
)

func Setup(conf *config.Config) error {
	apiConf := conf.AlarmServer.API

	log.WithFields(log.Fields{
		"bind":     apiConf.Bind,
	}).Info("api: starting alarm-server api server")
	
	grpcServer := grpc.NewServer()
	alsAPI := NewAlarmServerAPI()
	als.RegisterAlarmServerServiceServer(grpcServer, alsAPI)
	lis, err := net.Listen("tcp", "172.22.0.18:9000")
	if err != nil {
		log.Fatalf("Start api listener error: %v", err)
	}

	grpcServer.Serve(lis)
	
	return nil
}
