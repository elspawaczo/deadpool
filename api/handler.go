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

	report := Report{
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
	if _, err := col.Append(&report); err != nil {
		log.Fatal("Save data to database error: ", err)
		http.Error(w, "oops", http.StatusBadRequest)
		return
	}
	log.Info("data saved: ", report)

	doc, err := json.Marshal(report)
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
				log.Fatal("Panic: ", err, " ", Database_uri)
				http.Error(w, "oops", http.StatusInternalServerError)
			}
		}()

		var err error
		var sess db.Database
		var settings postgresql.ConnectionURL

		settings, err = postgresql.ParseURL(Database_uri)
		if err != nil {
			log.Error("Connection string cannot be parsed: ", Database_uri)
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
