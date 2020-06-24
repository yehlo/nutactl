// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package foreman

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var SubnetEndpointPrefix = fmt.Sprintf("%ss", strings.ToLower("Subnet"))

type QueryResponseSubnet struct {
	QueryResponse
	Results []Subnet `json:"results"`
}

func (c *Client) GetSubnet(ctx context.Context, idOrName string) (*Subnet, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetSubnetByID(ctx, int(id))
	}
	return c.GetSubnetByName(ctx, idOrName)
}

func (c *Client) GetSubnetByID(ctx context.Context, id int) (*Subnet, error) {
	response := new(Subnet)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s/%d", SubnetEndpointPrefix, id), http.MethodGet, nil, response)
	return response, err
}

func (c *Client) GetSubnetByName(ctx context.Context, name string) (*Subnet, error) {
	response := new(QueryResponseSubnet)
	err := c.requestSearchHelper(ctx, fmt.Sprintf("/%s", SubnetEndpointPrefix), http.MethodGet, "name", name, nil, response)
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, fmt.Errorf("Subnet not found: %s", name)

	}
	return &response.Results[0], err
}

func (c *Client) ListSubnet(ctx context.Context) (*QueryResponseSubnet, error) {
	response := new(QueryResponseSubnet)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s", SubnetEndpointPrefix), http.MethodGet, nil, response)
	return response, err
}

func (c *Client) CreateSubnet(ctx context.Context, createRequest interface{}) (*Subnet, error) {
	response := new(Subnet)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s", SubnetEndpointPrefix), http.MethodPost, createRequest, response)
	return response, err
}

func (c *Client) DeleteSubnet(ctx context.Context, id int) error {
	return c.requestHelper(ctx, fmt.Sprintf("/%s/%d", SubnetEndpointPrefix, id), http.MethodDelete, nil, nil)
}
