// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted to handle a subset of Pull Request webooks from Azure Devops
// Azure Devops Events docs: https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?view=azure-devops

package azuredevops

import (
	"time"
)

// VersionControlChangeType enum declaration
type VersionControlChangeType int

// VersionControlChangeType valid enum values
const (
	None VersionControlChangeType = iota
	Add
	Edit
	Encoding
	Rename
	Delete
	Undelete
	Branch
	Merge
	Lock
	Rollback
	SourceRename
	TargetRename
	Property
	All
)

func (d VersionControlChangeType) String() string {
	return [...]string{"none", "add", "edit", "encoding", "rename", "delete", "undelete", "branch", "merge", "lock", "rollback", "sourceRename", "targetRename", "property", "all"}[d]
}

// GitObjectType enum declaration
type GitObjectType int

// GitObjectType enum declaration
const (
	Bad GitObjectType = iota
	Commit
	Tree
	Blob
	Tag
	Ext2
	OfsDelta
	RefDelta
)

func (d GitObjectType) String() string {
	return [...]string{"bad", "commit", "tree", "blob", "tag", "ext2", "ofsDelta", "refDelta"}[d]
}

// ChangeCountDictionary maps the number of changes to each type
type ChangeCountDictionary *map[VersionControlChangeType]int

// GitChange describes file path and content changes
type GitChange struct {
	ChangeID           *int                      `json:"changeId,omitempty"`
	ChangeType         *VersionControlChangeType `json:"changeType,omitempty"`
	Item               *GitItem                  `json:"item,omitempty"`
	NewContent         *ItemContent              `json:"newContent,omitempty"`
	NewContentTemplate *GitTemplate              `json:"newContentTemplate,omitempty"`
	OriginalPath       *string                   `json:"originalPath,omitempty"`
	SourceServerItem   *string                   `json:"sourceServerItem,omitempty"`
	URL                *string                   `json:"url,omitempty"`
}

// GitCommitChanges is a list of GitCommitRefs and count of all changes describes in
// the response from the API
type GitCommitChanges struct {
	ChangeCounts *ChangeCountDictionary `json:"changeCounts,omitempty"`
	Changes      *[]GitChange           `json:"changes,omitempty"`
}

