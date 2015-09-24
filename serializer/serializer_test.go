package serializer

import (
	"net/http"
	"testing"
	"time"

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

func TestUnmarshalReport(t *testing.T) {
	rep, err := UnmarshalReport([]byte(testReport))
	assert.Nil(t, err, err)
	assert.Equal(t, rep.Origin, "192.168.37.1:62634")
	assert.Equal(t, rep.Method, "POST")
	assert.Equal(t, rep.Status, 200)
	assert.Equal(t, rep.ContentType, "text/plain; charset=utf-8")
	assert.Equal(t, rep.ContentLength, uint(18))
	assert.Equal(t, rep.Host, "192.168.37.10:5000")
	assert.Equal(t, rep.URL, "http://192.168.37.10:5000/sum?b=423")
	assert.Equal(t, rep.Scheme, "http")
	assert.Equal(t, rep.Path, "/sum")
	assert.Equal(t, rep.Header, Header{Header: http.Header{"Content-Type": []string{"application/json-hal"}}})
	assert.Equal(t, rep.Body, "ewogICJyZXN1bHQiOiA1Mwp9")
	assert.Equal(t, rep.RequestHeader, Header{Header: http.Header{"Content-Type": []string{"application/json"}}})
	assert.Equal(t, rep.RequestBody, "eyJhIjogIjQzIiwgImIiOiAiMTAifQ==")
	assert.Equal(t, rep.DateStart, time.Date(2015, 9, 22, 16, 45, 59, 479125723, time.UTC))
	assert.Equal(t, rep.DateStart, time.Date(2015, 9, 22, 16, 45, 59, 479125723, time.UTC))
	assert.Equal(t, rep.TimeTaken, time.Date(2015, 9, 22, 16, 45, 59, 479237627, time.UTC))
}

func TestUnmarshalCrippledJson(t *testing.T) {
	_, err := UnmarshalReport([]byte(`{"request": "request body}`))
	assert.NotNil(t, err, "Error is not nil")
}
