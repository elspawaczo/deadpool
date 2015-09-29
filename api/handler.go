package api

import (
	"bytes"
	"encoding/json"
	"flag"
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

func withDb(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal("Panic: ", err, Database_uri)
				http.Error(w, "oops", http.StatusInternalServerError)
			}
		}()

		var err error
		var sess db.Database
		var settings postgresql.ConnectionURL

		Database_uri := "postgres://postgres:mysecretpassword@172.17.0.1:5432/deadpool"
		settings, err = postgresql.ParseURL(Database_uri)
		if err != nil {
			log.Error("Connection string cannot be parsed ", Database_uri)
			http.Error(w, "oops", http.StatusBadRequest)
		}
		sess, err = db.Open(postgresql.Adapter, settings)
		defer sess.Close()

		if err != nil {
			log.Error("Canot connect to database: ", err)
			http.Error(w, "oops", http.StatusBadRequest)
		}
		context.Set(r, "db", sess)
		f(w, r)
	}
}

var httpReportHandler = withDb(
	func(w http.ResponseWriter, r *http.Request) {
		var err error
		var sess db.Database
		var col db.Collection

		sess = context.Get(r, "db").(db.Database)
		if sess == nil {
			log.Error("err?")
			http.Error(w, "oops", http.StatusBadRequest)
			return
		}

		col, err = sess.Collection("report")
		if err != nil {
			log.Fatal("Getting collection from db: ", err)
			http.Error(w, "oops", http.StatusBadRequest)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		rep, err := serializer.UnmarshalReport(buf.Bytes())
		if err != nil {
			log.Error("Unmarshal report: ", err)
			http.Error(w, "oops", 422)
			return
		}

		rrr := Report{
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
		}
		if _, err := col.Append(&rrr); err != nil {
			log.Fatal("Save data to database error: ", err)
			http.Error(w, "oops", http.StatusBadRequest)
			return
		}
		log.Info("data saved: ", rrr)

		doc, err := json.Marshal(rrr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(doc)
	},
)
