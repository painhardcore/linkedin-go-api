package linkedinapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

// MockTransport is a mock HTTP round tripper that returns a fixed response and error
type MockTransport struct {
	response *http.Response
	err      error
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.response, t.err
}

// TestCreateShare tests the CreateShare method
func TestCreateShare(t *testing.T) {
	// Set up test client
	testAccessToken := "test-access-token"
	testClient := &Client{AccessToken: testAccessToken, httpClient: http.DefaultClient}

	// Create a mock HTTP round tripper that returns a valid response
	mockTransport := &MockTransport{
		response: &http.Response{
			StatusCode: http.StatusCreated,
			Header:     http.Header{"X-Restli-Id": []string{"1234"}},
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		},
		err: nil,
	}

	// Set the mock HTTP round tripper for the default HTTP client
	http.DefaultClient.Transport = mockTransport

	// Set up test share request
	testShareRequest := ShareContentRequest{
		Author:         "urn:li:person:<person-id>",
		LifecycleState: "PUBLISHED",
		SpecificContent: SpecificContent{
			ShareContent: ShareContent{
				ShareCommentary: ShareCommentary{
					Text: "Test share",
				},
				ShareMediaCategory: "NONE",
			},
		},
		Visibility: Visibility{
			MemberNetworkVisibility: "PUBLIC",
		},
	}

	// Call the CreateShare method
	shareID, err := testClient.CreateShare(testShareRequest)
	if err != nil {
		t.Errorf("Failed to create share: %v", err)
	}

	// Check the share ID
	expectedShareID := "1234"
	if shareID != expectedShareID {
		t.Errorf("Invalid share ID: got %s, expected %s", shareID, expectedShareID)
	}
}
