package fhir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jbogarin/fhir-cli/internal/auth"
	"github.com/jbogarin/fhir-cli/internal/config"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Profile    *config.Profile
}

type Bundle struct {
	ResourceType string        `json:"resourceType"`
	Type         string        `json:"type"`
	Total        int           `json:"total,omitempty"`
	Link         []BundleLink  `json:"link,omitempty"`
	Entry        []BundleEntry `json:"entry,omitempty"`
}

type BundleLink struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}

type BundleEntry struct {
	FullURL  string          `json:"fullUrl,omitempty"`
	Resource json.RawMessage `json:"resource"`
	Search   *SearchInfo     `json:"search,omitempty"`
}

type SearchInfo struct {
	Mode  string  `json:"mode,omitempty"`
	Score float64 `json:"score,omitempty"`
}

type OperationOutcome struct {
	ResourceType string  `json:"resourceType"`
	Issue        []Issue `json:"issue"`
}

type Issue struct {
	Severity    string         `json:"severity"`
	Code        string         `json:"code"`
	Details     *CodeableConcept `json:"details,omitempty"`
	Diagnostics string         `json:"diagnostics,omitempty"`
	Location    []string       `json:"location,omitempty"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

func NewClient(profile *config.Profile) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(profile.FHIRBaseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		Profile: profile,
	}
}

// VerboseMode enables request/response logging
var VerboseMode bool

func (c *Client) doRequest(method, path string, body io.Reader) ([]byte, error) {
	accessToken, err := auth.GetAccessToken(c.Profile, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	fullURL := c.BaseURL + "/" + strings.TrimPrefix(path, "/")

	if VerboseMode {
		fmt.Printf("\n=== FHIR Request ===\n")
		fmt.Printf("%s %s\n", method, fullURL)
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/fhir+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/fhir+json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var outcome OperationOutcome
		if json.Unmarshal(respBody, &outcome) == nil && outcome.ResourceType == "OperationOutcome" {
			return nil, formatOperationOutcome(&outcome)
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func formatOperationOutcome(outcome *OperationOutcome) error {
	var messages []string
	for _, issue := range outcome.Issue {
		msg := fmt.Sprintf("[%s] %s", issue.Severity, issue.Code)
		if issue.Diagnostics != "" {
			msg += ": " + issue.Diagnostics
		}
		if issue.Details != nil && issue.Details.Text != "" {
			msg += " - " + issue.Details.Text
		}
		messages = append(messages, msg)
	}
	return fmt.Errorf("FHIR error: %s", strings.Join(messages, "; "))
}

// Get retrieves a single resource by ID
func (c *Client) Get(resourceType, id string) (json.RawMessage, error) {
	path := fmt.Sprintf("%s/%s", resourceType, id)
	return c.doRequest("GET", path, nil)
}

// Search performs a search operation on a resource type
func (c *Client) Search(resourceType string, params map[string]string) (*Bundle, error) {
	path := resourceType
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		path += "?" + values.Encode()
	}

	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse bundle: %w", err)
	}

	return &bundle, nil
}

// Create creates a new resource
func (c *Client) Create(resourceType string, resource interface{}) (json.RawMessage, error) {
	body, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	return c.doRequest("POST", resourceType, strings.NewReader(string(body)))
}

// Update updates an existing resource
func (c *Client) Update(resourceType, id string, resource interface{}) (json.RawMessage, error) {
	body, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	path := fmt.Sprintf("%s/%s", resourceType, id)
	return c.doRequest("PUT", path, strings.NewReader(string(body)))
}

// Delete removes a resource
func (c *Client) Delete(resourceType, id string) error {
	path := fmt.Sprintf("%s/%s", resourceType, id)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}

// GetNextPage follows a pagination link
func (c *Client) GetNextPage(bundle *Bundle) (*Bundle, error) {
	for _, link := range bundle.Link {
		if link.Relation == "next" {
			data, err := c.doRequest("GET", link.URL, nil)
			if err != nil {
				return nil, err
			}

			var nextBundle Bundle
			if err := json.Unmarshal(data, &nextBundle); err != nil {
				return nil, fmt.Errorf("failed to parse bundle: %w", err)
			}

			return &nextBundle, nil
		}
	}
	return nil, nil
}

// GetMetadata retrieves the CapabilityStatement
func (c *Client) GetMetadata() (json.RawMessage, error) {
	return c.doRequest("GET", "metadata", nil)
}

// Operation performs a FHIR operation
func (c *Client) Operation(resourceType, id, operation string, params map[string]string) (json.RawMessage, error) {
	var path string
	if id != "" {
		path = fmt.Sprintf("%s/%s/$%s", resourceType, id, operation)
	} else if resourceType != "" {
		path = fmt.Sprintf("%s/$%s", resourceType, operation)
	} else {
		path = fmt.Sprintf("$%s", operation)
	}

	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		path += "?" + values.Encode()
	}

	return c.doRequest("GET", path, nil)
}
