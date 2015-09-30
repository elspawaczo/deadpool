package serializer

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

type Header struct {
	http.Header
}

type Report struct {
	Origin        string    `json:"origin"`
	Method        string    `json:"method"`
	Status        int       `json:"status"`
	ContentType   string    `json:"content_type"`
	ContentLength uint      `json:"content_length"`
	Host          string    `json:"host"`
	URL           string    `json:"url"`
	Scheme        string    `json:"scheme"`
	Path          string    `json:"path",path"`
	Header        Header    `json:"header"`
	Body          string    `json:"body"`
	RequestHeader Header    `json:"request_header"`
	RequestBody   string    `json:"request_body"`
	DateStart     time.Time `json:"date_start"`
	DateEnd       time.Time `json:"date_end"`
	TimeTaken     time.Time `json:"time_taken"`
}

func UnmarshalReport(data []byte) (*Report, error) {
	r := &Report{}
	defer func() {
		if err := recover(); err != nil {
			log.Error("Incorrect JSON %v", data)
		}
	}()

	if err := json.Unmarshal(data, r); err != nil {
		return nil, err
	}

	return r, nil
}
