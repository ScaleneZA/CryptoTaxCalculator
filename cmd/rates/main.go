package main

import (
	pb "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/conversionratepb"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/server"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/db"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/di"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	b := di.SetupDI()

	//go syncCurrenciesForever(b)
	grpcServer(b)
}

func syncCurrenciesForever(b di.Backends) {
	for {
		err := sync.SyncAll(b)
		if err != nil {
			log.Println("Pair Syncing Failed:")
			log.Println(err.Error())
		}

		log.Println("Synced all currencies")

		time.Sleep(time.Hour * 12)
	}
}

func grpcServer(b di.Backends) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterConversionrateServer(s, &server.Server{
		B: b,
	})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func resetDB() {
	dbc := db.Connect()
	defer dbc.Close()

	db.ResetDB(dbc)
}
