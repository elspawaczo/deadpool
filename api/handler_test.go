package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nvellon/hal"
	"upper.io/db"
	"upper.io/db/postgresql"

	"github.com/stretchr/testify/assert"
)

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

func TestCreateNewReport(t *testing.T) {
	conn := os.Getenv("DATABASE_URI")

	r, _ := http.NewRequest(
		"POST",
		"/report",
		bytes.NewBuffer([]byte(testReport)))
	w := httptest.NewRecorder()

	httpReportHandler(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	settings, _ := postgresql.ParseURL(conn)
	sess, _ := db.Open(postgresql.Adapter, settings)

	reportCollection, _ := sess.Collection("report")
	defer reportCollection.Truncate()

	dbRes := reportCollection.Find()
	c, _ := dbRes.Count()
	assert.Equal(t, 1, int(c))
}

func TestGetReports(t *testing.T) {
	conn := os.Getenv("DATABASE_URI")

	settings, _ := postgresql.ParseURL(conn)
	sess, _ := db.Open(postgresql.Adapter, settings)

	reportCollection, _ := sess.Collection("report")
	reportCollection.Append(Report{})
	reportCollection.Append(Report{})
	defer reportCollection.Truncate()

	r, _ := http.NewRequest(
		"GET",
		"/report",
		bytes.NewBuffer([]byte(testReport)))
	w := httptest.NewRecorder()

	httpReportHandler(w, r)

	reps := reportCollection.Find()

	var reports []Report
	reps.All(&reports)
	assert.Equal(t, http.StatusOK, w.Code)

	c, _ := reps.Count()
	halDoc := hal.NewResource(Response{Count: len(reports), Total: int(c)}, "")
	for _, rep := range reports {
		halDoc.Embed("reports", hal.NewResource(rep, ""))
	}
	doc, _ := json.Marshal(halDoc)
	assert.Equal(t, string(doc), w.Body.String())
}
