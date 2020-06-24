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

var HostgroupEndpointPrefix = fmt.Sprintf("%ss", strings.ToLower("Hostgroup"))

type QueryResponseHostgroup struct {
	QueryResponse
	Results []Hostgroup `json:"results"`
}

func (c *Client) GetHostgroup(ctx context.Context, idOrName string) (*Hostgroup, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetHostgroupByID(ctx, int(id))
	}
	return c.GetHostgroupByName(ctx, idOrName)
}

func (c *Client) GetHostgroupByID(ctx context.Context, id int) (*Hostgroup, error) {
	response := new(Hostgroup)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s/%d", HostgroupEndpointPrefix, id), http.MethodGet, nil, response)
	return response, err
}

func (c *Client) GetHostgroupByName(ctx context.Context, name string) (*Hostgroup, error) {
	response := new(QueryResponseHostgroup)
	err := c.requestSearchHelper(ctx, fmt.Sprintf("/%s", HostgroupEndpointPrefix), http.MethodGet, "name", name, nil, response)
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, fmt.Errorf("Hostgroup not found: %s", name)

	}
	return &response.Results[0], err
}

func (c *Client) ListHostgroup(ctx context.Context) (*QueryResponseHostgroup, error) {
	response := new(QueryResponseHostgroup)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s", HostgroupEndpointPrefix), http.MethodGet, nil, response)
	return response, err
}

func (c *Client) CreateHostgroup(ctx context.Context, createRequest interface{}) (*Hostgroup, error) {
	response := new(Hostgroup)
	err := c.requestHelper(ctx, fmt.Sprintf("/%s", HostgroupEndpointPrefix), http.MethodPost, createRequest, response)
	return response, err
}

func (c *Client) DeleteHostgroup(ctx context.Context, id int) error {
	return c.requestHelper(ctx, fmt.Sprintf("/%s/%d", HostgroupEndpointPrefix, id), http.MethodDelete, nil, nil)
}
