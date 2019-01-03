package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/item"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/testdb"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/web"
	"github.com/google/go-cmp/cmp"
)

func Test_getItems(t *testing.T) {
	tests := []struct {
		Name         string
		ListID       int
		ExpectedBody []item.Record
		ExpectedCode int
	}{
		{
			Name:   "OK",
			ListID: testdb.SeedLists[0].ID,
			ExpectedBody: []item.Record{
				testdb.SeedItems[0],
				testdb.SeedItems[1],
			},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "NoContent",
			ListID:       testdb.SeedLists[2].ID,
			ExpectedBody: nil,
			ExpectedCode: http.StatusNoContent,
		},
		{
			Name:         "NotFound",
			ListID:       0, // postgres serial starts at 1, 0 will never exist
			ExpectedBody: nil,
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/list/%d/item", test.ListID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedBody != nil {
				var items []item.Record
				resp := web.Response{
					Results: &items,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if d := cmp.Diff(test.ExpectedBody, items); d != "" {
					t.Errorf("unexpected difference in response body:\n%v", d)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_createItem(t *testing.T) {
	// Test database needs reseeded after this test is ran because this test
	// adds items to the database
	defer ts.reseedDatabase(t)

	tests := []struct {
		Name         string
		ListID       int
		RequestBody  item.Record
		ExpectedCode int
	}{
		{
			Name:   "OK",
			ListID: testdb.SeedLists[0].ID,
			RequestBody: item.Record{
				Name:     "Foo",
				Quantity: 1,
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:   "NoName",
			ListID: testdb.SeedLists[0].ID,
			RequestBody: item.Record{
				Quantity: 1,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:   "LessThanOneQuantity",
			ListID: testdb.SeedLists[0].ID,
			RequestBody: item.Record{
				Name:     "Bar",
				Quantity: 0,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:   "NotFoundList",
			ListID: 0, // postgres serial starts at 1, 0 will never exist
			RequestBody: item.Record{
				Name:     "Bar",
				Quantity: 1,
			},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/list/%d/item", test.ListID), &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			defer req.Body.Close()

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode == http.StatusCreated {
				var i item.Record
				resp := web.Response{
					Results: &i,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := test.RequestBody.Name, i.Name; e != a {
					t.Errorf("expected item name: %v, got item name: %v", e, a)
				}

				if e, a := test.RequestBody.Quantity, i.Quantity; e != a {
					t.Errorf("expected item quantity: %v, got item quantity: %v", e, a)
				}

				if e, a := test.ListID, i.ListID; e != a {
					t.Errorf("expected item list id: %v, got item list id: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_getItem(t *testing.T) {
	tests := []struct {
		Name         string
		ListID       int
		ItemID       int
		ExpectedBody item.Record
		ExpectedCode int
	}{
		{
			Name:         "OK",
			ListID:       testdb.SeedLists[0].ID,
			ItemID:       testdb.SeedItems[0].ID,
			ExpectedBody: testdb.SeedItems[0],
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "NotFound",
			ListID:       testdb.SeedLists[0].ID,
			ItemID:       0, // postgres serial starts at 1, 0 will never exist
			ExpectedBody: item.Record{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/list/%d/item/%d", test.ListID, test.ItemID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var i item.Record
				resp := web.Response{
					Results: &i,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if d := cmp.Diff(test.ExpectedBody, i); d != "" {
					t.Errorf("unexpected difference in response body:\n%v", d)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_updateItem(t *testing.T) {
	// Test database needs reseeded after this test is ran because this test
	// changes items in the database
	defer ts.reseedDatabase(t)

	tests := []struct {
		Name         string
		ListID       int
		ItemID       int
		RequestBody  item.Record
		ExpectedCode int
	}{
		{
			Name:   "OK",
			ListID: testdb.SeedLists[0].ID,
			ItemID: testdb.SeedItems[0].ID,
			RequestBody: item.Record{
				Name:     "Foo",
				Quantity: 1,
			},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:   "NoName",
			ListID: testdb.SeedLists[0].ID,
			ItemID: testdb.SeedItems[0].ID,
			RequestBody: item.Record{
				Quantity: 1,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:   "LessThanOneQuantity",
			ListID: testdb.SeedLists[0].ID,
			ItemID: testdb.SeedItems[0].ID,
			RequestBody: item.Record{
				Name:     "Bar",
				Quantity: 0,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:   "NotFoundList",
			ListID: 0, // postgres serial starts at 1, 0 will never exist
			ItemID: testdb.SeedItems[0].ID,
			RequestBody: item.Record{
				Name:     "Bar",
				Quantity: 1,
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:   "NotFoundItem",
			ListID: testdb.SeedLists[0].ID,
			ItemID: 0, // postgres serial starts at 1, 0 will never exist
			RequestBody: item.Record{
				Name:     "Bar",
				Quantity: 1,
			},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/list/%d/item/%d", test.ListID, test.ItemID), &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			defer req.Body.Close()

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode == http.StatusOK {
				var i item.Record
				resp := web.Response{
					Results: &i,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := test.RequestBody.Name, i.Name; e != a {
					t.Errorf("expected item name: %v, got item name: %v", e, a)
				}

				if e, a := test.RequestBody.Quantity, i.Quantity; e != a {
					t.Errorf("expected item quantity: %v, got item quantity: %v", e, a)
				}

				if e, a := test.ItemID, i.ID; e != a {
					t.Errorf("expected item id: %v, got item id: %v", e, a)
				}

				if e, a := test.ListID, i.ListID; e != a {
					t.Errorf("expected item list id: %v, got item list id: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_deleteItem(t *testing.T) {
	// Test database needs reseeded after this test is ran because this test
	// deletes items in the database
	defer ts.reseedDatabase(t)

	tests := []struct {
		Name         string
		ListID       int
		ItemID       int
		ExpectedCode int
	}{
		{
			Name:         "OK",
			ListID:       testdb.SeedLists[0].ID,
			ItemID:       testdb.SeedItems[0].ID,
			ExpectedCode: http.StatusNoContent,
		},
		{
			Name:         "NotFound",
			ListID:       testdb.SeedLists[0].ID,
			ItemID:       0, // postgres serial starts at 1, 0 will never exist
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/list/%d/item/%d", test.ListID, test.ItemID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}
		}

		t.Run(test.Name, fn)
	}
}
