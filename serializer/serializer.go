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

func UnmarshalReport(data []byte, r *Report) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Incorrect JSON %v", data)
		}
	}()

	err := json.Unmarshal(data, r)
	if err != nil {
		return err
	}
	return nil
}
