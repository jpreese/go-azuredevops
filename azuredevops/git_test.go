package azuredevops_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/mcdafydd/go-azuredevops/azuredevops"
)

const (
	gitRefsListURL      = "/AZURE_DEVOPS_Project/_apis/git/repositories/vscode/refs/heads"
	gitRefsListResponse = `{
		"count": 6,
		"value": [
		  {
			"name": "refs/heads/develop",
			"objectId": "67cae2b029dff7eb3dc062b49403aaedca5bad8d",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/heads/develop"
		  },
		  {
			"name": "refs/heads/master",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/heads/master"
		  },
		  {
			"name": "refs/heads/npaulk/feature",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/heads/npaulk/feature"
		  },
		  {
			"name": "refs/tags/v1.0",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/tags/v1.0"
		  },
		  {
			"name": "refs/tags/v1.1",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/tags/v1.1"
		  },
		  {
			"name": "refs/tags/v2.0",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/tags/v2.0"
		  }
		]
		}`
	gitRepositoryURL         = "/AZURE_DEVOPS_Project/_apis/git/repositories/vscode"
	gitGetRepositoryResponse = `{
		"serverUrl": "https://dev.azure.com/fabrikam",
		"collection": {
			"id": "e22ddea7-989e-455d-b46a-67e991b04714",
			"name": "fabrikam",
			"url": "https://dev.azure.com/fabrikam/_apis/projectCollections/e22ddea7-989e-455d-b46a-67e991b04714"
		},
		"repository": {
			"id": "2f3d611a-f012-4b39-b157-8db63f380226",
			"name": "FabrikamCloud",
			"url": "https://dev.azure.com/fabrikam/_apis/git/repositories/2f3d611a-f012-4b39-b157-8db63f380226",
			"project": {
				"id": "3b3ae425-0079-421f-9101-bcf15d6df041",
				"name": "FabrikamCloud",
				"url": "https://dev.azure.com/fabrikam/_apis/projects/3b3ae425-0079-421f-9101-bcf15d6df041",
				"state": 1,
				"revision": 411518573
			},
			"remoteUrl": "https://dev.azure.com/fabrikam/FabrikamCloud/_git/FabrikamCloud"
		}
	}`
	gitCreateStatusURL      = "/AZURE_DEVOPS_Project/_apis/git/repositories/vscode/commits/67cae2b029dff7eb3dc062b49403aaedca5bad8d/statuses"
	gitCreateStatusResponse = `{
		"state": "succeeded",
		"description": "The build is successful",
		"context": {
			"name": "Build123",
			"genre": "continuous-integration"
		},
		"creationDate": "2016-01-27T09:33:07Z",
		"createdBy": {
			"id": "278d5cd2-584d-4b63-824a-2ba458937249",
			"displayName": "Norman Paulk",
			"uniqueName": "Fabrikamfiber16",
			"url": "https://dev.azure.com/fabrikam/_apis/Identities/278d5cd2-584d-4b63-824a-2ba458937249",
			"imageUrl": "https://dev.azure.com/fabrikam/_api/_common/identityImage?id=278d5cd2-584d-4b63-824a-2ba458937249"
		},
		"targetUrl": "https://ci.fabrikam.com/my-project/build/123 "
	}`
)

func TestGitService_ListRefs(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		response string
		count    int
		index    int
		refName  string
		refID    string
	}{
		{name: "return 6 refs", URL: gitRefsListURL, response: gitRefsListResponse, count: 6, index: 0, refName: "refs/heads/develop", refID: "67cae2b029dff7eb3dc062b49403aaedca5bad8d"},
		{name: "can handle no refs returned", URL: gitRefsListURL, response: "{}", count: 0, index: -1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(tc.URL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := tc.response
				fmt.Fprint(w, json)
			})

			opts := azuredevops.GitRefListOptions{}
			refs, count, err := c.Git.ListRefs("vscode", "heads", &opts)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if *refs[tc.index].Name != tc.refName {
					t.Fatalf("expected git ref name %s, got %s", tc.refName, *refs[tc.index].Name)
				}
				if *refs[tc.index].ObjectID != tc.refID {
					t.Fatalf("expected git ref object id %s, got %s", tc.refID, *refs[tc.index].ObjectID)
				}
			}

			if len(refs) != tc.count {
				t.Fatalf("expected length of git refs to be %d; got %d", tc.count, len(refs))
			}

			if count != tc.count {
				t.Fatalf("expected git ref count to be %d; got %d", tc.count, count)
			}
		})
	}
}

func TestGitService_Get(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		response string
		count    int
		repoName string
		id       string
	}{
		{name: "GetRepository() success", URL: gitRepositoryURL, response: gitGetRepositoryResponse, count: 1, repoName: "vscode", id: "2f3d611a-f012-4b39-b157-8db63f380226"},
		{name: "GetRepository() empty response", URL: gitRepositoryURL, response: "{}", count: 0, repoName: "", id: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(tc.URL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := tc.response
				fmt.Fprint(w, json)
			})

			resp, count, err := c.Git.GetRepository("vscode")
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if count > 0 {
				want := &azuredevops.GitRepository{}

				if !reflect.DeepEqual(resp, want) {
					t.Errorf("Repositories.Get returned %+v, want %+v", resp, want)
				}
			}

		})
	}
}

func TestGitService_CreateStatus(t *testing.T) {
	n := "Build123"
	g := "continuous-integration"
	context := azuredevops.GitStatusContext{
		Name:  &n,
		Genre: &g,
	}

	tt := []struct {
		name        string
		URL         string
		response    string
		count       int
		description string
		targetUrl   string
		state       azuredevops.GitStatusState
		context     *azuredevops.GitStatusContext
	}{
		{name: "CreateStatus() success", URL: gitCreateStatusURL, response: gitCreateStatusResponse, count: 1, description: "some", targetUrl: "", state: azuredevops.GitSucceeded, context: &context},
		{name: "CreateStatus() failed", URL: gitCreateStatusURL, response: "{}", count: 0, description: "some", targetUrl: "", state: azuredevops.GitSucceeded, context: &context},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(tc.URL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := tc.response
				fmt.Fprint(w, json)
			})

			// Build the example request payload
			// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/statuses/create?view#examples
			s := "The build is successful"
			state := azuredevops.GitSucceeded
			target := "https://ci.fabrikam.com/my-project/build/123"
			status := azuredevops.GitStatus{
				Context:     &context,
				Description: &s,
				State:       &state,
				TargetURL:   &target,
			}
			resp, count, err := c.Git.CreateStatus("repo", "67cae2b029dff7eb3dc062b49403aaedca5bad8d", status)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if count > 0 {
				if *resp.Description != tc.description {
					t.Fatalf("expected git ref name %s, got %s", tc.description, *resp.Description)
				}
				if !reflect.DeepEqual(resp.Context, tc.context) {
					t.Errorf("Git.GetRef returned %+v, want %+v", tc.context, resp.Context)
				}
				if *resp.State != tc.state {
					t.Fatalf("expected git ref name %s, got %s", tc.state, *resp.State)
				}
				if *resp.TargetURL != tc.targetUrl {
					t.Fatalf("expected git ref object id %s, got %s", tc.targetUrl, *resp.TargetURL)
				}
			}
		})
	}
}
