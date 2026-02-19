package jsonutil_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
)

func TestURL(t *testing.T) {
	type testURL struct {
		URL                 url.URL  `json:"url"`
		URLPointer          *url.URL `json:"urlPointer"`
		URLOmitZero        url.URL  `json:"urlOmitZero,omitzero"`
		URLPointerOmitEmpty *url.URL `json:"urlPointerOmitEmpty,omitempty"`
	}

	jsonOpts := json.JoinOptions(
		json.WithMarshalers(json.MarshalToFunc(jsonutil.URLMarshal)),
		json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.URLUnmarshal)),
	)

	t.Run("EOF", func(t *testing.T) {
		out := &testURL{}
		errSyn := &jsontext.SyntacticError{}

		if err := json.Unmarshal([]byte(`{"url":`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSyn) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if want := `unexpected EOF`; errSyn.Err.Error() != want {
			t.Fatalf("expected syntactic error be %s, got: %#v", want, errSyn.Err)
		}
	})

	t.Run("not a string", func(t *testing.T) {
		out := &testURL{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`{"url":3}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if tpInt := reflect.TypeFor[*url.URL](); errSem.GoType != tpInt {
			t.Fatalf("expected semantic error to have type %s, got: %s", tpInt, errSem.GoType)
		}
	})

	t.Run("parse error", func(t *testing.T) {
		out := &testURL{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`{"url":" http://example.org"}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if tpInt := reflect.TypeFor[*url.URL](); errSem.GoType != tpInt {
			t.Fatalf("expected semantic error to have type %s, got: %s", tpInt, errSem.GoType)
		}
	})

	t.Run("null", func(t *testing.T) {
		out := &testURL{}
		if err := json.Unmarshal([]byte(`{"url":null,"urlPointer":null,"urlOmitZero":null,"urlPointerOmitEmpty":null}`), out, jsonOpts); err != nil {
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
			testURL{URL: u, URLPointer: &u, URLOmitZero: u, URLPointerOmitEmpty: &u},
			`{"url":"https://example.com/path","urlPointer":"https://example.com/path","urlOmitZero":"https://example.com/path","urlPointerOmitEmpty":"https://example.com/path"}`,
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

			if got, want := out.URLOmitZero.String(), tc.in.URLOmitZero.String(); got != want {
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
