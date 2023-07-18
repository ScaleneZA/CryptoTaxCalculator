package server

import (
	"context"
	pb "github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/conversionratepb"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/marketvalue"
)

type Server struct {
	pb.UnimplementedConversionrateServer

	B Backends
}

func (s Server) ValueAtTime(ctx context.Context, req *pb.ValueAtTimeRequest) (*pb.ValueAtTimeResponse, error) {
	val, err := marketvalue.ValueAtTime(s.B, req.From, req.To, req.Timestamp)
	if err != nil {
		return nil, err
	}

	return &pb.ValueAtTimeResponse{
		Rate: val,
	}, nil
}
