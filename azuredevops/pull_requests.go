package azuredevops

import (
	"fmt"
	"time"
)

// Vote identifiers
const (
	VoteApproved                = 10
	VoteApprovedWithSuggestions = 5
	VoteNone                    = 0
	VoteWaitingForAuthor        = -5
	VoteRejected                = -10
)

// CommentType enum declaration
type CommentType int

// CommentType enum declaration
const (
	// The comment type is not known.
	CommentTypeUnknown CommentType = iota
	// This is a regular user comment.
	CommentTypeText
	// The comment comes as a result of a code change.
	CommentTypeCodeChange
	// The comment represents a system message.
	CommentTypeSystem
)

func (d CommentType) String() string {
	return [...]string{"unknown", "text", "codechange", "system"}[d]
}

// CommentThreadStatus enum declaration
type CommentThreadStatus int

// CommentThreadStatus enum declaration
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

// PullRequestAsyncStatus The current status of a pull request merge.
type PullRequestAsyncStatus int

// PullRequestAsyncStatus enum values
const (
	MergeConflicts PullRequestAsyncStatus = iota
	MergeFailure
	MergeNotSet
	MergeQueued
	MergeRejectedByPolicy
	MergeSucceeded
)

func (d PullRequestAsyncStatus) String() string {
	return [...]string{"conflicts", "failure", "notSet", "queued", "rejectedByPolicy", "succeeded"}[d]
}

// PullRequestMergeFailureType The specific type of merge request failure
type PullRequestMergeFailureType int

// PullRequestMergeFailureType enum values
const (
	NoFailure PullRequestMergeFailureType = iota
	UnknownFailure
	CaseSensitive
	ObjectTooLarge
)

func (d PullRequestMergeFailureType) String() string {
	return [...]string{"none", "unknown", "caseSensitive", "objectTooLarge"}[d]
}

// PullRequestStatus The current status of a pull request merge.
type PullRequestStatus int

// PullRequestStatus enum values
const (
	PullAbandoned PullRequestStatus = iota
	PullActive
	PullIncludeAll
	PullCompleted
	PullNotSet
)

func (d PullRequestStatus) String() string {
	return [...]string{"abandoned", "active", "all", "completed", "notSet"}[d]
}

// PullRequestsService handles communication with the pull requests methods on the API
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests
type PullRequestsService struct {
	client *Client
}

// PullRequestsCommitsResponse describes a pull requests commits response
type PullRequestsCommitsResponse struct {
	Count         int             `json:"count"`
	GitCommitRefs []*GitCommitRef `json:"value"`
}

// PullRequestsListResponse describes a pull requests list response
type PullRequestsListResponse struct {
	Count           int               `json:"count"`
	GitPullRequests []*GitPullRequest `json:"value"`
}

// PullRequestListOptions describes what the request to the API should look like
type PullRequestListOptions struct {
	// https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project#pullrequeststatus
	State string `url:"searchCriteria.status,omitempty"`
}

// List returns list of the pull requests
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project
func (s *PullRequestsService) List(opts *PullRequestListOptions) ([]*GitPullRequest, int, error) {
	URL := fmt.Sprintf("_apis/git/pullrequests?api-version=%s", APIVersion)
	URL, err := addOptions(URL, opts)

	req, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var resp PullRequestsListResponse
	_, err = s.client.Execute(req, &resp)

	return resp.GitPullRequests, resp.Count, err
}

// Get returns a single pull request
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project
func (s *PullRequestsService) Get(pullNum int, opts *PullRequestListOptions) (*GitPullRequest, int, error) {
	URL := fmt.Sprintf("_apis/git/pullrequests/%d?api-version=%s",
		pullNum,
		APIVersion,
	)
	URL, err := addOptions(URL, opts)

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
func (s *PullRequestsService) Merge(repoName string, pullNum int, id *IdentityRef, commitMsg string) (*GitPullRequest, int, error) {
	URL := fmt.Sprintf("_apis/git/repositories/%s/pullrequests?api-version=%s",
		repoName,
		APIVersion,
	)

	//	GitPullRequestMergeStrategy{
	//	NoFastForward: "true",
	//	rebase:        "false",
	//	rebaseMerge:   "false",
	//	squash:        "true",
	//}

	// Set default pull request completion options
	empty := ""
	mcm := NoFastForward.String()
	var twi *bool
	*twi = true
	prOptions := GitPullRequestCompletionOptions{
		BypassPolicy:            new(bool),
		BypassReason:            &empty,
		DeleteSourceBranch:      new(bool),
		MergeCommitMessage:      &commitMsg,
		MergeStrategy:           &mcm,
		SquashMerge:             new(bool),
		TransitionWorkItems:     twi,
		TriggeredByAutoComplete: new(bool),
	}

	pr := GitPullRequest{
		AutoCompleteSetBy: id,
		CompletionOptions: &prOptions,
		PullRequestID:     &pullNum,
	}

	// Now we're ready to make our API call to merge the pull request.
	request, err := s.client.NewRequest("PATCH", URL, pr)
	if err != nil {
		return nil, 0, err
	}
	var response GitPullRequest
	_, err = s.client.Execute(request, &response)

	return &response, 1, err
}

// Comment Represents a comment which is one of potentially many in a comment thread.
type Comment struct {
	Links                  *[]ReferenceLinks `json:"_links,omitempty"`
	Author                 *IdentityRef      `json:"author,omitempty"`
	CommentType            *string           `json:"commentType,omitempty"`
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
	Status                   *string                             `json:"status,omitempty"`
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

// ListCommits lists the commits in a pull request.
// Azure Devops API docs: https://docs.microsoft.com/en-us/rest/api/azure/devops/git/pull%20request%20commits/get%20pull%20request%20commits
//
func (s *PullRequestsService) ListCommits(repo string, pullNum int) ([]*GitCommitRef, int, error) {
	URL := fmt.Sprintf("_apis/git/repositories/%s/pullRequests/%d/commits?api-version=%s",
		repo,
		pullNum,
		APIVersion,
	)

	req, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}

	var resp PullRequestsCommitsResponse
	_, err = s.client.Execute(req, resp)

	return resp.GitCommitRefs, resp.Count, nil
}
