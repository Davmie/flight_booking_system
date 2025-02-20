package main

import (
	"flight_booking_system/ticketService/cmd/server"
	ticketDel "flight_booking_system/ticketService/internal/ticket/delivery"
	pgTicket "flight_booking_system/ticketService/internal/ticket/repository/postgres"
	ticketUseCase "flight_booking_system/ticketService/internal/ticket/usecase"
	"flight_booking_system/ticketService/pkg/middleware"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var prodCfgPg = postgres.Config{DSN: "host=postgres user=program password=test dbname=tickets port=5432"}

func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	ticketHandler := ticketDel.TicketHandler{
		TicketUseCase: ticketUseCase.New(pgTicket.New(logger, db)),
		Logger:        logger,
	}

	r := http.NewServeMux()

	r.Handle("GET /api/v1/tickets/{ticketId}", http.HandlerFunc(ticketHandler.Get))
	r.Handle("GET /api/v1/tickets", http.HandlerFunc(ticketHandler.GetAll))
	r.Handle("POST /api/v1/tickets", http.HandlerFunc(ticketHandler.Create))
	r.Handle("PATCH /api/v1/tickets", http.HandlerFunc(ticketHandler.Update))
	r.Handle("DELETE /api/v1/tickets/{ticketId}", http.HandlerFunc(ticketHandler.Delete))
	r.Handle("GET /api/v1/ticketsByUID", http.HandlerFunc(ticketHandler.GetByUID))

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
