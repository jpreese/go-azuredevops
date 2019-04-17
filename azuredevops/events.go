// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted for Azure Devops
package azuredevops

import (
	"encoding/json"
	"errors"
	"time"
)

// Message represents an Azure Devops webhook message property
type Message struct {
	Text     *string `json:"text,omitempty"`
	HTML     *string `json:"html,omitempty"`
	Markdown *string `json:"markdown,omitempty"`
}

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

// ParsePayload parses the event payload. For recognized event types,
// it returns the webhook payload with a parsed struct in the
// Event.Resource field
func (e *Event) ParsePayload() (payload interface{}, err error) {
	switch *e.EventType {
	case "git.pullrequest.created":
		payload = &GitPullRequest{}
	case "git.pullrequest.merged":
		payload = &GitPullRequest{}
	case "git.pullrequest.updated":
		payload = &GitPullRequest{}
	case "git.push":
		payload = &GitPush{}
	case "workitem.commented":
		payload = &WorkItem{}
	case "workitem.updated":
		payload = &WorkItemUpdate{}
	default:
		return payload, errors.New("Unknown EventType in webhook payload")
	}

	err = json.Unmarshal(*e.RawPayload, &payload)
	if err != nil {
		return payload, err
	}
	e.Resource = &payload
	return payload, nil
}
