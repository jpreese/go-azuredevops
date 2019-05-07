package azuredevops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://dev.azure.com"
	userAgent      = "go-azuredevops"
)

// Client for interacting with the Azure DevOps API
type Client struct {
	client *http.Client

	// BaseURL Comprised of baseURL and account
	BaseURL   url.URL
	UserAgent string

	// Account Required part of BaseURL
	Account string
	// Project Default project for api calls
	Project   string
	AuthToken string

	// Services used to proxy to other API endpoints
	Boards           *BoardsService
	BuildDefinitions *BuildDefinitionsService
	Builds           *BuildsService
	DeliveryPlans    *DeliveryPlansService
	Favourites       *FavouritesService
	Git              *GitService
	Iterations       *IterationsService
	PullRequests     *PullRequestsService
	Teams            *TeamsService
	Tests            *TestsService
	WorkItems        *WorkItemsService
}

// NewClient returns a new Azure DevOps API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
// The client's base URL is constructed from the supplied account and project.
// Token is a personal access token.
func NewClient(account string, project string, token string, httpClient *http.Client) (*Client, error) {
	if account == "" {
		return nil, fmt.Errorf("Missing valid account in call to NewClient(): account = %s", account)
	}

	if project == "" {
		return nil, fmt.Errorf("Missing valid project in call to NewClient(): project = %s", project)
	}

	if token == "" {
		return nil, fmt.Errorf("Missing personal access token in call to NewClient(): token = %s", token)
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// BaseURL
	baseURLstr := fmt.Sprintf("%s/%s/", defaultBaseURL, account)
	baseURL, _ := url.Parse(baseURLstr)

	c := &Client{
		client:    httpClient,
		BaseURL:   *baseURL,
		Account:   account,
		Project:   project,
		AuthToken: token,
	}

	c.Boards = &BoardsService{client: c}
	c.BuildDefinitions = &BuildDefinitionsService{client: c}
	c.Builds = &BuildsService{client: c}
	c.Favourites = &FavouritesService{client: c}
	c.Git = &GitService{client: c}
	c.Iterations = &IterationsService{client: c}
	c.PullRequests = &PullRequestsService{client: c}
	c.WorkItems = &WorkItemsService{client: c}
	c.Teams = &TeamsService{client: c}
	c.Tests = &TestsService{client: c}
	c.DeliveryPlans = &DeliveryPlansService{client: c}

	return c, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL.String())
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := new(url.URL)
	// If caller supplies an absolute urlStr, pass it unmodified to http.NewRequest
	if parsed.IsAbs() {
		u = parsed
	} else {
		// If caller supplies a relative URI in urlStr, prepend client project name
		s := fmt.Sprintf("%s/%s", url.PathEscape(c.Project), urlStr)
		u, err = c.BaseURL.Parse(s)
		if err != nil {
			return nil, err
		}
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Execute sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by r, or returned as an
// error if an API error has occurred. If r implements the io.Writer
// interface, the raw response body will be written to r, without attempting to
// first decode it. If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
// Execute runs all the http requests on the API
func (c *Client) Execute(ctx context.Context, request *http.Request, r interface{}) (*http.Response, error) {
	request = request.WithContext(ctx)
	request.SetBasicAuth("", c.AuthToken)

	//client := &http.Client{}
	response, err := c.client.Do(request)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return nil, e
			}
		}

		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 && response.StatusCode != 201 {
		return nil, fmt.Errorf("Request to %s responded with status %d", request.URL, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Decoding json response from %s failed: %v", request.URL, err)
	}

	return response, nil
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
// From: https://github.com/google/go-github/blob/master/github/github.go
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range u.Query() {
		qs[k] = v
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
