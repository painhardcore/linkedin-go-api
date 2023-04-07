package linkedinapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	baseURL                      = "https://api.linkedin.com/v2"
	mediaUploadBaseURL           = "https://api.linkedin.com/mediaUpload"
	ugcPostsEndpoint             = "/ugcPosts"
	assetsEndpoint               = "/assets"
	xRestliProtocolVersionHeader = "X-Restli-Protocol-Version"
	protocolVersion              = "2.0.0"
)

// ShareContentRequest represents the main structure of the share content request payload.
type ShareContentRequest struct {
	Author          string          `json:"author"`
	LifecycleState  string          `json:"lifecycleState"`
	SpecificContent SpecificContent `json:"specificContent"`
	Visibility      Visibility      `json:"visibility"`
}

// SpecificContent represents the content of the share.
type SpecificContent struct {
	ShareContent ShareContent `json:"com.linkedin.ugc.ShareContent"`
}

// ShareContent holds the details of the share commentary, media category, and media.
type ShareContent struct {
	ShareCommentary    ShareCommentary `json:"shareCommentary"`
	ShareMediaCategory string          `json:"shareMediaCategory"`
	Media              []ShareMedia    `json:"media,omitempty"`
}

// ShareCommentary provides the primary content for the share.
type ShareCommentary struct {
	Text string `json:"text"`
}

// ShareMedia represents the media assets attached to the share.
type ShareMedia struct {
	Status      string           `json:"status"`
	Description ShareDescription `json:"description,omitempty"`
	Media       string           `json:"media,omitempty"`
	OriginalURL string           `json:"originalUrl,omitempty"`
	Title       ShareTitle       `json:"title,omitempty"`
}

// ShareDescription is a short description for the image or article.
type ShareDescription struct {
	Text string `json:"text"`
}

// ShareTitle is the custom title of the image or article.
type ShareTitle struct {
	Text string `json:"text"`
}

// Visibility represents the visibility restrictions for the share.
type Visibility struct {
	MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
}

// RegisterUploadRequest represents the request payload for registering an image upload.
type RegisterUploadRequest struct {
	RegisterUpload RegisterUpload `json:"registerUploadRequest"`
}

// RegisterUpload holds the details of the image registration.
type RegisterUpload struct {
	Recipes              []string              `json:"recipes"`
	Owner                string                `json:"owner"`
	ServiceRelationships []ServiceRelationship `json:"serviceRelationships"`
}

// ServiceRelationship defines the relationship type and identifier for the image registration.
type ServiceRelationship struct {
	RelationshipType string `json:"relationshipType"`
	Identifier       string `json:"identifier"`
}

// RegisterUploadResponse represents the response payload for registering an image upload.
type RegisterUploadResponse struct {
	Value RegisterUploadValue `json:"value"`
}

// RegisterUploadValue holds the details of the upload mechanism, media artifact, and asset.
type RegisterUploadValue struct {
	UploadMechanism UploadMechanism `json:"uploadMechanism"`
	MediaArtifact   string          `json:"mediaArtifact"`
	Asset           string          `json:"asset"`
}

// UploadMechanism contains the details of the media upload HTTP request.
type UploadMechanism struct {
	MediaUploadHttpRequest MediaUploadHttpRequest `json:"com.linkedin.digitalmedia.uploading.MediaUploadHttpRequest"`
}

// MediaUploadHttpRequest holds the headers and upload URL for the image upload.
type MediaUploadHttpRequest struct {
	Headers   map[string]string `json:"headers"`
	UploadURL string            `json:"uploadUrl"`
}

func (c *Client) ShareText(personURN, text string) (string, error) {
	shareRequest := ShareContentRequest{
		Author:         personURN,
		LifecycleState: "PUBLISHED",
		SpecificContent: SpecificContent{
			ShareContent: ShareContent{
				ShareCommentary: ShareCommentary{
					Text: text,
				},
				ShareMediaCategory: "NONE",
			},
		},
		Visibility: Visibility{
			MemberNetworkVisibility: "PUBLIC",
		},
	}

	return c.CreateShare(shareRequest)
}

// Client struct for LinkedIn API
type Client struct {
	AccessToken string
	httpClient  *http.Client
}

// NewClient initializes a new LinkedIn API client
func NewClient(accessToken string) *Client {
	return &Client{
		AccessToken: accessToken,
		httpClient:  &http.Client{},
	}
}

func (c *Client) doRequest(method, url string, requestBody interface{}) (*http.Response, error) {
	// Convert the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) CreateShare(shareRequest ShareContentRequest) (string, error) {
	url := baseURL + ugcPostsEndpoint
	resp, err := c.doRequest(http.MethodPost, url, shareRequest)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Println(string(body))

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("failed to create share")
	}

	return resp.Header.Get("X-RestLi-Id"), nil
}

func (c *Client) RegisterUpload(registerUploadRequest RegisterUploadRequest) (RegisterUploadResponse, error) {
	url := baseURL + assetsEndpoint + "?action=registerUpload"
	resp, err := c.doRequest(http.MethodPost, url, registerUploadRequest)
	if err != nil {
		return RegisterUploadResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return RegisterUploadResponse{}, errors.New("failed to register upload")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RegisterUploadResponse{}, err
	}

	var registerUploadResponse RegisterUploadResponse
	err = json.Unmarshal(body, &registerUploadResponse)
	if err != nil {
		return RegisterUploadResponse{}, err
	}

	return registerUploadResponse, nil
}

func (c *Client) UploadImage(uploadURL, imagePath string) error {
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewBuffer(imgData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to upload image")
	}

	return nil
}

func (c *Client) ShareArticle(personURN, text, url, title, description string) (string, error) {
	shareRequest := ShareContentRequest{
		Author:         personURN,
		LifecycleState: "PUBLISHED",
		SpecificContent: SpecificContent{
			ShareContent: ShareContent{
				ShareCommentary: ShareCommentary{
					Text: text,
				},
				ShareMediaCategory: "ARTICLE",
				Media: []ShareMedia{
					{
						Status:      "READY",
						Description: ShareDescription{Text: description},
						OriginalURL: url,
						Title:       ShareTitle{Text: title},
					},
				},
			},
		},
		Visibility: Visibility{
			MemberNetworkVisibility: "PUBLIC",
		},
	}

	return c.CreateShare(shareRequest)
}
