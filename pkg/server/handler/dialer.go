package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.kostyadodich/demo/pkg/model"
	"github.kostyadodich/demo/pkg/repository"
	hm "github.kostyadodich/demo/pkg/server/model"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Dialer struct {
	dialerRepo *repository.Dialer
}

func NewDialer(dialerRepo *repository.Dialer) *Dialer {
	return &Dialer{dialerRepo: dialerRepo}
}

func (d *Dialer) CreateDialer(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handlDialer := hm.Dialer{}
	if err = json.Unmarshal(data, &handlDialer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	dialer := handlDialer.Convert()

	if err = d.dialerRepo.Create(dialer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Dialer) GetDialerByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dialer, err := d.dialerRepo.GetDialerByID(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(dialer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)

}

func (d *Dialer) GetDialers(w http.ResponseWriter, r *http.Request) {
	dialers, err := d.dialerRepo.DialerList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(dialers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write(data)

}

func (d *Dialer) UpdateDialer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	dialer := model.Dialer{}
	if err = json.Unmarshal(data, &dialer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	dialer.ID = int64(id)
	if err = d.dialerRepo.UpdateDialer(dialer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (d *Dialer) DeleteDialer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	if err = d.dialerRepo.DeleteDialer(int64(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
