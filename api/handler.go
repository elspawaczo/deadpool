package api

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"upper.io/db"
	"upper.io/db/postgresql"

	"github.com/thisissoon/deadpool/serializer"
)

var Database_uri string

func init() {
	databaseUrlUsage := `i.e.: postgres://postgres:mysecretpassword@postgres:5432/deadpool
	or you can use env variable: DATABASE_URI`

	flag.StringVar(&Database_uri, "d", os.Getenv("DATABASE_URI"), databaseUrlUsage)
	flag.Parse()
}

type Demo struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Bio       string `db:"bio"`
}

func withDb(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var sess db.Database
		var settings postgresql.ConnectionURL

		settings, err = postgresql.ParseURL(Database_uri)
		if err != nil {
			log.Error("Connection string cannot be parsed ", Database_uri)
			http.Error(w, "oops", 500)
		}
		sess, err = db.Open(postgresql.Adapter, settings)
		defer sess.Close()

		if err != nil {
			log.Error("Canot connect to database ", err)
			http.Error(w, "oops", 500)
		}
		context.Set(r, "db", sess)
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

var testReport string = `
{
    "origin": "192.168.37.1:62634",
    "method": "POST",
    "status": 200,
    "content_type": "text/plain; charset=utf-8",
    "content_length": 18,
    "host": "192.168.37.10:5000",
    "url": "http://192.168.37.10:5000/sum?b=423",
    "scheme": "http",
    "path": "/sum",
    "header": {
        "Header": {
            "Content-Type": [
                "application/json-hal"
            ]
        }
    },
    "body": "ewogICJyZXN1bHQiOiA1Mwp9",
    "request_header": {
        "Header": {
            "Content-Type": [
                "application/json"
            ]
        }
    },
    "request_body": "eyJhIjogIjQzIiwgImIiOiAiMTAifQ==",
    "date_start": "2015-09-22T16:45:59.479125723Z",
    "date_end": "2015-09-22T16:45:59.479237627Z",
    "time_taken": "2015-09-22T16:45:59.479237627Z"
}
`

var httpReportHandler = withDb(
	func(w http.ResponseWriter, r *http.Request) {
		var err error
		var sess db.Database
		var col db.Collection

		sess = context.Get(r, "db").(db.Database)
		if sess == nil {
			log.Error("err?")
			http.Error(w, "oops", 500)
		}

		col, err = sess.Collection("report")
		if err != nil {
			log.Fatal(err)
			http.Error(w, "oops", 500)
		}

		rep, err := serializer.UnmarshalReport([]byte(testReport))
		if err != nil {
			log.Error(err)
			http.Error(w, "oops", 422)
		}

		rrr, err := col.Append(Report{
			Origin:        rep.Origin,
			Method:        rep.Method,
			Status:        rep.Status,
			ContentType:   rep.ContentType,
			ContentLength: rep.ContentLength,
			Host:          rep.Host,
			URL:           rep.URL,
			Scheme:        rep.Scheme,
			Path:          rep.Path,
			Body:          rep.Body,
			RequestBody:   rep.RequestBody,
			DateStart:     rep.DateStart,
			DateEnd:       rep.DateEnd,
			TimeTaken:     rep.TimeTaken,
		})
		if err != nil {
			log.Fatal(err)
			http.Error(w, "oops", 500)
		}
		log.Info("data saved")
		respond(w, r, http.StatusOK, rrr)
	},
)
