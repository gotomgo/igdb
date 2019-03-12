package igdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// URL represents a URL as a string.
type URL string

// igdbURL is the base URL for the IGDB API.
const igdbURL string = "https://api-endpoint.igdb.com/"
const IgdbURL = igdbURL

// Errors returned when creating URLs for API calls.
var (
	// ErrNegativeID occurs when a negative ID is used as an argument in an API call.
	ErrNegativeID = errors.New("igdb.Client: negative ID")
	// ErrEmptyIDs occurs when an empty slice of IDs is used as an argument in an API call.
	ErrEmptyIDs = errors.New("igdb.Client: empty IDs")
	// ErrNoResults occurs when the IGDB returns no results
	ErrNoResults = errors.New("igdb.Client: no results")
)

// service is the underlying struct that handles
// all API calls for different IGDB endpoints.
type service struct {
	client *Client
}

// Client wraps an HTTP Client used to communicate with the IGDB,
// the root URL of the IGDB, and the user's IGDB API key.
// Client also initializes all the separate services to communicate
// with each individual IGDB API endpoint.
type Client struct {
	http    *http.Client
	rootURL string
	key     string
	isPro   bool

	common service

	// Services
	Characters   *CharacterService
	Collections  *CollectionService
	Companies    *CompanyService
	Credits      *CreditService
	Engines      *EngineService
	Feeds        *FeedService
	Franchises   *FranchiseService
	Games        *GameService
	GameModes    *GameModeService
	Genres       *GenreService
	Keywords     *KeywordService
	Pages        *PageService
	People       *PersonService
	Perspectives *PerspectiveService
	Platforms    *PlatformService
	Pulses       *PulseService
	PulseGroups  *PulseGroupService
	PulseSources *PulseSourceService
	ReleaseDates *ReleaseDateService
	Reviews      *ReviewService
	Themes       *ThemeService
	Titles       *TitleService
	Versions     *VersionService
}

// NewClientEx returns a new Client configured to communicate with the IGDB.
//
//	Notes
//		The provided apiKey will be used to make requests on your behalf. The
//		provided HTTP Client will be the client making requests to the IGDB.
//		If no HTTP Client is provided, a default HTTP client is used instead.
//
//		The isPro parameter can be used with pro subscriptions to use limit
//		values up to 3000.
//
//		If you need an IGDB API key, please visit: https://api.igdb.com/signup
func NewClientEx(apiURL string, apiKey string, isPro bool, custom *http.Client) *Client {
	if custom == nil {
		custom = http.DefaultClient
	}
	c := &Client{}
	c.http = custom
	c.key = apiKey
	c.isPro = isPro

	if isPro {
		c.rootURL = fmt.Sprintf("%spro/", apiURL)
	} else {
		c.rootURL = igdbURL
	}

	c.common.client = c
	c.Characters = (*CharacterService)(&c.common)
	c.Collections = (*CollectionService)(&c.common)
	c.Companies = (*CompanyService)(&c.common)
	c.Credits = (*CreditService)(&c.common)
	c.Engines = (*EngineService)(&c.common)
	c.Feeds = (*FeedService)(&c.common)
	c.Franchises = (*FranchiseService)(&c.common)
	c.Games = (*GameService)(&c.common)
	c.GameModes = (*GameModeService)(&c.common)
	c.Genres = (*GenreService)(&c.common)
	c.Keywords = (*KeywordService)(&c.common)
	c.Pages = (*PageService)(&c.common)
	c.People = (*PersonService)(&c.common)
	c.Perspectives = (*PerspectiveService)(&c.common)
	c.Platforms = (*PlatformService)(&c.common)
	c.Pulses = (*PulseService)(&c.common)
	c.PulseGroups = (*PulseGroupService)(&c.common)
	c.PulseSources = (*PulseSourceService)(&c.common)
	c.ReleaseDates = (*ReleaseDateService)(&c.common)
	c.Reviews = (*ReviewService)(&c.common)
	c.Themes = (*ThemeService)(&c.common)
	c.Titles = (*TitleService)(&c.common)
	c.Versions = (*VersionService)(&c.common)

	return c
}

// NewClient returns a new Client configured to communicate with the IGDB.
//
//	Notes
//		The provided apiKey will be used to make requests on your behalf. The
//		provided HTTP Client will be the client making requests to the IGDB.
//		If no HTTP Client is provided, a default HTTP client is used instead.
//
//		The isPro parameter is specified as false
//
//		If you need an IGDB API key, please visit: https://api.igdb.com/signup
//
func NewClient(apiKey string, custom *http.Client) *Client {
	return NewClientEx(igdbURL, apiKey, false, custom)
}

