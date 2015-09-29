package api

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/nvellon/hal"
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

type Response struct {
	Count int
	Total int
}

func (p Response) GetMap() hal.Entry {
	return hal.Entry{
		"count": p.Count,
		"total": p.Total,
	}
}

func (c Report) GetMap() hal.Entry {
	return hal.Entry{
		"id":             c.Id,
		"origin":         c.Origin,
		"method":         c.Method,
		"status":         c.Status,
		"content_type":   c.ContentType,
		"content_length": c.ContentLength,
		"host":           c.Host,
		"url":            c.URL,
		"scheme":         c.Scheme,
		"path":           c.Path,
		"body":           c.Body,
		"request_body":   c.RequestBody,
		"date_start":     c.DateStart,
		"date_end":       c.DateEnd,
		"time_taken":     c.TimeTaken,
		"created":        c.Ts,
	}
}

type ReportRest struct {
	sess db.Database
}

func (self ReportRest) Get(w http.ResponseWriter, r *http.Request) {
	col, err := self.sess.Collection("report")
	if err != nil {
		log.Fatal("Getting collection from db: ", err)
		http.Error(w, "oops", http.StatusBadRequest)
		return
	}

	limit := uint(20)
	reps := col.Find().Limit(limit)
	var reports []Report
	reps.All(&reports)
	c, _ := reps.Count()

	halDoc := hal.NewResource(Response{Count: len(reports), Total: int(c)}, "")
	for _, rep := range reports {
		halDoc.Embed("reports", hal.NewResource(rep, ""))
	}
	doc, _ := json.Marshal(halDoc)
	w.WriteHeader(http.StatusOK)
	w.Write(doc)
}

func (self ReportRest) Post(w http.ResponseWriter, r *http.Request) {
	col, err := self.sess.Collection("report")
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
		var sess db.Database

		sess = context.Get(r, "db").(db.Database)
		if sess == nil {
			log.Error("Cannot get `db` context")
			http.Error(w, "oops", http.StatusBadRequest)
			return
		}

		api := ReportRest{sess: sess}
		switch {
		case r.Method == "POST":
			api.Post(w, r)
		case r.Method == "GET":
			api.Get(w, r)
		default:
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	},
)
