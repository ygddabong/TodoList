package todo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPing(t *testing.T) {
	rec := do(t, NewHandler(NewStore()), http.MethodGet, "/ping", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "pong") {
		t.Fatalf("got %d %s", rec.Code, rec.Body)
	}
}

func TestCRUD(t *testing.T) {
	h := NewHandler(NewStore())

	rec := do(t, h, http.MethodPost, "/todos", `{"title":"  buy beans  "}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: got %d %s", rec.Code, rec.Body)
	}
	var created Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Title != "buy beans" || created.Done {
		t.Fatalf("unexpected todo: %+v", created)
	}
	if loc := rec.Header().Get("Location"); loc != "/todos/1" {
		t.Fatalf("Location = %q", loc)
	}

	rec = do(t, h, http.MethodPatch, "/todos/1", `{"done":true}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"done":true`) {
		t.Fatalf("update: got %d %s", rec.Code, rec.Body)
	}

	rec = do(t, h, http.MethodGet, "/todos", "")
	var list []Todo
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list) != 1 || !list[0].Done || list[0].Title != "buy beans" {
		t.Fatalf("list: %+v", list)
	}

	if rec = do(t, h, http.MethodDelete, "/todos/1", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: got %d", rec.Code)
	}
	if rec = do(t, h, http.MethodGet, "/todos/1", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("get after delete: got %d", rec.Code)
	}
}

func TestValidation(t *testing.T) {
	h := NewHandler(NewStore())
	cases := []struct {
		name, method, path, body string
		want                     int
	}{
		{"missing title", http.MethodPost, "/todos", `{}`, http.StatusBadRequest},
		{"blank title", http.MethodPost, "/todos", `{"title":"   "}`, http.StatusBadRequest},
		{"bad json", http.MethodPost, "/todos", `{`, http.StatusBadRequest},
		{"unknown field", http.MethodPost, "/todos", `{"title":"x","foo":1}`, http.StatusBadRequest},
		{"bad id", http.MethodGet, "/todos/abc", "", http.StatusBadRequest},
		{"missing id", http.MethodPatch, "/todos/42", `{"done":true}`, http.StatusNotFound},
		{"wrong method", http.MethodPut, "/todos/1", `{}`, http.StatusMethodNotAllowed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if rec := do(t, h, tc.method, tc.path, tc.body); rec.Code != tc.want {
				t.Fatalf("got %d, want %d (%s)", rec.Code, tc.want, rec.Body)
			}
		})
	}
}
