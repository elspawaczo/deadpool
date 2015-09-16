package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		log.Println("respond: ", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "#?!")
}
