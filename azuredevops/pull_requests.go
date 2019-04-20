package azuredevops

import (
	"fmt"
	"time"
)

// CommentType enum declaration
type CommentType int

// CommentType enum declaration
const (
	Unknown CommentType = iota
	Text
	CodeChange
	System
)

func (d CommentType) String() string {
	return [...]string{"unknown", "text", "codechange", "system"}[d]
}

// CommentType enum declaration
type CommentThreadStatus int

// CommentType enum declaration
const (
	StatusUnknown CommentThreadStatus = iota
	StatusActive
	Fixed
	WontFix
	Closed
	ByDesign
	Pending
)

func (d CommentThreadStatus) String() string {
	return [...]string{"unknown", "active", "fixed", "wontfix", "closed", "byDesign", "pending"}[d]
}

// PullRequestsService handles communication with the pull requests methods on the API
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests
type PullRequestsService struct {
	client *Client
}

// PullRequestsResponse describes the pull requests response
type PullRequestsResponse struct {
	GitPullRequests []GitPullRequest `json:"value"`
	Count           int              `json:"count"`
}

// PullRequestListOptions describes what the request to the API should look like
type PullRequestListOptions struct {
	// https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project#pullrequeststatus
	State string `url:"searchCriteria.status,omitempty"`
}

// List returns list of the pull requests
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project
func (s *PullRequestsService) List(opts *PullRequestListOptions) ([]GitPullRequest, int, error) {
	URL := fmt.Sprintf("/_apis/git/pullrequests?api-version=%s", APIVersion)
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response PullRequestsResponse
	_, err = s.client.Execute(request, &response)

	return response.GitPullRequests, response.Count, err
}

// List returns list of the pull requests
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project
func (s *PullRequestsService) ListOne(pullNum int, opts *PullRequestListOptions) (*GitPullRequest, int, error) {
	URL := fmt.Sprintf("/_apis/git/pullrequests/%d?api-version=%s",
		pullNum,
		APIVersion,
	)
	URL, err := addOptions(URL, opts)

	//opt := &RepositoryListByOrgOptions{"forks", ListOptions{Page: 2}}

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response GitPullRequest
	_, err = s.client.Execute(request, &response)

	return &response, 1, err
}

// Merge Completes a pull request
// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/pull%20requests/update?view=azure-devops-rest-5.1
// pullRequest = EnableAutoCompleteOnAnExistingPullRequest(gitHttpClient, pullRequest, mergeCommitMessage);
func (s *PullRequestsService) Merge(opts *PullRequestListOptions) (*GitPullRequest, int, error) {
	URL := fmt.Sprintf("/_apis/git/pullrequests?api-version=5.1-preview.1")
	URL, err := addOptions(URL, opts)

	var autoCompleteSetBy = IdentityRef{}

	mergeStrategy := GitPullRequestMergeStrategy{
		noFastForward: "true",
		//	rebase:        "false",
		//	rebaseMerge:   "false",
		//	squash:        "true",
	}

	pr := GitPullRequest{
		AutoCompleteSetBy: autoCompleteSetBy,
		CompletionOptions: prOptions,
		Status:            &PullRequestStatus{},
	}

	prOptions := GitPullRequestCompletionOptions{
		BypassPolicy:            false,
		BypassReason:            "",
		DeleteSourceBranch:      false,
		MergeCommitMessage:      "",
		MergeStrategy:           mergeStrategy,
		SquashMerge:             false,
		TransitionWorkItems:     true,
		TriggeredByAutoComplete: false,
	}

	// Now we're ready to make our API call to merge the pull request.
	options := &azuredevops.PullRequestOptions{
		MergeMethod: method,
	}

	var prOpts = GitPullRequestCompletionOptions{
		SquashMerge:        true,
		DeleteSourceBranch: true, // false if prefered otherwise
		MergeCommitMessage: mergeCommitMessages,
	}

	request, err := s.client.NewRequest("PATCH", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response GitPullRequest
	_, err = s.client.Execute(request, &response)

	return response.GitPullRequest, response.Count, err
}

// Comment Represents a comment which is one of potentially many in a comment thread.
type Comment struct {
	Links                  *[]ReferenceLinks `json:"_links,omitempty"`
	Author                 *IdentityRef      `json:"author,omitempty"`
	CommentType            *CommentType      `json:"commentType,omitempty"`
	Content                *string           `json:"content,omitempty"`
	ID                     *int              `json:"id,omitempty"`
	IsDeleted              *bool             `json:"isDeleted,omitempty"`
	LastContentUpdatedDate *time.Time        `json:"lastContentUpdatedDate,omitempty"`
	LastUpdatedDate        *time.Time        `json:"lastUpdatedDate,omitempty"`
	ParentCommentID        *int              `json:"parentCommentId,omitempty"`
	PublishedDate          *time.Time        `json:"publishedDate,omitempty"`
	UsersLiked             *[]IdentityRef    `json:"usersLiked,omitempty"`
}

type CommentPosition struct {
	Line   *int `json:"line,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

// GitPullRequestCommentThread Represents a comment thread of a pull request.
// A thread contains meta data about the file it was left on along with one or
// more comments (an initial comment and the subsequent replies).
type GitPullRequestCommentThread struct {
	Links                    *[]ReferenceLinks                   `json:"_links,omitempty"`
	Comments                 *[]Comment                          `json:"comments,omitempty"`
	ID                       *int                                `json:"id,omitempty"`
	Identities               *[]IdentityRef                      `json:"identities,omitempty"`
	IsDeleted                *bool                               `json:"isDeleted,omitempty"`
	LastUpdatedDate          *time.Time                          `json:"lastUpdatedDate,omitempty"`
	Properties               *[]int                              `json:"properties,omitempty"`
	PublishedDate            *time.Time                          `json:"publishedDate,omitempty"`
	Status                   *CommentThreadStatus                `json:"status,omitempty"`
	PullRequestThreadContext *GitPullRequestCommentThreadContext `json:"pullRequestThreadContext,omitempty"`
}

// GitPullRequestCommentThreadContext Comment thread context contains details about what
// diffs were being viewed at the time of thread creation and whether or not the thread
// has been tracked from that original diff.
type GitPullRequestCommentThreadContext struct {
	FilePath       *string          `json:"filePath,omitempty"`
	LeftFileEnd    *CommentPosition `json:"leftFileEnd,omitempty"`
	LeftFileStart  *CommentPosition `json:"leftFileStart,omitempty"`
	RightFileEnd   *CommentPosition `json:"rightFileEnd,omitempty"`
	RightFileStart *CommentPosition `json:"rightFileStart,omitempty"`
}
