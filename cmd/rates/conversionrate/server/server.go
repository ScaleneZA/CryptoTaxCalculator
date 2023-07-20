package server

import (
	"context"
	pb "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/conversionratepb"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/marketvalue"
	"github.com/luno/jettison/interceptors"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	pb.UnimplementedConversionrateServer

	b Backends
}

func Serve(b Backends) error {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryServerInterceptor),
		grpc.StreamInterceptor(interceptors.StreamServerInterceptor),
	)

	pb.RegisterConversionrateServer(s, &Server{
		b: b,
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s Server) ValueAtTime(ctx context.Context, req *pb.ValueAtTimeRequest) (*pb.ValueAtTimeResponse, error) {
	val, err := marketvalue.ValueAtTime(s.b, req.From, req.To, req.Timestamp)
	if err != nil {
		return nil, err
	}

	return &pb.ValueAtTimeResponse{
		Rate: val,
	}, nil
}
