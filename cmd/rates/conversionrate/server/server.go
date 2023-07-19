package server

import (
	"context"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/conversionratepb"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/marketvalue"
)

type Server struct {
	conversionratepb.UnimplementedConversionrateServer

	B Backends
}

func (s Server) ValueAtTime(ctx context.Context, req *conversionratepb.ValueAtTimeRequest) (*conversionratepb.ValueAtTimeResponse, error) {
	val, err := marketvalue.ValueAtTime(s.B, req.From, req.To, req.Timestamp)
	if err != nil {
		return nil, err
	}

	return &conversionratepb.ValueAtTimeResponse{
		Rate: val,
	}, nil
}
