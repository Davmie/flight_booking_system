package delivery

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"flight_booking_system/gatewayService/models"
	"flight_booking_system/gatewayService/pkg/logger"
	//"github.com/asaskevich/govalidator"
)

const (
	bonusHost  = "http://bonus_msv:8050"
	flightHost = "http://flight_msv:8060"
	ticketHost = "http://ticket_msv:8070"
)

type GatewayHandler struct {
	Logger logger.Logger
	Client http.Client
}

func makeFlightsInfoResponse(flightResponses []*models.FlightResponse, page, size int) models.FlightsInfo {
	flightsInfo := make([]models.FlightInfo, len(flightResponses))
	for i, response := range flightResponses {
		flightsInfo[i] = models.FlightInfo{
			FlightNumber: response.FlightNumber,
			FromAirport:  response.FromAirport,
			ToAirport:    response.ToAirport,
			Date:         response.Date,
			Price:        response.Price,
		}
	}

	return models.FlightsInfo{
		Page:    page,
		Size:    size,
		Total:   len(flightsInfo),
		Flights: flightsInfo,
	}
}

func makeTicketInfoResponse(ticketResponses []*models.TicketResponse, flightResponses []*models.FlightResponse) []models.TicketInfo {
	res := make([]models.TicketInfo, len(ticketResponses))
	for i, ticketResponse := range ticketResponses {
		res[i] = models.TicketInfo{
			UID:          ticketResponse.TicketUID,
			FlightNumber: ticketResponse.FlightNumber,
			FromAirport:  flightResponses[i].FromAirport,
			ToAirport:    flightResponses[i].ToAirport,
			Date:         flightResponses[i].Date,
			Price:        ticketResponse.Price,
			Status:       ticketResponse.Status,
		}
	}

	return res
}

func makeUserInfoResponse(ticketResponses []models.TicketResponse, flightResponses []*models.FlightResponse, privilegeResponse *models.PrivilegeResponse) models.UserInfoResponse {
	res := models.UserInfoResponse{}

	ticketInfos := make([]models.TicketInfo, len(ticketResponses))
	for i, ticketResponse := range ticketResponses {
		ticketInfos[i] = models.TicketInfo{
			UID:          ticketResponse.TicketUID,
			FlightNumber: ticketResponse.FlightNumber,
			FromAirport:  flightResponses[i].FromAirport,
			ToAirport:    flightResponses[i].ToAirport,
			Date:         flightResponses[i].Date,
			Price:        ticketResponse.Price,
			Status:       ticketResponse.Status,
		}
	}

	privilegeInfo := models.PrivilegeInfo{
		ID:      privilegeResponse.ID,
		Status:  privilegeResponse.Status,
		Balance: privilegeResponse.Balance,
	}

	res.Tickets = ticketInfos
	res.Privilege = privilegeInfo

	return res
}

func makePrivilegeFullResponse(privilegeResponse models.PrivilegeResponse, privilegeHistory []models.PrivilegeHistoryResponse) models.PrivilegeFullResponse {
	return models.PrivilegeFullResponse{
		Balance: privilegeResponse.Balance,
		Status:  privilegeResponse.Status,
		History: privilegeHistory,
	}
}

func (gh *GatewayHandler) makeUpdatePrivilegeRequest(privilegeInfo *models.PrivilegeResponse) error {
	jsonPrivilegeInfo, err := json.Marshal(privilegeInfo)
	if err != nil {
		return err
	}

	privReq, err := http.NewRequest("PATCH", bonusHost+"/api/v1/privileges", bytes.NewBuffer(jsonPrivilegeInfo))
	if err != nil {
		return err
	}

	_, err = gh.Client.Do(privReq)
	if err != nil {
		return err
	}

	return nil
}

func (gh *GatewayHandler) makeCreatePrivilegeHistoryRequest(privilege *models.PrivilegeResponse, operationType string, ticketUID string, balanceDiff int) error {
	privHistInfo := models.PrivilegeHistoryInfo{
		TicketUID:     ticketUID,
		OperationType: operationType,
		PrivilegeID:   privilege.ID,
		BalanceDiff:   balanceDiff,
		Date:          time.Now(),
	}

	jsonHistInfo, err := json.Marshal(privHistInfo)
	if err != nil {
		return err
	}

	privHistReq, err := http.NewRequest("POST", bonusHost+"/api/v1/privileges/history", bytes.NewBuffer(jsonHistInfo))
	if err != nil {
		return err
	}

	privHistReq.Header.Set("Content-Type", "application/json")

	_, err = gh.Client.Do(privHistReq)
	if err != nil {
		return err
	}

	return nil
}

