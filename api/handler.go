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

		col, err = sess.Collection("demo")
		if err != nil {
			log.Fatal(err)
			http.Error(w, "oops", 500)
		}

		col.Append(Demo{
			FirstName: "Hayao",
			LastName:  "Miyazaki",
			Bio:       "Japanese film director.",
		})
		log.Info("data saved")
		defer sess.Close()
		respond(w, r, http.StatusOK, "sess")
	},
)
