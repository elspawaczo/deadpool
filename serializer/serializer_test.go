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
	rep, err := UnmarshalReport([]byte(testReport))
	assert.Nil(t, err, err)
	assert.Equal(t, "request body", rep.HttpRequest)
	assert.Equal(t, "respond body", rep.HttpResponse)
	assert.Equal(t, "source name", rep.Source)
	assert.Equal(t, "destination name", rep.Destination)
	assert.Equal(t, "some description", rep.Description)
}

func TestUnmarshalCrippledJson(t *testing.T) {
	_, err := UnmarshalReport([]byte(`{"request": "request body}`))
	assert.NotNil(t, err, "Error is not nil")
}
