package serializer

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

type Report struct {
	HttpRequest  string `json:"request"`
	HttpResponse string `json:"respond"`
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	Description  string `json:"description"`
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
