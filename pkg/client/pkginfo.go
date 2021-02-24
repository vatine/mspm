package client

import (
	"context"

	log "github.com/sirupsen/logrus"

	pb "github.com/vatine/mspm/pkg/protos"
)

func (c *Client) GetPackageInformation(ctx context.Context, name string) ([]*pb.PackageInformation, error) {
	req := pb.PackageInformationRequest{
		PackageName: name,
	}

	resp, err := c.client.GetPackageInformation(ctx, &req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name":  name,
		}).Error("GetPackageInformation")
		return nil, err
	}

	return resp.PackageData, nil
}
