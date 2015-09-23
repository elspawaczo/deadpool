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
	Origin        string    `json:"origin",json"`
	Method        string    `json:"method",json"`
	Status        int       `json:"status",json"`
	ContentType   string    `json:"content_type",json"`
	ContentLength uint      `json:"content_length",json"`
	Host          string    `json:"host",json"`
	URL           string    `json:"url",json"`
	Scheme        string    `json:"scheme",json"`
	Path          string    `json:"path",path"`
	Header        Header    `json:"header",json"`
	Body          string    `json:"body",json"`
	RequestHeader Header    `json:"request_header",json"`
	RequestBody   string    `json:"request_body",json"`
	DateStart     time.Time `json:"date_start",json"`
	DateEnd       time.Time `json:"date_end",json"`
	TimeTaken     time.Time `json:"time_taken",json"`
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