// getgetWithCallback sends a GET request to the provided url and stores
// the response in the provided result empty interface.
//
//	Notes
//		The response will be checked and return any errors.
//
//		If callback is specified, it will be invoked with the http Response
//		prior to reading the body of the response. Note, the callback should
//		not read the body, but is free to access other response info such
//		as the headers
//
func (c *Client) getBody(url string, callback func(resp *http.Response) error) ([]byte, error) {
	fmt.Println(url)
	req, err := c.newRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if callback != nil {
		err = callback(resp)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()

	if err = checkResponse(resp); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = checkResults(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) getWithCallback(url string, callback func(resp *http.Response) error, result interface{}) (err error) {
	var body []byte

	if body, err = c.getBody(url, callback); err == nil {
		err = json.Unmarshal(body, &result)
	}

	return err
}

// get sends a GET request to the provided url and stores
// the response in the provided result empty interface.
// The response will be checked and return any errors.
func (c *Client) get(url string, result interface{}) error {
	return c.getWithCallback(url, nil, result)
}

// newRequest configures a new request for the provided URL and
// adds the necesarry headers to communicate with the IGDB.
func (c *Client) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("user-key", c.key)
	req.Header.Add("Accept", "application/json")

	return req, nil
}

// singleURL creates a URL configured to request a single IGDB object identified by
// its unique IGDB ID using the provided endpoint and options.
func (c *Client) singleURL(end endpoint, ID int, opts ...FuncOption) (string, error) {
	if ID < 0 {
		return "", ErrNegativeID
	}
	opt, err := newOpt(opts...)
	if err != nil {
		return "", err
	}

	url := c.rootURL + string(end) + strconv.Itoa(ID)
	url = encodeURL(&opt.Values, url)

	return url, nil
}

// multiURL creates a URL configured to request multiple IGDB objects identified
// by their unique IGDB IDs using the provided endpoint and options. An empty slice
// of IDs creates a URL configured to retrieve an index of IGDB objects from the given
// endpoint based solely on the provided options.
func (c *Client) multiURL(end endpoint, IDs []int, opts ...FuncOption) (string, error) {
	for _, ID := range IDs {
		if ID < 0 {
			return "", ErrNegativeID
		}
	}

	opt, err := newOpt(opts...)
	if err != nil {
		return "", err
	}

	url := c.rootURL + string(end) + intsToCommaString(IDs)
	url = encodeURL(&opt.Values, url)

	return url, nil
}

// paginatedURL creates a URL configured to request IDGB objects using the
// provided options, and to use pagination
func (c *Client) paginatedURL(end endpoint, limit int, opts ...FuncOption) (string, error) {
	opt, err := newOpt(opts...)
	if err != nil {
		return "", err
	}

	if c.isPro {
		if limit > 3000 {
			limit = 3000
		}
	} else {
		if limit > 50 {
			limit = 50
		}
	}

	// overwrite any limit option already specified
	opt.Values.Set("limit", fmt.Sprintf("%d", limit))
	// set the scroll param
	opt.Values.Set("scroll", "1")

	url := c.rootURL + string(end)
	url = encodeURL(&opt.Values, url)

	return url, nil
}

// enumURL creates a URL configured to retrieve an index of IGDB objects from
// the given endpoint based solely on the provided options
//
//	Notes
//		This is a short-hand form of multiUrl(end,nil,opts)
//
func (c *Client) enumURL(end endpoint, opts ...FuncOption) (string, error) {
	return c.multiURL(end, nil, opts...)
}

// searchURL creates a URL configured to search the IGDB based on the given query using
// the provided endpoint and options.
func (c *Client) searchURL(end endpoint, qry string, opts ...FuncOption) (string, error) {
	if strings.TrimSpace(qry) == "" {
		return "", ErrEmptyQuery
	}

	opts = append(opts, setSearch(qry))
	opt, err := newOpt(opts...)
	if err != nil {
		return "", err
	}

	url := c.rootURL + string(end)
	url = encodeURL(&opt.Values, url)

	return url, nil
}

// countURL creates a URL configured to retrieve the count of IGDB objects
// using the provided endpoint and options.
func (c *Client) countURL(end endpoint, opts ...FuncOption) (string, error) {
	opt, err := newOpt(opts...)
	if err != nil {
		return "", err
	}

	url := c.rootURL + string(end) + "count"
	url = encodeURL(&opt.Values, url)

	return url, nil
}

// Byte representations of ASCII characters. Used for empty result checks.
const (
	// openBracketASCII represents the ASCII code for an open bracket.
	openBracketASCII = 91
	// closedBracketASCII represents the ASCII code for a closed bracket.
	closedBracketASCII = 93
)

// checkResults checks if the results of an API call are an empty array.
// If they are, an error is returned. Otherwise, nil is returned.
func checkResults(r []byte) error {
	if len(r) != 2 {
		return nil
	}

	if r[0] == openBracketASCII && r[1] == closedBracketASCII {
		return ErrNoResults
	}

	return nil
}

// IntsToStrings is a helper function that converts a slice of ints to a
// slice of strings. Useful for functions that require a variadic number
// of strings instead of ints such as SetFilter.
func IntsToStrings(ints []int) []string {
	var str []string
	for _, i := range ints {
		str = append(str, strconv.Itoa(i))
	}
	return str
}

// intsToCommaString is a helper function that returns a comma separated
// list of ints as a single string.
func intsToCommaString(ints []int) string {
	s := IntsToStrings(ints)
	return strings.Join(s, ",")
}
