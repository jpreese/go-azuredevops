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

// FileContentMetadata Describes files referenced by a GitItem
type FileContentMetadata struct {
	ContentType *string `json:"contentType,omitempty"`
	Encoding    *int    `json:"encoding,omitempty"`
	Extension   *string `json:"extension,omitempty"`
	FileName    *string `json:"fileName,omitempty"`
	IsBinary    *bool   `json:"isBinary,omitempty"`
	IsImage     *bool   `json:"isImage,omitempty"`
	VSLink      *string `json:"vsLink,omitempty"`
}

// ItemContent
type ItemContent struct {
	Content     *string          `json:"content,omitempty"`
	ContentType *ItemContentType `json:"contentType,omitempty"`
}

// ItemContentType
type ItemContentType struct {
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
