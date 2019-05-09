package azuredevops_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/mcdafydd/go-azuredevops/azuredevops"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/testing"
)

// Pulled from https://github.com/google/go-github/blob/master/github/github_test.go
func setup() (client *azuredevops.Client, mux *http.ServeMux, serverURL string, teardown func()) {

	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// The client being tested and is configured to use test server.
	client, err := azuredevops.NewClient(nil)

	if err != nil {
		fmt.Errorf("Error requesting NewClient(): %v", err)
	}

	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = *url
	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

func testURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL; got.String() != want {
		t.Errorf("request URL is %s, want %s", got, want)
	}
}

func Test_NewClient(t *testing.T) {
	baseURL, _ := url.Parse(azuredevops.DefaultBaseURL)

	got, _ := azuredevops.NewClient(nil)
	want := azuredevops.Client{
		BaseURL:          *baseURL,
		UserAgent:        azuredevops.UserAgent,
		Account:          "",
		Project:          "",
		AuthToken:        "",
		Boards:           &azuredevops.BoardsService{},
		BuildDefinitions: &azuredevops.BuildDefinitionsService{},
		Builds:           &azuredevops.BuildsService{},
		DeliveryPlans:    &azuredevops.DeliveryPlansService{},
		Favourites:       &azuredevops.FavouritesService{},
		Git:              &azuredevops.GitService{},
		Iterations:       &azuredevops.IterationsService{},
		PullRequests:     &azuredevops.PullRequestsService{},
		Teams:            &azuredevops.TeamsService{},
		Tests:            &azuredevops.TestsService{},
		WorkItems:        &azuredevops.WorkItemsService{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatal("NewClient(): got = %#v, want = %#v", got, want)
	}
	return
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
