package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"

	"github.com/thisissoon/deadpool/api"
)

func main() {
	r := api.RouterFactory()

	log.Info("DeadPool is running")
	http.ListenAndServe(":8000", r)
}
