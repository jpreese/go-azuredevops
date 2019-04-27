package azuredevops

import (
	"fmt"
	"net/url"
)

// BoardsService handles communication with the boards methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/boards
type BoardsService struct {
	client *Client
}

// ListBoardsResponse describes the boards response
type ListBoardsResponse struct {
	Count           *int              `json:"count,omitempty"`
	BoardReferences []*BoardReference `json:"value,omitempty"`
}

// Board describes a board
type Board struct {
	BoardReference
	Links           *ReferenceLinks `json:"_links,omitempty"`
	AllowedMappings *string         `json:"allowedMappings,omitempty"`
	CanEdit         *bool           `json:"canEdit,omitempty"`
	Columns         []*BoardColumn  `json:"columns,omitempty"`
	Fields          *BoardFields    `json:"fields,omitempty"`
	IsValid         *bool           `json:"isvalid,omitempty"`
	Revision        *int            `json:"revision,omitempty"`
	Rows            []*BoardRow     `json:"rows,omitempty"`
}

// BoardColumn describes a column on the board
type BoardColumn struct {
	ColumnType    *string           `json:"columnType,omitempty"`
	Description   *string           `json:"description,omitempty"`
	ID            *string           `json:"id,omitempty"`
	IsSplit       *bool             `json:"isSplit,omitempty"`
	ItemLimit     *int              `json:"itemLimit,omitempty"`
	Name          *string           `json:"name,omitempty"`
	StateMappings map[string]string `json:"stateMappings,omitempty"`
}

// BoardColumnType describes a column on the board
type BoardColumnType int

// BoardColumnType Enum values
const (
	Incoming BoardColumnType = iota
	InProgress
	Outgoing
)

func (d BoardColumnType) String() string {
	return [...]string{"Incoming", "InProgress", "Outgoing"}[d]
}

// BoardFields describes a column on the board
type BoardFields struct {
	ColumnField *FieldReference `json:"columnField,omitempty"`
	DoneField   *FieldReference `json:"doneField,omitempty"`
	RowField    *FieldReference `json:"rowField,omitempty"`
}

// BoardReference Base object a board
type BoardReference struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	URL  *string `json:"URL,omitempty"`
}

// BoardRow describes a row on the board
type BoardRow struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// FieldReference describes a row on the board
type FieldReference struct {
	ReferenceName *string `json:"referenceName,omitempty"`
	URL           *string `json:"url,omitempty"`
}

// List returns list of the boards
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/boards/list
func (s *BoardsService) List(team string) ([]*BoardReference, error) {
	URL := fmt.Sprintf(
		"%s/_apis/work/boards?api-version=%s",
		url.PathEscape(team),
		APIVersion,
	)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response ListBoardsResponse
	_, err = s.client.Execute(request, &response)

	return response.BoardReferences, err
}

// Get returns a single board utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/boards/get
func (s *BoardsService) Get(team string, id string) (*Board, error) {
	URL := fmt.Sprintf(
		"%s/_apis/work/boards/%s?api-version=%s",
		url.PathEscape(team),
		id,
		APIVersion,
	)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response Board
	_, err = s.client.Execute(request, &response)

	return &response, err
}
