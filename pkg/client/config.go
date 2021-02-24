package client

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/vatine/mspm/pkg/protos"
)

type Client struct {
	client  pb.MspmClient
	conn    *grpc.ClientConn
	mspmDir string
}

// Create a new client Config, with a hooked-up gRPC client.
func New(backend, directory string) (*Client, error) {
	log.WithFields(log.Fields{
		"backend":   backend,
		"directory": directory,
	}).Debug("creating new client")
	var rv Client
	var err error

	rv.conn, err = grpc.Dial(backend, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("dialling server")
		return nil, err
	}
	rv.client = pb.NewMspmClient(rv.conn)
	rv.mspmDir = directory

	return &rv, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
