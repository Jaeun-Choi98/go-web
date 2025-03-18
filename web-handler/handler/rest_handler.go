package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"turker.web/model"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := model.User{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := model.CreateUser(newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if strId, exist := vars["id"]; exist {
		id, err := strconv.Atoi(strId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := model.GetUserById(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
		}
		w.Header().Set("content-type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	delUser := model.User{}
	if err := json.NewDecoder(r.Body).Decode(&delUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := model.UpdateUser(delUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateUserByEmail(w http.ResponseWriter, r *http.Request) {
	updateUser := model.User{}
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := model.UpdateUser(updateUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PutUserByEmail(w http.ResponseWriter, r *http.Request) {
	putUser := model.User{}
	if err := json.NewDecoder(r.Body).Decode(&putUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := model.PutUser(putUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
