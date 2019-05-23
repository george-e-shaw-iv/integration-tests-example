package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/testdb"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/web"
	"github.com/google/go-cmp/cmp"
)

func Test_getLists(t *testing.T) {
	// No Content (no seed data)
	{
		req, err := http.NewRequest(http.MethodGet, "/list", nil)
		if err != nil {
			t.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		ts.a.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			t.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var lists []list.List
		resp := web.Response{
			Results: &lists,
		}

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Errorf("error decoding response body: %v", err)
		}

		if len(lists) > 0 {
			t.Errorf("expected no lists to be returned, got %v lists", len(lists))
		}
	}

	// Ok (database has been seeded)
	{
		expectedLists := testdb.SeedLists(t, ts.a.db)
		defer testdb.Truncate(t, ts.a.db)

		req, err := http.NewRequest(http.MethodGet, "/list", nil)
		if err != nil {
			t.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		ts.a.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			t.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var lists []list.List
		resp := web.Response{
			Results: &lists,
		}

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Errorf("error decoding response body: %v", err)
		}

		if d := cmp.Diff(expectedLists, lists); d != "" {
			t.Errorf("unexpected difference in response body:\n%v", d)
		}
	}
}

func Test_createList(t *testing.T) {
	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		RequestBody  list.List
		ExpectedCode int
	}{
		{
			Name: "OK",
			RequestBody: list.List{
				Name: "Foo",
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Name: "BreakUniqueNameConstraint",
			RequestBody: list.List{
				Name: "Foo",
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "NoName",
			RequestBody:  list.List{},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/list", &b)
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

			if test.ExpectedCode != http.StatusBadRequest {
				var l list.List
				resp := web.Response{
					Results: &l,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := test.RequestBody.Name, l.Name; e != a {
					t.Errorf("expected list name: %v, got list name: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_getList(t *testing.T) {
	expectedLists := testdb.SeedLists(t, ts.a.db)
	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		ExpectedBody list.List
		ExpectedCode int
	}{
		{
			Name:         "OK",
			ListID:       expectedLists[0].ID,
			ExpectedBody: expectedLists[0],
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "NotFound",
			// Using 0 for ListID because postgres serial type starts at 1 so 0 will never exist.
			ListID:       0,
			ExpectedBody: list.List{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/list/%d", test.ListID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var l list.List
				resp := web.Response{
					Results: &l,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if d := cmp.Diff(test.ExpectedBody, l); d != "" {
					t.Errorf("unexpected difference in response body:\n%v", d)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_updateList(t *testing.T) {
	expectedLists := testdb.SeedLists(t, ts.a.db)
	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		RequestBody  list.List
		ExpectedCode int
	}{
		{
			Name:   "OK",
			ListID: expectedLists[0].ID,
			RequestBody: list.List{
				Name: "Foo",
			},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:   "BreakUniqueNameConstraint",
			ListID: expectedLists[1].ID,
			RequestBody: list.List{
				Name: "Foo",
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "NoName",
			ListID:       expectedLists[0].ID,
			RequestBody:  list.List{},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name: "NotFound",
			// Using 0 for ListID because postgres serial type starts at 1 so 0 will never exist.
			ListID: 0,
			RequestBody: list.List{
				Name: "Bar",
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

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/list/%d", test.ListID), &b)
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
				var l list.List
				resp := web.Response{
					Results: &l,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := test.RequestBody.Name, l.Name; e != a {
					t.Errorf("expected list name: %v, got list name: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_deleteList(t *testing.T) {
	expectedLists := testdb.SeedLists(t, ts.a.db)
	defer testdb.Truncate(t, ts.a.db)

	tests := []struct {
		Name         string
		ListID       int
		ExpectedCode int
	}{
		{
			Name:         "OK",
			ListID:       expectedLists[0].ID,
			ExpectedCode: http.StatusNoContent,
		},
		{
			Name: "NotFound",
			// Using 0 for ListID because postgres serial type starts at 1 so 0 will never exist.
			ListID:       0,
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/list/%d", test.ListID), nil)
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
