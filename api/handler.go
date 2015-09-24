package api

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

var Database_uri string

func init() {
	databaseUrlUsage := `i.e.: postgres://postgres:mysecretpassword@postgres:5432/deadpool
	or you can use env variable: DATABASE_URI`

	flag.StringVar(&Database_uri, "d", os.Getenv("DATABASE_URI"), databaseUrlUsage)
	flag.Parse()
}

func withDb(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("View using database ", Database_uri)
		f(w, r)
	}
}

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

var httpReportHandler = withDb(
	func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, http.StatusOK, "# fasf ds?!")
	},
)
