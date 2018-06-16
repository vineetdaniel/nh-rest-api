package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	log.Print(connectionString)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("connected with DB")
		log.Print(err)
		log.Print(a.DB)
		errd := a.DB.Ping()
		if errd != nil {
			panic(errd.Error())
		}
	}

	a.intitializeRoutes()
}

func (a *App) Run(addr string) {
	log.Print("server")
	log.Fatal(http.ListenAndServe(addr, a.Router))

}

var logger = log.New(os.Stdout, "[something shiny] ", 0)

func httpLogger(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Printf("Started %s %s", r.Method, r.URL.Path)
		fn(w, r)
		logger.Printf("Completed in %v", time.Since(start))
	}
}
func (a *App) intitializeRoutes() {
	a.Router = mux.NewRouter()
	log.Print("initializign routes")
	s := http.StripPrefix("/img/", http.FileServer(http.Dir("./img/")))

	a.Router.PathPrefix("/img/").Handler(s)
	a.Router.HandleFunc("/health/", httpLogger(HealthCheckHandler))
	a.Router.HandleFunc("/cities/", httpLogger(a.getCitiesHandler))
	a.Router.HandleFunc("/city/{id:[0-9]+}/", a.getCityHandler).Methods("GET")
	a.Router.HandleFunc("/locations/", httpLogger(a.getLocationsHandler))
	a.Router.HandleFunc("/location/{id:[0-9]+}/", httpLogger(a.getLocationHandler))
	a.Router.Use(loggingMiddleware)
	http.Handle("/", a.Router)
	log.Print(a.Router)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func (a *App) getCityHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *App) getCitiesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("in getCitys")
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	citys, err := getCities(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responsdWithJson(w, http.StatusOK, citys)
}

func (a *App) getLocationsHandler(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	locations, err := getLocations(a.DB, count, start)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responsdWithJson(w, http.StatusOK, locations)
}

func (a *App) getLocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Location ID")
		return
	}

	l := location{ID: id}

	if err := l.getLocation(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	responsdWithJson(w, http.StatusOK, l)
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
