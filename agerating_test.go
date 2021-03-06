package igdb

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

const (
	testAgeRatingGet  string = "test_data/agerating_get.json"
	testAgeRatingList string = "test_data/agerating_list.json"
)

func TestAgeRatingService_Get(t *testing.T) {
	f, err := ioutil.ReadFile(testAgeRatingGet)
	if err != nil {
		t.Fatal(err)
	}

	init := make([]*AgeRating, 1)
	json.Unmarshal(f, &init)

	var tests = []struct {
		name          string
		file          string
		id            int
		opts          []Option
		wantAgeRating *AgeRating
		wantErr       error
	}{
		{"Valid response", testAgeRatingGet, 9644, []Option{SetFields("name")}, init[0], nil},
		{"Invalid ID", testFileEmpty, -1, nil, nil, ErrNegativeID},
		{"Empty response", testFileEmpty, 9644, nil, nil, errInvalidJSON},
		{"Invalid option", testFileEmpty, 9644, []Option{SetOffset(-99999)}, nil, ErrOutOfRange},
		{"No results", testFileEmptyArray, 0, nil, nil, ErrNoResults},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c, err := testServerFile(http.StatusOK, test.file)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()

			age, err := c.AgeRatings.Get(test.id, test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if !reflect.DeepEqual(age, test.wantAgeRating) {
				t.Errorf("got: <%v>, \nwant: <%v>", age, test.wantAgeRating)
			}
		})
	}
}

func TestAgeRatingService_List(t *testing.T) {
	f, err := ioutil.ReadFile(testAgeRatingList)
	if err != nil {
		t.Fatal(err)
	}

	init := make([]*AgeRating, 0)
	json.Unmarshal(f, &init)

	var tests = []struct {
		name           string
		file           string
		ids            []int
		opts           []Option
		wantAgeRatings []*AgeRating
		wantErr        error
	}{
		{"Valid response", testAgeRatingList, []int{9644, 40}, []Option{SetLimit(5)}, init, nil},
		{"Zero IDs", testFileEmpty, nil, nil, nil, ErrEmptyIDs},
		{"Invalid ID", testFileEmpty, []int{-500}, nil, nil, ErrNegativeID},
		{"Empty response", testFileEmpty, []int{9644, 40}, nil, nil, errInvalidJSON},
		{"Invalid option", testFileEmpty, []int{9644, 40}, []Option{SetOffset(-99999)}, nil, ErrOutOfRange},
		{"No results", testFileEmptyArray, []int{0, 9999999}, nil, nil, ErrNoResults},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c, err := testServerFile(http.StatusOK, test.file)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()

			age, err := c.AgeRatings.List(test.ids, test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if !reflect.DeepEqual(age, test.wantAgeRatings) {
				t.Errorf("got: <%v>, \nwant: <%v>", age, test.wantAgeRatings)
			}
		})
	}
}

func TestAgeRatingService_Index(t *testing.T) {
	f, err := ioutil.ReadFile(testAgeRatingList)
	if err != nil {
		t.Fatal(err)
	}

	init := make([]*AgeRating, 0)
	json.Unmarshal(f, &init)

	tests := []struct {
		name           string
		file           string
		opts           []Option
		wantAgeRatings []*AgeRating
		wantErr        error
	}{
		{"Valid response", testAgeRatingList, []Option{SetLimit(5)}, init, nil},
		{"Empty response", testFileEmpty, nil, nil, errInvalidJSON},
		{"Invalid option", testFileEmpty, []Option{SetOffset(-99999)}, nil, ErrOutOfRange},
		{"No results", testFileEmptyArray, nil, nil, ErrNoResults},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c, err := testServerFile(http.StatusOK, test.file)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()

			age, err := c.AgeRatings.Index(test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if !reflect.DeepEqual(age, test.wantAgeRatings) {
				t.Errorf("got: <%v>, \nwant: <%v>", age, test.wantAgeRatings)
			}
		})
	}
}

func TestAgeRatingService_Count(t *testing.T) {
	var tests = []struct {
		name      string
		resp      string
		opts      []Option
		wantCount int
		wantErr   error
	}{
		{"Happy path", `{"count": 100}`, []Option{SetFilter("popularity", OpGreaterThan, "75")}, 100, nil},
		{"Empty response", "", nil, 0, errInvalidJSON},
		{"Invalid option", "", []Option{SetLimit(-100)}, 0, ErrOutOfRange},
		{"No results", "[]", nil, 0, ErrNoResults},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c := testServerString(http.StatusOK, test.resp)
			defer ts.Close()

			count, err := c.AgeRatings.Count(test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if count != test.wantCount {
				t.Fatalf("got: <%v>, want: <%v>", count, test.wantCount)
			}
		})
	}
}

func TestAgeRatingService_Fields(t *testing.T) {
	var tests = []struct {
		name       string
		resp       string
		wantFields []string
		wantErr    error
	}{
		{"Happy path", `["name", "slug", "url"]`, []string{"url", "slug", "name"}, nil},
		{"Dot operator", `["logo.url", "background.id"]`, []string{"background.id", "logo.url"}, nil},
		{"Asterisk", `["*"]`, []string{"*"}, nil},
		{"Empty response", "", nil, errInvalidJSON},
		{"No results", "[]", nil, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c := testServerString(http.StatusOK, test.resp)
			defer ts.Close()

			fields, err := c.AgeRatings.Fields()
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if !equalSlice(fields, test.wantFields) {
				t.Fatalf("Expected fields '%v', got '%v'", test.wantFields, fields)
			}
		})
	}
}
