package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/testdb"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/web"
	"github.com/google/go-cmp/cmp"
)

func Test_getLists(t *testing.T) {
	// Test database needs reseeded after this test is ran because this test
	// removes lists from the database in order to finish testing
	defer ts.reseedDatabase(t)

	tests := []struct {
		Name         string
		ExpectedBody []list.Record
		ExpectedCode int
	}{
		{
			Name:         "OK",
			ExpectedBody: testdb.SeedLists,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "NoContent",
			ExpectedBody: nil,
			ExpectedCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		if test.ExpectedBody == nil {
			if err := testdb.Truncate(ts.a.db); err != nil {
				t.Errorf("error encountered truncating database: %v", err)
			}
		}

		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/list", nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			ts.a.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedBody != nil {
				var lists []list.Record
				resp := web.Response{
					Results: &lists,
				}

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if d := cmp.Diff(test.ExpectedBody, lists); d != "" {
					t.Errorf("unexpected difference in response body:\n%v", d)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}
