package igdb

// CollectionService handles all the API
// calls for the IGDB Collection endpoint.
type CollectionService service

// Collection contains information on an IGDB
// entry for a particular video game series.
//
// For more information, visit: https://igdb.github.io/api/endpoints/collection/
type Collection struct {
	BaseEntity

	Games []int `json:"games"`
}

// Get returns a single Collection identified by the provided IGDB ID. Provide
// the SetFields functional option if you need to specify which fields to
// retrieve. If the ID does not match any Collections, an error is returned.
func (cs *CollectionService) Get(id int, opts ...FuncOption) (*Collection, error) {
	url, err := cs.client.singleURL(CollectionEndpoint, id, opts...)
	if err != nil {
		return nil, err
	}

	var col []Collection

	err = cs.client.get(url, &col)
	if err != nil {
		return nil, err
	}

	return &col[0], nil
}

// List returns a list of Collections identified by the provided list of IGDB IDs.
// Provide functional options to sort, filter, and paginate the results. Omitting
// IDs will instead retrieve an index of Collections based solely on the provided
// options. Any ID that does not match a Collection is ignored. If none of the IDs
// match a Collection, an error is returned.
func (cs *CollectionService) List(ids []int, opts ...FuncOption) ([]*Collection, error) {
	url, err := cs.client.multiURL(CollectionEndpoint, ids, opts...)
	if err != nil {
		return nil, err
	}

	var col []*Collection

	err = cs.client.get(url, &col)
	if err != nil {
		return nil, err
	}

	return col, nil
}

func (cs *CollectionService) ListPaginated(limit int, opts ...FuncOption) (*Pagination, error) {
	return NewPaginationForEndpoint(cs.client, CollectionEndpoint, limit, opts...)
}

func (cs *CollectionService) GetPaginated(p *Pagination) ([]*Collection, bool, error) {
	var c []*Collection

	moreItems, err := p.Get(&c)
	if err != nil {
		return nil, false, err
	}

	return c, moreItems, nil
}

// Search returns a list of Collections found by searching the IGDB using the provided
// query. Provide functional options to sort, filter, and paginate the results. If
// no Collections are found using the provided query, an error is returned.
func (cs *CollectionService) Search(qry string, opts ...FuncOption) ([]*Collection, error) {
	url, err := cs.client.searchURL(CollectionEndpoint, qry, opts...)
	if err != nil {
		return nil, err
	}

	var col []*Collection

	err = cs.client.get(url, &col)
	if err != nil {
		return nil, err
	}

	return col, nil
}

// Count returns the number of Collections available in the IGDB.
// Provide the SetFilter functional option if you need to filter
// which Collections to count.
func (cs *CollectionService) Count(opts ...FuncOption) (int, error) {
	ct, err := cs.client.getEndpointCount(CollectionEndpoint, opts...)
	if err != nil {
		return 0, err
	}

	return ct, nil
}

// ListFields returns the up-to-date list of fields in an
// IGDB Collection object.
func (cs *CollectionService) ListFields() ([]string, error) {
	fl, err := cs.client.getEndpointFieldList(CollectionEndpoint)
	if err != nil {
		return nil, err
	}

	return fl, nil
}
