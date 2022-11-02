package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// IncludeVersions contains operations available on IncludeVersion resource
	IncludeVersions interface {
		// CreateIncludeVersion creates a new include version based on any previous version
		CreateIncludeVersion(context.Context, CreateIncludeVersionRequest) (*CreateIncludeVersionResponse, error)

		// GetIncludeVersion polls the state of a specific include version, for example to check its activation status
		GetIncludeVersion(context.Context, GetIncludeVersionRequest) (*IncludeVersionResponse, error)

		// ListIncludeVersions lists the include versions, with results limited to the 500 most recent versions
		ListIncludeVersions(context.Context, ListIncludeVersionsRequest) (*IncludeVersionResponse, error)

		// ListIncludeVersionAvailableCriteria lists available criteria for the include version
		ListIncludeVersionAvailableCriteria(context.Context, ListAvailableCriteriaRequest) (*AvailableCriteriaResponse, error)

		// ListIncludeVersionAvailableBehaviors lists available behaviors for the include version
		ListIncludeVersionAvailableBehaviors(context.Context, ListAvailableBehaviorsRequest) (*AvailableBehaviorsResponse, error)
	}

	// CreateIncludeVersionRequest contains parameters used to create a new include version
	CreateIncludeVersionRequest struct {
		IncludeID string
		IncludeVersionRequest
	}

	// IncludeVersionRequest contains body parameters used to create a new include version
	IncludeVersionRequest struct {
		CreateFromVersion     int    `json:"createFromVersion"`
		CreateFromVersionEtag string `json:"createFromVersionEtag,omitempty"`
	}

	// CreateIncludeVersionResponse represents a response object returned by CreateIncludeVersion
	CreateIncludeVersionResponse struct {
		VersionLink string
	}

	// GetIncludeVersionRequest contains parameters used to get the include version
	GetIncludeVersionRequest struct {
		IncludeID  string
		Version    int
		ContractID string
		GroupID    string
	}

	// ListIncludeVersionsRequest contains parameters used to list the include versions
	ListIncludeVersionsRequest struct {
		IncludeID  string
		ContractID string
		GroupID    string
	}

	// IncludeVersionResponse represents a response object returned by GetIncludeVersion
	IncludeVersionResponse struct {
		IncludeID       string      `json:"includeId"`
		IncludeName     string      `json:"includeName"`
		AccountID       string      `json:"accountId"`
		ContractID      string      `json:"contractId"`
		GroupID         string      `json:"groupId"`
		AssetID         string      `json:"assetId"`
		IncludeType     IncludeType `json:"includeType"`
		IncludeVersions Versions    `json:"versions"`
	}

	// Versions represents IncludeVersions object
	Versions struct {
		Items []IncludeVersion `json:"items"`
	}

	// IncludeVersion represents an include version object
	IncludeVersion struct {
		UpdatedByUser    string     `json:"updatedByUser"`
		UpdatedDate      string     `json:"updatedDate"`
		ProductionStatus StatusType `json:"productionStatus"`
		Etag             string     `json:"etag"`
		ProductID        string     `json:"productId"`
		Note             string     `json:"note,omitempty"`
		RuleFormat       string     `json:"ruleFormat,omitempty"`
		IncludeVersion   int        `json:"includeVersion"`
		StagingStatus    StatusType `json:"stagingStatus"`
	}

	// StatusType is type of staging status, whether the include version has been activated to the test network
	StatusType string

	// ListAvailableCriteriaRequest contains parameters used to get available include version criteria
	ListAvailableCriteriaRequest struct {
		IncludeID string
		Version   int
	}

	// AvailableCriteriaResponse represents a response object returned by ListIncludeVersionAvailableCriteria
	AvailableCriteriaResponse struct {
		ContractID        string            `json:"contractId"`
		GroupID           string            `json:"groupId"`
		ProductID         string            `json:"productId"`
		RuleFormat        string            `json:"ruleFormat"`
		AvailableCriteria AvailableCriteria `json:"criteria"`
	}

	// AvailableCriteria represents list of available criteria for the include version
	AvailableCriteria struct {
		Items []Criteria `json:"items"`
	}

	// Criteria represents available criteria object
	Criteria struct {
		Name       string `json:"name"`
		SchemaLink string `json:"schemaLink"`
	}

	// ListAvailableBehaviorsRequest contains parameters used to get available include version behaviors
	ListAvailableBehaviorsRequest struct {
		IncludeID string
		Version   int
	}

	// AvailableBehaviorsResponse represents a response object returned by GetIncludeVersionAvailableBehavior
	AvailableBehaviorsResponse struct {
		ContractID         string             `json:"contractId"`
		GroupID            string             `json:"groupId"`
		ProductID          string             `json:"productId"`
		RuleFormat         string             `json:"ruleFormat"`
		AvailableBehaviors AvailableBehaviors `json:"behaviors"`
	}

	// AvailableBehaviors represents list of available behaviors for the include version
	AvailableBehaviors struct {
		Items []Behavior `json:"items"`
	}

	// Behavior represents available behavior object
	Behavior struct {
		Name       string `json:"name"`
		SchemaLink string `json:"schemaLink"`
	}
)

