package delivery

import (
	"encoding/json"
	privilegeUseCase "flight_booking_system/bonusService/internal/privilege/usecase"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"flight_booking_system/bonusService/models"
	"flight_booking_system/bonusService/pkg/logger"
	//"github.com/asaskevich/govalidator"
)

type PrivilegeHandler struct {
	PrivilegeUseCase privilegeUseCase.PrivilegeUseCaseI
	Logger           logger.Logger
}

func (ph *PrivilegeHandler) Create(w http.ResponseWriter, r *http.Request) {
	privilege := models.Privilege{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ph.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ph.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &privilege)
	if err != nil {
		ph.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(privilege)
	//if err != nil {
	//	ph.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	err = ph.PrivilegeUseCase.Create(&privilege)
	if err != nil {
		ph.Logger.Infow("can`t create privilege",
			"err:", err.Error())
		http.Error(w, "can`t create privilege", http.StatusBadRequest)
		return
	}

	//resp, err := json.Marshal(privilege)
	//
	//if err != nil {
	//	ph.Logger.Errorw("can`t marshal privilege",
	//		"err:", err.Error())
	//	http.Error(w, "can`t make privilege", http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/privileges/%d", privilege.ID))
	w.WriteHeader(http.StatusCreated)

	//_, err = w.Write(resp)
	//if err != nil {
	//	ph.Logger.Errorw("can`t write response",
	//		"err:", err.Error())
	//	http.Error(w, "can`t write response", http.StatusInternalServerError)
	//	return
	//}
}

func (ph *PrivilegeHandler) Get(w http.ResponseWriter, r *http.Request) {
	privilegeIdString := r.PathValue("privilegeId")
	if privilegeIdString == "" {
		ph.Logger.Errorw("no privilegeId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	privilegeId, err := strconv.Atoi(privilegeIdString)
	if err != nil {
		ph.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	privilege, err := ph.PrivilegeUseCase.Get(privilegeId)
	if err != nil {
		ph.Logger.Infow("can`t get privilege",
			"err:", err.Error())
		http.Error(w, "can`t get privilege", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(privilege)

	if err != nil {
		ph.Logger.Errorw("can`t marshal privilege",
			"err:", err.Error())
		http.Error(w, "can`t make privilege", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ph.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ph *PrivilegeHandler) Update(w http.ResponseWriter, r *http.Request) {
	privilege := &models.Privilege{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ph.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ph.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, privilege)
	if err != nil {
		ph.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(privilege)
	//if err != nil {
	//	ph.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}
	err = ph.PrivilegeUseCase.Update(privilege)
	if err != nil {
		ph.Logger.Infow("can`t update privilege",
			"err:", err.Error())
		http.Error(w, "can`t update privilege", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(privilege)

	if err != nil {
		ph.Logger.Errorw("can`t marshal privilege",
			"err:", err.Error())
		http.Error(w, "can`t make privilege", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ph.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ph *PrivilegeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	privilegeIdString := r.PathValue("privilegeId")
	if privilegeIdString == "" {
		ph.Logger.Errorw("no privilegeId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	privilegeId, err := strconv.Atoi(privilegeIdString)
	if err != nil {
		ph.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = ph.PrivilegeUseCase.Delete(privilegeId)
	if err != nil {
		ph.Logger.Infow("can`t delete privilege",
			"err:", err.Error())
		http.Error(w, "can`t delete privilege", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ph *PrivilegeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	userName := r.Header["X-User-Name"][0]
	if userName != "" {
		var privilege *models.Privilege

		privilege, err = ph.PrivilegeUseCase.GetByUserName(userName)

		if err != nil {
			ph.Logger.Infow("can`t get privilege by username",
				"err:", err.Error())
			http.Error(w, "can`t get privilege by username", http.StatusInternalServerError)
			return
		}

		dto := models.PrivilegeToDTO(*privilege)

		resp, err = json.Marshal(dto)
		if err != nil {
			ph.Logger.Errorw("can`t marshal privilege",
				"err:", err.Error())
			http.Error(w, "can`t make privilege", http.StatusInternalServerError)
			return
		}
	} else {
		var privileges []*models.Privilege

		privileges, err = ph.PrivilegeUseCase.GetAll()
		if err != nil {
			ph.Logger.Infow("can`t get all privileges",
				"err:", err.Error())
			http.Error(w, "can`t get all privileges", http.StatusInternalServerError)
			return
		}

		resp, err = json.Marshal(privileges)
		if err != nil {
			ph.Logger.Errorw("can`t marshal privilege",
				"err:", err.Error())
			http.Error(w, "can`t make privilege", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ph.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ph *PrivilegeHandler) CreateHistory(w http.ResponseWriter, r *http.Request) {
	privilegeHistory := models.PrivilegeHistory{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ph.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ph.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &privilegeHistory)
	if err != nil {
		ph.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(privilegeHistory)
	//if err != nil {
	//	ph.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	err = ph.PrivilegeUseCase.CreateHistory(&privilegeHistory)
	if err != nil {
		ph.Logger.Infow("can`t create privilegeHistory",
			"err:", err.Error())
		http.Error(w, "can`t create privilegeHistory", http.StatusBadRequest)
		return
	}

	//resp, err := json.Marshal(privilegeHistory)
	//
	//if err != nil {
	//	ph.Logger.Errorw("can`t marshal privilegeHistory",
	//		"err:", err.Error())
	//	http.Error(w, "can`t make privilegeHistory", http.StatusInternalServerError)
	//	return
	//}

	//w.Header().Set("Location", fmt.Sprintf("/api/v1/privileges/%d", privilegeHistory.ID))
	w.WriteHeader(http.StatusCreated)

	//_, err = w.Write(resp)
	//if err != nil {
	//	ph.Logger.Errorw("can`t write response",
	//		"err:", err.Error())
	//	http.Error(w, "can`t write response", http.StatusInternalServerError)
	//	return
	//}
}

func (ph *PrivilegeHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	var privilegeHistory []*models.PrivilegeHistory
	var err error
	var resp []byte

	userName := r.Header["X-User-Name"][0]
	privilegeHistory, err = ph.PrivilegeUseCase.GetHistory(userName)
	if err != nil {
		ph.Logger.Infow("can`t get privilege by username",
			"err:", err.Error())
		http.Error(w, "can`t get privilege by username", http.StatusInternalServerError)
		return
	}

	dto := models.PrivilegeHistoryToDTOs(privilegeHistory)

	resp, err = json.Marshal(dto)
	if err != nil {
		ph.Logger.Errorw("can`t marshal privilege",
			"err:", err.Error())
		http.Error(w, "can`t make privilege", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ph.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
