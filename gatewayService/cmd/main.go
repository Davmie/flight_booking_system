package main

import (
	"context"
	"flight_booking_system/gatewayService/cmd/server"
	gatewayDel "flight_booking_system/gatewayService/internal/delivery"
	"flight_booking_system/gatewayService/pkg/middleware"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func ticketUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/tickets/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "Invalid ticket TicketUID", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "ticketUID", parts[0])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//var prodCfgPg = postgres.Config{DSN: "host=dpg-crn8te08fa8c738bekog-a.frankfurt-postgres.render.com user=program password=ZbuzttiIWI6DpKrwhYryoHhYm7NjeLQ9 dbname=persons_1jwg port=5432"}

// postgresql://program:ZbuzttiIWI6DpKrwhYryoHhYm7NjeLQ9@dpg-crn8te08fa8c738bekog-a.frankfurt-postgres.render.com/persons_1jwg
func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	//db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	//if err != nil {
	//	log.Fatal(err)
	//}

	gatewayHandler := gatewayDel.GatewayHandler{
		//ServerUseCase: personUseCase.New(pgPerson.New(logger, db)),
		Logger: logger,
	}

	r := http.NewServeMux()

	r.Handle("GET /api/v1/flights", http.HandlerFunc(gatewayHandler.GetFlights))
	r.Handle("GET /api/v1/me", http.HandlerFunc(gatewayHandler.GetMe))
	r.Handle("GET /api/v1/tickets", http.HandlerFunc(gatewayHandler.GetTickets))
	r.Handle("POST /api/v1/tickets", http.HandlerFunc(gatewayHandler.BuyTicket))
	r.Handle("GET /api/v1/tickets/", ticketUIDMiddleware(http.HandlerFunc(gatewayHandler.GetTicketByUID)))
	r.Handle("DELETE /api/v1/tickets/", ticketUIDMiddleware(http.HandlerFunc(gatewayHandler.ReturnTicket)))
	r.Handle("GET /api/v1/privilege", http.HandlerFunc(gatewayHandler.GetPrivilege))

	r.Handle("GET /manage/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	router := middleware.AccessLog(logger, r)
	router = middleware.Panic(logger, router)

	s := server.NewServer(router)
	if err := s.Start(); err != nil {
		logger.Fatal(err)
	}

	err := zapLogger.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
