package main

import (
	"net/http"

	"github.com/thisissoon/deadpool/api"
)

func main() {
	r := api.RouterFactory()
	http.ListenAndServe(":8000", r)
}
