package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"upper.io/db"
	"upper.io/db/postgresql"

	"github.com/stretchr/testify/assert"
)

// var testReport string = `{"origin":"192.168.37.1:62634","method":"POST","status":200,"content_type":"text/plain; charset=utf-8","content_length":18,"host":"192.168.37.10:5000","url":"http://192.168.37.10:5000/sum?b=423","scheme":"http","path":"/sum","header":{"Header":{"Content-Type":["application/json-hal"]}},"body":"ewogICJyZXN1bHQiOiA1Mwp9","request_header":{"Header":{"Content-Type":["application/json"]}},"request_body":"eyJhIjogIjQzIiwgImIiOiAiMTAifQ==","date_start":"2015-09-22T16:45:59.479125723Z","date_end":"2015-09-22T16:45:59.479237627Z","time_taken":"2015-09-22T16:45:59.479237627Z"}`
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

func TestHttpResponseHandler(t *testing.T) {
	conn := "postgres://postgres:mysecretpassword@172.17.0.1:5432/deadpool"
	os.Setenv("DATABASE_URI", conn)

	r, _ := http.NewRequest(
		"POST",
		"/report",
		bytes.NewBuffer([]byte(testReport)))
	w := httptest.NewRecorder()

	// Call our handler
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