const (
	// StagingStatusTypeActive indicates that the include version is read-only
	StagingStatusTypeActive StatusType = "ACTIVE"
	// StagingStatusTypeInactive indicates that the include version is inactive
	StagingStatusTypeInactive StatusType = "INACTIVE"
	// StagingStatusTypePending indicates that the include version is pending
	StagingStatusTypePending StatusType = "PENDING"
	// StagingStatusTypeDeactivated indicates that the include version is deactivated
	StagingStatusTypeDeactivated StatusType = "DEACTIVATED"
)

// Validate validates CreateIncludeVersionRequest
func (i CreateIncludeVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":         validation.Validate(i.IncludeID, validation.Required),
		"CreateFromVersion": validation.Validate(i.CreateFromVersion, validation.Required),
	})
}

// Validate validates GetIncludeVersionRequest
func (i GetIncludeVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":  validation.Validate(i.IncludeID, validation.Required),
		"Version":    validation.Validate(i.Version, validation.Required),
		"ContractID": validation.Validate(i.ContractID, validation.Required),
		"GroupID":    validation.Validate(i.GroupID, validation.Required),
	})
}

// Validate validates ListIncludeVersionsRequest
func (i ListIncludeVersionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":  validation.Validate(i.IncludeID, validation.Required),
		"ContractID": validation.Validate(i.ContractID, validation.Required),
		"GroupID":    validation.Validate(i.GroupID, validation.Required),
	})
}

// Validate validates ListAvailableCriteriaRequest
func (i ListAvailableCriteriaRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID": validation.Validate(i.IncludeID, validation.Required),
		"Version":   validation.Validate(i.Version, validation.Required),
	})
}

// Validate validates ListAvailableBehaviorsRequest
func (i ListAvailableBehaviorsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID": validation.Validate(i.IncludeID, validation.Required),
		"Version":   validation.Validate(i.Version, validation.Required),
	})
}

var (
	// ErrCreateIncludeVersion is returned in case an error occurs on CreateIncludeVersion operation
	ErrCreateIncludeVersion = errors.New("create an include version")
	// ErrGetIncludeVersion is returned in case an error occurs on GetIncludeVersion operation
	ErrGetIncludeVersion = errors.New("get an include version")
	// ErrListIncludeVersions is returned in case an error occurs on ListIncludeVersions operation
	ErrListIncludeVersions = errors.New("list include versions")
	// ErrListIncludeVersionAvailableCriteria is returned in case an error occurs on ListIncludeVersionAvailableCriteria operation
	ErrListIncludeVersionAvailableCriteria = errors.New("list include version available criteria")
	// ErrListIncludeVersionAvailableBehaviors is returned in case an error occurs on ListIncludeVersionAvailableBehaviors operation
	ErrListIncludeVersionAvailableBehaviors = errors.New("list include version available behaviors")
)

func (p *papi) CreateIncludeVersion(ctx context.Context, params CreateIncludeVersionRequest) (*CreateIncludeVersionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateIncludeVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateIncludeVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/versions", params.IncludeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateIncludeVersion, err)
	}

	var result CreateIncludeVersionResponse
	resp, err := p.Exec(req, &result, params.IncludeVersionRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateIncludeVersion, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateIncludeVersion, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) GetIncludeVersion(ctx context.Context, params GetIncludeVersionRequest) (*IncludeVersionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetIncludeVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetIncludeVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/versions/%d", params.IncludeID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetIncludeVersion, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetIncludeVersion, err)
	}

	var result IncludeVersionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetIncludeVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetIncludeVersion, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) ListIncludeVersions(ctx context.Context, params ListIncludeVersionsRequest) (*IncludeVersionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListIncludeVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListIncludeVersions, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/versions", params.IncludeID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListIncludeVersions, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListIncludeVersions, err)
	}

	var result IncludeVersionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListIncludeVersions, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListIncludeVersions, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) ListIncludeVersionAvailableCriteria(ctx context.Context, params ListAvailableCriteriaRequest) (*AvailableCriteriaResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListIncludeVersionAvailableCriteria")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListIncludeVersionAvailableCriteria, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/versions/%d/available-criteria", params.IncludeID, params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListIncludeVersionAvailableCriteria, err)
	}

	var result AvailableCriteriaResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListIncludeVersionAvailableCriteria, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListIncludeVersionAvailableCriteria, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) ListIncludeVersionAvailableBehaviors(ctx context.Context, params ListAvailableBehaviorsRequest) (*AvailableBehaviorsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListIncludeVersionAvailableBehaviors")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListIncludeVersionAvailableBehaviors, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/versions/%d/available-behaviors", params.IncludeID, params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListIncludeVersionAvailableBehaviors, err)
	}

	var result AvailableBehaviorsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListIncludeVersionAvailableBehaviors, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListIncludeVersionAvailableBehaviors, p.Error(resp))
	}

	return &result, nil
}