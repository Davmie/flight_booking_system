package main

import (
	"flight_booking_system/flightService/cmd/server"
	flightDel "flight_booking_system/flightService/internal/flight/delivery"
	pgFlight "flight_booking_system/flightService/internal/flight/repository/postgres"
	flightUseCase "flight_booking_system/flightService/internal/flight/usecase"
	"flight_booking_system/flightService/pkg/middleware"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var prodCfgPg = postgres.Config{DSN: "host=postgres user=program password=test dbname=flights port=5432"}

func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	flightHandler := flightDel.FlightHandler{
		FlightUseCase: flightUseCase.New(pgFlight.New(logger, db)),
		Logger:        logger,
	}

	r := http.NewServeMux()

	r.Handle("GET /api/v1/flights/{flightId}", http.HandlerFunc(flightHandler.Get))
	r.Handle("GET /api/v1/flights", http.HandlerFunc(flightHandler.GetAll))
	r.Handle("POST /api/v1/flights", http.HandlerFunc(flightHandler.Create))
	r.Handle("PATCH /api/v1/flights/{flightId}", http.HandlerFunc(flightHandler.Update))
	r.Handle("DELETE /api/v1/flights/{flightId}", http.HandlerFunc(flightHandler.Delete))
	r.Handle("GET /api/v1/flightsPaginate", http.HandlerFunc(flightHandler.GetAllPaginate))

	r.Handle("GET /manage/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	router := middleware.AccessLog(logger, r)
	router = middleware.Panic(logger, router)

	s := server.NewServer(router)
	if err := s.Start(); err != nil {
		logger.Fatal(err)
	}

	err = zapLogger.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
