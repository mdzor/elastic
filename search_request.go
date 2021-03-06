package elastic

import (
	"strings"
)

// SearchRequest combines a search request and its
// query details (see SearchSource).
// It is used in combination with MultiSearch.
type SearchRequest struct {
	searchType string // default in ES is "query_then_fetch"
	indices    []string
	types      []string
	routing    *string
	preference *string
	source     interface{}
}

// NewSearchRequest creates a new search request.
func NewSearchRequest() *SearchRequest {
	return &SearchRequest{
		indices: make([]string, 0),
		types:   make([]string, 0),
	}
}

// SearchRequest must be one of "query_then_fetch", "query_and_fetch",
// "scan", "count", "dfs_query_then_fetch", or "dfs_query_and_fetch".
// Use one of the constants defined via SearchType.
func (r *SearchRequest) SearchType(searchType string) *SearchRequest {
	r.searchType = searchType
	return r
}

func (r *SearchRequest) SearchTypeDfsQueryThenFetch() *SearchRequest {
	return r.SearchType("dfs_query_then_fetch")
}

func (r *SearchRequest) SearchTypeDfsQueryAndFetch() *SearchRequest {
	return r.SearchType("dfs_query_and_fetch")
}

func (r *SearchRequest) SearchTypeQueryThenFetch() *SearchRequest {
	return r.SearchType("query_then_fetch")
}

func (r *SearchRequest) SearchTypeQueryAndFetch() *SearchRequest {
	return r.SearchType("query_and_fetch")
}

func (r *SearchRequest) SearchTypeScan() *SearchRequest {
	return r.SearchType("scan")
}

func (r *SearchRequest) SearchTypeCount() *SearchRequest {
	return r.SearchType("count")
}

func (r *SearchRequest) Index(index string) *SearchRequest {
	r.indices = append(r.indices, index)
	return r
}

func (r *SearchRequest) Indices(indices ...string) *SearchRequest {
	r.indices = append(r.indices, indices...)
	return r
}

func (r *SearchRequest) HasIndices() bool {
	return len(r.indices) > 0
}

func (r *SearchRequest) Type(typ string) *SearchRequest {
	r.types = append(r.types, typ)
	return r
}

func (r *SearchRequest) Types(types ...string) *SearchRequest {
	r.types = append(r.types, types...)
	return r
}

func (r *SearchRequest) Routing(routing string) *SearchRequest {
	r.routing = &routing
	return r
}

func (r *SearchRequest) Routings(routings ...string) *SearchRequest {
	if routings != nil {
		routings := strings.Join(routings, ",")
		r.routing = &routings
	} else {
		r.routing = nil
	}
	return r
}

func (r *SearchRequest) Preference(preference string) *SearchRequest {
	r.preference = &preference
	return r
}

func (r *SearchRequest) Source(source interface{}) *SearchRequest {
	switch v := source.(type) {
	case *SearchSource:
		r.source = v.Source()
	default:
		r.source = source
	}
	return r
}

// header is used by MultiSearch to get information about the search header
// of one SearchRequest.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-multi-search.html
func (r *SearchRequest) header() interface{} {
	h := make(map[string]interface{})
	if r.searchType != "" {
		h["search_type"] = r.searchType
	}

	switch len(r.indices) {
	case 0:
	case 1:
		h["index"] = r.indices[0]
	default:
		h["indices"] = r.indices
	}

	switch len(r.types) {
	case 0:
	case 1:
		h["types"] = r.types[0]
	default:
		h["type"] = r.types
	}

	if r.routing != nil && *r.routing != "" {
		h["routing"] = *r.routing
	}

	if r.preference != nil && *r.preference != "" {
		h["preference"] = *r.preference
	}

	return h
}

// bidy is used by MultiSearch to get information about the search body
// of one SearchRequest.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-multi-search.html
func (r *SearchRequest) body() interface{} {
	return r.source
}
