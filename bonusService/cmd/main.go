package main

import (
	"flight_booking_system/bonusService/cmd/server"
	privilegeDel "flight_booking_system/bonusService/internal/privilege/delivery"
	pgPrivilege "flight_booking_system/bonusService/internal/privilege/repository/postgres"
	privilegeUseCase "flight_booking_system/bonusService/internal/privilege/usecase"
	"flight_booking_system/bonusService/pkg/middleware"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var prodCfgPg = postgres.Config{DSN: "host=postgres user=program password=test dbname=privileges port=5432"}

func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	privilegeHandler := privilegeDel.PrivilegeHandler{
		PrivilegeUseCase: privilegeUseCase.New(pgPrivilege.New(logger, db)),
		Logger:           logger,
	}

	r := http.NewServeMux()

	r.Handle("GET /api/v1/privilegeHistory", http.HandlerFunc(privilegeHandler.GetHistory))
	r.Handle("GET /api/v1/privileges/{privilegeId}", http.HandlerFunc(privilegeHandler.Get))
	r.Handle("GET /api/v1/privileges", http.HandlerFunc(privilegeHandler.GetAll))
	r.Handle("POST /api/v1/privileges", http.HandlerFunc(privilegeHandler.Create))
	r.Handle("PATCH /api/v1/privileges", http.HandlerFunc(privilegeHandler.Update))
	r.Handle("DELETE /api/v1/privileges/{privilegeId}", http.HandlerFunc(privilegeHandler.Delete))
	r.Handle("POST /api/v1/privileges/history", http.HandlerFunc(privilegeHandler.CreateHistory))

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
