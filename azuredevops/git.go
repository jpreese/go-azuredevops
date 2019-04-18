package azuredevops

import (
	"context"
	"fmt"
	"net/url"
)

// GitService handles communication with the git methods on the API
// See: https://docs.microsoft.com/en-us/rest/api/vsts/git/
type GitService struct {
	client *Client
}

// GitRefsResponse describes the git refs list response
type GitRefsResponse struct {
	Count int         `json:"count"`
	Refs  []GitStatus `json:"value"`
}

// GitRef describes what the git reference looks like
type GitRef struct {
	Name     string      `json:"name,omitempty"`
	ObjectID string      `json:"objectId,omitempty"`
	URL      string      `json:"url,omitempty"`
	Statuses []GitStatus `json:"statuses,omitempty"`
}

// GitRefListOptions describes what the request to the API should look like
type GitRefListOptions struct {
	Filter             string `url:"filter,omitempty"`
	IncludeStatuses    bool   `url:"includeStatuses,omitempty"`
	LatestStatusesOnly bool   `url:"latestStatusesOnly,omitempty"`
}

// ListRefs returns a list of the references for a git repo
func (s *GitService) ListRefs(repo, refType string, opts *GitRefListOptions) ([]GitRef, int, error) {
	URL := fmt.Sprintf(
		"/_apis/git/repositories/%s/refs/%s?api-version=5.1-preview.1",
		repo,
		refType,
	)

	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response GitRefsResponse
	_, err = s.client.Execute(request, &response)

	return response.Refs, response.Count, err
}

// CreateStatus creates a new status for a repository at the specified
// reference. Ref can be a SHA, a branch name, or a tag name.
// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/statuses/create?view=azure-devops-rest-5.0
func (s *GitService) CreateStatus(ctx context.Context, owner, repo, ref string, status *GitStatus) (*[]GitStatus, int, error) {
	URL := fmt.Sprintf(
		"/_apis/git/repositories/%s/commits/%s/statuses?api-version=5.1-preview.1",
		url.QueryEscape(ref),
		ref,
	)

	request, err := s.client.NewRequest("POST", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response GitRefsResponse
	_, err = s.client.Execute(request, &response)

	return &response.Refs, response.Count, err
	// return repoStatus, resp, nil
}