// GitCommitRef describes a single git commit reference
type GitCommitRef struct {
	Links            *[]ReferenceLinks      `json:"_links,omitempty"`
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

// GitForkRef provides information about a fork ref.
type GitForkRef struct {
	Links          *[]ReferenceLinks `json:"_links,omitempty"`
	Creator        *IdentityRef      `json:"creator,omitempty"`
	IsLocked       *bool             `json:"isLocked,omitempty"`
	IsLockedBy     *IdentityRef      `json:"isLockedBy,omitempty"`
	Name           *string           `json:"name,omitempty"`
	ObjectID       *string           `json:"objectId,omitempty"`
	PeeledObjectID *string           `json:"peeledObjectId,omitempty"`
	Repository     *GitRepository    `json:"repository,omitempty"`
	Statuses       *[]GitStatus      `json:"statuses,omitempty"`
	URL            *string           `json:"url,omitempty"`
}

// GitItem describes a single git item
type GitItem struct {
	CommitID              *string        `json:"commitId,omitempty"`
	GitObjectType         *GitObjectType `json:"gitObjectType,omitempty"`
	LatestProcessedChange *GitCommitRef  `json:"latestProcessedChange,omitempty"`
	ObjectID              *string        `json:"objectId,omitempty"`
	OriginalObjectID      *string        `json:"originalObjectId,omitempty"`
}

// GitPullRequest represents all the data associated with a pull request.
type GitPullRequest struct {
	Links                 *[]ReferenceLinks                `json:"_links,omitempty"`
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

// GitPullRequestCompletionOptions describes preferences about how the pull
// request should be completed.
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

// GitPullRequestMergeOptions describes the options which are used when a pull
// request merge is created.
type GitPullRequestMergeOptions struct {
	DetectRenameFalsePositives *bool `json:"detectRenameFalsePositives,omitempty"`
	DisableRenames             *bool `json:"disableRenames,omitempty"`
}

// GitPullRequestMergeStrategy specifies the strategy used to merge the pull request
// during completion. If MergeStrategy is not set to any value, a no-FF merge will be
// created if SquashMerge == false. If MergeStrategy is not set to any value, the pull
// request commits will be squash if SquashMerge == true. The SquashMerge member is
// deprecated. It is recommended that you explicitly set MergeStrategy in all cases.
// If an explicit value is provided for MergeStrategy, the SquashMerge member will
// be ignored.
type GitPullRequestMergeStrategy struct {
	NoFastForward *string `json:"noFastForward,omitempty"`
	Rebase        *string `json:"rebase,omitempty"`
	RebaseMerge   *string `json:"rebaseMerge,omitempty"`
	Squash        *string `json:"squash,omitempty"`
}

// GitPush describes a code push request event.
type GitPush struct {
	Links      *[]ReferenceLinks `json:"_links,omitempty"`
	Commits    *[]GitCommitRef   `json:"commits,omitempty"`
	Date       *time.Time        `json:"date,omitempty"`
	PushID     *int              `json:"pushId,omitempty"`
	PushedBy   *IdentityRef      `json:"pushedBy,omitempty"`
	RefUpdates *[]GitRefUpdate   `json:"refUpdates,omitempty"`
	Repository *GitRepository    `json:"repository,omitempty"`
	URL        *string           `json:"url,omitempty"`
}

// GitRefUpdate
type GitRefUpdate struct {
	IsLocked     *bool   `json:"isLocked,omitempty"`
	Name         *string `json:"name,omitempty"`
	NewObjectID  *string `json:"newObjectId,omitempty"`
	OldObjectID  *string `json:"oldObjectId,omitempty"`
	RepositoryID *string `json:"repositoryId,omitempty"`
}

// GitRepository describes an Azure Devops Git repository.
type GitRepository struct {
	Links            *[]ReferenceLinks `json:"_links,omitempty"`
	DefaultBranch    *string           `json:"defaultBranch,omitempty"`
	ID               *string           `json:"id,omitempty"`
	IsFork           *bool             `json:"isFork,omitempty"`
	Name             *string           `json:"name,omitempty"`
	ParentRepository *GitRepositoryRef `json:"parentRepository,omitempty"`
	Project          *TeamProjectRef   `json:"project,omitempty"`
	RemoteURL        *string           `json:"remoteUrl,omitempty"`
	Size             *int              `json:"size,omitempty"`
	SSHURL           *string           `json:"sshUrl,omitempty"`
	URL              *string           `json:"url,omitempty"`
	ValidRemoteURLs  *[]string         `json:"validRemoteUrls,omitempty"`
}

// GitRepositoryRef
type GitRepositoryRef struct {
	Collection *TeamProjectCollectionReference `json:"collection,omitempty"`
	ID         *string                         `json:"id,omitempty"`
	IsFork     *bool                           `json:"isFork,omitempty"`
	Name       *string                         `json:"name,omitempty"`
	Project    *TeamProjectRef                 `json:"project,omitempty"`
	RemoteURL  *string                         `json:"remoteUrl,omitempty"`
	SSHURL     *string                         `json:"sshUrl,omitempty"`
	URL        *string                         `json:"url,omitempty"`
}

// GitStatus contains the metadata of a service/extension posting a status.
type GitStatus struct {
	Links        *[]ReferenceLinks `json:"_links,omitempty"`
	Context      *GitStatusContext `json:"context,omitempty"`
	CreatedBy    *IdentityRef      `json:"createdBy,omitempty"`
	CreationDate *time.Time        `json:"creationDate,omitempty"`
	Description  *string           `json:"description,omitempty"`
	ID           *int              `json:"id,omitempty"`
	TargetURL    *string           `json:"targetUrl,omitempty"`
	UpdatedDate  *time.Time        `json:"updatedDate,omitempty"`
}

// GitStatusContext Status context that uniquely identifies the status.
type GitStatusContext struct {
	Genre *string `json:"genre,omitempty"`
	Name  *string `json:"name,omitempty"`
}

// GitStatusState State of the status.
type GitStatusState struct {
	Error         *string `json:"error,omitempty"`
	Failed        *string `json:"failed,omitempty"`
	NotApplicable *string `json:"notApplicable,omitempty"`
	NotSet        *string `json:"notSet,omitempty"`
	Pending       *string `json:"pending,omitempty"`
	Succeeded     *string `json:"succeeded,omitempty"`
}

// GitTemplate
type GitTemplate struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

// GitUserDate User info and date for Git operations.
type GitUserDate struct {
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`
	Date  *time.Time `json:"date,omitempty"`
}

// IdentityRef describes an Azure Devops identity
type IdentityRef struct {
	Links             *[]ReferenceLinks `json:"_links,omitempty"`
	Descriptor        *string           `json:"descriptor,omitempty"`
	DirectoryAlias    *string           `json:"directoryAlias,omitempty"`
	DisplayName       *string           `json:"displayName,omitempty"`
	ID                *string           `json:"id,omitempty"`
	ImageURL          *string           `json:"imageUrl,omitempty"`
	Inactive          *bool             `json:"inactive,omitempty"`
	IsAadIdentity     *bool             `json:"isAadIdentity,omitempty"`
	IsContainer       *bool             `json:"isContainer,omitempty"`
	IsDeletedInOrigin *bool             `json:"isDeletedInOrigin,omitempty"`
	ProfileURL        *string           `json:"profileUrl,omitempty"`
	URL               *string           `json:"url,omitempty"`
	UniqueName        *string           `json:"uniqueName,omitempty"`
}

// IdentityRefWithVote Identity information including a vote on a pull request.
type IdentityRefWithVote struct {
	IdentityRef
	IsRequired  *bool                  `json:"isRequired,omitempty"`
	ReviewerURL *string                `json:"reviewerUrl,omitempty"`
	Vote        *int                   `json:"vote,omitempty"`
	VotedFor    *[]IdentityRefWithVote `json:"votedFor,omitempty"`
}

// ItemContent
type ItemContent struct {
	Content     *string          `json:"content,omitempty"`
	ContentType *ItemContentType `json:"contentType,omitempty"`
}

// ItemContentType
type ItemContItemContentTypeent struct {
	Base64Encoded *string `json:"base64Encoded,omitempty"`
	RawText       *string `json:"rawText,omitempty"`
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

// PullRequestAsyncStatus The current status of the pull request merge.
type PullRequestAsyncStatus struct {
	Conflicts        *string `json:"conflicts,omitempty"`
	Failure          *string `json:"failure,omitempty"`
	NotSet           *string `json:"notSet,omitempty"`
	Queued           *string `json:"queued,omitempty"`
	RejectedByPolicy *string `json:"rejectedByPolicy,omitempty"`
	Succeeded        *string `json:"succeeded,omitempty"`
}

// PullRequestMergeFailureType The type of failure (if any) of the pull request
// merge.
type PullRequestMergeFailureType struct {
	CaseSensitive  *string `json:"caseSensitive,omitempty"`
	None           *string `json:"none,omitempty"`
	ObjectTooLarge *string `json:"objectTooLarge,omitempty"`
	Unknown        *string `json:"Unknown,omitempty"`
}

// PullRequestStatus The current status of the pull request merge.
type PullRequestStatus struct {
	Abandoned *string `json:"abandoned,omitempty"`
	Active    *string `json:"active,omitempty"`
	All       *string `json:"all,omitempty"`
	Completed *string `json:"completed,omitempty"`
	NotSet    *string `json:"notSet,omitempty"`
}

// ReferenceLinks The class to represent a collection of REST reference links.
type ReferenceLinks struct {
	Links *map[string]Link `json:",omitempty"`
}

// Link A single item in a collection of ReferenceLinks.
type Link struct {
	Href *string `json:"href,omitempty"`
}

// ResourceContainers provides information related to the Resources in a payload
type ResourceContainers struct {
	Collection *ResourceRef `json:"text,omitempty"`
	Account    *ResourceRef `json:"html,omitempty"`
	Project    *ResourceRef `json:"markdown,omitempty"`
}

// ResourceRef Describes properties to identify a resource
type ResourceRef struct {
	ID      *string `json:"id,omitempty"`
	BaseURL *string `json:"baseUrl,omitempty"`
	URL     *string `json:"url,omitempty"`
}

// TeamProjectCollectionReference Reference object for a TeamProjectCollection.
type TeamProjectCollectionReference struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

// TeamProjectReference Represents a shallow reference to a TeamProject.
type TeamProjectReference struct {
	Abbreviation        *string `json:"abbreviation,omitempty"`
	DefaultTeamImageUrl *string `json:"defaultTeamImageUrl,omitempty"`
	Description         *string `json:"description,omitempty"`
	ID                  *string `json:"id,omitempty"`
	Name                *string `json:"name,omitempty"`
	Revision            *string `json:"revision,omitempty"`
	State               *string `json:"state,omitempty"`
	URL                 *string `json:"url,omitempty"`
	Visibility          *string `json:"visibility,omitempty"`
}

// WebAPITagDefinition The representation of a tag definition which is sent across
// the wire.
type WebAPITagDefinition struct {
	Active *bool   `json:"active,omitempty"`
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	URL    *string `json:"url,omitempty"`
}