package igdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type Pagination struct {
	client   *Client
	startURL string
	limit    int

	pageQuery  string
	itemCount  int
	totalRead  int
	totalPages int
}

func NewPaginationForEndpoint(client *Client, endpoint endpoint, limit int, opts ...FuncOption) (*Pagination, error) {
	startURL, err := client.paginatedURL(endpoint, limit, opts...)
	if err != nil {
		return nil, err
	}

	return NewPagination(client, startURL, limit), nil
}

func NewPagination(client *Client, startURL string, limit int) *Pagination {
	return &Pagination{client: client, startURL: startURL, limit: limit}
}

func (p *Pagination) RemainingItems() int {
	return p.itemCount - p.totalRead
}

func (p *Pagination) HasRemainingItems() bool {
	return p.itemCount > p.totalRead
}

func (p *Pagination) RemainingPages() int {
	return p.totalPages - p.getPagesRead()
}

func (p *Pagination) TotalPages() int {
	return p.totalPages
}

func (p *Pagination) setItemCount(itemCount int) {
	p.itemCount = itemCount
	p.totalPages = (itemCount + (p.limit - 1)) / p.limit
}

func (p *Pagination) getPagesRead() int {
	return p.getPageCount(p.itemCount - p.totalRead)
}

func (p *Pagination) getPageCount(itemCount int) int {
	return (itemCount + (p.limit - 1)) / p.limit
}

func (p *Pagination) updateTotalRead(result interface{}) bool {
	if result == nil {
		return false
	}

	// result is a *[] of some sort, get the underlying elem which is the slice
	rv := reflect.ValueOf(result).Elem()

	if rv.Kind() == reflect.Slice {
		p.totalRead += rv.Len()
	} else {
		p.totalRead += p.limit
		if p.totalRead > p.itemCount {
			p.totalRead = p.itemCount
		}

		// panic(fmt.Errorf("Result type not slice, it is => %v", reflect.TypeOf(rv).Kind()))
	}

	return p.HasRemainingItems()
}

func (p *Pagination) start() (body []byte, err error) {
	// on the 1st request we expect the x-next-page header, and an x-count header
	body, err = p.client.getBody(p.startURL, func(resp *http.Response) error {
		// the pagination query never changes and is good for about 3 minutes
		// (each call thru the query resets the timer)
		p.pageQuery = resp.Header.Get("x-next-page")

		xcount := resp.Header.Get("x-count")

		if len(xcount) > 0 {
			// We get the total count of the items to be paged
			itemCount, err2 := strconv.ParseInt(xcount, 10, 32)
			if err2 != nil {
				return err2
			}

			p.setItemCount(int(itemCount))
		} else {
			p.setItemCount(0)
		}

		return nil
	})

	return
}

func (p *Pagination) Get(result interface{}) (moreItems bool, err error) {
	// convet the []byte by marshaling it as JSON
	result, moreItems, err = p.GetWithCallback(func(b []byte) interface{} {
		err = json.Unmarshal(b, &result)
		return result
	})

	return
}

func (p *Pagination) GetRaw() (result []byte, moreItems bool, err error) {
	var temp interface{}

	// by not specifying a callback, the result is an []byte
	temp, moreItems, err = p.GetWithCallback(nil)
	if err != nil {
		return nil, false, err
	}

	return temp.([]byte), moreItems, nil
}

func (p *Pagination) GetWithCallback(callback func(body []byte) interface{}) (result interface{}, moreItems bool, err error) {
	var body []byte

	if len(p.pageQuery) == 0 {
		body, err = p.start()
	} else {
		if p.HasRemainingItems() {
			// rootUrl contains trailing '/', pageQuery starts with '/'.
			// Remove one of them
			pageURL := fmt.Sprintf("%s%s", p.client.rootURL, p.pageQuery[1:])
			body, err = p.client.getBody(pageURL, nil)
		} else {
			err = ErrNoResults
		}
	}

	if err != nil {
		return nil, false, err
	}

	if callback != nil {
		if len(body) > 0 {
			result = callback(body)
		}
	} else {
		if len(body) > 0 {
			result = body
		}
	}

	moreItems = p.updateTotalRead(result)

	return
}
