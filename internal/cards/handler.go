package cards

import (
	"encoding/json"
	"github.com/rotscher/cruds/internal/route"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	all := selectAll()
	jsonData, err := json.Marshal(all)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}

func GetById(w http.ResponseWriter, r *http.Request) {

	vars := route.Vars(r)
	if vars == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
		return
	}

	cardId := vars["cardId"]
	if cardId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("missing query param: cardId"))
		return
	}

	id, err := strconv.ParseInt(cardId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	card := selectById(id)
	if card == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	jsonData, err := json.Marshal(card)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var card Card
	err = json.Unmarshal(body, &card)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := insert(&card)
	if result {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
