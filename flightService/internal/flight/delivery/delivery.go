package delivery

import (
	"encoding/json"
	flightUseCase "flight_booking_system/flightService/internal/flight/usecase"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"flight_booking_system/flightService/models"
	"flight_booking_system/flightService/pkg/logger"
	//"github.com/asaskevich/govalidator"
)

type FlightHandler struct {
	FlightUseCase flightUseCase.FlightUseCaseI
	Logger        logger.Logger
}

func (ah *FlightHandler) Create(w http.ResponseWriter, r *http.Request) {
	flight := models.Flight{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &flight)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(flight)
	//if err != nil {
	//	ah.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	err = ah.FlightUseCase.Create(&flight)
	if err != nil {
		ah.Logger.Infow("can`t create flight",
			"err:", err.Error())
		http.Error(w, "can`t create flight", http.StatusBadRequest)
		return
	}

	//resp, err := json.Marshal(flight)
	//
	//if err != nil {
	//	ah.Logger.Errorw("can`t marshal flight",
	//		"err:", err.Error())
	//	http.Error(w, "can`t make flight", http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/flights/%d", flight.ID))
	w.WriteHeader(http.StatusCreated)

	//_, err = w.Write(resp)
	//if err != nil {
	//	ah.Logger.Errorw("can`t write response",
	//		"err:", err.Error())
	//	http.Error(w, "can`t write response", http.StatusInternalServerError)
	//	return
	//}
}

func (ah *FlightHandler) Get(w http.ResponseWriter, r *http.Request) {
	flightIdString := r.PathValue("flightId")
	if flightIdString == "" {
		ah.Logger.Errorw("no flightId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	flightId, err := strconv.Atoi(flightIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	flight, err := ah.FlightUseCase.Get(flightId)
	if err != nil {
		ah.Logger.Infow("can`t get flight",
			"err:", err.Error())
		http.Error(w, "can`t get flight", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(flight)

	if err != nil {
		ah.Logger.Errorw("can`t marshal flight",
			"err:", err.Error())
		http.Error(w, "can`t make flight", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *FlightHandler) Update(w http.ResponseWriter, r *http.Request) {
	flightIdString := r.PathValue("flightId")
	if flightIdString == "" {
		ah.Logger.Errorw("no flightId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	flightId, err := strconv.Atoi(flightIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	flight := &models.Flight{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, flight)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(flight)
	//if err != nil {
	//	ah.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	flight.ID = flightId
	err = ah.FlightUseCase.Update(flight)
	if err != nil {
		ah.Logger.Infow("can`t update flight",
			"err:", err.Error())
		http.Error(w, "can`t update flight", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(flight)

	if err != nil {
		ah.Logger.Errorw("can`t marshal flight",
			"err:", err.Error())
		http.Error(w, "can`t make flight", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *FlightHandler) Delete(w http.ResponseWriter, r *http.Request) {
	flightIdString := r.PathValue("flightId")
	if flightIdString == "" {
		ah.Logger.Errorw("no flightId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	flightId, err := strconv.Atoi(flightIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = ah.FlightUseCase.Delete(flightId)
	if err != nil {
		ah.Logger.Infow("can`t delete flight",
			"err:", err.Error())
		http.Error(w, "can`t delete flight", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ah *FlightHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var flights []*models.FlightDTO
	var err error
	var resp []byte

	flightNumber := r.Header.Get("flightNumber")
	flights, err = ah.FlightUseCase.GetAll(flightNumber)
	if err != nil {
		ah.Logger.Infow("can`t get all flights",
			"err:", err.Error())
		http.Error(w, "can`t get all flights", http.StatusInternalServerError)
		return
	}

	resp, err = json.Marshal(flights)
	if err != nil {
		ah.Logger.Errorw("can`t marshal flight",
			"err:", err.Error())
		http.Error(w, "can`t make flight", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *FlightHandler) GetAllPaginate(w http.ResponseWriter, r *http.Request) {
	var flights []*models.FlightDTO
	var err error
	var resp []byte

	q := r.URL.Query()
	strPage := q.Get("page")
	strSize := q.Get("size")
	page, _ := strconv.Atoi(strPage)
	if page <= 0 {
		page = 1
	}

	size, _ := strconv.Atoi(strSize)
	switch {
	case size <= 0:
		size = 1
	case size > 100:
		size = 100
	}

	offset := (page - 1) * size

	flights, err = ah.FlightUseCase.GetAllPaginate(offset, size)
	if err != nil {
		ah.Logger.Infow("can`t get all flights",
			"err:", err.Error())
		http.Error(w, "can`t get all flights", http.StatusInternalServerError)
		return
	}

	resp, err = json.Marshal(flights)
	if err != nil {
		ah.Logger.Errorw("can`t marshal flight",
			"err:", err.Error())
		http.Error(w, "can`t make flight", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
