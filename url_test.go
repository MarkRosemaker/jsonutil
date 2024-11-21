package jsonutil_test

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/go-json-experiment/json"
)

func TestURL(t *testing.T) {
	type testURL struct {
		URL                 url.URL  `json:"url"`
		URLPointer          *url.URL `json:"urlPointer"`
		URLOmitEmpty        url.URL  `json:"urlOmitEmpty,omitempty"`
		URLPointerOmitEmpty *url.URL `json:"urlPointerOmitEmpty,omitempty"`
	}

	jsonOpts := json.JoinOptions(
		json.WithMarshalers(json.MarshalFuncV2(jsonutil.URLMarshal)),
		json.WithUnmarshalers(json.UnmarshalFuncV2(jsonutil.URLUnmarshal)),
	)

	t.Run("errors", func(t *testing.T) {
		tpURL := reflect.TypeOf(&url.URL{})

		for _, data := range []string{
			`{"url":`,                       // EOF
			`{"url":3}`,                     // not a string
			`{"url":" http://example.org"}`, // parse error
		} {
			t.Run(data, func(t *testing.T) {
				out := &testURL{}
				errSem := &json.SemanticError{}
				if err := json.Unmarshal([]byte(data), out, jsonOpts); err == nil {
					t.Fatalf("expected error")
				} else if !errors.As(err, &errSem) {
					t.Fatalf("expected error to be a semantic error, got: %v", err)
				} else if errSem.GoType != tpURL {
					t.Fatalf("expected semantic error to have type %s, got: %s", tpURL, errSem.GoType)
				}
			})
		}
	})

	t.Run("null", func(t *testing.T) {
		out := &testURL{}
		if err := json.Unmarshal([]byte(`{"url":null,"urlPointer":null,"urlOmitEmpty":null,"urlPointerOmitEmpty":null}`), out, jsonOpts); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	u := url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/path",
	}

	for i, tc := range []struct {
		in  testURL
		out string
	}{
		{testURL{}, `{"url":"","urlPointer":null}`},
		{
			testURL{URL: u, URLPointer: &u, URLOmitEmpty: u, URLPointerOmitEmpty: &u},
			`{"url":"https://example.com/path","urlPointer":"https://example.com/path","urlOmitEmpty":"https://example.com/path","urlPointerOmitEmpty":"https://example.com/path"}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := json.Marshal(tc.in, jsonOpts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(b) != tc.out {
				t.Fatalf("want: %s, got: %s", tc.out, string(b))
			}

			var out testURL
			if err := json.Unmarshal(b, &out, jsonOpts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := out.URL.String(), tc.in.URL.String(); got != want {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			if got, want := out.URLOmitEmpty.String(), tc.in.URLOmitEmpty.String(); got != want {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			for _, tt := range []struct{ got, want *url.URL }{
				{out.URLPointer, tc.in.URLPointer},
				{out.URLPointerOmitEmpty, tc.in.URLPointerOmitEmpty},
			} {
				if tt.got == nil && tt.want != nil {
					t.Fatalf("expected non-nil URLPointer")
				} else if tt.got != nil && tt.want == nil {
					t.Fatalf("expected nil URLPointer")
				} else if tt.got != nil && tt.want != nil {
					if tt.got.String() != tt.want.String() {
						t.Fatalf("want: %s, got: %s", tt.want, tt.got)
					}
				}
			}
		})
	}
}
