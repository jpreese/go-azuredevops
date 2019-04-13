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

// PullRequestResource represents the resource field in a pull request event webhook
type PullRequestResource struct {
	Repo                  *Repo          `json:"repository,omitempty"`
	PullRequestID         *int           `json:"pullRequestId,omitempty"`
	CodeReviewID          *int           `json:"codeReviewId,omitempty"`
	Status                *string        `json:"status,omitempty"`
	CreatedBy             *CreatedBy     `json:"createdBy,omitempty"`
	CreationDate          *time.Time     `json:"creationDate,omitempty"`
	Title                 *string        `json:"title,omitempty"`
	Description           *string        `json:"description,omitempty"`
	SourceRefName         *string        `json:"sourceRefName,omitempty"`
	TargetRefName         *string        `json:"targetRefName,omitempty"`
	MergeStatus           *string        `json:"mergeStatus,omitempty"`
	IsDraft               *bool          `json:"isDraft,omitempty"`
	MergeID               *string        `json:"mergeId,omitempty"`
	LastMergeSourceCommit *MergeCommit   `json:"lastMergeSourceCommit,omitempty"`
	LastMergeTargetCommit *MergeCommit   `json:"lastMergeTargetCommit,omitempty"`
	LastMergeCommit       *MergeCommit   `json:"lastMergeCommit,omitempty"`
	Reviewers             *[]interface{} `json:"reviewers,omitempty"`
	URL                   *string        `json:"url,omitempty"`
	Links                 *Links         `json:"_links,omitempty"`
	SupportsIterations    *bool          `json:"supportsIterations,omitempty"`
	ArtifactID            *string        `json:"artifactId,omitempty"`
}

type Repo struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	URL       string   `json:"url,omitempty"`
	Project   *Project `json:"project,omitempty"`
	Size      int      `json:"size,omitempty"`
	RemoteURL string   `json:"remoteUrl,omitempty"`
	SSHURL    string   `json:"sshUrl,omitempty"`
}

type Project struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	URL            string    `json:"url,omitempty"`
	State          string    `json:"state,omitempty"`
	Revision       int       `json:"revision,omitempty"`
	Visibility     string    `json:"visibility,omitempty"`
	LastUpdateTime time.Time `json:"lastUpdateTime,omitempty"`
}

type CreatedBy struct {
	DisplayName string `json:"displayName,omitempty"`
	URL         string `json:"url,omitempty"`
	Links       *Links `json:"_links,omitempty"`
	ID          string `json:"id,omitempty"`
	UniqueName  string `json:"uniqueName,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty"`
	Descriptor  string `json:"descriptor,omitempty"`
}

// Links
type Links struct {
	Avatar *Avatar `json:"avatar,omitempty"`
}

// Avatar
type Avatar struct {
	Href string `json:"href,omitempty"`
}

// MergeCommit
type MergeCommit struct {
	CommitID  string     `json:"commitId,omitempty"`
	Author    *Committer `json:"author,omitempty"`
	Committer *Committer `json:"committer,omitempty"`
	Comment   string     `json:"comment,omitempty"`
	URL       string     `json:"url,omitempty"`
}

// Committer
type Committer struct {
	Name  string    `json:"name,omitempty"`
	Email string    `json:"email,omitempty"`
	Date  time.Time `json:"date,omitempty"`
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
}
