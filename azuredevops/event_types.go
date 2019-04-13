// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted to handle a subset of Pull Request webooks from Azure Devops
// Azure Devops Events docs: https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?view=azure-devops

package azuredevops

import (
	"encoding/json"
	"time"
)

// Event - Describes an Azure Devops webhook payload parent
// For now, delay parsing Resource using *json.RawMessage
// until we know EventType
type Event struct {
	SubscriptionID     *string             `json:"subscriptionId,omitempty"`
	NotificationID     *int                `json:"notificationId,omitempty"`
	ID                 *string             `json:"id,omitempty"`
	EventType          *string             `json:"eventType,omitempty"`
	Message            *Message            `json:"message,omitempty"`
	DetailedMessage    *Message            `json:"detailedMessage,omitempty"`
	RawPayload         *json.RawMessage    `json:"resource,omitempty"`
	ResourceVersion    *string             `json:"resourceVersion,omitempty"`
	ResourceContainers *ResourceContainers `json:"resourceContainers,omitempty"`
	CreatedDate        *time.Time          `json:"createdDate,omitempty"`
	Resource           *interface{}
}

// Message represents an Azure Devops webhook message property
type Message struct {
	Text     *string `json:"text,omitempty"`
	HTML     *string `json:"html,omitempty"`
	Markdown *string `json:"markdown,omitempty"`
}

// GitPullRequest represents the resource field in a pull request event webhook
type GitPullRequest struct {
	Links                 *ReferenceLinks                  `json:"_links,omitempty"`
	ArtifactID            *string                          `json:"artifactId,omitempty"`
	AutoCompleteSetBy     *IdentityRef                     `json:"autoCompleteSetBy,omitempty"`
	ClosedBy              *IdentityRef                     `json:"closedBy,omitempty"`
	ClosedDate            *time.Time                       `json:"closedDate,omitempty"`
	CodeReviewID          *int                             `json:"codeReviewId,omitempty"`
	Commits               *[]GitCommitRef                  `json:"commits,omitempty"`
	CompletionOptions     *GitPullRequestCompletionOptions `json:"completionOptions,omitempty"`
	CompletionQueueTime   *time.Time                       `json:"completionQueueTime,	omitempty"`
	CreatedBy             *IdentityRef                     `json:"createdBy,omitempty"`
	CreationDate          *time.Time                       `json:"creationDate,omitempty"`
	Description           *string                          `json:"description,omitempty"`
	CreationDate          *time.Time                       `json:"creationDate,omitempty"`
	ForkSource            *GitForkRef                      `json:"forkSource,omitempty"`
	IsDraft               *bool                            `json:"isDraft,omitempty"`
	Labels                *WebApiTagDefinition             `json:"labels,omitempty"`
	LastMergeCommit       *GitCommitRef                    `json:"lastMergeCommit,omitempty"`
	LastMergeSourceCommit *GitCommitRef                    `json:"lastMergeSourceCommit,omitempty"`
	LastMergeTargetCommit *GitCommitRef                    `json:"lastMergeTargetCommit,omitempty"`
	MergeFailureMessage   *string                          `json:"mergeFailureMessage,omitempty"`
	MergeFailureType      *PullRequestMergeFailureType     `json:"mergeFailureType,omitempty"`
	MergeID               *string                          `json:"mergeId,omitempty"`
	MergeOptions          *GitPullRequestMergeOptions      `json:"mergeOptions,omitempty"`
	MergeStatus           *PullRequestAsyncStatus          `json:"mergeStatus,omitempty"`
	PullRequestID         *int                             `json:"pullRequestId,omitempty"`
	Repository            *GitRepository                   `json:"repository,omitempty"`
	Reviewers             *[]IdentityRefWithVote           `json:"reviewers,omitempty"`
	RemoteURL             *string                          `json:"remoteUrl,omitempty"`
	SourceRefName         *string                          `json:"sourceRefName,omitempty"`
	Status                *PullRequestStatus               `json:"status,omitempty"`
	SupportsIterations    *bool                            `json:"supportsIterations,omitempty"`
	TargetRefName         *string                          `json:"targetRefName,omitempty"`
	Title                 *string                          `json:"title,omitempty"`
	URL                   *string                          `json:"url,omitempty"`
	WorkItemRefs          *[]ResourceRef                   `json:"workItemRefs,omitempty"`
}

// GitForkRef
// Information about a fork ref.
type GitForkRef struct {
	Links          *ReferenceLinks `json:"_links,omitempty"`
	Creator        *IdentityRef    `json:"creator,omitempty"`
	IsLocked       *bool           `json:"isLocked,omitempty"`
	IsLockedBy     *IdentityRef    `json:"isLockedBy,omitempty"`
	Name           *string         `json:"name,omitempty"`
	ObjectID       *string         `json:"objectId,omitempty"`
	PeeledObjectID *string         `json:"peeledObjectId,omitempty"`
	Repository     *GitRepository  `json:"repository,omitempty"`
	Statuses       *[]GitStatus    `json:"statuses,omitempty"`
	URL            *string         `json:"url,omitempty"`
}

// GitPullRequestCompletionOptions
// Preferences about how the pull request should be completed.
type GitPullRequestCompletionOptions struct {
	BypassPolicy            *bool                        `json:"bypassPolicy,omitempty"`
	BypassReason            *string                      `json:"bypassReason,omitempty"`
	DeleteSourceBranch      *bool                        `json:"deleteSourceBranch,omitempty"`
	MergeCommitMessage      *int                         `json:"mergeCommitMessage,omitempty"`
	MergeStrategy           *GitPullRequestMergeStrategy `json:"mergeStrategy,omitempty"`
	SquashMerge             *bool                        `json:"squashMerge,omitempty"`
	TransitionWorkItems     *bool                        `json:"transitionWorkItems,omitempty"`
	TriggeredByAutoComplete *bool                        `json:"triggeredByAutoComplete,omitempty"`
}

