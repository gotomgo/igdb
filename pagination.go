package igdb

import (
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
	// result is a *[] of some sort, get the underlying elem which is the slice
	rv := reflect.ValueOf(result).Elem()

	if rv.Kind() == reflect.Slice {
		p.totalRead += rv.Len()
	} else {
		panic(fmt.Errorf("Result type not slice, it is => %v", reflect.TypeOf(rv).Kind()))
	}

	return p.HasRemainingItems()
}

func (p *Pagination) start(result interface{}) (moreItems bool, err error) {
	err = p.client.getWithCallback(p.startURL, func(resp *http.Response) error {
		p.pageQuery = resp.Header.Get("x-next-page")
		fmt.Println(p.pageQuery)

		itemCount, err2 := strconv.ParseInt(resp.Header.Get("x-count"), 10, 32)
		if err2 != nil {
			return err2
		}

		p.setItemCount(int(itemCount))

		return nil
	}, result)

	if err == nil {
		moreItems = p.updateTotalRead(result)
		// fmt.Printf("itemCount: %d,totalRead: %d\n", p.itemCount, p.totalRead)
	}

	return
}

func (p *Pagination) Get(result interface{}) (moreItems bool, err error) {
	if len(p.pageQuery) == 0 {
		moreItems, err = p.start(result)
	} else {
		if p.HasRemainingItems() {
			// rootUrl contains trailing '/', pageQuery starts with '/'.
			// Remove one of them
			pageURL := fmt.Sprintf("%s%s", p.client.rootURL, p.pageQuery[1:])
			if err = p.client.get(pageURL, result); err == nil {
				moreItems = p.updateTotalRead(result)
				// fmt.Printf("itemCount: %d,totalRead: %d\n", p.itemCount, p.totalRead)
			}
		} else {
			err = ErrNoResults
		}
	}

	return
}
