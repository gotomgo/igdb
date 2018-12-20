package igdb

// PulseService handles all the API
// calls for the IGDB Pulse endpoint.
type PulseService service

// Pulse contains information on an IGDB
// entry for a single news article.
//
// For more information, visit: https://igdb.github.io/api/endpoints/pulse/
type Pulse struct {
	ID        int `json:"id"`
	URL       URL `json:"url"`
	CreatedAt int `json:"created_at"` // Unix time in milliseconds
	UpdatedAt int `json:"updated_at"` // Unix time in milliseconds

	PulseSource int          `json:"pulse_source"`
	Category    int          `json:"category"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	UID         string       `json:"uid"`          //perhaps switch to ID
	PublishedAt int          `json:"published_at"` // Unix time in milliseconds
	ImageURL    URL          `json:"image"`
	PulseImage  Image        `json:"pulse_image"`
	Videos      []PulseVideo `json:"videos"`
	Author      string       `json:"author"`
	Tags        []Tag        `json:"tags"`
	Ignored     interface{}  `json:"ignored"`
}

// PulseVideo contains the ID and category
// for a video related to a pulse.
type PulseVideo struct {
	Category int    `json:"category"`
	ID       string `json:"video_id"`
}

// Get returns a single Pulse identified by the provided IGDB ID. Provide
// the SetFields functional option if you need to specify which fields to
// retrieve. If the ID does not match any Pulses, an error is returned.
func (ps *PulseService) Get(id int, opts ...FuncOption) (*Pulse, error) {
	url, err := ps.client.singleURL(PulseEndpoint, id, opts...)
	if err != nil {
		return nil, err
	}

	var p []Pulse

	err = ps.client.get(url, &p)
	if err != nil {
		return nil, err
	}

	return &p[0], nil
}

// List returns a list of Pulses identified by the provided list of IGDB IDs.
// Provide functional options to sort, filter, and paginate the results. Omitting
// IDs will instead retrieve an index of Pulses based solely on the provided
// options. Any ID that does not match a Pulse is ignored. If none of the IDs
// match a Pulse, an error is returned.
func (ps *PulseService) List(ids []int, opts ...FuncOption) ([]*Pulse, error) {
	url, err := ps.client.multiURL(PulseEndpoint, ids, opts...)
	if err != nil {
		return nil, err
	}

	var p []*Pulse

	err = ps.client.get(url, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Search returns a list of Pulses found by searching the IGDB using the provided
// query. Provide functional options to sort, filter, and paginate the results. If
// no Pulses are found using the provided query, an error is returned.
func (ps *PulseService) Search(qry string, opts ...FuncOption) ([]*Pulse, error) {
	url, err := ps.client.searchURL(PulseEndpoint, qry, opts...)
	if err != nil {
		return nil, err
	}

	var p []*Pulse

	err = ps.client.get(url, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Count returns the number of Pulses available in the IGDB.
// Provide the SetFilter functional option if you need to filter
// which Pulses to count.
func (ps *PulseService) Count(opts ...FuncOption) (int, error) {
	ct, err := ps.client.getEndpointCount(PulseEndpoint, opts...)
	if err != nil {
		return 0, err
	}

	return ct, nil
}

// ListFields returns the up-to-date list of fields in an
// IGDB Pulse object.
func (ps *PulseService) ListFields() ([]string, error) {
	fl, err := ps.client.getEndpointFieldList(PulseEndpoint)
	if err != nil {
		return nil, err
	}

	return fl, nil
}