// GitPullRequestMergeStrategy
// Specify the strategy used to merge the pull request during completion. If MergeStrategy
// is not set to any value, a no-FF merge will be created if SquashMerge == false. If
// MergeStrategy is not set to any value, the pull request commits will be squash if
// SquashMerge == true. The SquashMerge member is deprecated. It is recommended that you
// explicitly set MergeStrategy in all cases. If an explicit value is provided for
// MergeStrategy, the SquashMerge member will be ignored.
type GitPullRequestMergeStrategy struct {
	NoFastForward *string `json:"noFastForward,omitempty"`
	Rebase        *string `json:"rebase,omitempty"`
	RebaseMerge   *string `json:"rebaseMerge,omitempty"`
	Squash        *string `json:"squash,omitempty"`
}

// PushResource represents the resource field in a code push request event webhook
type PushResource struct {
	ChangeSetID *int       `json:"changeSetID,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Author      *AuthorRef `json:"author,omitempty"`
	CheckedInBy *AuthorRef `json:"checkedInBy,omitempty"`
	CreatedDate *time.Time `json:"createdDate,omitempty"`
	Comment     *string    `json:"comment,omitempty"`
}

type Repo struct {
	ID        *string  `json:"id,omitempty"`
	Name      *string  `json:"name,omitempty"`
	URL       *string  `json:"url,omitempty"`
	Project   *Project `json:"project,omitempty"`
	Size      *int     `json:"size,omitempty"`
	RemoteURL *string  `json:"remoteUrl,omitempty"`
	SSHURL    *string  `json:"sshUrl,omitempty"`
}

type Project struct {
	ID             *string    `json:"id,omitempty"`
	Name           *string    `json:"name,omitempty"`
	Description    *string    `json:"description,omitempty"`
	URL            *string    `json:"url,omitempty"`
	State          *string    `json:"state,omitempty"`
	Revision       *int       `json:"revision,omitempty"`
	Visibility     *string    `json:"visibility,omitempty"`
	LastUpdateTime *time.Time `json:"lastUpdateTime,omitempty"`
}

// IdentityRef
// Note that many of these fields are being deprecated in the 5.1 API
// This struct follows the 5.0 API
type IdentityRef struct {
	Links             *ReferenceLinks `json:"_links,omitempty"`
	Descriptor        *string         `json:"descriptor,omitempty"`
	DirectoryAlias    *string         `json:"directoryAlias,omitempty"`
	DisplayName       *string         `json:"displayName,omitempty"`
	ID                *string         `json:"id,omitempty"`
	ImageURL          *string         `json:"imageUrl,omitempty"`
	Inactive          *bool           `json:"inactive,omitempty"`
	IsAadIdentity     *bool           `json:"isAadIdentity,omitempty"`
	IsContainer       *bool           `json:"isContainer,omitempty"`
	IsDeletedInOrigin *bool           `json:"isDeletedInOrigin,omitempty"`
	ProfileURL        *string         `json:"profileUrl,omitempty"`
	URL               *string         `json:"url,omitempty"`
	UniqueName        *string         `json:"uniqueName,omitempty"`
}

// ReferenceLinks
type ReferenceLinks struct {
	Links *interface{} `json:"links,omitempty"`
}

// Avatar
type Avatar struct {
	Href *string `json:"href,omitempty"`
}

// GitCommitRef
type GitCommitRef struct {
	Links            *ReferenceLinks        `json:"_links,omitempty"`
	CommitID         *string                `json:"commitId,omitempty"`
	Author           *GitUserDate           `json:"author,omitempty"`
	Committer        *GitUserDate           `json:"committer,omitempty"`
	Comment          *string                `json:"comment,omitempty"`
	CommentTruncated *bool                  `json:"commentTruncated,omitempty"`
	URL              *string                `json:"url,omitempty"`
	ChangeCounts     *ChangeCountDictionary `json:"changeCounts,omitempty"`
	Changes          *GitChange             `json:"changes,omitempty"`
	Parents          *[]string              `json:"parents,omitempty"`
	Push             *GitPushRef            `json:"push,omitempty"`
	RemoteURL        *string                `json:"remoteUrl,omitempty"`
	Statuses         *[]GitStatus           `json:"statuses,omitempty"`
	WorkItems        *ResourceRef           `json:"workItems,omitempty"`
}

type GitChange struct {
}

type GitStatus struct {
}

type GitRepository struct {
}

// GitUserDate
type GitUserDate struct {
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`
	Date  *time.Time `json:"date,omitempty"`
}

// ResourceContainers provides information related to the Resources in a payload
type ResourceContainers struct {
	Collection *ResourceRef `json:"text,omitempty"`
	Account    *ResourceRef `json:"html,omitempty"`
	Project    *ResourceRef `json:"markdown,omitempty"`
}

// Describes properties to identify a resource
type ResourceRef struct {
	ID      *string `json:"id,omitempty"`
	BaseURL *string `json:"baseUrl,omitempty"`
	URL     *string `json:"url,omitempty"`
}

// WebApiTagDefinition
// The representation of a tag definition which is sent across the wire.
type WebApiTagDefinition struct {
	Active *bool   `json:"active,omitempty"`
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	URL    *string `json:"url,omitempty"`
}
