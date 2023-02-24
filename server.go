package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"byoa/redis"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var store = make(map[string]string)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/texts", PostHandler).Methods("POST")
	r.HandleFunc("/texts/{uuid}", GetHandler).Methods("GET")
	r.HandleFunc("/texts/{uuid}", DeleteHandler).Methods("DELETE")
	r.HandleFunc("/texts/{uuid}/search", SearchHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

type Data struct {
	Data string `json:"data"`
}

type PostResp struct {
	Id string `json:"id"`
}

type Error struct {
	Error string `json:"error"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	var data Data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	}

	encoder := json.NewEncoder(w)

	// generate id
	id := uuid.New().String()

	// store text
	//store[id] = data.Data
	err = redis.Store(id, data.Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, Error{Error: err.Error()})
	}

	resp := PostResp{
		Id: id,
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(resp)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)

	// get text from store
	//data, ok := store[uuid]
	data, found, err := redis.GetText(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, Error{Error: err.Error()})
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(Error{Error: "not found"})
		return
	}

	encoder.Encode(Data{Data: data})
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	encoder := json.NewEncoder(w)

	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := store[uuid]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(Error{Error: "not found"})
		return
	}

	//delete(store, uuid)
	err := redis.Delete(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, Error{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type Result struct {
	Found bool `json:"found"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	encoder := json.NewEncoder(w)

	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(Error{Error: "invalid uuid"})
		return
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(Error{Error: "invalid term"})
		return
	}

	// get text from store
	data, ok := store[uuid]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(Error{Error: "not found"})
		return
	}

	var res Result
	if strings.Contains(data, term) {
		res.Found = true
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(res)

}
