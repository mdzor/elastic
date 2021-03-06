// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

// IndexResult is the result of indexing a document in Elasticsearch.
type IndexResult struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	Id      string `json:"_id"`
	Version int    `json:"_version"`
	Created bool   `json:"created"`
}

// IndexService adds documents to Elasticsearch.
type IndexService struct {
	client      *Client
	index       string
	_type       string
	id          string
	routing     string
	parent      string
	opType      string
	refresh     *bool
	version     *int64
	versionType string
	timestamp   string
	ttl         string
	timeout     string
	bodyString  string
	bodyJson    interface{}
	pretty      bool
	debug       bool
}

func NewIndexService(client *Client) *IndexService {
	builder := &IndexService{
		client: client,
	}
	return builder
}

func (b *IndexService) Index(name string) *IndexService {
	b.index = name
	return b
}

func (b *IndexService) Type(_type string) *IndexService {
	b._type = _type
	return b
}

func (b *IndexService) Id(id string) *IndexService {
	b.id = id
	return b
}

func (b *IndexService) Routing(routing string) *IndexService {
	b.routing = routing
	return b
}

func (b *IndexService) Parent(parent string) *IndexService {
	b.parent = parent
	return b
}

// OpType is either "create" or "index" (the default).
func (b *IndexService) OpType(opType string) *IndexService {
	b.opType = opType
	return b
}

func (b *IndexService) Refresh(refresh bool) *IndexService {
	b.refresh = &refresh
	return b
}

func (b *IndexService) Version(version int64) *IndexService {
	b.version = &version
	return b
}

// VersionType is either "internal" (default), "external",
// "external_gt", "external_gte", or "force".
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/docs-index_.html#_version_types
// for details.
func (b *IndexService) VersionType(versionType string) *IndexService {
	b.versionType = versionType
	return b
}

func (b *IndexService) Timestamp(timestamp string) *IndexService {
	b.timestamp = timestamp
	return b
}

func (b *IndexService) TTL(ttl string) *IndexService {
	b.ttl = ttl
	return b
}

func (b *IndexService) Timeout(timeout string) *IndexService {
	b.timeout = timeout
	return b
}

func (b *IndexService) BodyString(body string) *IndexService {
	b.bodyString = body
	return b
}

func (b *IndexService) BodyJson(json interface{}) *IndexService {
	b.bodyJson = json
	return b
}

func (b *IndexService) Pretty(pretty bool) *IndexService {
	b.pretty = pretty
	return b
}

func (b *IndexService) Debug(debug bool) *IndexService {
	b.debug = debug
	return b
}

func (b *IndexService) Do() (*IndexResult, error) {
	// Build url
	var urls, method string
	if b.id != "" {
		// Create document with manual id
		method = "PUT"
		urls = "/{index}/{type}/{id}"
	} else {
		// Automatic ID generation
		// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/docs-index_.html#index-creation
		method = "POST"
		urls = "/{index}/{type}/"
	}
	urls, err := uritemplates.Expand(urls, map[string]string{
		"index": b.index,
		"type":  b._type,
		"id":    b.id,
	})
	if err != nil {
		return nil, err
	}

	// Parameters
	params := make(url.Values)
	if b.pretty {
		params.Set("pretty", "true")
	}
	if b.routing != "" {
		params.Set("routing", b.routing)
	}
	if b.parent != "" {
		params.Set("parent", b.parent)
	}
	if b.opType != "" {
		params.Set("op_type", b.opType)
	}
	if b.refresh != nil && *b.refresh {
		params.Set("refresh", "true")
	}
	if b.version != nil {
		params.Set("version", fmt.Sprintf("%d", *b.version))
	}
	if b.versionType != "" {
		params.Set("version_type", b.versionType)
	}
	if b.timestamp != "" {
		params.Set("timestamp", b.timestamp)
	}
	if b.ttl != "" {
		params.Set("ttl", b.ttl)
	}
	if b.timeout != "" {
		params.Set("timeout", b.timeout)
	}

	/*
		routing string
		parent string
		opType string
		refresh *bool
		version *int64
		versionType string
		timestamp string
		ttl string
	*/

	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := b.client.NewRequest(method, urls)
	if err != nil {
		return nil, err
	}

	// Set body
	if b.bodyJson != nil {
		req.SetBodyJson(b.bodyJson)
	} else {
		req.SetBodyString(b.bodyString)
	}

	if b.debug {
		b.client.dumpRequest((*http.Request)(req))
	}

	// Get response
	res, err := b.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if b.debug {
		b.client.dumpResponse(res)
	}

	ret := new(IndexResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}
