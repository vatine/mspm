package main

// The MSPM main entry-point

import (
	"flag"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/vatine/mspm/pkg/protos"
	"github.com/vatine/mspm/pkg/server"
)

func main() {
	var debug bool
	var ssl bool
	var playground, store string
	var port string

	flag.BoolVar(&debug, "debug", false, "Enable debug logging.")
	flag.BoolVar(&ssl, "ssl", false, "Serve requests on an SSL port.")
	flag.StringVar(&playground, "playground", "/var/mspm/tempstore", "Path to temporary storage.")
	flag.StringVar(&store, "store", "/var/mspm/store", "Path to more permanent storage.")
	flag.StringVar(&port, "listen", ":10240", "Host:Port for the gRPC communication.")

	flag.Parse()
	log.SetLevel(log.InfoLevel)
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("Debug logging enabled.")

	log.Debug("Creating gRPC server")
	s := grpc.NewServer()
	log.WithFields(log.Fields{
		"store":      store,
		"playground": playground,
	}).Debug("Creating MSPM server")
	mspmServer := server.NewServer(playground, store)

	log.Debug("Registering MSPM server")
	pb.RegisterMspmServer(s, mspmServer)
	log.Debug("Registering gRPC service reflection")
	reflection.Register(s)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"port":  port,
		}).Fatal("opening gRPC port")
	}

	err = s.Serve(listener)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"port":  port,
		}).Fatal("starting gRPC server")
	}
}
