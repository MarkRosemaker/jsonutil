package jsonutil_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"maps"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
)

func TestHTTPHeader(t *testing.T) {
	jsonOpts := json.JoinOptions(
		json.WithMarshalers(json.MarshalToFunc(jsonutil.HTTPHeaderMarshal)),
		json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.HTTPHeaderUnmarshal)),
	)

	t.Run("EOF", func(t *testing.T) {
		out := &http.Header{}
		errSyn := &jsontext.SyntacticError{}

		if err := json.Unmarshal([]byte(`{"foo":`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSyn) {
			t.Fatalf("expected error to be a syntactic error, got: %v", err)
		} else if want := `unexpected EOF`; errSyn.Err.Error() != want {
			t.Fatalf("expected syntactic error be %s, got: %#v", want, errSyn.Err)
		}
	})

	t.Run("not a map", func(t *testing.T) {
		out := &http.Header{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`"foo"`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %T", err)
		} else if want := `expected begin object, got string`; errSem.Err.Error() != want {
			t.Fatalf("expected semantic error be %s, got: %#v", want, errSem.Err)
		}
	})

	t.Run("not a string", func(t *testing.T) {
		out := &http.Header{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`{"foo": 3}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if tpHeader := reflect.TypeFor[*http.Header](); errSem.GoType != tpHeader {
			t.Fatalf("expected semantic error to have type %s, got: %s", tpHeader, errSem.GoType)
		}
	})

	for i, tc := range []struct {
		in  http.Header
		out string
	}{
		{nil, `null`},
		{http.Header{}, `{}`},
		{http.Header{"foo": []string{"bar"}, "baz": []string{"quux"}}, `{"Baz":"quux","Foo":"bar"}`},
		{http.Header{"foo": nil}, `{}`},
		{http.Header{"foo": []string{}}, `{}`},
		{http.Header{"foo": []string{""}}, `{}`},
		{http.Header{"foo": []string{"bar", "baz"}}, `{"Foo":"bar"}`},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := json.Marshal(tc.in, jsonOpts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(b) != tc.out {
				t.Fatalf("want: %s, got: %s", tc.out, string(b))
			}

			out := http.Header{}
			if err := json.Unmarshal(b, &out, jsonOpts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.in == nil && out != nil {
				t.Fatalf("if in is nil, out should also be nil, got: %v", out)
			}

			want := http.Header{}
			for key, val := range tc.in {
				if len(val) > 0 && val[0] != "" {
					want.Set(key, val[0])
				}
			}

			if !maps.EqualFunc(want, out, slices.Equal) {
				t.Fatalf("want: %v, got: %v", want, out)
			}
		})
	}
}
