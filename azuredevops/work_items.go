package azuredevops

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// WorkItemsService handles communication with the work items methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/wit/work%20items
type WorkItemsService struct {
	client *Client
}

// IterationWorkItems Represents work items in an iteration backlog
type IterationWorkItems struct {
	Links             *map[string]Link `json:"_links,omitempty"`
	WorkItemRelations []*WorkItemLink  `json:"workItemRelations"`
	URL               *string          `json:"url,omitempty"`
}

// WorkItemComment Describes a response to CreateComment()
type WorkItemComment struct {
	CreatedBy    *IdentityRef `json:"createdBy,omitempty"`
	CreatedDate  *time.Time   `json:"createdDate,omitempty"`
	ID           *int         `json:"id,omitempty"`
	ModifiedBy   *IdentityRef `json:"modifiedBy,omitempty"`
	ModifiedDate *time.Time   `json:"modifiedDate,omitempty"`
	Text         *string      `json:"text,omitempty"`
	URL          *string      `json:"url,omitempty"`
	Version      *int         `json:"version,omitempty"`
	WorkItemID   *int         `json:"workItemId,omitempty"`
}

// WorkItemLink A link between two work items.
type WorkItemLink struct {
	Rel    *string            `json:"rel,omitempty"`
	Source *WorkItemReference `json:"source,omitempty"`
	Target *WorkItemReference `json:"target,omitempty"`
}

// WorkItemListResponse describes the list response for work items
type WorkItemListResponse struct {
	Count     int         `json:"count,omitempty"`
	WorkItems []*WorkItem `json:"value,omitempty"`
}

// WorkItem describes an individual work item in TFS
type WorkItem struct {
	Links             *map[string]Link        `json:"_links,omitempty"`
	CommentVersionRef *CommentVersionRef      `json:"commentVersionRef,omitempty"`
	Fields            *map[string]interface{} `json:"fields,omitempty"`
	ID                *int                    `json:"id,omitempty"`
	Relations         []*WorkItemRelation     `json:"relations,omitempty"`
	Rev               *int                    `json:"rev,omitempty"`
	URL               *string                 `json:"url,omitempty"`
}

// WorkItemFieldUpdate Describes an update to a work item field.
type WorkItemFieldUpdate struct {
	NewValue interface{} `json:"newValue,omitempty"`
	OldValue interface{} `json:"oldValue,omitempty"`
}

// WorkItemRelationUpdates Describes updates to a work item's relations.
type WorkItemRelationUpdates struct {
	Added   []*WorkItemRelation `json:"added,omitempty"`
	Removed []*WorkItemRelation `json:"removed,omitempty"`
	Updated []*WorkItemRelation `json:"updated,omitempty"`
}

// CommentVersionRef refers to the specific version of a comment
type CommentVersionRef struct {
	CommentID *int    `json:"commentId,omitempty"`
	Version   *int    `json:"version,omitempty"`
	URL       *string `json:"url,omitempty"`
}

// WorkItemReference Contains reference to a work item.
type WorkItemReference struct {
	ID  *int    `json:"id,omitempty,string"`
	URL *string `json:"url,omitempty"`
}

// WorkItemRelation describes an intermediary between iterations and work items
type WorkItemRelation struct {
	Attributes *map[string]interface{} `json:"attributes,omitempty"`
	Rel        *string                 `json:"rel,omitempty"`
	URL        *string                 `json:"url,omitempty"`
}

// WorkItemUpdate Describes an update to a work item.
type WorkItemUpdate struct {
	Links       *map[string]interface{}         `json:"attributes,omitempty"`
	Fields      *map[string]WorkItemFieldUpdate `json:"fields,omitempty"`
	ID          *int                            `json:"id,omitempty"`
	Relations   *WorkItemRelationUpdates        `json:"relations,omitempty"`
	Rev         *int                            `json:"rev,omitempty"`
	RevisedBy   *IdentityRef                    `json:"revisedBy,omitempty"`
	RevisedDate *time.Time                      `json:"revisedDate,omitempty"`
	WorkItemID  *int                            `json:"workItemId,omitempty"`
	URL         *string                         `json:"url,omitempty"`
}

// GetForIteration will get a list of work items based on an iteration name
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/wit/work%20items/list
func (s *WorkItemsService) GetForIteration(ctx context.Context, team string, iteration Iteration) ([]*WorkItem, error) {
	queryIds, err := s.GetIdsForIteration(ctx, team, iteration)
	if err != nil {
		return nil, err
	}

	var workIds []string
	for index := 0; index < len(queryIds); index++ {
		workIds = append(workIds, strconv.Itoa(queryIds[index]))
	}

	// https://docs.microsoft.com/en-us/rest/api/vsts/wit/work%20item%20types%20field/list
	fields := []string{
		"System.Id", "System.Title", "System.State", "System.WorkItemType",
		"Microsoft.VSTS.Scheduling.StoryPoints", "System.BoardColumn",
		"System.CreatedBy", "System.AssignedTo", "System.Tags",
	}

	// Now we want to pad out the fields for the work items
	URL := fmt.Sprintf(
		"_apis/wit/workitems?ids=%s&fields=%s&api-version=5.1-preview.1",
		strings.Join(workIds, ","),
		strings.Join(fields, ","),
	)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	var response WorkItemListResponse
	_, err = s.client.Execute(ctx, request, &response)

	return response.WorkItems, err
}

// GetIdsForIteration will return an array of ids for a given iteration
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/iterations/get%20iteration%20work%20items
func (s *WorkItemsService) GetIdsForIteration(ctx context.Context, team string, iteration Iteration) ([]int, error) {
	URL := fmt.Sprintf(
		"%s/_apis/work/teamsettings/iterations/%s/workitems?api-version=5.1-preview.1",
		url.PathEscape(team),
		iteration.ID,
	)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	var response IterationWorkItems

	_, err = s.client.Execute(ctx, request, &response)

	var queryIds []int
	for index := 0; index < len(response.WorkItemRelations); index++ {
		relationship := (response.WorkItemRelations)[index]
		queryIds = append(queryIds, *relationship.Target.ID)
	}

	return queryIds, err
}

// CreateComment Posts a comment to a work item
// https://docs.microsoft.com/en-us/rest/api/azure/devops/wit/comments/add
func (s *WorkItemsService) CreateComment(ctx context.Context, workItemID int, comment *WorkItemComment) (*WorkItemComment, *WorkItemComment, error) {
	URL := fmt.Sprintf(
		"_apis/wit/workItems/%d/comments?api-version=5.1-preview.3",
		workItemID,
	)

	response := WorkItemComment{}

	request, err := s.client.NewRequest("POST", URL, &response)
	if err != nil {
		return nil, &response, err
	}

	_, err = s.client.Execute(ctx, request, &response)

	return comment, &response, nil
}
