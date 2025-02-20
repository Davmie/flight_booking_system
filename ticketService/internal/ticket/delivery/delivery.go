package delivery

import (
	"encoding/json"
	ticketUseCase "flight_booking_system/ticketService/internal/ticket/usecase"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"

	"flight_booking_system/ticketService/models"
	"flight_booking_system/ticketService/pkg/logger"
)

type TicketHandler struct {
	TicketUseCase ticketUseCase.TicketUseCaseI
	Logger        logger.Logger
}

func (th *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	ticket := models.Ticket{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		th.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		th.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &ticket)
	if err != nil {
		th.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(ticket)
	//if err != nil {
	//	th.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	if ticket.TicketUID == "" {
		ticket.TicketUID = uuid.New().String()
	}

	err = th.TicketUseCase.Create(&ticket)
	if err != nil {
		th.Logger.Infow("can`t create ticket",
			"err:", err.Error())
		http.Error(w, "can`t create ticket", http.StatusBadRequest)
		return
	}

	//resp, err := json.Marshal(ticket)
	//
	//if err != nil {
	//	th.Logger.Errorw("can`t marshal ticket",
	//		"err:", err.Error())
	//	http.Error(w, "can`t make ticket", http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/tickets/%d", ticket.ID))
	w.Header().Set("X-Ticket-UID", ticket.TicketUID)
	w.WriteHeader(http.StatusCreated)

	//_, err = w.Write(resp)
	//if err != nil {
	//	th.Logger.Errorw("can`t write response",
	//		"err:", err.Error())
	//	http.Error(w, "can`t write response", http.StatusInternalServerError)
	//	return
	//}
}

func (th *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	ticketIdString := r.PathValue("ticketId")
	if ticketIdString == "" {
		th.Logger.Errorw("no ticketId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	ticketId, err := strconv.Atoi(ticketIdString)
	if err != nil {
		th.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	ticket, err := th.TicketUseCase.Get(ticketId)
	if err != nil {
		th.Logger.Infow("can`t get ticket",
			"err:", err.Error())
		http.Error(w, "can`t get ticket", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(ticket)

	if err != nil {
		th.Logger.Errorw("can`t marshal ticket",
			"err:", err.Error())
		http.Error(w, "can`t make ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		th.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (th *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	ticket := &models.Ticket{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		th.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		th.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, ticket)
	if err != nil {
		th.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(ticket)
	//if err != nil {
	//	th.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	err = th.TicketUseCase.Update(ticket)
	if err != nil {
		th.Logger.Infow("can`t update ticket",
			"err:", err.Error())
		http.Error(w, "can`t update ticket", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(ticket)

	if err != nil {
		th.Logger.Errorw("can`t marshal ticket",
			"err:", err.Error())
		http.Error(w, "can`t make ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		th.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (th *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ticketIdString := r.PathValue("ticketId")
	if ticketIdString == "" {
		th.Logger.Errorw("no ticketId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	ticketId, err := strconv.Atoi(ticketIdString)
	if err != nil {
		th.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = th.TicketUseCase.Delete(ticketId)
	if err != nil {
		th.Logger.Infow("can`t delete ticket",
			"err:", err.Error())
		http.Error(w, "can`t delete ticket", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TicketHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userName := r.Header["X-User-Name"][0]

	tickets, err := th.TicketUseCase.GetAll(userName)
	if err != nil {
		th.Logger.Infow("can`t get all tickets",
			"err:", err.Error())
		http.Error(w, "can`t get all tickets", http.StatusInternalServerError)
		return
	}

	ticketsDTO := make([]*models.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		dto := models.TicketToDTO(*ticket)
		ticketsDTO[i] = dto
	}

	resp, err := json.Marshal(ticketsDTO)
	if err != nil {
		th.Logger.Errorw("can`t marshal ticketDTO",
			"err:", err.Error())
		http.Error(w, "can`t make ticketDTO", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		th.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (th *TicketHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	userName := r.Header["X-User-Name"][0]
	ticketUid := r.Header["X-Ticket-Uid"][0]

	ticket, err := th.TicketUseCase.GetByUID(ticketUid, userName)
	if err != nil {
		th.Logger.Infow("can`t get ticket by uid",
			"err:", err.Error())
		http.Error(w, "can`t get ticket by uid", http.StatusInternalServerError)
		return
	}

	dto := models.TicketToDTO(*ticket)

	resp, err := json.Marshal(dto)
	if err != nil {
		th.Logger.Errorw("can`t marshal ticketDTO",
			"err:", err.Error())
		http.Error(w, "can`t make ticketDTO", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		th.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
