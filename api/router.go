package api

import (
	"github.com/gorilla/mux"
)

func RouterFactory() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/report", httpSaveReportHandler).Methods("POST")
	return r
}
