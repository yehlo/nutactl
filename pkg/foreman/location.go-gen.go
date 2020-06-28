// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// maybe someone find a better solution. For now this is more than a workaround.
//
package foreman

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var LocationEndpointPrefix = fmt.Sprintf("%s", strings.ToLower("Locations"))

type QueryResponseLocation struct {
	QueryResponse
	Results []Location `json:"results"`
}

func (c *Client) GetLocation(ctx context.Context, idOrName string) (*Location, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetLocationByID(ctx, int(id))
	}
	return c.GetLocationByName(ctx, idOrName)
}

func (c *Client) GetLocationByID(ctx context.Context, id int) (*Location, error) {
	response := new(Location)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s/%d", LocationEndpointPrefix, id), http.MethodGet, nil, response)
	return response, err
}

func (c *Client) GetLocationByName(ctx context.Context, name string) (*Location, error) {
	response := new(QueryResponseLocation)
	filter := fmt.Sprintf("%s=\"%s\"", strings.ToLower("Name"), name)
	err := c.requestSearchHelper(ctx, fmt.Sprintf("/%s", LocationEndpointPrefix), http.MethodGet, filter, nil, response)
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, fmt.Errorf("Location not found")

	}
	return &response.Results[0], err
}

func (c *Client) ListLocation(ctx context.Context) (*QueryResponseLocation, error) {
	response := new(QueryResponseLocation)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s", LocationEndpointPrefix), http.MethodGet, nil, response)
	return response, err
}

func (c *Client) SearchLocation(ctx context.Context, filter string) (*QueryResponseLocation, error) {
	response := new(QueryResponseLocation)
	err := c.requestSearchHelper(ctx, fmt.Sprintf("/%s", LocationEndpointPrefix), http.MethodGet, filter, nil, response)
	return response, err
}

func (c *Client) CreateLocation(ctx context.Context, createRequest interface{}) (*Location, error) {
	response := new(Location)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s", LocationEndpointPrefix), http.MethodPost, createRequest, response)
	return response, err
}

func (c *Client) DeleteLocation(ctx context.Context, id int) error {
	return c.requestHelper(ctx, fmt.Sprintf("/%s/%d", LocationEndpointPrefix, id), http.MethodDelete, nil, nil)
}
