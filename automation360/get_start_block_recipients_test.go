package automation360

import (
	"fmt"
	"github.com/dimuska139/sendpulse-sdk-go"
	"github.com/dimuska139/sendpulse-sdk-go/client"
	"github.com/icrowley/fake"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAutomation360_GetStartBlockRecipients(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	blockID := 1
	responseBody, _ := ioutil.ReadFile("./testdata/startBlockRecipients.json")
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/a360/stats/main-trigger/%d/addresses", client.ApiBaseUrl, blockID),
		httpmock.NewBytesResponder(http.StatusOK, responseBody),
	)

	config := sendpulse.Config{
		UserID: fake.CharactersN(50),
		Secret: fake.CharactersN(50),
		Token:  fake.Word(),
	}

	spClient := New(http.DefaultClient, &config)

	recipients, total, err := spClient.GetStartBlockRecipients(blockID, 10, 0, "desc", "id")
	assert.NoError(t, err)
	assert.Equal(t, 12345, int(*recipients[0].FlowID))
	assert.Equal(t, 1, *total)
}

func TestAutomation360_GetStartBlockRecipients_HttpError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	blockID := 1
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/a360/stats/main-trigger/%d/addresses", client.ApiBaseUrl, blockID),
		httpmock.NewStringResponder(http.StatusInternalServerError, ""),
	)

	config := sendpulse.Config{
		UserID: fake.CharactersN(50),
		Secret: fake.CharactersN(50),
		Token:  fake.Word(),
	}

	spClient := New(http.DefaultClient, &config)

	recipients, total, err := spClient.GetStartBlockRecipients(blockID, 10, 0, "desc", "id")
	assert.Error(t, err)
	assert.Equal(t, 0, len(recipients))
	assert.Nil(t, total)
}

func TestAutomation360_GetStartBlockRecipients_InvalidJson(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	blockID := 1
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/a360/stats/main-trigger/%d/addresses", client.ApiBaseUrl, blockID),
		httpmock.NewStringResponder(http.StatusOK, ""),
	)

	config := sendpulse.Config{
		UserID: fake.CharactersN(50),
		Secret: fake.CharactersN(50),
		Token:  fake.Word(),
	}

	spClient := New(http.DefaultClient, &config)

	recipients, total, err := spClient.GetStartBlockRecipients(blockID, 10, 0, "desc", "id")
	assert.Error(t, err)
	assert.Equal(t, 0, len(recipients))
	assert.Nil(t, total)
}

func TestAutomation360_GetStartBlockRecipients_NoData(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	blockID := 1
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/a360/stats/main-trigger/%d/addresses", client.ApiBaseUrl, blockID),
		httpmock.NewStringResponder(http.StatusOK, "{}"),
	)

	config := sendpulse.Config{
		UserID: fake.CharactersN(50),
		Secret: fake.CharactersN(50),
		Token:  fake.Word(),
	}

	spClient := New(http.DefaultClient, &config)

	recipients, total, err := spClient.GetStartBlockRecipients(blockID, 10, 0, "desc", "id")
	assert.Error(t, err)
	assert.Equal(t, 0, len(recipients))
	assert.Nil(t, total)
}
