package main

import (
	"fmt"
	pb "github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/conversionratepb"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/server"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/di"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/webhandlers"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	b := di.SetupDI()

	go syncCurrenciesForever(b)

	grpcServers(b)
	httpServers()
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

func grpcServers(b di.Backends) {
	// TODO: Abstract this a bit to allow for other servers
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

func httpServers() {
	http.HandleFunc("/", webhandlers.Home)
	http.HandleFunc("/ajax/upload", webhandlers.UploadTransform)
	http.HandleFunc("/ajax/calculate", webhandlers.Calculate)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
