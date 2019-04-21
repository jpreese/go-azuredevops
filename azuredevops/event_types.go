// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted to handle a subset of Pull Request webooks from Azure Devops
// Azure Devops Events docs: https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?view=azure-devops

package azuredevops

// ChangeCountDictionary maps the number of changes to each type
type ChangeCountDictionary *map[VersionControlChangeType]int

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
	Links map[string]Link `json:",omitempty"`
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

// WebAPITagDefinition The representation of a tag definition which is sent across
// the wire.
type WebAPITagDefinition struct {
	Active *bool   `json:"active,omitempty"`
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	URL    *string `json:"url,omitempty"`
}