func (gh *GatewayHandler) makeCreateTicketRequest(userName string, flightNumber string, price int) (string, error) {
	ticketInfo := models.TicketInfoRequest{
		UID:          "",
		FlightNumber: flightNumber,
		Username:     userName,
		Price:        price,
		Status:       "PAID",
	}

	jsonTicketInfo, err := json.Marshal(ticketInfo)
	if err != nil {
		return "", err
	}

	ticketReq, err := http.NewRequest("POST", ticketHost+"/api/v1/tickets", bytes.NewBuffer(jsonTicketInfo))
	if err != nil {
		return "", err
	}

	ticketReq.Header.Set("Content-Type", "application/json")

	ticketResp, err := gh.Client.Do(ticketReq)
	if err != nil {
		return "", err
	}

	ticketUID := ticketResp.Header.Get("X-Ticket-UID")

	return ticketUID, nil
}

func (gh *GatewayHandler) makeGetFlightRequest(flightNumber string) (*models.FlightResponse, error) {
	flightReq, err := http.NewRequest("GET", flightHost+"/api/v1/flights", nil)
	if err != nil {
		return nil, err
	}

	flightReq.Header.Set("flightNumber", flightNumber)

	flightResp, err := gh.Client.Do(flightReq)
	if err != nil {
		return nil, err
	}

	flightResponse := make([]*models.FlightResponse, 0)
	body, err := io.ReadAll(flightResp.Body)
	if err != nil {
		return nil, err
	}

	err = flightResp.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &flightResponse)
	if err != nil {
		return nil, err
	}

	return flightResponse[0], nil
}

