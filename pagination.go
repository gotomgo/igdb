package igdb

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type Pagination struct {
	client   *Client
	endpoint endpoint
	limit    int
	options  []Option
	offset   int
}

// NewPaginationForEndpoint is @deprecated
func NewPaginationForEndpoint(client *Client, end endpoint, limit int, opts ...Option) (*Pagination, error) {
	return NewPagination(client, end, limit, opts...), nil
}

func NewPagination(client *Client, end endpoint, limit int, opts ...Option) *Pagination {
	return &Pagination{client: client, endpoint: end, limit: limit, options: opts}
}

func (p *Pagination) updateTotalRead(result interface{}) bool {
	if result == nil {
		return false
	}

	// result is a *[] of some sort, get the underlying elem which is the slice
	rv := reflect.ValueOf(result).Elem()

	var itemCount int
	if rv.Kind() == reflect.Slice {
		itemCount = rv.Len()
	} else {
		itemCount = p.limit
	}

	p.offset += itemCount

	fmt.Printf("pagination => items read: %d, total read/offset: %d\n", itemCount, p.offset)

	return itemCount >= p.limit
}

func (p *Pagination) Get(result interface{}) (moreItems bool, err error) {
	options := make([]Option, len(p.options)+2)
	options[0] = SetLimit(p.limit)
	options[1] = SetOffset(p.offset)
	for i, opt := range p.options {
		options[i+2] = opt
	}

	if err = p.client.get(p.endpoint, result, options...); err == nil {
		moreItems = p.updateTotalRead(result)
	} else if errors.Cause(err) == ErrNoResults {
		err = nil
	}

	return
}
