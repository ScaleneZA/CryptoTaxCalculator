package grpc

import (
	"context"
	pb "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/conversionratepb"
	"github.com/luno/jettison/interceptors"
	"google.golang.org/grpc"
)

type Client struct {
	grpcConnection *grpc.ClientConn
	grpcClient     pb.ConversionrateClient
}

func New() (*Client, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		grpcConnection: conn,
		grpcClient:     pb.NewConversionrateClient(conn),
	}, nil
}

func (c Client) ValueAtTime(from, to string, timestamp int64) (float64, error) {
	resp, err := c.grpcClient.ValueAtTime(context.TODO(), &pb.ValueAtTimeRequest{
		From:      from,
		To:        to,
		Timestamp: timestamp,
	})
	if err != nil {
		return 0, err
	}

	return resp.Rate, nil
}