func (gh *GatewayHandler) GetFlights(w http.ResponseWriter, r *http.Request) {
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

	flightReq, err := http.NewRequest("GET", flightHost+"/api/v1/flightsPaginate", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create flight request")
		http.Error(w, "can`t create flight request", http.StatusInternalServerError)
		return
	}
	newQ := flightReq.URL.Query()
	newQ.Add("page", strPage)
	newQ.Add("size", strSize)
	flightReq.URL.RawQuery = newQ.Encode()

	flightResp, err := gh.Client.Do(flightReq)
	if err != nil {
		gh.Logger.Errorw("can`t get flights", "err:", err.Error())
		http.Error(w, "can`t get flights", http.StatusBadRequest)
		return
	}

	flightResponses := make([]*models.FlightResponse, 0)
	body, err := io.ReadAll(flightResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = flightResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &flightResponses)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	res := makeFlightsInfoResponse(flightResponses, page, size)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := json.Marshal(res)
	if err != nil {
		gh.Logger.Errorw("can`t marshal response", "err:", err.Error())
		http.Error(w, "can`t marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		gh.Logger.Errorw("can`t write response", "err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userName := r.Header["X-User-Name"]
	if len(userName) == 0 {
		gh.Logger.Errorw("no X-User-Name header found")
		http.Error(w, "unknown error", http.StatusBadRequest)
		return
	}

	ticketReq, err := http.NewRequest("GET", ticketHost+"/api/v1/tickets", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create ticket request")
		http.Error(w, "can`t create ticket request", http.StatusInternalServerError)
		return
	}
	ticketReq.Header["X-User-Name"] = userName

	ticketResp, err := gh.Client.Do(ticketReq)
	if err != nil {
		gh.Logger.Errorw("can`t get tickets", "err:", err.Error())
		http.Error(w, "can`t get tickets", http.StatusBadRequest)
		return
	}

	ticketResponses := make([]models.TicketResponse, 0)
	body, err := io.ReadAll(ticketResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = ticketResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &ticketResponses)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	flightResponses := make([]*models.FlightResponse, len(ticketResponses))

	for i, ticketResponse := range ticketResponses {
		flightResponse, err := gh.makeGetFlightRequest(ticketResponse.FlightNumber)
		if err != nil {
			gh.Logger.Errorw("can`t create flight request", "err:", err.Error())
			http.Error(w, "can`t create flight request", http.StatusBadRequest)
			return
		}

		flightResponses[i] = flightResponse
	}

	privilegeReq, err := http.NewRequest("GET", bonusHost+"/api/v1/privileges", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create privilege request", "err:", err.Error())
		http.Error(w, "can`t create privilege request", http.StatusInternalServerError)
		return
	}
	privilegeReq.Header["X-User-Name"] = userName

	privilegeResp, err := gh.Client.Do(privilegeReq)
	if err != nil {
		gh.Logger.Errorw("can`t get privilege", "err:", err.Error())
		http.Error(w, "can`t get privilege", http.StatusInternalServerError)
		return
	}

	privilegeResponse := &models.PrivilegeResponse{}
	body, err = io.ReadAll(privilegeResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = privilegeResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, privilegeResponse)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	res := makeUserInfoResponse(ticketResponses, flightResponses, privilegeResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := json.Marshal(res)
	if err != nil {
		gh.Logger.Errorw("can`t marshal response", "err:", err.Error())
		http.Error(w, "can`t marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		gh.Logger.Errorw("can`t write response", "err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	userName := r.Header["X-User-Name"]
	if len(userName) == 0 {
		gh.Logger.Errorw("no X-User-Name header found")
		http.Error(w, "unknown error", http.StatusBadRequest)
		return
	}

	ticketReq, err := http.NewRequest("GET", ticketHost+"/api/v1/tickets", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create ticket request")
		http.Error(w, "can`t create ticket request", http.StatusInternalServerError)
		return
	}
	ticketReq.Header["X-User-Name"] = userName

	ticketResp, err := gh.Client.Do(ticketReq)
	if err != nil {
		gh.Logger.Errorw("can`t get tickets", "err:", err.Error())
		http.Error(w, "can`t get tickets", http.StatusBadRequest)
		return
	}

	ticketResponses := make([]*models.TicketResponse, 0)
	body, err := io.ReadAll(ticketResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = ticketResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &ticketResponses)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	flightResponses := make([]*models.FlightResponse, len(ticketResponses))

	for i, ticketResponse := range ticketResponses {
		flightResponse, err := gh.makeGetFlightRequest(ticketResponse.FlightNumber)
		if err != nil {
			gh.Logger.Errorw("can`t create flight request", "err:", err.Error())
			http.Error(w, "can`t create flight request", http.StatusBadRequest)
			return
		}

		flightResponses[i] = flightResponse
	}

	ticketInfoResponse := makeTicketInfoResponse(ticketResponses, flightResponses)

	resp, err := json.Marshal(ticketInfoResponse)

	if err != nil {
		gh.Logger.Errorw("can`t marshal ticketInfoResponse",
			"err:", err.Error())
		http.Error(w, "can`t make ticketInfoResponse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		gh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
	buyInfo := models.BuyTicketInfo{}
	buyInfoResponse := models.BuyTicketResponse{}

	userName := r.Header["X-User-Name"]
	if len(userName) == 0 {
		gh.Logger.Errorw("no X-User-Name header found")
		http.Error(w, "unknown error", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &buyInfo)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	// Проверка, что рейс существует
	flightResponse, err := gh.makeGetFlightRequest(buyInfo.FlightNumber)
	if err != nil {
		gh.Logger.Errorw("can`t create flight request", "err:", err.Error())
		http.Error(w, "can`t create flight request", http.StatusBadRequest)
		return
	}

	// Инфа о бонусах
	privilegeReq, err := http.NewRequest("GET", bonusHost+"/api/v1/privileges", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create privilege request", "err:", err.Error())
		http.Error(w, "can`t create privilege request", http.StatusInternalServerError)
		return
	}
	privilegeReq.Header["X-User-Name"] = userName

	privilegeResp, err := gh.Client.Do(privilegeReq)
	if err != nil {
		gh.Logger.Errorw("can`t get privilege", "err:", err.Error())
		http.Error(w, "can`t get privilege", http.StatusInternalServerError)
		return
	}

	privilegeResponse := &models.PrivilegeResponse{}
	body, err = io.ReadAll(privilegeResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = privilegeResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, privilegeResponse)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	balanceDiff := 0
	operationType := "FILL_IN_BALANCE"

	if buyInfo.PaidFromBalance {
		operationType = "DEBIT_THE_ACCOUNT"
		if privilegeResponse.Balance >= buyInfo.Price {
			buyInfoResponse.PaidByBonuses = buyInfo.Price
			buyInfoResponse.PaidByMoney = 0
			balanceDiff = buyInfo.Price
		} else {
			buyInfoResponse.PaidByBonuses = privilegeResponse.Balance
			buyInfoResponse.PaidByMoney = buyInfo.Price - privilegeResponse.Balance
			balanceDiff = privilegeResponse.Balance
		}
	} else {
		buyInfoResponse.PaidByBonuses = 0
		buyInfoResponse.PaidByMoney = buyInfo.Price
		balanceDiff = int(math.Round(0.1 * float64(buyInfo.Price)))
		privilegeResponse.Balance += balanceDiff
	}

	privilegeResponse.Balance -= buyInfoResponse.PaidByBonuses

	// Создание билета
	ticketUID, err := gh.makeCreateTicketRequest(userName[0], buyInfo.FlightNumber, buyInfo.Price)
	if err != nil {
		gh.Logger.Errorw("can`t create ticket", "err:", err.Error())
		http.Error(w, "can`t create ticket", http.StatusBadRequest)
		return
	}

	err = gh.makeCreatePrivilegeHistoryRequest(privilegeResponse, operationType, ticketUID, balanceDiff)

	err = gh.makeUpdatePrivilegeRequest(privilegeResponse)

	buyInfoResponse.UID = ticketUID
	buyInfoResponse.FlightNumber = buyInfo.FlightNumber
	buyInfoResponse.FromAirport = flightResponse.FromAirport
	buyInfoResponse.ToAirport = flightResponse.ToAirport
	buyInfoResponse.Date = flightResponse.Date
	buyInfoResponse.Price = buyInfo.Price
	buyInfoResponse.Status = "PAID"
	buyInfoResponse.Privilege = *privilegeResponse

	resp, err := json.Marshal(buyInfoResponse)
	if err != nil {
		gh.Logger.Errorw("can`t marshal buyInfoResponse",
			"err:", err.Error())
		http.Error(w, "can`t marshal buyInfoResponse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		gh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandler) GetTicketByUID(w http.ResponseWriter, r *http.Request) {
	userName := r.Header["X-User-Name"]
	if len(userName) == 0 {
		gh.Logger.Errorw("no X-User-Name header found")
		http.Error(w, "unknown error", http.StatusBadRequest)
		return
	}

	ticketUid := r.Context().Value("ticketUID").(string)

	ticketReq, err := http.NewRequest("GET", ticketHost+"/api/v1/ticketsByUID", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create ticket request")
		http.Error(w, "can`t create ticket request", http.StatusInternalServerError)
		return
	}
	ticketReq.Header["X-User-Name"] = userName
	ticketReq.Header["X-Ticket-Uid"] = []string{ticketUid}

	ticketResp, err := gh.Client.Do(ticketReq)
	if err != nil {
		gh.Logger.Errorw("can`t get tickets", "err:", err.Error())
		http.Error(w, "can`t get tickets", http.StatusBadRequest)
		return
	}

	ticketResponse := &models.TicketResponse{}
	body, err := io.ReadAll(ticketResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = ticketResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &ticketResponse)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	flightResponse := &models.FlightResponse{}

	flightResponse, err = gh.makeGetFlightRequest(ticketResponse.FlightNumber)
	if err != nil {
		gh.Logger.Errorw("can`t create flight request", "err:", err.Error())
		http.Error(w, "can`t create flight request", http.StatusBadRequest)
		return
	}

	ticketInfoResponse := makeTicketInfoResponse([]*models.TicketResponse{ticketResponse}, []*models.FlightResponse{flightResponse})

	resp, err := json.Marshal(ticketInfoResponse[0])

	if err != nil {
		gh.Logger.Errorw("can`t marshal ticketInfoResponse",
			"err:", err.Error())
		http.Error(w, "can`t make ticketInfoResponse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		gh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandler) ReturnTicket(w http.ResponseWriter, r *http.Request) {
	userName := r.Header["X-User-Name"]
	if len(userName) == 0 {
		gh.Logger.Errorw("no X-User-Name header found")
		http.Error(w, "unknown error", http.StatusBadRequest)
		return
	}

	ticketUid := r.Context().Value("ticketUID").(string)

	ticketReq, err := http.NewRequest("GET", ticketHost+"/api/v1/ticketsByUID", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create ticket request")
		http.Error(w, "can`t create ticket request", http.StatusInternalServerError)
		return
	}
	ticketReq.Header["X-User-Name"] = userName
	ticketReq.Header["X-Ticket-Uid"] = []string{ticketUid}

	ticketResp, err := gh.Client.Do(ticketReq)
	if err != nil {
		gh.Logger.Errorw("can`t get tickets", "err:", err.Error())
		http.Error(w, "can`t get tickets", http.StatusBadRequest)
		return
	}

	ticketResponse := models.TicketResponse{}
	body, err := io.ReadAll(ticketResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = ticketResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &ticketResponse)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	ticketResponse.Status = "CANCELED"

	//TODO: Сделать запрос в PrivilegeHistory, чтобы отменить эту операцию

	jsonTicketResponse, err := json.Marshal(ticketResponse)
	if err != nil {
		gh.Logger.Errorw("can`t marshal ticketResponse",
			"err:", err.Error())
		http.Error(w, "can`t make ticketResponse", http.StatusInternalServerError)
		return
	}

	ticketReq, err = http.NewRequest("PATCH", ticketHost+"/api/v1/tickets", bytes.NewBuffer(jsonTicketResponse))
	if err != nil {
		gh.Logger.Errorw("can`t create ticket request")
		http.Error(w, "can`t create ticket request", http.StatusInternalServerError)
		return
	}

	_, err = gh.Client.Do(ticketReq)
	if err != nil {
		gh.Logger.Errorw("can`t cancel ticket", "err:", err.Error())
		http.Error(w, "can`t cancel ticket", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (gh *GatewayHandler) GetPrivilege(w http.ResponseWriter, r *http.Request) {
	userName := r.Header["X-User-Name"]
	if len(userName) == 0 {
		gh.Logger.Errorw("no X-User-Name header found")
		http.Error(w, "unknown error", http.StatusBadRequest)
		return
	}

	privilegeReq, err := http.NewRequest("GET", bonusHost+"/api/v1/privileges", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create privilege request", "err:", err.Error())
		http.Error(w, "can`t create privilege request", http.StatusInternalServerError)
		return
	}
	privilegeReq.Header["X-User-Name"] = userName

	privilegeResp, err := gh.Client.Do(privilegeReq)
	if err != nil {
		gh.Logger.Errorw("can`t get privilege", "err:", err.Error())
		http.Error(w, "can`t get privilege", http.StatusInternalServerError)
		return
	}

	privilegeResponse := &models.PrivilegeResponse{}
	body, err := io.ReadAll(privilegeResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = privilegeResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, privilegeResponse)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	privilegeHistoryReq, err := http.NewRequest("GET", bonusHost+"/api/v1/privilegeHistory", nil)
	if err != nil {
		gh.Logger.Errorw("can`t create privilege request", "err:", err.Error())
		http.Error(w, "can`t create privilege request", http.StatusInternalServerError)
		return
	}
	privilegeHistoryReq.Header["X-User-Name"] = userName

	privilegeHistoryResp, err := gh.Client.Do(privilegeHistoryReq)
	if err != nil {
		gh.Logger.Errorw("can`t get privilege", "err:", err.Error())
		http.Error(w, "can`t get privilege", http.StatusInternalServerError)
		return
	}

	privilegeHistoryResponse := make([]models.PrivilegeHistoryResponse, 0)
	body, err = io.ReadAll(privilegeHistoryResp.Body)
	if err != nil {
		gh.Logger.Errorw("can`t read body of request", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = privilegeHistoryResp.Body.Close()
	if err != nil {
		gh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &privilegeHistoryResponse)
	if err != nil {
		gh.Logger.Infow("can`t unmarshal form", "err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	privilegeFullResponse := makePrivilegeFullResponse(*privilegeResponse, privilegeHistoryResponse)

	resp, err := json.Marshal(privilegeFullResponse)

	if err != nil {
		gh.Logger.Errorw("can`t marshal privilegeFullResponse",
			"err:", err.Error())
		http.Error(w, "can`t make privilegeFullResponse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		gh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
