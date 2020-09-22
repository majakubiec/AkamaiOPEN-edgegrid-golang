package papi

import (
	"context"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
	"net/http"
)

type (
	// Search contains SearchProperty method used for fetching properties
	// https://developer.akamai.com/api/core_features/property_manager/v1.html#searchgroup
	Search interface {
		// Search earches properties by name, or by the hostname or edge hostname for which it’s currently active
		// https://developer.akamai.com/api/core_features/property_manager/v1.html#postfindbyvalue
		SearchProperties(context.Context, SearchRequest) (*SearchResponse, error)
	}

	// SearchResponse contains response body of POST /search request
	SearchResponse struct {
		Versions SearchItems `json:"versions"`
	}

	// SearchItems contains a list of search results
	SearchItems struct {
		Items []SearchItem `json:"items"`
	}

	// SearchItem contains details of a search result
	SearchItem struct {
		AccountID        string `json:"accountId"`
		AssetID          string `json:"assetId"`
		ContractID       string `json:"contractId"`
		EdgeHostname     string `json:"edgeHostname"`
		GroupID          string `json:"groupId"`
		Hostname         string `json:"hostname"`
		ProductionStatus string `json:"productionStatus"`
		PropertyID       string `json:"propertyId"`
		PropertyName     string `json:"propertyName"`
		PropertyVersion  int    `json:"propertyVersion"`
		StagingStatus    string `json:"stagingStatus"`
		UpdatedByUser    string `json:"updatedByUser"`
		UpdatedDate      string `json:"updatedDate"`
	}

	// SearchRequest contains key-value pair for search request
	// Key must have one of three values: "edgeHostname", "hostname" or "propertyName"
	SearchRequest struct {
		key   string
		value string
	}
)

const (
	// SearchKeyEdgeHostname search request key
	SearchKeyEdgeHostname = "edgeHostname"
	// SearchKeyHostname search request key
	SearchKeyHostname = "hostname"
	// SearchKeyPropertyName search request key
	SearchKeyPropertyName = "propertyName"
)

// Validate validate SearchRequest struct
func (s SearchRequest) Validate() error {
	return validation.Errors{
		"SearchKey": validation.Validate(s.key,
			validation.Required,
			validation.In(SearchKeyEdgeHostname, SearchKeyHostname, SearchKeyPropertyName)),
		"SearchValue": validation.Validate(s.value, validation.Required),
	}.Filter()
}

func (p *papi) SearchProperties(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("SearchProperties")

	searchURL := "/papi/v1/search/find-by-value"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create SearchProperties request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var search SearchResponse
	resp, err := p.Exec(req, &search, map[string]string{request.key: request.value})
	if err != nil {
		return nil, fmt.Errorf("SearchProperties request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, searchURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &search, nil
}