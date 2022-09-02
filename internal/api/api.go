package api

import (
	"net"

	"github.com/caarlos0/log"
	"github.com/ibrahimozekici/chirpstack-api/go/v5/als"
	"github.com/yurttasutkan/alarmservice/internal/api/alarm"
	"github.com/yurttasutkan/alarmservice/internal/config"
	"google.golang.org/grpc"
)


//Sets up the AlarmServer.
func Setup(conf *config.Config) error {

	//apiConf defines the socket which AlarmServer will be listening to.
	apiConf := conf.AlarmServer.API

	log.WithFields(log.Fields{
		"bind": apiConf.Bind,
	}).Info("api: starting alarm-server api server")

	//Initialize the gRPC server.
	grpcServer := grpc.NewServer()
	alsAPI := alarm.NewAlarmServerAPI()
	als.RegisterAlarmServerServiceServer(grpcServer, alsAPI)

	//Listen on the given socket.
	lis, err := net.Listen("tcp", apiConf.Bind)
	if err != nil {
		log.Fatalf("Start api listener error: %v", err)
	}

	//Starts the connection.
	grpcServer.Serve(lis)

	return nil
}
