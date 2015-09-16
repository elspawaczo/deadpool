package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testReport string = `
{
    "request": "request body",
    "respond": "respond body",
    "source": "source name",
    "destination": "destination name",
    "description": "some description"
}
`

func TestUnmarshalReport(t *testing.T) {
	var rep Report
	err := UnmarshalReport([]byte(testReport), &rep)
	assert.Nil(t, err, err)
	assert.Equal(t, "request body", rep.HttpRequest)
	assert.Equal(t, "respond body", rep.HttpResponse)
	assert.Equal(t, "source name", rep.Source)
	assert.Equal(t, "destination name", rep.Destination)
	assert.Equal(t, "some description", rep.Description)
}

func TestUnmarshalCrippledJson(t *testing.T) {
	var rep Report
	err := UnmarshalReport([]byte(`{"request": "request body}`), &rep)
	assert.NotNil(t, err, "Error is not nil")
}
