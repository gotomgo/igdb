package igdb

// CompanyService handles all the API
// calls for the IGDB Company endpoint.
type CompanyService service

// Company contains information on an IGDB entry
// for a particular video game company, including
// both publishers and developers.
//
// For more information, visit: https://igdb.github.io/api/endpoints/company/
type Company struct {
	ID                 int          `json:"ID,omitempty"`
	Name               string       `json:"name,omitempty"`
	Slug               string       `json:"slug,omitempty"`
	URL                URL          `json:"url,omitempty"`
	CreatedAt          int          `json:"created_at,omitempty"` // Unix time in milliseconds
	UpdatedAt          int          `json:"updated_at,omitempty"` // Unix time in milliseconds
	Logo               Image        `json:"logo,omitempty"`
	Description        string       `json:"description,omitempty"`
	Country            CountryCode  `json:"country,omitempty"`
	Website            string       `json:"website,omitempty"`
	StartDate          int          `json:"start_date,omitempty"` // Unix time in milliseconds
	StartDateCategory  DateCategory `json:"start_date_category,omitempty"`
	ChangedID          int          `json:"changed_company_id,omitempty"`
	ChangeDate         int          `json:"change_date,omitempty"` // Unix time in milliseconds
	ChangeDateCategory DateCategory `json:"change_date_category,omitempty"`
	Twitter            string       `json:"twitter,omitempty"`
	Published          []int        `json:"published,omitempty"`
	Developed          []int        `json:"developed,omitempty"`
	Parent             int          `json:"parent,omitempty"`
	Facebook           string       `json:"facebook,omitempty"`
}

// Get returns a single Company identified by the provided IGDB ID. Provide
// the SetFields functional option if you need to specify which fields to
// retrieve. If the ID does not match any Companies, an error is returned.
func (cs *CompanyService) Get(id int, opts ...FuncOption) (*Company, error) {
	url, err := cs.client.singleURL(CompanyEndpoint, id, opts...)
	if err != nil {
		return nil, err
	}

	var com []Company

	err = cs.client.get(url, &com)
	if err != nil {
		return nil, err
	}

	return &com[0], nil
}

// List returns a list of Companies identified by the provided list of IGDB IDs.
// Provide functional options to sort, filter, and paginate the results. Omitting
// IDs will instead retrieve an index of Companies based solely on the provided
// options. Any ID that does not match a Company is ignored. If none of the IDs
// match a Company, an error is returned.
func (cs *CompanyService) List(ids []int, opts ...FuncOption) ([]*Company, error) {
	url, err := cs.client.multiURL(CompanyEndpoint, ids, opts...)
	if err != nil {
		return nil, err
	}

	var com []*Company

	err = cs.client.get(url, &com)
	if err != nil {
		return nil, err
	}

	return com, nil
}

func (cs *CompanyService) ListPaginated(limit int, opts ...FuncOption) (*Pagination, error) {
	return NewPaginationForEndpoint(cs.client, CompanyEndpoint, limit, opts...)
}

func (cs *CompanyService) GetPaginated(p *Pagination) ([]*Company, bool, error) {
	var c []*Company

	moreItems, err := p.Get(&c)
	if err != nil {
		return nil, false, err
	}

	return c, moreItems, nil
}

// Search returns a list of Companies found by searching the IGDB using the provided
// query. Provide functional options to sort, filter, and paginate the results. If
// no Companies are found using the provided query, an error is returned.
func (cs *CompanyService) Search(qry string, opts ...FuncOption) ([]*Company, error) {
	url, err := cs.client.searchURL(CompanyEndpoint, qry, opts...)
	if err != nil {
		return nil, err
	}

	var com []*Company

	err = cs.client.get(url, &com)
	if err != nil {
		return nil, err
	}

	return com, nil
}

// Count returns the number of Companies available in the IGDB.
// Provide the SetFilter functional option if you need to filter
// which Companies to count.
func (cs *CompanyService) Count(opts ...FuncOption) (int, error) {
	ct, err := cs.client.getEndpointCount(CompanyEndpoint, opts...)
	if err != nil {
		return 0, err
	}

	return ct, nil
}

// ListFields returns the up-to-date list of fields in an
// IGDB Company object.
func (cs *CompanyService) ListFields() ([]string, error) {
	fl, err := cs.client.getEndpointFieldList(CompanyEndpoint)
	if err != nil {
		return nil, err
	}

	return fl, nil
}
