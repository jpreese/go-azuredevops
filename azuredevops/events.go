// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted for Azure Devops
package azuredevops

import (
	"encoding/json"
	"errors"
)

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
		payload = &PushResource{}
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
