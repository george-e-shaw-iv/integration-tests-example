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
	expectedLists := testdb.SeedLists(t, ts.a.db)
	expectedItems := testdb.SeedItems(t, ts.a.db, expectedLists)

	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		ExpectedBody []item.Item
		ExpectedCode int
	}{
		{
			Name:   "OK",
			ListID: expectedLists[0].ID,
			ExpectedBody: []item.Item{
				expectedItems[0],
				expectedItems[1],
			},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "NoContent",
			ListID:       expectedLists[2].ID,
			ExpectedBody: []item.Item{},
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "NotFound",
			// Using 0 for ListID because postgres serial type starts at 1 so 0 will never exist.
			ListID:       0,
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
				var items []item.Item
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
	expectedLists := testdb.SeedLists(t, ts.a.db)
	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		RequestBody  item.Item
		ExpectedCode int
	}{
		{
			Name:   "OK",
			ListID: expectedLists[0].ID,
			RequestBody: item.Item{
				Name:     "Foo",
				Quantity: 1,
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:   "NoName",
			ListID: expectedLists[0].ID,
			RequestBody: item.Item{
				Quantity: 1,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:   "LessThanOneQuantity",
			ListID: expectedLists[0].ID,
			RequestBody: item.Item{
				Name:     "Bar",
				Quantity: 0,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name: "NotFoundList",
			// Using 0 for ListID because postgres serial type starts at 1 so 0 will never exist.
			ListID: 0,
			RequestBody: item.Item{
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

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode == http.StatusCreated {
				var i item.Item
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
	expectedLists := testdb.SeedLists(t, ts.a.db)
	expectedItems := testdb.SeedItems(t, ts.a.db, expectedLists)

	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		ItemID       int
		ExpectedBody item.Item
		ExpectedCode int
	}{
		{
			Name:         "OK",
			ListID:       expectedLists[0].ID,
			ItemID:       expectedItems[0].ID,
			ExpectedBody: expectedItems[0],
			ExpectedCode: http.StatusOK,
		},
		{
			Name:   "NotFound",
			ListID: expectedLists[0].ID,
			// Using 0 for ItemID because postgres serial type starts at 1 so 0 will never exist.
			ItemID:       0,
			ExpectedBody: item.Item{},
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
				var i item.Item
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
	expectedLists := testdb.SeedLists(t, ts.a.db)
	expectedItems := testdb.SeedItems(t, ts.a.db, expectedLists)

	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		ItemID       int
		RequestBody  item.Item
		ExpectedCode int
	}{
		{
			Name:   "OK",
			ListID: expectedLists[0].ID,
			ItemID: expectedItems[0].ID,
			RequestBody: item.Item{
				Name:     "Foo",
				Quantity: 1,
			},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:   "NoName",
			ListID: expectedLists[0].ID,
			ItemID: expectedItems[0].ID,
			RequestBody: item.Item{
				Quantity: 1,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:   "LessThanOneQuantity",
			ListID: expectedLists[0].ID,
			ItemID: expectedItems[0].ID,
			RequestBody: item.Item{
				Name:     "Bar",
				Quantity: 0,
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name: "NotFoundList",
			// Using 0 for ListID because postgres serial type starts at 1 so 0 will never exist.
			ListID: 0,
			ItemID: expectedItems[0].ID,
			RequestBody: item.Item{
				Name:     "Bar",
				Quantity: 1,
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:   "NotFoundItem",
			ListID: expectedLists[0].ID,
			// Using 0 for ItemID because postgres serial type starts at 1 so 0 will never exist.
			ItemID: 0,
			RequestBody: item.Item{
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

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode == http.StatusOK {
				var i item.Item
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
	expectedLists := testdb.SeedLists(t, ts.a.db)
	expectedItems := testdb.SeedItems(t, ts.a.db, expectedLists)

	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		ItemID       int
		ExpectedCode int
	}{
		{
			Name:         "OK",
			ListID:       expectedLists[0].ID,
			ItemID:       expectedItems[0].ID,
			ExpectedCode: http.StatusNoContent,
		},
		{
			Name:   "NotFound",
			ListID: expectedLists[0].ID,
			// Using 0 for ItemID because postgres serial type starts at 1 so 0 will never exist.
			ItemID:       0,
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
