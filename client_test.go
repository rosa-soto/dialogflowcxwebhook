package dialogflowcxwebhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryDialogFlow(t *testing.T) {
	url := "https://dialogflow.googleapis.com/v3/projects"
	client := NewClient(url)
	resp, err := client.Post("how do I change my payout bank account?")
	assert.NotNil(t, resp)
	assert.NoError(t, err)
}
