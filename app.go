package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("connected with DB")
		log.Print(err)
	}
	a.Router = mux.NewRouter()
	a.intitializeRoutes()
}

func (a *App) Run(addr string) {
	log.Print("server")
	log.Fatal(http.ListenAndServe(addr, a.Router))

}
func (a *App) intitializeRoutes() {
	a.Router.HandleFunc("/city", a.getCitys).Methods("GET")
	a.Router.HandleFunc("/city/{id:[0-9]+}", a.getCity).Methods("GET")
}

func (a *App) getCity(w http.ResponseWriter, r *http.Request) {
	log.Printf("in city")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid City ID")
		return
	}

	c := city{ID: id}

	if err := c.getCity(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	responsdWithJson(w, http.StatusOK, c)
}

func (a *App) getCitys(w http.ResponseWriter, r *http.Request) {
	log.Printf("in getCitys")
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	citys, err := getCitys(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responsdWithJson(w, http.StatusOK, citys)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	responsdWithJson(w, code, map[string]string{"error": message})
}

func responsdWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
